package service

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"strconv"
	"strings"
)

// Status represents a status of a service execution.
type Status int

const (
	// New indicates that the service execution was created but not yet started.
	New = Status(0)
	// Down indicates that the service execution was run and is now done.
	Down = Status(1)
	// Running indicates that the service execution is still running.
	Running = Status(2)
	// Stopped indicates that the service execution is still running but a stop was initiated.
	Stopped = Status(3)
	// Killed indicates that the service execution is still running but a kill was initiated.
	Killed = Status(4)
	// Unknown indicates a status that should never happen. Please consult the log for more information.
	Unknown = Status(5)
)

// AllStatus contains all possible variants of Status.
var AllStatus = []Status{
	New,
	Down,
	Running,
	Stopped,
	Killed,
	Unknown,
}

func (instance Status) String() string {
	s, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return s
}

// CheckedString is like String but also returns an optional error if there are any
// validation errors.
func (instance Status) CheckedString() (string, error) {
	switch instance {
	case New:
		return "new", nil
	case Down:
		return "down", nil
	case Running:
		return "running", nil
	case Stopped:
		return "stopped", nil
	case Killed:
		return "killed", nil
	case Unknown:
		return "unknown", nil
	}
	return "", errors.New("Illegal status: %d", instance)
}

// Set sets the given string to the current object from a string.
// Returns an error object if there are any problems while transforming the string.
func (instance *Status) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllStatus {
			if int(candidate) == valueAsInt {
				*instance = candidate
				return nil
			}
		}
		return errors.New("Illegal status: " + value)
	}
	lowerValue := strings.ToLower(value)
	for _, candidate := range AllStatus {
		if strings.ToLower(candidate.String()) == lowerValue {
			*instance = candidate
			return nil
		}
	}
	return errors.New("Illegal status: " + value)
}

// MarshalJSON is used until json marshalling. Do not call this method directly.
func (instance Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(instance.String())
}

// UnmarshalJSON is used until json unmarshalling. Do not call this method directly.
func (instance *Status) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate validates actions on this object and returns an error object if there are any.
func (instance Status) Validate() error {
	_, err := instance.CheckedString()
	return err
}

// IsGoDownRequest returns "true" if the status indicates that the service must be stopped.
func (instance Status) IsGoDownRequest() bool {
	switch instance {
	case Stopped:
		fallthrough
	case Killed:
		return true
	}
	return false
}
