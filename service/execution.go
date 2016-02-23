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
	cmd := generateServiceBasedCmd(instance, a, (*instance).config.Command)
	lock := syncGroup.NewMutex()
	condition := syncGroup.NewCondition(lock)
	return &Execution{
		service:   instance,
		logger:    instance.logger,
		cmd:       cmd,
		status:    New,
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

func getServiceBasedRunArgumentsFor(s *Service, ai *access.Access, command []String) []string {
	args := []string{}
	for i := 1; i < len(command); i++ {
		args = append(args, s.expandValue(ai, command[i].String()))
	}
	return args
}

func generateServiceBasedCmd(s *Service, ai *access.Access, command []String) *exec.Cmd {
	logger := (*s).logger
	config := (*s).config
	executable := s.expandValue(ai, command[0].String())
	cmd := exec.Command(executable, getServiceBasedRunArgumentsFor(s, ai, command)...)
	cmd.Stdout = logger.Stdout()
	cmd.Stderr = logger.Stderr()
	cmd.Stdin = logger.Stdin()
	cmd.SysProcAttr = s.createSysProcAttr()
	if !config.Directory.IsTrimmedEmpty() {
		cmd.Dir = s.expandValue(ai, config.Directory.String())
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
	serviceHandleUsersFor(s, cmd)
	return cmd
}

func (instance *Execution) generateCmd(command []String) *exec.Cmd {
	return generateServiceBasedCmd(instance.service, instance.access, command)
}

func (instance *Execution) extractCommandProperties(command []String) (cleanCommand []String, handleErrors bool) {
	if len(command) > 0 && command[0] == "-" {
		return command[1:], false
	} else {
		return command, true
	}
}

func (instance *Execution) commandLineOf(cmd *exec.Cmd) string {
	result := ""
	for i, arg := range cmd.Args {
		if i != 0 {
			result += " "
		}
		if strings.Contains(arg, "\"") || strings.Contains(arg, "\\") || strings.Contains(arg, " ") {
			result += strconv.Quote(arg)
		} else {
			result += arg
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

func (instance *Execution) preExecution() (ExitCode, error) {
	preCommands := instance.service.config.PreCommands
	for _, preCommand := range preCommands {
		command, handleErrors := instance.extractCommandProperties(preCommand)
		if len(command) > 0 {
			cmd := instance.generateCmd(command)
			instance.logger.Log(logger.Debug, "Execute pre command: %s", instance.commandLineOf(cmd))
			exitCode, err := instance.runCommand(cmd)
			if handleErrors {
				if err != nil {
					instance.logger.LogProblem(err, logger.Error, "Pre command failed.")
					return exitCode, err
				} else if exitCode != 0 {
					instance.logger.Log(logger.Error, "Pre command failed. Exit with unexpected exit code: %d", exitCode)
					return exitCode, err
				}
			}
		}
	}
	return ExitCode(0), nil
}

func (instance *Execution) postExecution() {
	postCommands := instance.service.config.PostCommands
	for _, preCommand := range postCommands {
		command, handleErrors := instance.extractCommandProperties(preCommand)
		if len(command) > 0 {
			cmd := instance.generateCmd(command)
			instance.logger.Log(logger.Debug, "Execute post command: %s", instance.commandLineOf(cmd))
			exitCode, err := instance.runCommand(cmd)
			if handleErrors {
				if err != nil {
					instance.logger.LogProblem(err, logger.Warning, "Post command failed.")
				} else if exitCode != 0 {
					instance.logger.Log(logger.Warning, "Post command failed. Exit with unexpected exit code: %d", exitCode)
				}
			}
		}
	}
}

func (instance *Execution) Run() (ExitCode, error) {
	err := instance.handleBeforeRun()
	if err != nil {
		return ExitCode(1), err
	}
	exitCode, err := instance.preExecution()
	if err != nil || exitCode != 0 {
		return exitCode, err
	}
	instance.logger.Log(logger.Debug, "Start service '%s' with command: %s", instance.Name(), instance.commandLineOf(instance.cmd))
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
	instance.postExecution()
	return exitCode, err
}

type UnrecoverableError struct {
	error
}

type StoppedOrKilledError struct {
	error
}

func (instance *Execution) runCommand(cmd *exec.Cmd) (ExitCode, error) {
	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			exitSignal := waitStatus.Signal()
			if exitSignal > 0 {
				return ExitCode(int(exitSignal) + 128), nil
			} else {
				return ExitCode(waitStatus.ExitStatus()), nil
			}
		} else {
			return ExitCode(0), UnrecoverableError{error: err}
		}
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		return ExitCode(waitStatus.ExitStatus()), nil
	}
}

func (instance *Execution) runBare() (ExitCode, error, Status) {
	if instance.doTrySetRunningState() {
		defer instance.doSetDownState()
		exitCode, err := instance.runCommand((*instance).cmd)
		return exitCode, err, instance.status
	} else {
		return ExitCode(0), UnrecoverableError{error: errors.New("Cannot run service. Already in status: %v", instance.status)}, instance.status
	}
}

func (instance *Execution) doTrySetRunningState() bool {
	if instance.doLock() != nil {
		return false
	}
	defer instance.doUnlock()
	if (*instance).status == New {
		(*instance).status = Running
		return true
	}
	return false
}

func (instance *Execution) doSetDownState() {
	instance.setStateSyncedTo(Down)
}

func (instance *Execution) setStateSyncedTo(s Status) bool {
	if instance.doLock() != nil {
		return false
	}
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
	if instance.doLock() != nil {
		return
	}
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
		c := (*instance).service.config
		stopCommand, handleErrors := instance.extractCommandProperties(c.StopCommand)
		if len(stopCommand) > 0 {
			cmd := instance.generateCmd(stopCommand)
			instance.logger.Log(logger.Debug, "Execute stop command: %s", instance.commandLineOf(cmd))
			exitCode, err := instance.runCommand(cmd)
			if handleErrors {
				if err != nil {
					instance.logger.LogProblem(err, logger.Warning, "Stop command failed.")
				} else if exitCode != 0 {
					instance.logger.Log(logger.Warning, "Stop command failed. Exit with unexpected exit code: %d", exitCode)
				}
			}
		} else {
			instance.sendSignal(c.StopSignal)
		}
	}
}

func (instance *Execution) Kill() error {
	instance.syncGroup.Interrupt()
	if err := instance.doLock(); err != nil {
		return err
	}
	defer instance.doUnlock()
	instance.logger.Log(logger.Debug, "Killing '%s'...", instance.Name())
	instance.sendKill()
	return nil
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
	if err := instance.doLock() ; err != nil {
		return err
	}
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

func (instance *Execution) doLock() error {
	return instance.lock.Lock()
}

func (instance *Execution) doUnlock() {
	instance.lock.Unlock()
}

func (instance *Execution) Pid() int {
	if instance.doLock() != nil {
		return 0
	}
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
	if instance.doLock() != nil {
		return Unknown
	}
	defer instance.doUnlock()
	return instance.status
}

func (instance Execution) Service() *Service {
	return instance.service
}

func (instance Execution) String() string {
	return instance.service.String()
}

func (instance *Execution) SyncGroup() *sync.SyncGroup {
	return instance.syncGroup
}

