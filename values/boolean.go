package values

import (
	"strings"
)

// Boolean represents a boolean with more features as the primitive type.
// @inline
type Boolean bool

func (instance Boolean) String() string {
	s, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return s
}

// CheckedString is like String but return also an optional error if there are some
// validation errors.
func (instance Boolean) CheckedString() (string, error) {
	if instance == Boolean(true) {
		return "true", nil
	}
	return "false", nil
}

// Set the given string to current object from a string.
func (instance *Boolean) Set(value string) error {
	switch strings.ToLower(value) {
	case "1":
		fallthrough
	case "on":
		fallthrough
	case "yes":
		fallthrough
	case "y":
		fallthrough
	case "true":
		return instance.SetFromBool(true)
	}
	return instance.SetFromBool(false)
}

// SetFromBool set given boolean to current object.
func (instance *Boolean) SetFromBool(value bool) error {
	(*instance) = Boolean(value)
	return nil
}

// MarshalYAML is used until yaml marshalling. Do not call directly.
func (instance Boolean) MarshalYAML() (interface{}, error) {
	return bool(instance), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call directly.
func (instance *Boolean) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value bool
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.SetFromBool(value)
}

// Validate do validate action on this object and return an error object if any.
func (instance Boolean) Validate() error {
	_, err := instance.CheckedString()
	return err
}
