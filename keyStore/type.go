package keyStore

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"strconv"
	"strings"
)

// # Description
//
// Represents the type of the keyStore.
type Type int

const (
	// @id generated
	// Indicates that caretakerd have to generate its own keyStore on startup.
	// This is the best solution in most cases.
	Generated Type = 0

	// @id fromFile
	// Load keyStore from a provided PEM file.
	// If instance type is selected instance file have to be provided.
	FromFile Type = 1

	// @id fromEnvironment
	// Load keyStore from the environment variable ``CTD_PEM`` in PEM format.
	// If instance type is selected instance variable have to be provided.
	FromEnvironment Type = 2
)

// AllTypes contains all possible variants of Type.
var AllTypes = []Type{
	Generated,
	FromFile,
	FromEnvironment,
}

func (instance Type) String() string {
	result, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return result
}

// CheckedString is like String but return also an optional error if there are some
// validation errors.
func (instance Type) CheckedString() (string, error) {
	switch instance {
	case Generated:
		return "generated", nil
	case FromFile:
		return "fromFile", nil
	case FromEnvironment:
		return "fromEnvironment", nil
	}
	return "", errors.New("Illegal keyStore type: %d", instance)
}

// Set the given string to current object from a string.
// Return an error object if there are some problems while transforming the string.
func (instance *Type) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllTypes {
			if int(candidate) == valueAsInt {
				(*instance) = candidate
				return nil
			}
		}
		return errors.New("Illegal keyStore type: " + value)
	}
	lowerValue := strings.ToLower(value)
	for _, candidate := range AllTypes {
		if strings.ToLower(candidate.String()) == lowerValue {
			(*instance) = candidate
			return nil
		}
	}
	return errors.New("Illegal keyStore type: " + value)
}

// MarshalYAML is used until yaml marshalling. Do not call directly.
func (instance Type) MarshalYAML() (interface{}, error) {
	return instance.String(), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call directly.
func (instance *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call directly.
func (instance Type) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

// UnmarshalJSON is used until json unmarshalling. Do not call directly.
func (instance *Type) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// IsTakingFilename returns true if the KeyStore instance created with this type is created from file.
func (instance Type) IsTakingFilename() bool {
	return instance == FromFile
}

// IsGenerating returns true if the KeyStore instance created with this type will be generated.
func (instance Type) IsGenerating() bool {
	return instance == Generated
}

// IsConsumingCAFile returns true if the KeyStore instance created with this type can consumes a CA bundle file.
func (instance Type) IsConsumingCAFile() bool {
	return instance == FromFile || instance == FromEnvironment
}

// Validate do validate action on this object and return an error object if any.
func (instance Type) Validate() error {
	_, err := instance.CheckedString()
	return err
}
