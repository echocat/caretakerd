package values

import (
    "strings"
)

// @id Boolean
// @type simple
//
// This represents a value that could only be ``true`` or ``false``.
type Boolean bool

func (this Boolean) String() string {
    s, err := this.CheckedString()
    if err != nil {
        panic(err)
    }
    return s
}

func (this Boolean) CheckedString() (string, error) {
    if this == Boolean(true) {
        return "true", nil
    } else {
        return "false", nil
    }
}

func (this *Boolean) Set(value string) error {
    switch strings.ToLower(value) {
    case "1": fallthrough
    case "on": fallthrough
    case "yes": fallthrough
    case "y": fallthrough
    case "true":
        return this.SetFromBool(true)
    }
    return this.SetFromBool(false)
}

func (this *Boolean) SetFromBool(value bool) error {
    (*this) = Boolean(value)
    return nil
}

func (this Boolean) MarshalYAML() (interface{}, error) {
    return bool(this), nil
}

func (this *Boolean) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var value bool
    if err := unmarshal(&value); err != nil {
        return err
    }
    return this.SetFromBool(value)
}

func (this Boolean) Validate() error {
    _, err := this.CheckedString()
    return err
}
