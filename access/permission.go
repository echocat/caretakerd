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

var AllPermissions []Permission = []Permission{
    Forbidden,
    ReadOnly,
    ReadWrite,
}

func (this Permission) String() string {
    switch this {
    case Forbidden:
        return "forbidden"
    case ReadOnly:
        return "readOnly"
    case ReadWrite:
        return "readWrite"
    }
    panic(panics.New("Illegal permission: %d", this))
}

func (this Permission) CheckedString() (string, error) {
    switch this {
    case Forbidden:
        return "forbidden", nil
    case ReadOnly:
        return "readOnly", nil
    case ReadWrite:
        return "readWrite", nil
    }
    return "", errors.New("Illegal permission: %d", this)
}

func (this *Permission) Set(value string) error {
    if valueAsInt, err := strconv.Atoi(value); err == nil {
        for _, candidate := range AllPermissions {
            if int(candidate) == valueAsInt {
                (*this) = candidate
                return nil
            }
        }
        return errors.New("Illegal permission: " + value)
    } else {
        lowerValue := strings.ToLower(value)
        for _, candidate := range AllPermissions {
            if strings.ToLower(candidate.String()) == lowerValue {
                (*this) = candidate
                return nil
            }
        }
        return errors.New("Illegal permission: " + value)
    }
}

func (this Permission) MarshalYAML() (interface{}, error) {
    return this.CheckedString()
}

func (this *Permission) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var value string
    if err := unmarshal(&value); err != nil {
        return err
    }
    return this.Set(value)
}

func (this Permission) MarshalJSON() ([]byte, error) {
    s, err := this.CheckedString()
    if err != nil {
        return []byte{}, err
    }
    return json.Marshal(s)
}

func (this *Permission) UnmarshalJSON(b []byte) error {
    var value string
    if err := json.Unmarshal(b, &value); err != nil {
        return err
    }
    return this.Set(value)
}

func (this Permission) Validate() error {
    _, err := this.CheckedString()
    return err
}

