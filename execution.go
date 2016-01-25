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

func (this *Execution) Run() (ExitCode, error) {
    for _, service := range this.executable.Services().GetAllAutoStartable() {
        err := this.Start(service)
        if err != nil {
            this.executable.Logger().LogProblem(err, logger.Error, "Could not start execution of service '%v'.", service)
        }
    }
    this.wg.Wait()
    return *this.masterExitCode, this.masterError
}

func (this *Execution) GetCountOfActiveExecutions() int {
    this.doRLock()
    defer this.doRUnlock()
    return len(this.executions)
}

func (this *Execution) Start(target *service.Service) error {
    execution, err := this.createAndRegisterNotExistingExecutionFor(target)
    if err != nil {
        if sare, ok := err.(service.ServiceAlreadyRunningError); ok {
            return sare
        } else {
            return errors.New("Could not start service '%v'.", target).CausedBy(err)
        }
    }
    this.wg.Add(1)
    go this.drive(execution)
    return nil
}

func (this *Execution) drive(target *service.Execution) {
    var exitCode ExitCode
    var err error
    defer this.doAfterExecution(target, exitCode, err)
    respectDelay := true
    doRun := true
    for run := 1; doRun && target != nil; run++ {
        if respectDelay {
            if ! this.delayedStartIfNeeded(target.Service(), run) {
                break
            }
        } else {
            run = 1
        }
        exitCode, err = target.Run()
        doRun, respectDelay = this.checkAfterExecutionStates(target, exitCode, err)
        if doRun {
            newTarget, err := this.recreateExecution(target)
            if err != nil {
                this.executable.Logger().LogProblem(err, logger.Error, "Could not retrigger execution of '%v'.", target)
            } else {
                target = newTarget
            }
        }
    }
}

func (this *Execution) recreateExecution(target *service.Execution) (*service.Execution, error) {
    this.doWLock()
    defer this.doWUnlock()
    s := target.Service()
    newTarget, err := s.NewExecution(this.executable.KeyStore())
    if err != nil {
        delete(this.executions, s)
        delete(this.restartRequests, s)
    } else {
        this.executions[s] = newTarget
        this.restartRequests[s] = false
    }
    return newTarget, err
}

func (this *Execution) checkAfterExecutionStates(target *service.Execution, exitCode ExitCode, err error) (doRestart bool, respectDelay bool) {
    if _, ok := err.(service.StoppedOrKilledError); ok {
        doRestart = false
    } else if _, ok := err.(service.UnrecoverableError); ok {
        doRestart = target.Service().Config().CronExpression.IsEnabled() && this.masterExitCode == nil
    } else if this.checkRestartRequestedAndClean(target.Service()) {
        doRestart = true
        respectDelay = false
    } else if target.Service().Config().SuccessExitCodes.Contains(exitCode) {
        doRestart = target.Service().Config().CronExpression.IsEnabled() && this.masterExitCode == nil
    } else {
        doRestart = true
        respectDelay = true
    }
    return doRestart, respectDelay
}

func (this *Execution) doAfterExecution(target *service.Execution, exitCode ExitCode, err error) {
    defer this.doUnregisterExecution(target)
    if target.Service().Config().Type == service.Master {
        this.masterExitCode = &exitCode
        this.masterError = err
        others := this.allExecutionsButMaster()
        if len(others) > 0 {
            this.executable.Logger().Log(logger.Debug, "Master '%s' is down. Stopping all left services...", target.Name())
            for _, other := range others {
                go this.Stop(other.Service())
            }
        }
    }
}

func (this *Execution) delayedStartIfNeeded(s *service.Service, currentRun int) bool {
    if currentRun == 1 {
        return this.delayedStartIfNeededFor(s, s.Config().StartDelayInSeconds, "Wait %d seconds before start...")
    } else {
        return this.delayedStartIfNeededFor(s, s.Config().RestartDelayInSeconds, "Wait %d seconds before restart...")
    }
}

func (this *Execution) delayedStartIfNeededFor(s *service.Service, delayInSeconds NonNegativeInteger, messagePattern string) bool {
    if s.Config().StartDelayInSeconds > 0 {
        s.Logger().Log(logger.Debug, messagePattern, delayInSeconds)
        return this.syncGroup.Sleep(time.Duration(delayInSeconds) * time.Second) == nil
    } else {
        return true
    }
}

func (this *Execution) checkRestartRequestedAndClean(target *service.Service) bool {
    this.doWLock()
    defer this.doWUnlock()
    result := this.restartRequests[target]
    delete(this.restartRequests, target)
    return result
}

func (this *Execution) StopAll() {
    for _, execution := range this.allExecutions() {
        this.Stop(execution.Service())
    }
    this.wg.Wait()
}

func (this *Execution) Restart(target *service.Service) error {
    this.doRLock()
    execution, ok := this.executions[target]
    if ! ok {
        this.doRUnlock()
        return this.Start(target)
    }
    this.restartRequests[target] = true
    this.doRUnlock()
    execution.Stop()
    return nil
}

func (this *Execution) Stop(target *service.Service) error {
    this.doRLock()
    execution, ok := this.executions[target]
    if ! ok {
        this.doRUnlock()
        return service.ServiceDownError{Name: target.Name()}
    }
    this.doRUnlock()
    execution.Stop()
    return nil
}

func (this *Execution) Kill(target *service.Service) error {
    this.doRLock()
    execution, ok := this.executions[target]
    if ! ok {
        this.doRUnlock()
        return service.ServiceDownError{Name: target.Name()}
    }
    this.doRUnlock()
    execution.Kill()
    return nil
}

func (this *Execution) Signal(target *service.Service, what Signal) error {
    this.doRLock()
    execution, ok := this.executions[target]
    if ! ok {
        this.doRUnlock()
        return service.ServiceDownError{Name: target.Name()}
    }
    this.doRUnlock()
    execution.Signal(what)
    return nil
}

func (this *Execution) createAndRegisterNotExistingExecutionFor(target *service.Service) (*service.Execution, error) {
    this.doWLock()
    defer this.doWUnlock()
    result, err := target.NewExecution(this.executable.KeyStore())
    if err != nil {
        return nil, err
    }
    if _, ok := this.executions[target]; ok {
        return nil, service.ServiceAlreadyRunningError{Name: target.Name()}
    }
    this.executions[target] = result
    return result, nil
}

func (this *Execution) allExecutionsButMaster() []*service.Execution {
    this.doRLock()
    defer this.doRUnlock()
    result := []*service.Execution{}
    for s, candidate := range this.executions {
        if s.Config().Type != service.Master {
            result = append(result, candidate)
        }
    }
    return result
}

func (this *Execution) allExecutions() []*service.Execution {
    this.doRLock()
    defer this.doRUnlock()
    result := []*service.Execution{}
    for _, candidate := range this.executions {
        result = append(result, candidate)
    }
    return result
}

func (this *Execution) doUnregisterExecution(target *service.Execution) {
    this.doWLock()
    defer this.doWUnlock()
    delete(this.executions, target.Service())
    delete(this.restartRequests, target.Service())
    this.wg.Done()
}

func (this *Execution) doWLock() {
    this.lock.Lock()
}

func (this *Execution) doWUnlock() {
    this.lock.Unlock()
}

func (this *Execution) doRLock() {
    this.lock.RLock()
}

func (this *Execution) doRUnlock() {
    this.lock.RUnlock()
}

func (this *Execution) GetFor(s *service.Service) (*service.Execution, bool) {
    result, ok := this.executions[s]
    return result, ok
}

func (this *Execution) Information() map[string]service.Information {
    result := map[string]service.Information{}
    for _, service := range (*this.executable.Services()) {
        result[service.Name()] = this.InformationFor(service)
    }
    return result
}

func (this *Execution) InformationFor(s *service.Service) service.Information {
    if result, ok := this.GetFor(s); ok {
        return service.NewInformationForExecution(result)
    } else {
        return service.NewInformationForService(s)
    }
}
