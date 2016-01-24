package sync

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
)

func TestSleep_Interrupt(t *testing.T) {
    sg := NewSyncGroup().NewSyncGroup()
    start := time.Now()
    go func() {
        time.Sleep(10 * time.Millisecond)
        sg.Interrupt()
    }()
    sg.Sleep(10 * time.Second)
    duration := time.Since(start) / time.Millisecond
    assert.Equal(t, true, duration < 50, "Expected was not longer then 50ms but was %d", duration)
}
