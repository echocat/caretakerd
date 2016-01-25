package values

import (
    "strconv"
    "github.com/echocat/caretakerd/errors"
)

// @id NonNegativeInteger
// @type simple
//
// Same like {@ref Integer} but could not be negative.
type NonNegativeInteger int

func (i NonNegativeInteger) String() string {
    result, err := i.CheckedString()
    if err != nil {
        panic(err)
    }
    return result
}

func (i NonNegativeInteger) CheckedString() (string, error) {
    return strconv.Itoa(int(i)), nil
}

func (i *NonNegativeInteger) Set(value string) error {
    valueAsInt, err := strconv.Atoi(value)
    if err != nil {
        return errors.New("Illegal integer value: " + value)
    }
    return i.SetFromInt(valueAsInt)
}

func (i *NonNegativeInteger) SetFromInt(value int) error {
    if value < 0 {
        return errors.New("This intger value should not be negative. But got: %v", value)
    }
    (*i) = NonNegativeInteger(value)
    return nil
}

func (i NonNegativeInteger) MarshalYAML() (interface{}, error) {
    return int(i), nil
}

func (i *NonNegativeInteger) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var value int
    if err := unmarshal(&value); err != nil {
        return err
    }
    return i.SetFromInt(value)
}

func (i NonNegativeInteger) Int() int {
    return int(i)
}

func (i NonNegativeInteger) Validate() error {
    _, err := i.CheckedString()
    return err
}
