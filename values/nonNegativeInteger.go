package values

import (
	"github.com/echocat/caretakerd/errors"
	"strconv"
)

// NonNegativeInteger represents an int with more features as the primitive type and which could not be negative.
// @inline
type NonNegativeInteger int

func (i NonNegativeInteger) String() string {
	result, err := i.CheckedString()
	if err != nil {
		panic(err)
	}
	return result
}

// CheckedString is like String but return also an optional error if there are some
// validation errors.
func (i NonNegativeInteger) CheckedString() (string, error) {
	return strconv.Itoa(int(i)), nil
}

// Set the given string to current object from a string.
// Return an error object if there are some problems while transforming the string.
func (i *NonNegativeInteger) Set(value string) error {
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("Illegal integer value: " + value)
	}
	return i.SetFromInt(valueAsInt)
}

// SetFromInt try to set the given int value to this instance.
// Return an error object if there are some problems while transforming the plain int.
func (i *NonNegativeInteger) SetFromInt(value int) error {
	if value < 0 {
		return errors.New("This intger value should not be negative. But got: %v", value)
	}
	(*i) = NonNegativeInteger(value)
	return nil
}

// MarshalYAML is used until yaml marshalling. Do not call directly.
func (i NonNegativeInteger) MarshalYAML() (interface{}, error) {
	return int(i), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call directly.
func (i *NonNegativeInteger) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value int
	if err := unmarshal(&value); err != nil {
		return err
	}
	return i.SetFromInt(value)
}

// Int return this value as int.
func (i NonNegativeInteger) Int() int {
	return int(i)
}

// Validate do validate action on this object and return an error object if any.
func (i NonNegativeInteger) Validate() error {
	_, err := i.CheckedString()
	return err
}
