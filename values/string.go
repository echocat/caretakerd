package values

import "strings"

// @id String
// @type simple
//
// This represents a slice of characters.
type String string

func (instance String) String() string {
	s, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return s
}

func (instance String) CheckedString() (string, error) {
	return string(instance), nil
}

func (instance *String) Set(value string) error {
	(*instance) = String(value)
	return nil
}

func (instance String) MarshalYAML() (interface{}, error) {
	return string(instance), nil
}

func (instance *String) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

func (instance String) Validate() error {
	_, err := instance.CheckedString()
	return err
}

func (instance String) IsEmpty() bool {
	return len(instance) <= 0
}

func (instance String) IsTrimmedEmpty() bool {
	return len(strings.TrimSpace(instance.String())) <= 0
}
