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

// CheckedString is like String but also returns an optional error if there are any
// validation errors.
func (instance Boolean) CheckedString() (string, error) {
	if instance == Boolean(true) {
		return "true", nil
	}
	return "false", nil
}

// Sets the given string to current object from a string.
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

// SetFromBool sets the given boolean value to the current object.
func (instance *Boolean) SetFromBool(value bool) error {
	(*instance) = Boolean(value)
	return nil
}

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (instance Boolean) MarshalYAML() (interface{}, error) {
	return bool(instance), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *Boolean) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value bool
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.SetFromBool(value)
}

// Validate validates actions on this object and returns an error object if there are any.
func (instance Boolean) Validate() error {
	_, err := instance.CheckedString()
	return err
}
