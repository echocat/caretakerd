package access

import (
	"strconv"
	"strings"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/panics"
	"encoding/json"
)

type Permission int

const (
	Forbidden Permission = 0
	ReadOnly Permission = 1
	ReadWrite Permission = 2
)

var AllPermissions = []Permission{
	Forbidden,
	ReadOnly,
	ReadWrite,
}

func (instance Permission) String() string {
	switch instance {
	case Forbidden:
		return "forbidden"
	case ReadOnly:
		return "readOnly"
	case ReadWrite:
		return "readWrite"
	}
	panic(panics.New("Illegal permission: %d", instance))
}

func (instance Permission) CheckedString() (string, error) {
	switch instance {
	case Forbidden:
		return "forbidden", nil
	case ReadOnly:
		return "readOnly", nil
	case ReadWrite:
		return "readWrite", nil
	}
	return "", errors.New("Illegal permission: %d", instance)
}

func (instance *Permission) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllPermissions {
			if int(candidate) == valueAsInt {
				(*instance) = candidate
				return nil
			}
		}
		return errors.New("Illegal permission: " + value)
	} else {
		lowerValue := strings.ToLower(value)
		for _, candidate := range AllPermissions {
			if strings.ToLower(candidate.String()) == lowerValue {
				(*instance) = candidate
				return nil
			}
		}
		return errors.New("Illegal permission: " + value)
	}
}

func (instance Permission) MarshalYAML() (interface{}, error) {
	return instance.CheckedString()
}

func (instance *Permission) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

func (instance Permission) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

func (instance *Permission) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

func (instance Permission) Validate() error {
	_, err := instance.CheckedString()
	return err
}

