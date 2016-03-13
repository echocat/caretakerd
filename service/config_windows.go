// +build windows

package service

import (
	"github.com/echocat/caretakerd/values"
)

func defaultStopSignal() values.Signal {
	return values.TERM
}
