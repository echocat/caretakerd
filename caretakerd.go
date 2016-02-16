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
	. "github.com/echocat/caretakerd/values"
	osignal "os/signal"
	"runtime"
	"sync"
	"syscall"
	"os"
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
	log, err := logger.NewLogger(conf.Logger, "caretakerd", syncGroup)
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

func (instance Caretakerd) IsOpen() bool {
	return instance.open
}

func (instance *Caretakerd) Close() {
	defer func() {
		instance.open = false
	}()
	instance.Stop()
	instance.services.Close()
	instance.logger.Close()
}

func (instance Caretakerd) Logger() *logger.Logger {
	return instance.logger
}

func (instance *Caretakerd) Control() *control.Control {
	return instance.control
}

func (instance *Caretakerd) Services() *service.Services {
	return instance.services
}

func (instance *Caretakerd) KeyStore() *keyStore.KeyStore {
	return instance.keyStore
}

func (instance *Caretakerd) ConfigObject() interface{} {
	return instance.config
}

func (instance *Caretakerd) Run() (ExitCode, error) {
	var r *rpc.Rpc
	defer func() {
		instance.uninstallTerminationNotificationHandler()
		if r != nil {
			r.Stop()
		}
	}()

	execution := NewExecution(instance)
	if instance.config.Rpc.Enabled == Boolean(true) {
		r = rpc.NewRpc(instance.config.Rpc, execution, instance, instance.logger)
		r.Start()
	}
	instance.installTerminationNotificationHandler()
	instance.execution = execution
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
		osignal.Notify(i.signalChannel, syscall.SIGINT, syscall.SIGTERM)
		go i.terminationNotificationHandler()
	}
}

func (i *Caretakerd) terminationNotificationHandler() {
	defer panics.DefaultPanicHandler()
	for {
		osSignal, channelReady := <-i.signalChannel
		if channelReady {
			signal := Signal(osSignal.(syscall.Signal))
			i.Logger().Log(logger.Debug, "Received shudown signal: %v", signal)
			i.Stop()
		} else {
			break
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
