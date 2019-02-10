package logger

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"strconv"
	"strings"
)

// # Description
//
// Represents a level for logging with a {@ref .Config Logger}
type Level int

const (
	// @id debug
	// Used for debugging purposes. This level is only required if something goes wrong and you need more information.
	Debug Level = 200

	// @id info
	// This is the regular level. Every normal message will be logged with instance level.
	Info Level = 300

	// @id warning
	// If a problem appears but the program is still able to continue its work, this instance level is used.
	Warning Level = 400

	// @id error
	// If a problem appears and the program is not longer able to continue its work, this instance level is used.
	Error Level = 500

	// @id fatal
	// This level is used on fatal problems.
	Fatal Level = 600
)

// AllLevels contains all possible variants of Level.
var AllLevels = []Level{
	Debug,
	Info,
	Warning,
	Error,
	Fatal,
}

func (instance Level) String() string {
	result, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return result
}

// CheckedString is like String but also returns an optional error if there are any
// validation errors.
func (instance Level) CheckedString() (string, error) {
	switch instance {
	case Debug:
		return "debug", nil
	case Info:
		return "info", nil
	case Warning:
		return "warning", nil
	case Error:
		return "error", nil
	case Fatal:
		return "fatal", nil
	}
	return strconv.Itoa(int(instance)), nil
}

// DisplayForLogging returns a string that can be used to display this level in log messages.
func (instance Level) DisplayForLogging() string {
	if instance == Warning {
		return "WARN"
	}
	return strings.ToUpper(instance.String())
}

// Set the given string to the current object from a string.
// Returns an error object if there are any problems while transforming the string.
func (instance *Level) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllLevels {
			if int(candidate) == valueAsInt {
				*instance = candidate
				return nil
			}
		}
		return errors.New("Illegal level: " + value)
	}
	lowerValue := strings.ToLower(value)
	switch lowerValue {
	case "warn":
		*instance = Warning
		return nil
	case "err":
		*instance = Error
		return nil
	}
	for _, candidate := range AllLevels {
		if candidate.String() == lowerValue {
			*instance = candidate
			return nil
		}
	}
	return errors.New("Illegal level: " + value)
}

// IsIndicatingProblem returns "true" if this level indicates a problem.
func (instance Level) IsIndicatingProblem() bool {
	return instance == Warning || instance == Error || instance == Fatal
}

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (instance Level) MarshalYAML() (interface{}, error) {
	return instance.CheckedString()
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *Level) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call this method directly.
func (instance Level) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

// UnmarshalJSON is used until json unmarshalling. Do not call this method directly.
func (instance *Level) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate validates actions on this object and returns an error object if there any.
func (instance Level) Validate() error {
	_, err := instance.CheckedString()
	return err
}
