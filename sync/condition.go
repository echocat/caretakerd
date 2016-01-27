package sync

import (
	"time"
	"runtime"
	"errors"
)

type Condition struct {
	sg      *SyncGroup
	channel chan bool
	mutex   *Mutex
}

func (sg *SyncGroup) NewCondition(mutex *Mutex) *Condition {
	result := &Condition{
		sg: sg,
		channel: make(chan bool),
		mutex: mutex,
	}
	runtime.SetFinalizer(result, finalizeConditionInstance)
	return result
}

func finalizeConditionInstance(s *Condition) {
	closeChannel(s.channel)
}

func (i *Condition) Wait(duration time.Duration) error {
	sg := (*i).sg
	sg.append(i)
	i.doUnlock()
	defer i.doLock()
	select {
	case <-time.After(duration):
		return sg.removeAndReturn(i, TimeoutError{})
	case c := <-i.channel:
		var result error
		if ! c {
			result = InterruptedError{}
		}
		return sg.removeAndReturn(i, result)
	}
	return sg.removeAndReturn(i, InterruptedError{})
}

func (i *Condition) doLock() {
	i.mutex.Lock()
}

func (i *Condition) doUnlock() {
	i.mutex.Unlock()
}

func (i *Condition) send() (bool, error) {
	var err error
	defer func() {
		p := recover()
		if p != nil {
			if s, ok := p.(string); ok {
				if s != "send on closed channel" {
					panic(p)
				} else {
					err = errors.New("Signal interrupted.")
				}
			} else {
				panic(p)
			}
		}
	}()
	sent := true
	select {
	case i.channel <- true:
		sent = true
	default:
		sent = false
	}
	return sent, err
}

func (i *Condition) Send() error {
	_, err := i.send()
	return err
}

func (i *Condition) Broadcast() error {
	var err error
	doSend := true
	for doSend {
		doSend, err = i.send()
	}
	return err
}

func (i *Condition) Interrupt() {
	closeChannel(i.channel)
	i.mutex.Interrupt()
}

func (i *Condition) Mutex() *Mutex {
	return i.mutex
}
