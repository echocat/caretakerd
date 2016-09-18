package values

import "strings"

// String represents a string with more features as the primitive type.
// @inline
type String string

func (instance String) String() string {
	s, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return s
}

// CheckedString is like String but also returns an optional error if there are any
// validation errors.
func (instance String) CheckedString() (string, error) {
	return string(instance), nil
}

// Set sets the given string to the current object from a string.
// Returns an error object if there are any problems while transforming the string.
func (instance *String) Set(value string) error {
	(*instance) = String(value)
	return nil
}

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (instance String) MarshalYAML() (interface{}, error) {
	return string(instance), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *String) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate validates actions on this object and returns an error object if there are any.
func (instance String) Validate() error {
	_, err := instance.CheckedString()
	return err
}

// IsEmpty returns "true" if the current string instance has no content.
func (instance String) IsEmpty() bool {
	return len(instance) <= 0
}

// IsTrimmedEmpty returns "true" if the current string instance has no trimmed content.
func (instance String) IsTrimmedEmpty() bool {
	return len(strings.TrimSpace(instance.String())) <= 0
}
