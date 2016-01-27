// +build linux darwin

package service

import (
	"os/user"
	"syscall"
	"os/exec"
	"github.com/echocat/caretakerd/values"
	"github.com/echocat/caretakerd/panics"
)

func serviceHandleUsersFor(service *Service, cmd *exec.Cmd) {
	config := (*service).config
	userName := config.User
	if !userName.IsTrimmedEmpty() {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		userId, err := user.Lookup(userName.String())
		if err != nil {
			panics.New("Could not run as user '%v'.", userName).CausedBy(err).Throw()
		}
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: userId.Uid, Gid: userId.Gid}
	}
}

func sendSignalToService(service *Service, process *os.Process, what values.Signal) error {
	signal.RecordSendSignal(what)
	process.Signal(syscall.Signal(what))
}
