package values

import (
	"github.com/echocat/caretakerd/errors"
	"strconv"
)

// ExitCode represents an exitCode of a command.
// @inline
type ExitCode int

func (i ExitCode) String() string {
	return strconv.Itoa(int(i))
}

// Set sets the given string to current object from a string.
// Returns an error object if there are any problems while transforming the string.
func (i *ExitCode) Set(value string) error {
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("Illegal exit Code value: " + value)
	}
	(*i) = ExitCode(valueAsInt)
	return nil
}

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (i ExitCode) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (i *ExitCode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return i.Set(value)
}

// Validate validates actions on this object and returns an error object if there are any.
func (i ExitCode) Validate() {}
