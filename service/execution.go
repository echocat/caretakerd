package service

import (
	"github.com/echocat/caretakerd/access"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/keyStore"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/sync"
	. "github.com/echocat/caretakerd/values"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Execution struct {
	service   *Service
	logger    *logger.Logger
	cmd       *exec.Cmd
	status    Status
	lock      *sync.Mutex
	condition *sync.Condition
	access    *access.Access
	syncGroup *sync.SyncGroup
}

func (instance *Service) NewExecution(sec *keyStore.KeyStore) (*Execution, error) {
	a, err := access.NewAccess(instance.config.Access, instance.name, sec)
	if err != nil {
		return nil, errors.New("Could not create caretakerd base execution.").CausedBy(err)
	}
	syncGroup := instance.syncGroup.NewSyncGroup()
	cmd := instance.generateCmd(a)
	lock := syncGroup.NewMutex()
	condition := syncGroup.NewCondition(lock)
	return &Execution{
		service:   instance,
		logger:    instance.logger,
		cmd:       cmd,
		status:    Down,
		lock:      lock,
		condition: condition,
		access:    a,
		syncGroup: syncGroup,
	}, nil
}

func (instance *Service) expandValue(ai *access.Access, in string) string {
	return os.Expand(in, func(key string) string {
		if value, ok := (*instance).config.Environment[key]; ok {
			return value
		} else if key == "CTD_PEM" {
			if ai.Type() == access.GenerateToEnvironment {
				return string(ai.Pem())
			} else {
				return ""
			}
		} else {
			return os.Getenv(key)
		}
	})
}

func (instance *Service) getRunArgumentsFor(ai *access.Access) []string {
	args := []string{}
	config := (*instance).config
	command := config.Command
	for i := 1; i < len(command); i++ {
		args = append(args, instance.expandValue(ai, command[i].String()))
	}
	return args
}

func (instance *Service) generateCmd(ai *access.Access) *exec.Cmd {
	logger := (*instance).logger
	config := (*instance).config
	executable := instance.expandValue(ai, config.Command[0].String())
	cmd := exec.Command(executable, instance.getRunArgumentsFor(ai)...)
	cmd.Stdout = logger.Stdout()
	cmd.Stderr = logger.Stderr()
	cmd.Stdin = logger.Stdin()
	cmd.SysProcAttr = instance.createSysProcAttr()
	if !config.Directory.IsTrimmedEmpty() {
		cmd.Dir = instance.expandValue(ai, config.Directory.String())
	}
	for key, value := range config.Environment {
		cmd.Env = append(cmd.Env, key+"="+value)
	}
	if ai.Type() == access.GenerateToEnvironment {
		cmd.Env = append(cmd.Env, "CTD_PEM="+string(ai.Pem()))
	} else {
		cmd.Env = append(cmd.Env, "CTD_PEM=")
	}
	if config.InheritEnvironment {
		cmd.Env = append(cmd.Env, os.Environ()...)
	}
	serviceHandleUsersFor(instance, cmd)
	return cmd
}

func (instance *Execution) CommandLine() string {
	service := (*instance).service
	result := ""
	for i, arg := range (*service).config.Command {
		argAsString := arg.String()
		if i != 0 {
			result += " "
		}
		if strings.Contains(argAsString, "\"") || strings.Contains(argAsString, "\\") || strings.Contains(argAsString, " ") {
			result += strconv.Quote(argAsString)
		} else {
			result += argAsString
		}
	}
	return result
}

func (instance *Execution) handleBeforeRun() error {
	cronExpression := instance.service.Config().CronExpression
	startAt := cronExpression.Next(time.Now())
	if startAt != nil {
		waitDuration := startAt.Sub(time.Now())
		instance.logger.Log(logger.Debug, "Start of service '%s' is timed for %v (in %v).", instance.Name(), startAt, waitDuration)
		if err := instance.syncGroup.Sleep(waitDuration); err != nil {
			return StoppedOrKilledError{error: errors.New("Process was stopped before start.")}
		}
	}
	return nil
}

func (instance *Execution) Run() (ExitCode, error) {
	err := instance.handleBeforeRun()
	if err != nil {
		return ExitCode(1), err
	}
	instance.logger.Log(logger.Debug, "Start service '%s' with command: %s", instance.Name(), instance.CommandLine())
	exitCode, err, lastState := instance.runBare()
	if lastState == Killed {
		err = StoppedOrKilledError{error: errors.New("Process was killed.")}
		instance.logger.Log(logger.Debug, "Service '%s' ended after kill: %d", instance.Name(), exitCode)
	} else if lastState == Stopped {
		err = StoppedOrKilledError{error: errors.New("Process was stopped.")}
		instance.logger.Log(logger.Debug, "Service '%s' ended successful after stop: %d", instance.Name(), exitCode)
	} else if err != nil {
		instance.logger.Log(logger.Fatal, err)
	} else if instance.service.config.SuccessExitCodes.Contains(exitCode) {
		instance.logger.Log(logger.Debug, "Service '%s' ended successful: %d", instance.Name(), exitCode)
	} else {
		instance.logger.Log(logger.Error, "Service '%s' ended with unexpected code: %d", instance.Name(), exitCode)
		err = errors.New("Unexpected error code %d generated by service '%s'", exitCode, instance.Name())
	}
	return exitCode, err
}

type UnrecoverableError struct {
	error
}

type StoppedOrKilledError struct {
	error
}

func (instance *Execution) runBare() (ExitCode, error, Status) {
	cmd := (*instance).cmd
	var waitStatus syscall.WaitStatus
	instance.doSetRunningState()
	defer instance.doSetDownState()
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			return ExitCode(waitStatus.ExitStatus()), nil, instance.status
		} else {
			return ExitCode(0), UnrecoverableError{error: err}, instance.status
		}
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		return ExitCode(waitStatus.ExitStatus()), nil, instance.status
	}
}

func (instance *Execution) doSetRunningState() {
	instance.doLock()
	defer instance.doUnlock()
	(*instance).status = Running
}

func (instance *Execution) doSetDownState() {
	instance.setStateSyncedTo(Down)
}

func (instance *Execution) setStateSyncedTo(s Status) bool {
	instance.doLock()
	defer instance.doUnlock()
	return instance.setStateTo(s)
}

func (instance *Execution) setStateTo(ns Status) bool {
	cs := instance.status
	if cs != Down || (ns != Killed && ns != Stopped) {
		(*instance).status = ns
		instance.condition.Send()
		if cs == Down {
			instance.access.Cleanup()
		}
		return true
	}
	instance.condition.Send()
	return false
}

func (instance *Execution) Name() string {
	return (*instance).service.name
}

func (instance *Execution) Stop() {
	instance.syncGroup.Interrupt()
	instance.doLock()
	defer instance.doUnlock()
	if instance.status != Down {
		instance.sendStop()
		if instance.status != Down {
			instance.logger.Log(logger.Debug, "Stopping '%s'...", instance.Name())
			instance.condition.Wait(time.Duration(instance.service.config.StopWaitInSeconds) * time.Second)
			if instance.status != Down {
				instance.logger.Log(logger.Warning, "Service '%s' does not respond after %d seconds. Going to kill it now...", instance.Name(), instance.service.config.StopWaitInSeconds)
				instance.sendKill()
			}
		}
	}
}

func (instance *Execution) sendStop() {
	if instance.status != Killed && instance.status != Stopped && instance.setStateTo(Stopped) {
		instance.sendSignal((*instance).service.config.StopSignal)
	}
}

func (instance *Execution) Kill() {
	instance.syncGroup.Interrupt()
	instance.doLock()
	defer instance.doUnlock()
	instance.logger.Log(logger.Debug, "Killing '%s'...", instance.Name())
	instance.sendKill()
}

func (instance *Execution) sendKill() {
	if instance.status != Killed && instance.setStateTo(Killed) {
		for instance.status != Down {
			if err := instance.sendSignal(KILL); err != nil {
				instance.logger.LogProblem(err, logger.Warning, "Could not kill: %v", instance.service.Name())
			}
			if instance.status != Down {
				(*instance).condition.Wait(1 * time.Second)
			}
		}
	}
}

func (instance *Execution) Signal(what Signal) error {
	instance.doLock()
	defer instance.doUnlock()
	instance.logger.Log(logger.Debug, "Sending signal %v to '%s'...", what, instance.Name())
	return instance.sendSignal(what)
}

func (instance *Execution) sendSignal(s Signal) error {
	if instance.isKillSignal(s) {
		if !instance.setStateTo(Killed) {
			return errors.New("Service '%v' is not running.", instance)
		}
	} else if instance.isStopSignal(s) {
		if !instance.setStateTo(Stopped) {
			return errors.New("Service '%v' is not running.", instance)
		}
	}
	cmd := (*instance).cmd
	process := cmd.Process
	ps := cmd.ProcessState
	if process == nil || ps != nil {
		instance.setStateTo(Down)
		return errors.New("Service '%v' is not running.", instance)
	}
	if s != NOOP {
		return sendSignalToService((*instance).service, process, s)
	}
	return nil
}

func (instance Execution) isStopSignal(s Signal) bool {
	return s == instance.service.config.StopSignal
}

func (instance Execution) isKillSignal(s Signal) bool {
	return s == KILL
}

func (instance *Execution) doLock() {
	instance.lock.Lock()
}

func (instance *Execution) doUnlock() {
	instance.lock.Unlock()
}

func (instance *Execution) Pid() int {
	instance.doLock()
	defer instance.doUnlock()
	cmd := instance.cmd
	if cmd != nil {
		process := cmd.Process
		if process != nil {
			return process.Pid
		}
	}
	return 0
}

func (instance *Execution) Status() Status {
	instance.doLock()
	defer instance.doUnlock()
	return instance.status
}

func (instance Execution) Service() *Service {
	return instance.service
}
