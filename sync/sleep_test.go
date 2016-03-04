package sync

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSleep_Interrupt(t *testing.T) {
	sg := NewGroup().NewGroup()
	start := time.Now()
	go func() {
		time.Sleep(10 * time.Millisecond)
		sg.Interrupt()
	}()
	sg.Sleep(10 * time.Second)
	duration := time.Since(start) / time.Millisecond
	assert.Equal(t, true, duration < 50, "Expected was not longer then 50ms but was %d", duration)
}
