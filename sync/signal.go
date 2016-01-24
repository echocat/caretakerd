package sync

import (
    "time"
    "runtime"
    "github.com/echocat/caretakerd/errors"
)

type Signal struct {
    sg      *SyncGroup
    channel chan bool
}

func (sg *SyncGroup) NewSignal() *Signal {
    result := &Signal{
        sg: sg,
        channel: make(chan bool),
    }
    runtime.SetFinalizer(result, finalizeSignalInstance)
    return result
}

func finalizeSignalInstance(s *Signal) {
    closeChannel(s.channel)
}

func (i *Signal) Wait(duration time.Duration) error {
    sg := (*i).sg
    sg.append(i)
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

func (i *Signal) Send() error {
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
    send := true
    for ; send ; {
        select {
        case i.channel <- true:
            send = true
        default:
            send = false
        }
    }
    return err
}

func (i *Signal) Interrupt() {
    closeChannel(i.channel)
}
