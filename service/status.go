package service

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/panics"
	"strconv"
	"strings"
)

type Status int

const (
	New     = Status(0)
	Down    = Status(1)
	Running = Status(2)
	Stopped = Status(3)
	Killed  = Status(4)
	Unknown = Status(5)
)

var AllStatus = []Status{
	New,
	Down,
	Running,
	Stopped,
	Killed,
	Unknown,
}

func (instance Status) String() string {
	switch instance {
	case New:
		return "new"
	case Down:
		return "down"
	case Running:
		return "running"
	case Stopped:
		return "stopped"
	case Killed:
		return "killed"
	case Unknown:
		return "unknown"
	}
	panic(panics.New("Illegal status: %d", instance))
}

// Set the given string to current object from a string.
// Return an error object if there are some problems while transforming the string.
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

// MarshalJSON is used until json marshalling. Do not call directly.
func (instance Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(instance.String())
}

// UnmarshalJSON is used until json unmarshalling. Do not call directly.
func (instance *Status) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate do validate action on this object and return an error object if any.
func (instance Status) Validate() {
	instance.String()
}

func (instance Status) IsGoDownRequest() bool {
	switch instance {
	case Stopped:
		fallthrough
	case Killed:
		return true
	}
	return false
}
