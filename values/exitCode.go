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

// Set the given string to current object from a string.
// Return an error object if there are some problems while transforming the string.
func (i *ExitCode) Set(value string) error {
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("Illegal exit Code value: " + value)
	}
	(*i) = ExitCode(valueAsInt)
	return nil
}

// MarshalYAML is used until yaml marshalling. Do not call directly.
func (i ExitCode) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call directly.
func (i *ExitCode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return i.Set(value)
}

// Validate do validate action on this object and return an error object if any.
func (i ExitCode) Validate() {
	i.String()
}
