// +build windows

package service

import (
    "os/exec"
    "os"
    "syscall"
    "github.com/echocat/caretakerd/panics"
    "github.com/echocat/caretakerd/service/signal"
    "github.com/echocat/caretakerd/errors"
    "github.com/echocat/caretakerd/logger"
)

func serviceHandleUsersFor(service *Service, cmd *exec.Cmd) {
    if ! (*service).config.User.IsTrimmedEmpty() {
        panics.New("Could not handle users under windows. Please remove it from service '%s'.", service.Name()).Throw()
    }
}

func sendSignalToService(service *Service, process *os.Process, what signal.Signal) error {
    signal.RecordSendSignal(what)
    if what == signal.TERM {
        sendSpecialSignal(process, syscall.CTRL_BREAK_EVENT)
    } else if what == signal.INT {
        sendSpecialSignal(process, syscall.CTRL_C_EVENT)
    } else {
        sSignal := syscall.Signal(what)
        service.Logger().Log(logger.Debug, "Send %v to %v with PID %d", sSignal, service, process.Pid)
        if err := process.Signal(sSignal); err != nil {
            ignore := false
            if se, ok := err.(syscall.Errno); ok {
                ignore = se.Error() == "invalid argument"
            }
            if ! ignore {
                return errors.New("Could no signal '%v': %v", service, what).CausedBy(err)
            }
        }
    }
    return nil
}

func sendSpecialSignal(process *os.Process, what uintptr) {
    d, e := syscall.LoadDLL("kernel32.dll")
    if e != nil {
        panics.New("Could not set terminate to process %v.", process).CausedBy(e).Throw()
    }
    p, e := d.FindProc("GenerateConsoleCtrlEvent")
    if e != nil {
        panics.New("Could not set terminate to process %v.", process).CausedBy(e).Throw()
    }
    p.Call(what, uintptr(process.Pid))
}
