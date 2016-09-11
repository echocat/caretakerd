// +build linux darwin

package service

import (
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/panics"
	"github.com/echocat/caretakerd/values"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

func serviceHandleUsersFor(service *Service, cmd *exec.Cmd) {
	config := (*service).config
	userName := config.User
	if !userName.IsTrimmedEmpty() {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		uid, gid, err := lookupUser(userName.String())
		if err != nil {
			panics.New("Could not run as user '%v'.", userName).CausedBy(err).Throw()
		}
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	}
}

func sendSignalToService(service *Service, process *os.Process, what values.Signal, signalTarget values.SignalTarget) error {
	switch signalTarget {
	case values.Process:
		return sendSignalToProcess(process, what, false)
	case values.ProcessGroup:
		return sendSignalToProcess(process, what, true)
	case values.Mixed:
		return sendSignalToProcess(process, what, what == values.KILL || what == values.STOP)
	}
	return errors.New("Could not signal '%v' (#%v) %v because signalTarget %v is not supported.", service, process.Pid, what, signalTarget)
}

func sendSignalToProcess(process *os.Process, what values.Signal, tryGroup bool) error {
	if tryGroup {
		pgid, err := syscall.Getpgid(process.Pid)
		if err == nil {
			if syscall.Kill(-pgid, syscall.Signal(what)) == nil {
				return nil
			}
		}
	}
	process.Signal(syscall.Signal(what))
	return nil
}

func (instance *Service) createSysProcAttr() *syscall.SysProcAttr {
	// Prevents that a created process receives signals for caretakerd
	return &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}
}

func lookupUser(username string) (uid, gid int, err error) {
	u, err := lookupUserInPasswd(username)
	if err != nil {
		usernameParts := strings.SplitN(username, ":", 2)
		if len(usernameParts) == 2 {
			u = &user.User{
				Uid: usernameParts[0],
				Gid: usernameParts[1],
			}
		} else {
			return -1, -1, err
		}
	}
	uid, err = strconv.Atoi(u.Uid)
	if err != nil {
		return -1, -1, err
	}
	gid, err = strconv.Atoi(u.Gid)
	if err != nil {
		return -1, -1, err
	}
	return uid, gid, nil
}

func lookupUserInPasswd(uid string) (*user.User, error) {
	file, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(file), "\n") {
		data := strings.Split(line, ":")
		if len(data) > 5 && (data[0] == uid || data[2] == uid) {
			return &user.User{
				Uid:      data[2],
				Gid:      data[3],
				Username: data[0],
				Name:     data[4],
				HomeDir:  data[5],
			}, nil
		}
	}
	return nil, errors.New("User not found in /etc/passwd")
}
