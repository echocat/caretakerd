package values

import (
	"github.com/echocat/caretakerd/errors"
	"strconv"
)

// Integer represents an int with more features as the primitive type.
// @inline
type Integer int

func (instance Integer) String() string {
	s, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return s
}

// CheckedString is like String but return also an optional error if there are some
// validation errors.
func (instance Integer) CheckedString() (string, error) {
	return strconv.Itoa(int(instance)), nil
}

// Set the given string to current object from a string.
// Return an error object if there are some problems while transforming the string.
func (instance *Integer) Set(value string) error {
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("Illegal integer value: " + value)
	}
	return instance.SetFromInt(valueAsInt)
}

// SetFromInt try to set the given int value to this instance.
// Return an error object if there are some problems while transforming the plain int.
func (instance *Integer) SetFromInt(value int) error {
	(*instance) = Integer(value)
	return nil
}

// MarshalYAML is used until yaml marshalling. Do not call directly.
func (instance Integer) MarshalYAML() (interface{}, error) {
	return int(instance), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call directly.
func (instance *Integer) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value int
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.SetFromInt(value)
}

// Int return this value as int.
func (instance Integer) Int() int {
	return int(instance)
}

// Validate do validate action on this object and return an error object if any.
func (instance Integer) Validate() error {
	_, err := instance.CheckedString()
	return err
}
