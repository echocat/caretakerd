package sync

import (
	"runtime"
	"sync"
)

// Interruptable represents an object that could be interrupted.
type Interruptable interface {
	Interrupt()
}

// Group is a couple of tools (like sleep, locks, conditions, ...) that are grouped
// together and could be interrupted by calling Interrupt() method.
type Group struct {
	interruptables map[Interruptable]int
	lock           *sync.Mutex
}

// TimeoutError occurs if a timeout condition is reached.
type TimeoutError struct{}

func (instance TimeoutError) Error() string {
	return "Timeout."
}

// InterruptedError occurs if someone has called the Interrupt() method.
type InterruptedError struct{}

func (instance InterruptedError) Error() string {
	return "Interrupted."
}

// NewGroup creates a new SyncGroup instance.
func NewGroup() *Group {
	result := &Group{
		interruptables: map[Interruptable]int{},
		lock:           new(sync.Mutex),
	}
	runtime.SetFinalizer(result, finalizeSyncGroup)
	return result
}

func finalizeSyncGroup(instance *Group) {
	instance.Interrupt()
}

// Interrupt interrupts every action on this SyncGroup.
// After calling this method the instance is no longer usable anymore.
func (instance *Group) Interrupt() {
	for interruptable := range instance.interruptables {
		interruptable.Interrupt()
	}
}

func (instance Group) doUnlock() {
	instance.lock.Unlock()
}

// NewGroup creates a new sub instance of this instance.
func (instance *Group) NewGroup() *Group {
	result := NewGroup()
	instance.append(result)
	return result
}

func (instance *Group) append(what Interruptable) {
	instance.lock.Lock()
	defer instance.doUnlock()
	if existing, ok := instance.interruptables[what]; ok {
		instance.interruptables[what] = existing + 1
	} else {
		instance.interruptables[what] = 1
	}
}

func (instance *Group) removeAndReturn(what Interruptable, result error) error {
	instance.lock.Lock()
	defer instance.doUnlock()
	if existing, ok := instance.interruptables[what]; ok {
		if existing <= 1 {
			delete(instance.interruptables, what)
		} else {
			instance.interruptables[what] = existing - 1
		}
	}
	return result
}

func closeChannel(c chan bool) {
	defer func() {
		p := recover()
		if p != nil {
			if s, ok := p.(string); ok {
				if s != "close of closed channel" {
					panic(p)
				}
			} else {
				panic(p)
			}
		}

	}()
	close(c)
}
