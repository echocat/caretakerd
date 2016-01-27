package sync

import (
	"time"
)

func (sg *SyncGroup) Sleep(duration time.Duration) error {
	signal := sg.NewSignal()
	err := signal.Wait(duration)
	if _, ok := err.(TimeoutError); ok {
		return nil
	}
	return err
}
