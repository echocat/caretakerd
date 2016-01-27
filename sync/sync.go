package sync

import (
	"sync"
	"runtime"
)

type Interruptable interface {
	Interrupt()
}

type SyncGroup struct {
	interruptables map[Interruptable]int
	lock           *sync.Mutex
}

type TimeoutError struct{}

func (instance TimeoutError) Error() string {
	return "Timeout."
}

type InterruptedError struct{}

func (instance InterruptedError) Error() string {
	return "Interrupted."
}

func NewSyncGroup() *SyncGroup {
	result := &SyncGroup{
		interruptables: map[Interruptable]int{},
		lock: new(sync.Mutex),
	}
	runtime.SetFinalizer(result, finalizeSyncGroup)
	return result
}

func finalizeSyncGroup(sg *SyncGroup) {
	sg.Interrupt()
}

func (sg *SyncGroup) Interrupt() {
	for interruptable, _ := range sg.interruptables {
		interruptable.Interrupt()
	}
}

func (sg SyncGroup) doUnlock() {
	sg.lock.Unlock()
}

func (sg *SyncGroup) NewSyncGroup() *SyncGroup {
	result := NewSyncGroup()
	sg.append(result)
	return result
}

func (sg *SyncGroup) append(what Interruptable) {
	sg.lock.Lock()
	defer sg.doUnlock()
	if existing, ok := sg.interruptables[what]; ok {
		sg.interruptables[what] = existing + 1
	} else {
		sg.interruptables[what] = 1
	}
}

func (sg *SyncGroup) removeAndReturn(what Interruptable, result error) (error) {
	sg.lock.Lock()
	defer sg.doUnlock()
	if existing, ok := sg.interruptables[what]; ok {
		if existing <= 1 {
			delete(sg.interruptables, what)
		} else {
			sg.interruptables[what] = existing - 1
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
