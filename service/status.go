package service

import (
	"github.com/echocat/caretakerd/panics"
	"github.com/echocat/caretakerd/errors"
	"strconv"
	"strings"
	"encoding/json"
)

type Status int

const (
	Down = Status(0)
	Running = Status(1)
	Stopped = Status(2)
	Killed = Status(3)
)

var AllStatus = []Status{
	Down,
	Running,
	Stopped,
	Killed,
}

func (instance Status) String() string {
	switch instance {
	case Down:
		return "down"
	case Running:
		return "running"
	case Stopped:
		return "stopped"
	case Killed:
		return "killed"
	}
	panic(panics.New("Illegal status: %d", instance))
}

func (instance *Status) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllStatus {
			if int(candidate) == valueAsInt {
				(*instance) = candidate
				return nil
			}
		}
		return errors.New("Illegal status: " + value)
	} else {
		lowerValue := strings.ToLower(value)
		for _, candidate := range AllStatus {
			if strings.ToLower(candidate.String()) == lowerValue {
				(*instance) = candidate
				return nil
			}
		}
		return errors.New("Illegal status: " + value)
	}
}

func (instance Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(instance.String())
}

func (instance *Status) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

func (instance Status) Validate() {
	instance.String()
}

func (instance Status) IsGoDownRequest() bool {
	switch instance {
	case Stopped: fallthrough
	case Killed:
		return true
	}
	return false
}
