// +build linux,darwin

package service

import (
    "os/user"
    "syscall"
    "os/exec"
    "strings"
    "github.com/echocat/caretakerd/values"
    "github.com/echocat/caretakerd/panics"
)

func serviceHandleUsersFor(service *Service, cmd *exec.Cmd) {
    config := (*service).config
    u := config.User
    if len(strings.TrimSpace(u)) > 0 {
        cmd.SysProcAttr = &syscall.SysProcAttr{}
        u, err := user.Lookup(u)
        if err != nil {
            panics.New("Could not run as user '%s'.", u).CausedBy(err).Throw()
        }
        cmd.SysProcAttr.Credential = &syscall.Credential{Uid: u.Uid, Gid: u.Gid}
    }
}

func sendSignalToService(service *Service, process *os.Process, what values.Signal) error {
    signal.RecordSendSignal(what)
    process.Signal(syscall.Signal(what))
}
