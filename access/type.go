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
	None                  Type = 0
	// @id trusted
	Trusted               Type = 1
	// @id generateToEnvironment
	GenerateToEnvironment Type = 2
	// @id generateToFile
	GenerateToFile        Type = 3
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

func (instance Type) MarshalYAML() (interface{}, error) {
	return instance.CheckedString()
}

func (instance *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

func (instance Type) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

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

func (instance Type) Validate() error {
	_, err := instance.CheckedString()
	return err
}
