package access

import (
	"strconv"
	"strings"
	"github.com/echocat/caretakerd/errors"
	"encoding/json"
)

type Type int

const (
	None Type = 0
	Trusted Type = 1
	GenerateToEnvironment Type = 2
	GenerateToFile Type = 3
)

var AllTypes []Type = []Type{
	None,
	Trusted,
	GenerateToEnvironment,
	GenerateToFile,
}

func (this Type) String() string {
	s, err := this.CheckedString()
	if err != nil {
		panic(err)
	}
	return s
}

func (this Type) CheckedString() (string, error) {
	switch this {
	case None:
		return "none", nil
	case Trusted:
		return "trusted", nil
	case GenerateToEnvironment:
		return "generateToEnvironment", nil
	case GenerateToFile:
		return "generateToFile", nil
	}
	return "", errors.New("Illegal access type: %d", this)
}

func (this *Type) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllTypes {
			if int(candidate) == valueAsInt {
				(*this) = candidate
				return nil
			}
		}
		return errors.New("Illegal access type: " + value)
	} else {
		lowerValue := strings.ToLower(value)
		for _, candidate := range AllTypes {
			if strings.ToLower(candidate.String()) == lowerValue {
				(*this) = candidate
				return nil
			}
		}
		return errors.New("Illegal access type: " + value)
	}
}

func (this Type) MarshalYAML() (interface{}, error) {
	return this.CheckedString()
}

func (this *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return this.Set(value)
}

func (this Type) MarshalJSON() ([]byte, error) {
	s, err := this.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

func (this *Type) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return this.Set(value)
}

func (this Type) IsTakingFilename() bool {
	return this == GenerateToFile
}

func (this Type) IsTakingFilePermission() bool {
	return this == GenerateToFile
}

func (this Type) IsTakingFileUser() bool {
	return this == GenerateToFile
}

func (this Type) IsTakingGroup() bool {
	return this == GenerateToFile
}

func (this Type) IsGenerating() bool {
	return this == GenerateToFile || this == GenerateToEnvironment
}

func (this Type) Validate() error {
	_, err := this.CheckedString()
	return err
}

