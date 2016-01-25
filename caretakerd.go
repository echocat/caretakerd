package caretakerd

import (
    "runtime"
    "sync"
    "os"
    "syscall"
    osignal "os/signal"
    . "github.com/echocat/caretakerd/values"
    "github.com/echocat/caretakerd/service"
    "github.com/echocat/caretakerd/logger"
    usync "github.com/echocat/caretakerd/sync"
    "github.com/echocat/caretakerd/keyStore"
    "github.com/echocat/caretakerd/errors"
    "github.com/echocat/caretakerd/control"
    "github.com/echocat/caretakerd/panics"
    "github.com/echocat/caretakerd/rpc"
)

type Caretakerd struct {
    config        Config
    logger        *logger.Logger
    control       *control.Control
    services      *service.Services
    lock          *sync.Mutex
    syncGroup     *usync.SyncGroup
    execution     *Execution
    signalChannel chan os.Signal
    open          bool
    keyStore      *keyStore.KeyStore
}

func finalize(what *Caretakerd) {
    what.Close()
}

func NewCaretakerd(conf Config, syncGroup *usync.SyncGroup) (*Caretakerd, error) {
    err := conf.Validate()
    if err != nil {
        return nil, err
    }
    log, err := logger.NewLogger(conf.Logger, "caretakerd", syncGroup.NewSyncGroup())
    if err != nil {
        return nil, errors.New("Could not create logger for caretakerd.").CausedBy(err)
    }
    ks, err := keyStore.NewKeyStore(bool(conf.Rpc.Enabled), conf.KeyStore)
    if err != nil {
        return nil, err
    }
    ctl, err := control.NewControl(conf.Control, ks)
    if err != nil {
        return nil, err
    }
    services, err := service.NewServices(conf.Services, syncGroup.NewSyncGroup(), ks)
    if err != nil {
        return nil, err
    }
    result := Caretakerd{
        open: true,
        config: conf,
        logger: log,
        control: ctl,
        keyStore: ks,
        services: services,
        lock: new(sync.Mutex),
        syncGroup: syncGroup,
        signalChannel: nil,
    }
    runtime.SetFinalizer(&result, finalize)
    return &result, nil
}

func (this Caretakerd) IsOpen() bool {
    return this.open
}

func (this *Caretakerd) Close() {
    defer func() {
        this.open = false
    }()
    this.Stop()
    this.services.Close()
    this.logger.Close()
}

func (this Caretakerd) Logger() *logger.Logger {
    return this.logger
}

func (this *Caretakerd) Control() *control.Control {
    return this.control
}

func (this *Caretakerd) Services() *service.Services {
    return this.services
}

func (this *Caretakerd) KeyStore() *keyStore.KeyStore {
    return this.keyStore
}

func (this *Caretakerd) ConfigObject() interface{} {
    return this.config
}

func (this *Caretakerd) Run() (ExitCode, error) {
    var r *rpc.Rpc
    defer func() {
        this.uninstallTerminationNotificationHandler()
        if r != nil {
            r.Stop()
        }
    }()

    execution := NewExecution(this, this.syncGroup.NewSyncGroup())
    if this.config.Rpc.Enabled == Boolean(true) {
        r = rpc.NewRpc(this.config.Rpc, execution, this, this.logger)
        r.Start()
    }
    this.installTerminationNotificationHandler()
    this.execution = execution
    return execution.Run()
}

func (i *Caretakerd) Stop() {
    defer func() {
        i.execution = nil
    }()
    execution := i.execution
    if execution != nil {
        execution.StopAll()
    }
    i.syncGroup.Interrupt()
}

func (i *Caretakerd) installTerminationNotificationHandler() {
    i.lock.Lock()
    defer func() {
        i.lock.Unlock()
    }()
    if i.signalChannel == nil {
        i.signalChannel = make(chan os.Signal, 1)
        osignal.Notify(i.signalChannel, os.Interrupt, os.Kill, syscall.SIGTERM)
        go i.terminationNotificationHandler()
    }
}

func (i *Caretakerd) terminationNotificationHandler() {
    defer panics.DefaultPanicHandler()
    for {
        plainSignal := <-i.signalChannel
        s := Signal(plainSignal.(syscall.Signal))
        if s == NOOP {
            break
        }
        if !IsHandlingOfSignalIgnoreable(s) {
            i.Stop()
        }
    }
}

func (i *Caretakerd) uninstallTerminationNotificationHandler() {
    i.lock.Lock()
    defer func() {
        i.signalChannel = nil
        i.lock.Unlock()
    }()
    if i.signalChannel != nil {
        osignal.Stop(i.signalChannel)
    }
}
