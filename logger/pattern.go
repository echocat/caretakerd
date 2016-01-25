package logger

import "strings"

type Pattern string

func (this Pattern) String() string {
    s, err := this.CheckedString()
    if err != nil {
        panic(err)
    }
    return s
}

func (this Pattern) CheckedString() (string, error) {
    return string(this), nil
}

func (this *Pattern) Set(value string) error {
    (*this) = Pattern(value)
    return nil
}

func (this Pattern) MarshalYAML() (interface{}, error) {
    return string(this), nil
}

func (this *Pattern) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var value string
    if err := unmarshal(&value); err != nil {
        return err
    }
    return this.Set(value)
}

func (this Pattern) Validate() error {
    _, err := this.CheckedString()
    return err
}

func (this Pattern) IsEmpty() bool {
    return len(this) <= 0
}

func (this Pattern) IsTrimmedEmpty() bool {
    return len(strings.TrimSpace(this.String())) <= 0
}
