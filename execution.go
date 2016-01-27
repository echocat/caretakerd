package caretakerd

import (
	ssync "sync"
	"time"
	. "github.com/echocat/caretakerd/values"
	"github.com/echocat/caretakerd/service"
	"github.com/echocat/caretakerd/keyStore"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/sync"
	"github.com/echocat/caretakerd/errors"
)

type Executable interface {
	Services() *service.Services
	KeyStore() *keyStore.KeyStore
	Logger() *logger.Logger
}

type Execution struct {
	executable      Executable
	executions      map[*service.Service]*service.Execution
	restartRequests map[*service.Service]bool
	masterExitCode  *ExitCode
	masterError     error
	lock            *ssync.RWMutex
	wg              *ssync.WaitGroup
	syncGroup       *sync.SyncGroup
}

func NewExecution(executable Executable, syncGroup *sync.SyncGroup) *Execution {
	return &Execution{
		executable: executable,
		executions: map[*service.Service]*service.Execution{},
		restartRequests: map[*service.Service]bool{},
		lock: new(ssync.RWMutex),
		wg: new(ssync.WaitGroup),
		syncGroup: syncGroup,
	}
}

func (instance *Execution) Run() (ExitCode, error) {
	for _, service := range instance.executable.Services().GetAllAutoStartable() {
		err := instance.Start(service)
		if err != nil {
			instance.executable.Logger().LogProblem(err, logger.Error, "Could not start execution of service '%v'.", service)
		}
	}
	instance.wg.Wait()
	return *instance.masterExitCode, instance.masterError
}

func (instance *Execution) GetCountOfActiveExecutions() int {
	instance.doRLock()
	defer instance.doRUnlock()
	return len(instance.executions)
}

func (instance *Execution) Start(target *service.Service) error {
	execution, err := instance.createAndRegisterNotExistingExecutionFor(target)
	if err != nil {
		if sare, ok := err.(service.ServiceAlreadyRunningError); ok {
			return sare
		} else {
			return errors.New("Could not start service '%v'.", target).CausedBy(err)
		}
	}
	instance.wg.Add(1)
	go instance.drive(execution)
	return nil
}

func (instance *Execution) drive(target *service.Execution) {
	var exitCode ExitCode
	var err error
	defer instance.doAfterExecution(target, exitCode, err)
	respectDelay := true
	doRun := true
	for run := 1; doRun && target != nil; run++ {
		if respectDelay {
			if ! instance.delayedStartIfNeeded(target.Service(), run) {
				break
			}
		} else {
			run = 1
		}
		exitCode, err = target.Run()
		doRun, respectDelay = instance.checkAfterExecutionStates(target, exitCode, err)
		if doRun {
			newTarget, err := instance.recreateExecution(target)
			if err != nil {
				instance.executable.Logger().LogProblem(err, logger.Error, "Could not retrigger execution of '%v'.", target)
			} else {
				target = newTarget
			}
		}
	}
}

func (instance *Execution) recreateExecution(target *service.Execution) (*service.Execution, error) {
	instance.doWLock()
	defer instance.doWUnlock()
	s := target.Service()
	newTarget, err := s.NewExecution(instance.executable.KeyStore())
	if err != nil {
		delete(instance.executions, s)
		delete(instance.restartRequests, s)
	} else {
		instance.executions[s] = newTarget
		instance.restartRequests[s] = false
	}
	return newTarget, err
}

func (instance *Execution) checkAfterExecutionStates(target *service.Execution, exitCode ExitCode, err error) (doRestart bool, respectDelay bool) {
	if _, ok := err.(service.StoppedOrKilledError); ok {
		doRestart = false
	} else if _, ok := err.(service.UnrecoverableError); ok {
		doRestart = target.Service().Config().CronExpression.IsEnabled() && instance.masterExitCode == nil
	} else if instance.checkRestartRequestedAndClean(target.Service()) {
		doRestart = true
		respectDelay = false
	} else if target.Service().Config().SuccessExitCodes.Contains(exitCode) {
		doRestart = target.Service().Config().CronExpression.IsEnabled() && instance.masterExitCode == nil
	} else {
		doRestart = true
		respectDelay = true
	}
	return doRestart, respectDelay
}

func (instance *Execution) doAfterExecution(target *service.Execution, exitCode ExitCode, err error) {
	defer instance.doUnregisterExecution(target)
	if target.Service().Config().Type == service.Master {
		instance.masterExitCode = &exitCode
		instance.masterError = err
		others := instance.allExecutionsButMaster()
		if len(others) > 0 {
			instance.executable.Logger().Log(logger.Debug, "Master '%s' is down. Stopping all left services...", target.Name())
			for _, other := range others {
				go instance.Stop(other.Service())
			}
		}
	}
}

func (instance *Execution) delayedStartIfNeeded(s *service.Service, currentRun int) bool {
	if currentRun == 1 {
		return instance.delayedStartIfNeededFor(s, s.Config().StartDelayInSeconds, "Wait %d seconds before start...")
	} else {
		return instance.delayedStartIfNeededFor(s, s.Config().RestartDelayInSeconds, "Wait %d seconds before restart...")
	}
}

func (instance *Execution) delayedStartIfNeededFor(s *service.Service, delayInSeconds NonNegativeInteger, messagePattern string) bool {
	if s.Config().StartDelayInSeconds > 0 {
		s.Logger().Log(logger.Debug, messagePattern, delayInSeconds)
		return instance.syncGroup.Sleep(time.Duration(delayInSeconds) * time.Second) == nil
	} else {
		return true
	}
}

func (instance *Execution) checkRestartRequestedAndClean(target *service.Service) bool {
	instance.doWLock()
	defer instance.doWUnlock()
	result := instance.restartRequests[target]
	delete(instance.restartRequests, target)
	return result
}

func (instance *Execution) StopAll() {
	for _, execution := range instance.allExecutions() {
		instance.Stop(execution.Service())
	}
	instance.wg.Wait()
}

func (instance *Execution) Restart(target *service.Service) error {
	instance.doRLock()
	execution, ok := instance.executions[target]
	if ! ok {
		instance.doRUnlock()
		return instance.Start(target)
	}
	instance.restartRequests[target] = true
	instance.doRUnlock()
	execution.Stop()
	return nil
}

func (instance *Execution) Stop(target *service.Service) error {
	instance.doRLock()
	execution, ok := instance.executions[target]
	if ! ok {
		instance.doRUnlock()
		return service.ServiceDownError{Name: target.Name()}
	}
	instance.doRUnlock()
	execution.Stop()
	return nil
}

func (instance *Execution) Kill(target *service.Service) error {
	instance.doRLock()
	execution, ok := instance.executions[target]
	if ! ok {
		instance.doRUnlock()
		return service.ServiceDownError{Name: target.Name()}
	}
	instance.doRUnlock()
	execution.Kill()
	return nil
}

func (instance *Execution) Signal(target *service.Service, what Signal) error {
	instance.doRLock()
	execution, ok := instance.executions[target]
	if ! ok {
		instance.doRUnlock()
		return service.ServiceDownError{Name: target.Name()}
	}
	instance.doRUnlock()
	execution.Signal(what)
	return nil
}

func (instance *Execution) createAndRegisterNotExistingExecutionFor(target *service.Service) (*service.Execution, error) {
	instance.doWLock()
	defer instance.doWUnlock()
	result, err := target.NewExecution(instance.executable.KeyStore())
	if err != nil {
		return nil, err
	}
	if _, ok := instance.executions[target]; ok {
		return nil, service.ServiceAlreadyRunningError{Name: target.Name()}
	}
	instance.executions[target] = result
	return result, nil
}

func (instance *Execution) allExecutionsButMaster() []*service.Execution {
	instance.doRLock()
	defer instance.doRUnlock()
	result := []*service.Execution{}
	for s, candidate := range instance.executions {
		if s.Config().Type != service.Master {
			result = append(result, candidate)
		}
	}
	return result
}

func (instance *Execution) allExecutions() []*service.Execution {
	instance.doRLock()
	defer instance.doRUnlock()
	result := []*service.Execution{}
	for _, candidate := range instance.executions {
		result = append(result, candidate)
	}
	return result
}

func (instance *Execution) doUnregisterExecution(target *service.Execution) {
	instance.doWLock()
	defer instance.doWUnlock()
	delete(instance.executions, target.Service())
	delete(instance.restartRequests, target.Service())
	instance.wg.Done()
}

func (instance *Execution) doWLock() {
	instance.lock.Lock()
}

func (instance *Execution) doWUnlock() {
	instance.lock.Unlock()
}

func (instance *Execution) doRLock() {
	instance.lock.RLock()
}

func (instance *Execution) doRUnlock() {
	instance.lock.RUnlock()
}

func (instance *Execution) GetFor(s *service.Service) (*service.Execution, bool) {
	result, ok := instance.executions[s]
	return result, ok
}

func (instance *Execution) Information() map[string]service.Information {
	result := map[string]service.Information{}
	for _, service := range (*instance.executable.Services()) {
		result[service.Name()] = instance.InformationFor(service)
	}
	return result
}

func (instance *Execution) InformationFor(s *service.Service) service.Information {
	if result, ok := instance.GetFor(s); ok {
		return service.NewInformationForExecution(result)
	} else {
		return service.NewInformationForService(s)
	}
}
