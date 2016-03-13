// +build windows

package service

import (
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/panics"
	"github.com/echocat/caretakerd/values"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

const processAllAccess = 0x001F0FFF

func serviceHandleUsersFor(service *Service, cmd *exec.Cmd) {
	if !(*service).config.User.IsTrimmedEmpty() {
		panics.New("Could not handle users under windows. Please remove it from service '%s'.", service.Name()).Throw()
	}
}

func sendSignalToService(service *Service, process *os.Process, what values.Signal, signalTarget values.SignalTarget) error {
	switch signalTarget {
	case values.Process:
		return sendSignalToProcess(process.Pid, service, what)
	case values.ProcessGroup:
		return sendSignalToProcessGroup(process.Pid, service, what)
	case values.Mixed:
		if what == values.KILL || what == values.STOP {
			return sendSignalToProcessGroup(process.Pid, service, what)
		} else {
			return sendSignalToProcess(process.Pid, service, what)
		}
	}
	return errors.New("Could not signal '%v' (#%v) %v because signalTarget %v is not supported.", service, process.Pid, what, signalTarget)
}

func sendSignalToProcess(pid int, service *Service, what values.Signal) error {
	switch what {
	case values.TERM:
		return sendSpecialSignal(pid, syscall.CTRL_BREAK_EVENT)
	case values.INT:
		return sendSpecialSignal(pid, syscall.CTRL_C_EVENT)
	case values.STOP:
		fallthrough
	case values.KILL:
		processHandle, err := syscall.OpenProcess(processAllAccess, false, uint32(pid))
		if err != nil {
			if se, ok := err.(syscall.Errno); ok {
				if se.Error() != "invalid argument" {
					return errors.New("Could not signal '%v' (#%v): %v", service, pid, what).CausedBy(err)
				}
			} else {
				return errors.New("Could not signal '%v' (#%v): %v", service, pid, what).CausedBy(err)
			}
		}
		if err == nil {
			err := syscall.TerminateProcess(processHandle, 1)
			if err != nil {
				if se, ok := err.(syscall.Errno); ok {
					if se.Error() != "invalid argument" {
						return errors.New("Could not signal '%v' (#%v): %v", service, pid, what).CausedBy(err)
					}
				} else {
					return errors.New("Could not signal '%v' (#%v): %v", service, pid, what).CausedBy(err)
				}
			}
			syscall.CloseHandle(processHandle)
		}
		return nil
	}
	return errors.New("Could not signal '%v' (#%v): %v ... because this signal is not supported on widows.", service, pid, what)
}

func sendSignalToProcessGroup(pid int, service *Service, what values.Signal) error {
	pe := syscall.ProcessEntry32{}
	pe.Size = uint32(unsafe.Sizeof(pe))

	hSnap, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return errors.New("Could iterate over children of service %v (#%v).", service, pid).CausedBy(err)
	}

	err = syscall.Process32First(hSnap, &pe)
	if err != nil {
		return err
	}

	tryNext := true
	for tryNext {
		if pe.ParentProcessID == uint32(pid) {
			err := sendSignalToProcess(int(pe.ProcessID), service, what)
			if err != nil {
				return err
			}
		}
		tryNext = syscall.Process32Next(hSnap, &pe) == nil
	}

	return sendSignalToProcess(pid, service, what)
}

func sendSpecialSignal(pid int, what uintptr) error {
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		panics.New("Could not set terminate to process #%v.", pid).CausedBy(e).Throw()
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		panics.New("Could not set terminate to process #%v.", pid).CausedBy(e).Throw()
	}
	_, _, err := p.Call(what, uintptr(pid))
	return err
}

func (instance *Service) createSysProcAttr() *syscall.SysProcAttr {
	// Prevent that a created process receives signals for caretakerd
	return &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}
