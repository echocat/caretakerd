package sync

import (
	"errors"
	"runtime"
	"time"
)

// Condition could be used to wait for a synced change of something.
type Condition struct {
	sg      *Group
	channel chan bool
	mutex   *Mutex
}

// NewCondition creates a new condition in the current SyncGroup with the given Mutex.
func (instance *Group) NewCondition(mutex *Mutex) *Condition {
	result := &Condition{
		sg:      instance,
		channel: make(chan bool),
		mutex:   mutex,
	}
	runtime.SetFinalizer(result, finalizeConditionInstance)
	return result
}

func finalizeConditionInstance(condition *Condition) {
	closeChannel(condition.channel)
}

// Wait waits for someone that calls Send() or Broadcast() on this Condition instance for the given maximum duration.
// If someone calls Interrupt() or there is no trigger received within the maximum duration an error will be returned.
func (instance *Condition) Wait(duration time.Duration) error {
	return instance.wait(duration, true)
}

func (instance *Condition) wait(duration time.Duration, guarded bool) error {
	sg := (*instance).sg
	sg.append(instance)
	if guarded {
		instance.doUnlock()
		defer instance.doLock()
	}
	select {
	case <-time.After(duration):
		return sg.removeAndReturn(instance, TimeoutError{})
	case c := <-instance.channel:
		var result error
		if !c {
			result = InterruptedError{}
		}
		return sg.removeAndReturn(instance, result)
	}
}

func (instance *Condition) doLock() error {
	return instance.mutex.Lock()
}

func (instance *Condition) doUnlock() {
	instance.mutex.Unlock()
}

func (instance *Condition) send() (bool, error) {
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
	case instance.channel <- true:
		sent = true
	default:
		sent = false
	}
	return sent, err
}

// Send sends ONE trigger to the condition.
// If there are more then one listener only one of them will receive the trigger.
// This method is not blocking.
func (instance *Condition) Send() error {
	_, err := instance.send()
	return err
}

// Broadcast broadcast trigger to the condition.
// If there are more then one listener every of them will receive the trigger.
// This method is not blocking.
func (instance *Condition) Broadcast() error {
	var err error
	doSend := true
	for doSend {
		doSend, err = instance.send()
	}
	return err
}

// Interrupt interrupts every possible current running Wait() method of this instance.
// Nobody is able to call Wait() from this moment on anymore of this instance.
func (instance *Condition) Interrupt() {
	closeChannel(instance.channel)
	instance.mutex.Interrupt()
}

// Mutex returns the instance of this Condition.
func (instance *Condition) Mutex() *Mutex {
	return instance.mutex
}
