package caretakerd

import (
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/keyStore"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/service"
	"github.com/echocat/caretakerd/values"
	ssync "sync"
	"time"
)

// Executable indicates an object that could be executed.
type Executable interface {
	Services() *service.Services
	KeyStore() *keyStore.KeyStore
	Logger() *logger.Logger
}

// Execution is an instance of an execution of every service of caretakerd.
type Execution struct {
	executable      Executable
	executions      map[*service.Service]*service.Execution
	restartRequests map[*service.Service]bool
	stopRequests    map[*service.Service]bool
	masterExitCode  *values.ExitCode
	masterError     error
	lock            *ssync.RWMutex
	wg              *ssync.WaitGroup
}

// NewExecution creates a new Execution instance of caretakerd.
//
// Hint: A caretakerd Execution instance could only be called once.
func NewExecution(executable Executable) *Execution {
	return &Execution{
		executable:      executable,
		masterExitCode:  nil,
		executions:      map[*service.Service]*service.Execution{},
		restartRequests: map[*service.Service]bool{},
		stopRequests:    map[*service.Service]bool{},
		lock:            new(ssync.RWMutex),
		wg:              new(ssync.WaitGroup),
	}
}

// Run starts the caretakerd execution, every service and required resource.
// This is a blocking method.
func (instance *Execution) Run() (values.ExitCode, error) {
	autoStartableServices := instance.executable.Services().GetAllAutoStartable()
	// Start all non-master services first to start the master properly.
	for _, target := range autoStartableServices {
		if target.Config().Type != service.Master {
			instance.startAndLogProblemsIfNeeded(target)
		}
	}
	// Now start the master.
	masterStarted := false
	for _, target := range autoStartableServices {
		if target.Config().Type == service.Master {
			err := instance.Start(target)
			if err != nil {
				(*instance).masterError = err
			}
			masterStarted = true
		}
	}
	if !masterStarted {
		(*instance).masterError = errors.New("No master was started. There are no master configured?")
	}
	if (*instance).masterError != nil {
		instance.stopOthers()
	}
	instance.wg.Wait()
	exitCode := values.ExitCode(-1)
	if (*instance).masterExitCode != nil {
		exitCode = *(*instance).masterExitCode
	}
	return exitCode, (*instance).masterError
}

// GetCountOfActiveExecutions returns the number of active executions of all services.
func (instance *Execution) GetCountOfActiveExecutions() int {
	instance.doRLock()
	defer instance.doRUnlock()
	return len(instance.executions)
}

func (instance *Execution) startAndLogProblemsIfNeeded(target *service.Service) {
	err := instance.Start(target)
	if err != nil {
		instance.executable.Logger().LogProblem(err, logger.Error, "Could not start execution of service '%v'.", target)
	}
}

// Start starts the given service.
// This is not a blocking method.
func (instance *Execution) Start(target *service.Service) error {
	execution, err := instance.createAndRegisterNotExistingExecutionFor(target)
	if err != nil {
		if sare, ok := err.(service.AlreadyRunningError); ok {
			return sare
		}
		return errors.New("Could not start service '%v'.", target).CausedBy(err)
	}
	instance.wg.Add(1)
	go instance.drive(execution)
	return nil
}

func (instance *Execution) drive(target *service.Execution) {
	var exitCode values.ExitCode
	var err error
	defer instance.doAfterExecution(target, exitCode, err)
	respectDelay := true
	doRun := true
	for run := 1; doRun && target != nil && !instance.isAlreadyStopRequested(target); run++ {
		if respectDelay {
			if !instance.delayedStartIfNeeded(target, run) {
				break
			}
		} else {
			run = 1
		}
		if !instance.isAlreadyStopRequested(target) {
			exitCode, err = target.Run()
			doRun, respectDelay = instance.checkAfterExecutionStates(target, exitCode, err)
			if doRun && !instance.isAlreadyStopRequested(target) {
				newTarget, err := instance.recreateExecution(target)
				if err != nil {
					instance.executable.Logger().LogProblem(err, logger.Error, "Could not retrigger execution of '%v'.", target)
				} else {
					target = newTarget
				}
			}
		}
	}
}

func (instance *Execution) isAlreadyStopRequested(target *service.Execution) bool {
	instance.doRLock()
	defer instance.doRUnlock()
	stopRequested, ok := instance.stopRequests[target.Service()]
	return ok && stopRequested
}

func (instance *Execution) recreateExecution(target *service.Execution) (*service.Execution, error) {
	instance.doWLock()
	defer instance.doWUnlock()
	s := target.Service()
	newTarget, err := s.NewExecution(instance.executable.KeyStore())
	if err != nil {
		delete(instance.executions, s)
	} else {
		instance.executions[s] = newTarget
	}
	delete(instance.restartRequests, s)
	return newTarget, err
}

func (instance *Execution) checkAfterExecutionStates(target *service.Execution, exitCode values.ExitCode, err error) (doRestart bool, respectDelay bool) {
	if _, ok := err.(service.StoppedOrKilledError); ok {
		doRestart = false
	} else if _, ok := err.(service.UnrecoverableError); ok {
		doRestart = target.Service().Config().CronExpression.IsEnabled() && instance.masterExitCode == nil
	} else if instance.checkRestartRequestedAndClean(target.Service()) {
		doRestart = true
		respectDelay = false
	} else if target.Service().Config().SuccessExitCodes.Contains(exitCode) {
		doRestart = (target.Service().Config().CronExpression.IsEnabled() && instance.masterExitCode == nil) || target.Service().Config().AutoRestart.OnSuccess()
	} else {
		doRestart = target.Service().Config().AutoRestart.OnFailures()
		respectDelay = true
	}
	return doRestart, respectDelay
}

func (instance *Execution) doAfterExecution(target *service.Execution, exitCode values.ExitCode, err error) {
	defer instance.doUnregisterExecution(target)
	if target.Service().Config().Type == service.Master {
		instance.masterExitCode = &exitCode
		instance.masterError = err
		instance.stopOthers()
	}
}

func (instance *Execution) stopOthers() {
	others := instance.allExecutionsButMaster()
	if len(others) > 0 {
		instance.executable.Logger().Log(logger.Debug, "Master is down. Stopping all remaining services...")
		for _, other := range others {
			go instance.Stop(other.Service())
		}
	}
}

func (instance *Execution) delayedStartIfNeeded(target *service.Execution, currentRun int) bool {
	config := target.Service().Config()
	if currentRun == 1 {
		return instance.delayedStartIfNeededFor(target, config.StartDelayInSeconds, "Wait %d seconds before starting...")
	}
	return instance.delayedStartIfNeededFor(target, config.RestartDelayInSeconds, "Wait %d seconds before restarting...")
}

func (instance *Execution) delayedStartIfNeededFor(target *service.Execution, delayInSeconds values.NonNegativeInteger, messagePattern string) bool {
	s := target.Service()
	if s.Config().StartDelayInSeconds > 0 {
		s.Logger().Log(logger.Debug, messagePattern, delayInSeconds)
		return target.SyncGroup().Sleep(time.Duration(delayInSeconds)*time.Second) == nil
	}
	return true
}

func (instance *Execution) checkRestartRequestedAndClean(target *service.Service) bool {
	instance.doWLock()
	defer instance.doWUnlock()
	result := instance.restartRequests[target]
	delete(instance.restartRequests, target)
	return result
}

func (instance *Execution) registerStopRequestsFor(executions ...*service.Execution) {
	instance.doWLock()
	defer instance.doWUnlock()
	for _, execution := range executions {
		instance.stopRequests[execution.Service()] = true
		delete(instance.restartRequests, execution.Service())
	}
}

// StopAll stop all running serivces.
func (instance *Execution) StopAll() {
	for _, execution := range instance.allExecutions() {
		if execution.Service().Config().Type == service.Master {
			// Hint: Stopping the master should also trigger the shutdown of all other services.
			// This is the reason why at this point only the master is shut down.
			instance.Stop(execution.Service())
		}
	}
	instance.wg.Wait()
}

// Restart stops and restarts the given service.
func (instance *Execution) Restart(target *service.Service) error {
	instance.doRLock()
	if stopRequested, ok := instance.stopRequests[target]; ok && stopRequested {
		instance.doRUnlock()
		return service.AlreadyStoppedError{Name: target.Name()}
	}
	execution, ok := instance.executions[target]
	if !ok {
		instance.doRUnlock()
		return instance.Start(target)
	}
	instance.restartRequests[target] = true
	instance.doRUnlock()
	execution.Stop()
	return nil
}

// Stop stops the given service.
func (instance *Execution) Stop(target *service.Service) error {
	instance.doRLock()
	execution, ok := instance.executions[target]
	if !ok {
		instance.doRUnlock()
		return service.AlreadyStoppedError{Name: target.Name()}
	}
	instance.doRUnlock()
	instance.registerStopRequestsFor(execution)
	execution.Stop()
	return nil
}

// Kill kills the given service.
func (instance *Execution) Kill(target *service.Service) error {
	instance.doRLock()
	execution, ok := instance.executions[target]
	if !ok {
		instance.doRUnlock()
		return service.AlreadyStoppedError{Name: target.Name()}
	}
	instance.doRUnlock()
	instance.registerStopRequestsFor(execution)
	return execution.Kill()
}

// Signal sends the given signal to the given service.
func (instance *Execution) Signal(target *service.Service, what values.Signal) error {
	instance.doRLock()
	execution, ok := instance.executions[target]
	if !ok {
		instance.doRUnlock()
		return service.AlreadyStoppedError{Name: target.Name()}
	}
	instance.doRUnlock()
	return execution.Signal(what)
}

func (instance *Execution) createAndRegisterNotExistingExecutionFor(target *service.Service) (*service.Execution, error) {
	instance.doWLock()
	defer instance.doWUnlock()
	result, err := target.NewExecution(instance.executable.KeyStore())
	if err != nil {
		return nil, err
	}
	if _, ok := instance.executions[target]; ok {
		return nil, service.AlreadyRunningError{Name: target.Name()}
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

// GetFor queries the current active service execution for the given service.
// Returns "false" if no current execution matches.
func (instance *Execution) GetFor(s *service.Service) (*service.Execution, bool) {
	result, ok := instance.executions[s]
	return result, ok
}

// Information returns an information object that contains information for every
// configured service.
func (instance *Execution) Information() map[string]service.Information {
	result := map[string]service.Information{}
	for _, service := range *instance.executable.Services() {
		result[service.Name()] = instance.InformationFor(service)
	}
	return result
}

// InformationFor returns an information object for the given service.
func (instance *Execution) InformationFor(s *service.Service) service.Information {
	if result, ok := instance.GetFor(s); ok {
		return service.NewInformationForExecution(result)
	}
	return service.NewInformationForService(s)
}
