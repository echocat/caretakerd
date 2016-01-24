// +build linux,darwin

package service

import (
    "github.com/echocat/caretakerd/service/signal"
)

func defaultStopSignal() signal.Signal {
    return signal.TERM
}
