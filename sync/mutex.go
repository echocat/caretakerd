package sync

import (
	"github.com/echocat/caretakerd/errors"
	"runtime"
	"time"
)

// Mutex is a lock in a SyncGroup.
type Mutex struct {
	sg      *Group
	channel chan bool
}

// NewMutex creates a new instance of Mutex in the current SyncGroup.
func (instance *Group) NewMutex() *Mutex {
	result := &Mutex{
		sg:      instance,
		channel: make(chan bool, 1),
	}
	runtime.SetFinalizer(result, finalizeMutexInstance)
	return result
}

func finalizeMutexInstance(mutex *Mutex) {
	closeChannel(mutex.channel)
}

// Lock locks the current thread to this mutex.
// If this is not possible an error will be returned.
// This method is blocking until locking is possible.
func (instance *Mutex) Lock() error {
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
	case instance.channel <- true:
		return nil
	default:
		if err != nil {
			return err
		}
		return errors.New("Lock interrupted.")
	}
}

// Unlock unlocks the current thread from this Mutex.
func (instance *Mutex) Unlock() {
	select {
	case <-instance.channel:
		return
	}
}

// TryLock tries to locks the current thread to this mutex.
// This method will wait for a maximum of the given duration to get the lock
// - in this case "true" is returned.
func (instance *Mutex) TryLock(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer func() {
		timer.Stop()
	}()
	select {
	case instance.channel <- true:
		return true
	case <-timer.C:
	}
	return false
}

// Interrupt interrupts every possible current running Lock() and TryLock() method of this instance.
// In this instance, nobody will be able to call Lock() and TryLock() from this moment on.
func (instance *Mutex) Interrupt() {
	closeChannel(instance.channel)
}
