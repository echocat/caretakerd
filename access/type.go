package access

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/echocat/caretakerd/errors"
)

// Type indicates the validation type for the access of a service/node to caretakerd.

type Type int

const (
	// @id none
	//
	// No ID given
	None Type = 0
	// @id trusted
	//
	// caretakerd trusts the remote connection based on the remote name and the configured {@ref github.com/echocat/caretakerd/keyStore.Config#CaFile}.
	// or if the {@ref github.com/echocat/caretakerd/access.Config#PemFile} is specified to expect exactly this identity.
	Trusted Type = 1
	// @id generateToEnvironment
	//
	// Generates a new certificate to the environment variable ``CTD_PEM`` and trusts it.
	GenerateToEnvironment Type = 2
	// @id generateToFile
	//
	// Generates a new certificate to the configured {@ref github.com/echocat/caretakerd/access.Config#PemFile} and trusts it.
	GenerateToFile Type = 3
)

// AllTypes contains all possible variants of Type.
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

// CheckedString is like String but also returns an optional error if there are
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

// Set sets the given string to the current object from a string.
// Returns an error object if there are problems while transforming the string.
func (instance *Type) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllTypes {
			if int(candidate) == valueAsInt {
				*instance = candidate
				return nil
			}
		}
		return fmt.Errorf("illegal access type: %v", value)
	}
	lowerValue := strings.ToLower(value)
	for _, candidate := range AllTypes {
		if strings.ToLower(candidate.String()) == lowerValue {
			*instance = candidate
			return nil
		}
	}
	return fmt.Errorf("illegal permission type: %v", value)
}

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (instance Type) MarshalYAML() (interface{}, error) {
	return instance.CheckedString()
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call this method directly.
func (instance Type) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

// UnmarshalJSON is used until json unmarshalling. Do not call this method directly.
func (instance *Type) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// IsTakingFilename returns true if this Type indicates that it accepts a file name.
func (instance Type) IsTakingFilename() bool {
	return instance == GenerateToFile
}

// IsTakingFilePermission returns true if this Type indicates that it accepts a file permission.
func (instance Type) IsTakingFilePermission() bool {
	return instance == GenerateToFile
}

// IsTakingFileUser returns true if this Type indicates that it accepts a file user.
func (instance Type) IsTakingFileUser() bool {
	return instance == GenerateToFile
}

// IsTakingFileGroup returns true if this Type indicates that it accepts a file group.
func (instance Type) IsTakingFileGroup() bool {
	return instance == GenerateToFile
}

// IsGenerating returns true if this Type indicates that it will create a key.
func (instance Type) IsGenerating() bool {
	return instance == GenerateToFile || instance == GenerateToEnvironment
}

// Validate validates an action on this object and returns an error object if there are any.
func (instance Type) Validate() error {
	_, err := instance.CheckedString()
	return err
}
