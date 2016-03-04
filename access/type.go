package access

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"strconv"
	"strings"
)

type Type int

const (
	// @id none
	//
	// Like the name said.
	None Type = 0
	// @id trusted
	//
	// Means that caretakerd trust the remote connection based on remote name and configured {@ref github.com/echocat/caretakerd/keyStore.Config#CaFile}.
	// Or if {@ref github.com/echocat/caretakerd/access.Config#PemFile} is specified expect exact this identity.
	Trusted Type = 1
	// @id generateToEnvironment
	//
	// Generate a new certificate to environment variable ``CTD_PEM`` and trust it.
	GenerateToEnvironment Type = 2
	// @id generateToFile
	//
	// Generate a new certificate to configured {@ref github.com/echocat/caretakerd/access.Config#PemFile} and trust it.
	GenerateToFile Type = 3
)

var AllTypes = []Type{
	None,
	Trusted,
	GenerateToEnvironment,
	GenerateToFile,
}

func (instance Type) String() string {
	s, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return s
}

// CheckedString is like String but return also an optional error if there are some
// validation errors.
func (instance Type) CheckedString() (string, error) {
	switch instance {
	case None:
		return "none", nil
	case Trusted:
		return "trusted", nil
	case GenerateToEnvironment:
		return "generateToEnvironment", nil
	case GenerateToFile:
		return "generateToFile", nil
	}
	return "", errors.New("Illegal access type: %d", instance)
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
		return errors.New("Illegal access type: " + value)
	} else {
		lowerValue := strings.ToLower(value)
		for _, candidate := range AllTypes {
			if strings.ToLower(candidate.String()) == lowerValue {
				(*instance) = candidate
				return nil
			}
		}
		return errors.New("Illegal access type: " + value)
	}
}

// MarshalYAML is used until yaml marshalling. Do not call directly.
func (instance Type) MarshalYAML() (interface{}, error) {
	return instance.CheckedString()
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

func (instance Type) IsTakingFilename() bool {
	return instance == GenerateToFile
}

func (instance Type) IsTakingFilePermission() bool {
	return instance == GenerateToFile
}

func (instance Type) IsTakingFileUser() bool {
	return instance == GenerateToFile
}

func (instance Type) IsTakingGroup() bool {
	return instance == GenerateToFile
}

func (instance Type) IsGenerating() bool {
	return instance == GenerateToFile || instance == GenerateToEnvironment
}

// Validate do validate action on this object and return an error object if any.
func (instance Type) Validate() error {
	_, err := instance.CheckedString()
	return err
}
