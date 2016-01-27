// +build linux darwin

package service

import (
	"github.com/echocat/caretakerd/panics"
	"github.com/echocat/caretakerd/values"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func lookupUser(username string) (uid, gid int, err error) {
	u, err := user.Lookup(username)
	if err != nil {
		return -1, -1, err
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

func serviceHandleUsersFor(service *Service, cmd *exec.Cmd) {
	config := (*service).config
	userName := config.User
	if !userName.IsTrimmedEmpty() {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		//		u, err := userLookup(userName.String())
		uid, gid, err := lookupUser(userName.String())
		if err != nil {
			panics.New("Could not run as user '%v'.", userName).CausedBy(err).Throw()
		}
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	}
}

func sendSignalToService(service *Service, process *os.Process, what values.Signal) error {
	values.RecordSendSignal(what)
	process.Signal(syscall.Signal(what))
	return nil
}
