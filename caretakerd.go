package caretakerd

import (
	"github.com/echocat/caretakerd/control"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/keyStore"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/panics"
	"github.com/echocat/caretakerd/rpc"
	"github.com/echocat/caretakerd/service"
	usync "github.com/echocat/caretakerd/sync"
	"github.com/echocat/caretakerd/values"
	"os"
	osignal "os/signal"
	"runtime"
	"sync"
	"syscall"
)

// Caretakerd instance structure
type Caretakerd struct {
	config        Config
	logger        *logger.Logger
	control       *control.Control
	services      *service.Services
	lock          *sync.Mutex
	syncGroup     *usync.Group
	execution     *Execution
	signalChannel chan os.Signal
	open          bool
	keyStore      *keyStore.KeyStore
}

func finalize(what *Caretakerd) {
	what.Close()
}

// NewCaretakerd create a new Caretakerd instance from given config
func NewCaretakerd(conf Config, syncGroup *usync.Group) (*Caretakerd, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}
	log, err := logger.NewLogger(conf.Logger, "caretakerd", syncGroup)
	if err != nil {
		return nil, errors.New("Could not create logger for caretakerd.").CausedBy(err)
	}
	ks, err := keyStore.NewKeyStore(bool(conf.RPC.Enabled), conf.KeyStore)
	if err != nil {
		return nil, err
	}
	ctl, err := control.NewControl(conf.Control, ks)
	if err != nil {
		return nil, err
	}
	services, err := service.NewServices(conf.Services, syncGroup, ks)
	if err != nil {
		return nil, err
	}
	result := Caretakerd{
		open:          true,
		config:        conf,
		logger:        log,
		control:       ctl,
		keyStore:      ks,
		services:      services,
		lock:          new(sync.Mutex),
		syncGroup:     syncGroup,
		signalChannel: nil,
	}
	runtime.SetFinalizer(&result, finalize)
	return &result, nil
}

// IsOpen return true if caretakerd is still open. This should be false after Close() was called.
func (instance Caretakerd) IsOpen() bool {
	return instance.open
}

// Close close the caretakerd instance and free resoruces.
// After call of this method it is not longer possible to use this instance.
func (instance *Caretakerd) Close() {
	defer func() {
		instance.open = false
	}()
	instance.Stop()
	instance.services.Close()
	instance.logger.Close()
}

// Logger return instantiated logger that belongs to this instance.
func (instance Caretakerd) Logger() *logger.Logger {
	return instance.logger
}

// Control return instantiated control that belongs to this instance.
func (instance *Caretakerd) Control() *control.Control {
	return instance.control
}

// Services return instantiated services that belongs to this instance.
func (instance *Caretakerd) Services() *service.Services {
	return instance.services
}

// KeyStore return instantiated keyStore that belongs to this instance.
func (instance *Caretakerd) KeyStore() *keyStore.KeyStore {
	return instance.keyStore
}

// ConfigObject return config this instances was created with.
func (instance *Caretakerd) ConfigObject() interface{} {
	return instance.config
}

// Run starts every services and requird resources of caretakerd.
// This method is blocking.
func (instance *Caretakerd) Run() (values.ExitCode, error) {
	var r *rpc.RPC
	defer func() {
		instance.uninstallTerminationNotificationHandler()
		if r != nil {
			r.Stop()
		}
	}()

	execution := NewExecution(instance)
	if instance.config.RPC.Enabled == values.Boolean(true) {
		r = rpc.NewRPC(instance.config.RPC, execution, instance, instance.logger)
		r.Start()
	}
	instance.installTerminationNotificationHandler()
	instance.execution = execution
	return execution.Run()
}

// Stop stops this instance (if running).
// This method is blocking until every service and resource is stopped.
func (instance *Caretakerd) Stop() {
	defer func() {
		instance.execution = nil
	}()
	execution := instance.execution
	if execution != nil {
		execution.StopAll()
	}
	instance.syncGroup.Interrupt()
}

func (instance *Caretakerd) installTerminationNotificationHandler() {
	instance.lock.Lock()
	defer func() {
		instance.lock.Unlock()
	}()
	if instance.signalChannel == nil {
		instance.signalChannel = make(chan os.Signal, 1)
		osignal.Notify(instance.signalChannel, syscall.SIGINT, syscall.SIGTERM)
		go instance.terminationNotificationHandler()
	}
}

func (instance *Caretakerd) terminationNotificationHandler() {
	defer panics.DefaultPanicHandler()
	for {
		osSignal, channelReady := <-instance.signalChannel
		if channelReady {
			signal := values.Signal(osSignal.(syscall.Signal))
			instance.Logger().Log(logger.Debug, "Received shudown signal: %v", signal)
			instance.Stop()
		} else {
			break
		}
	}
}

func (instance *Caretakerd) uninstallTerminationNotificationHandler() {
	instance.lock.Lock()
	defer func() {
		instance.signalChannel = nil
		instance.lock.Unlock()
	}()
	if instance.signalChannel != nil {
		osignal.Stop(instance.signalChannel)
	}
}
