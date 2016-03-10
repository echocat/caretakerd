// +build windows

package service

import (
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/panics"
	"github.com/echocat/caretakerd/values"
	"os"
	"os/exec"
	"syscall"
)

func serviceHandleUsersFor(service *Service, cmd *exec.Cmd) {
	if !(*service).config.User.IsTrimmedEmpty() {
		panics.New("Could not handle users under windows. Please remove it from service '%s'.", service.Name()).Throw()
	}
}

func sendSignalToService(service *Service, process *os.Process, what values.Signal) error {
	if service.config.StopSignalTarget != values.Process {
		return errors.New("On windows only stopSignalTarget == 'process' is supported but got: %v", service.config.StopSignalTarget)
	}
	if what == values.TERM {
		sendSpecialSignal(process, syscall.CTRL_BREAK_EVENT)
	} else if what == values.INT {
		sendSpecialSignal(process, syscall.CTRL_C_EVENT)
	} else {
		sSignal := syscall.Signal(what)
		service.Logger().Log(logger.Debug, "Send %v to %v with PID %d", sSignal, service, process.Pid)
		if err := process.Signal(sSignal); err != nil {
			ignore := false
			if se, ok := err.(syscall.Errno); ok {
				ignore = se.Error() == "invalid argument"
			}
			if !ignore {
				return errors.New("Could not signal '%v': %v", service, what).CausedBy(err)
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

func (instance *Service) createSysProcAttr() *syscall.SysProcAttr {
	// Prevent that a created process receives signals for caretakerd
	return &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}
