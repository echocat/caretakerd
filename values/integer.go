package values

import (
	"strconv"
	"github.com/echocat/caretakerd/errors"
)

// @id Integer
// @type simple
//
// This represents a natural number. It could be negative and positive.
type Integer int

func (this Integer) String() string {
	s, err := this.ChekcedString()
	if err != nil {
		panic(err)
	}
	return s
}

func (this Integer) ChekcedString() (string, error) {
	return strconv.Itoa(int(this)), nil
}

func (this *Integer) Set(value string) error {
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("Illegal integer value: " + value)
	}
	return this.SetFromInt(valueAsInt)
}

func (this *Integer) SetFromInt(value int) error {
	(*this) = Integer(value)
	return nil
}

func (this Integer) MarshalYAML() (interface{}, error) {
	return int(this), nil
}

func (this *Integer) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value int
	if err := unmarshal(&value); err != nil {
		return err
	}
	return this.SetFromInt(value)
}

func (this Integer) Int() int {
	return int(this)
}

func (this Integer) Validate() error {
	_, err := this.ChekcedString()
	return err
}
