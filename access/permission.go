package access

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Permission represents the service's/node's permissions in caretakerd.
type Permission int

const (
	// @id forbidden
	//
	// The remote control/service does not have any permissions in caretakerd.
	Forbidden Permission = 0
	// @id readOnly
	//
	// The remote control/service does only have read permissions in caretakerd.
	ReadOnly Permission = 1
	// @id readWrite
	//
	// The remote control/service does have read and write permissions in caretakerd.
	ReadWrite Permission = 2
)

// AllPermissions contains all possible variants of Permission.
var AllPermissions = []Permission{
	Forbidden,
	ReadOnly,
	ReadWrite,
}

func (instance Permission) String() string {
	s, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return s
}

// CheckedString - Same as String but also returns an optional error message if errors occur.
// validation errors.
func (instance Permission) CheckedString() (string, error) {
	switch instance {
	case Forbidden:
		return "forbidden", nil
	case ReadOnly:
		return "readOnly", nil
	case ReadWrite:
		return "readWrite", nil
	}
	return "", fmt.Errorf("illegal permission: %d", instance)
}

// Set the given string to current object from a string.
// Return an error object if there are some problems while transforming the string.
func (instance *Permission) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllPermissions {
			if int(candidate) == valueAsInt {
				*instance = candidate
				return nil
			}
		}
		return fmt.Errorf("illegal permission: %v", value)
	}
	lowerValue := strings.ToLower(value)
	for _, candidate := range AllPermissions {
		if strings.ToLower(candidate.String()) == lowerValue {
			*instance = candidate
			return nil
		}
	}
	return fmt.Errorf("illegal permission: %v", value)
}

// MarshalYAML is used until yaml marshalling. Do not call directly.
func (instance Permission) MarshalYAML() (interface{}, error) {
	return instance.CheckedString()
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call directly.
func (instance *Permission) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call directly.
func (instance Permission) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

// UnmarshalJSON is used until json unmarshalling. Do not call directly.
func (instance *Permission) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate validates actions on the given object and returns an error object if errors occur.
func (instance Permission) Validate() error {
	_, err := instance.CheckedString()
	return err
}
