// +build windows

package config

import (
    "github.com/echocat/caretakerd/service/signal"
)

func defaultStopSignal() signal.Signal {
    return signal.KILL
}
