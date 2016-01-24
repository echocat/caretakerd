package values

import "strings"

type String string

func (this String) String() string {
    s, err := this.CheckedString()
    if err != nil {
        panic(err)
    }
    return s
}

func (this String) CheckedString() (string, error) {
    return string(this), nil
}

func (this *String) Set(value string) error {
    (*this) = String(value)
    return nil
}

func (this String) MarshalYAML() (interface{}, error) {
    return string(this), nil
}

func (this *String) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var value string
    if err := unmarshal(&value); err != nil {
        return err
    }
    return this.Set(value)
}

func (this String) Validate() error {
    _, err := this.CheckedString()
    return err
}

func (this String) IsEmpty() bool {
    return len(this) <= 0
}

func (this String) IsTrimmedEmpty() bool {
    return len(strings.TrimSpace(this.String())) <= 0
}
