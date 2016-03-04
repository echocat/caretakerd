package sync

import (
	"github.com/echocat/caretakerd/errors"
	"runtime"
	"time"
)

type Mutex struct {
	sg      *SyncGroup
	channel chan bool
}

func (sg *SyncGroup) NewMutex() *Mutex {
	result := &Mutex{
		sg:      sg,
		channel: make(chan bool, 1),
	}
	runtime.SetFinalizer(result, finalizeMutexInstance)
	return result
}

func finalizeMutexInstance(s *Mutex) {
	closeChannel(s.channel)
}

func (i *Mutex) Lock() error {
	var err error
	defer func() {
		p := recover()
		if p != nil {
			if s, ok := p.(string); ok {
				if s != "send on closed channel" {
					panic(p)
				} else {
					err = errors.New("Lock interrupted.")
				}
			} else {
				panic(p)
			}
		}
	}()
	select {
	case i.channel <- true:
		return nil
	default:
		return errors.New("Lock interrupted.")
	}
	return err
}

func (i *Mutex) Unlock() {
	select {
	case <-i.channel:
		return
	}
	return
}

func (i *Mutex) TryLock(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer func() {
		timer.Stop()
	}()
	select {
	case i.channel <- true:
		return true
	case <-timer.C:
	}
	return false
}

func (i *Mutex) Interrupt() {
	closeChannel(i.channel)
}
