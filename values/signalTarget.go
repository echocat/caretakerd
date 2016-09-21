package values

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"strconv"
	"strings"
)

// SignalTarget defines who have to receives a signal.
type SignalTarget int

const (
	// @id process
	//
	// Send a signal to the process only.
	Process SignalTarget = 1
	// @id processGroup
	//
	// Send a signal to the whole process group.
	ProcessGroup SignalTarget = 2
	// @id mixed
	//
	// Send every signal only to a process - except KILL and STOP. This are send to the processGroup.
	Mixed SignalTarget = 3
)

// AllSignalTargets contains all possible variants of SignalTarget.
var AllSignalTargets = []SignalTarget{
	Process,
	ProcessGroup,
	Mixed,
}

func (instance SignalTarget) String() string {
	result, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return result
}

// CheckedString is like String but also returns an optional error if there are any
// validation errors.
func (instance SignalTarget) CheckedString() (string, error) {
	switch instance {
	case Process:
		return "process", nil
	case ProcessGroup:
		return "processGroup", nil
	}
	return "", errors.New("Illegal signal target: %d", instance)
}

// Set sets the given string to current object from a string.
// Returns an error object if there are any problems while transforming the string.
func (instance *SignalTarget) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllSignalTargets {
			if int(candidate) == valueAsInt {
				(*instance) = candidate
				return nil
			}
		}
		return errors.New("Illegal signal target: " + value)
	}
	lowerValue := strings.ToLower(value)
	for _, candidate := range AllSignalTargets {
		if strings.ToLower(candidate.String()) == lowerValue {
			(*instance) = candidate
			return nil
		}
	}
	return errors.New("Illegal signal target: " + value)
}

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (instance SignalTarget) MarshalYAML() (interface{}, error) {
	return instance.String(), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *SignalTarget) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call this method directly.
func (instance SignalTarget) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

// UnmarshalJSON is used until json unmarshalling. Do not call this method directly.
func (instance *SignalTarget) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate validates actions on this object and returns an error object if there are any.
func (instance SignalTarget) Validate() error {
	_, err := instance.CheckedString()
	return err
}
