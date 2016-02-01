package values

import (
	"strings"
)

// @inline
type Boolean bool

func (instance Boolean) String() string {
	s, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return s
}

func (instance Boolean) CheckedString() (string, error) {
	if instance == Boolean(true) {
		return "true", nil
	} else {
		return "false", nil
	}
}

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

func (instance *Boolean) SetFromBool(value bool) error {
	(*instance) = Boolean(value)
	return nil
}

func (instance Boolean) MarshalYAML() (interface{}, error) {
	return bool(instance), nil
}

func (instance *Boolean) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value bool
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.SetFromBool(value)
}

func (instance Boolean) Validate() error {
	_, err := instance.CheckedString()
	return err
}
