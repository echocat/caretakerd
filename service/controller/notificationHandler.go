package controller
import (
    "sync"
    "os"
    osignal "os/signal"
    "github.com/echocat/caretakerd/service/signal"
    "syscall"
    "github.com/echocat/caretakerd/panics"
)

type NotificationHandler struct {
    lock          *sync.Mutex
    signalChannel chan os.Signal
    onSignal      func(signal.Signal)
}

func InstallNotificationHandler(onSignal func(signal.Signal)) *NotificationHandler {
    result := &NotificationHandler{
        lock: new(sync.Mutex),
        signalChannel: make(chan os.Signal, 1),
        onSignal: onSignal,
    }
    osignal.Notify(result.signalChannel, os.Interrupt, os.Kill, syscall.SIGTERM)
    go result.terminationNotificationHandler()
    return &result;
}

func (i *NotificationHandler) terminationNotificationHandler() {
    defer panics.DefaultPanicHandler()
    for {
        plainSignal := <-i.signalChannel
        s := signal.Signal(plainSignal.(syscall.Signal))
        if s == signal.NOOP {
            break
        }
        if !signal.IsHandlingOfSignalIgnoreable(s) {
            i.onSignal(signal)
        }
    }
}


func (i *NotificationHandler) Uninstall() {
    i.lock.Lock()
    defer func() {
        i.signalChannel = nil
        i.lock.Unlock()
    }()
    if i.signalChannel != nil {
        osignal.Stop(i.signalChannel)
    }
}
