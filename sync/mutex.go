package sync

import (
	"time"
	"runtime"
)

type Mutex struct {
	sg      *SyncGroup
	channel chan bool
}

func (sg *SyncGroup) NewMutex() *Mutex {
	result := &Mutex{
		sg: sg,
		channel: make(chan bool, 1),
	}
	runtime.SetFinalizer(result, finalizeMutexInstance)
	return result
}

func finalizeMutexInstance(s *Mutex) {
	closeChannel(s.channel)
}

func (i *Mutex) Lock() {
	i.channel <- true
}

func (i *Mutex) Unlock() {
	<-i.channel
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
