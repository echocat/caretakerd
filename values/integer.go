package values

import (
	"github.com/echocat/caretakerd/errors"
	"strconv"
)

// @inline
type Integer int

func (instance Integer) String() string {
	s, err := instance.ChekcedString()
	if err != nil {
		panic(err)
	}
	return s
}

func (instance Integer) ChekcedString() (string, error) {
	return strconv.Itoa(int(instance)), nil
}

func (instance *Integer) Set(value string) error {
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("Illegal integer value: " + value)
	}
	return instance.SetFromInt(valueAsInt)
}

func (instance *Integer) SetFromInt(value int) error {
	(*instance) = Integer(value)
	return nil
}

func (instance Integer) MarshalYAML() (interface{}, error) {
	return int(instance), nil
}

func (instance *Integer) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value int
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.SetFromInt(value)
}

func (instance Integer) Int() int {
	return int(instance)
}

func (instance Integer) Validate() error {
	_, err := instance.ChekcedString()
	return err
}
