package keyStore

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"strconv"
	"strings"
)

// ## Description
//
// Represents the type of the keyStore.
type Type int

const (
	// Indicates that caretakerd have to generate its own keyStore on startup.
	// This is the best solution in most cases.
	Generated Type = 0

	// Load keyStore from a provided PEM file.
	// If instance type is selected instance file have to be provided.
	FromFile Type = 1

	// Load keyStore from the environment variable ``CTD_PEM`` in PEM format.
	// If instance type is selected instance variable have to be provided.
	FromEnvironment Type = 2
)

var AllTypes = []Type{
	Generated,
	FromFile,
	FromEnvironment,
}

func (i Type) String() string {
	result, err := i.CheckedString()
	if err != nil {
		panic(err)
	}
	return result
}

func (i Type) CheckedString() (string, error) {
	switch i {
	case Generated:
		return "generated", nil
	case FromFile:
		return "fromFile", nil
	case FromEnvironment:
		return "fromEnvironment", nil
	}
	return "", errors.New("Illegal keyStore type: %d", i)
}

func (i *Type) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllTypes {
			if int(candidate) == valueAsInt {
				(*i) = candidate
				return nil
			}
		}
		return errors.New("Illegal keyStore type: " + value)
	} else {
		lowerValue := strings.ToLower(value)
		for _, candidate := range AllTypes {
			if strings.ToLower(candidate.String()) == lowerValue {
				(*i) = candidate
				return nil
			}
		}
		return errors.New("Illegal keyStore type: " + value)
	}
}

func (i Type) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

func (i *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return i.Set(value)
}

func (i Type) MarshalJSON() ([]byte, error) {
	s, err := i.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

func (i *Type) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return i.Set(value)
}

func (i Type) IsTakingFilename() bool {
	return i == FromFile
}

func (i Type) IsGenerating() bool {
	return i == Generated
}

func (i Type) IsConsumingCaFile() bool {
	return i == FromFile || i == FromEnvironment
}

func (i Type) Validate() error {
	_, err := i.CheckedString()
	return err
}
