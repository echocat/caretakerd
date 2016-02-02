package values

import (
	"github.com/echocat/caretakerd/errors"
	"strconv"
)

// This represents an exit value of a process.
type ExitCode int

func (i ExitCode) String() string {
	return strconv.Itoa(int(i))
}

func (i *ExitCode) Set(value string) error {
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("Illegal exit Code value: " + value)
	}
	if valueAsInt < 0 {
		return errors.New("ExitCode value have to be larger or equal to 0. But got: " + value)
	}
	(*i) = ExitCode(valueAsInt)
	return nil
}

func (i ExitCode) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

func (i *ExitCode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return i.Set(value)
}

func (i ExitCode) Validate() {
	i.String()
}
