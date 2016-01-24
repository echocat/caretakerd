package security

import (
    "strconv"
    "strings"
    "github.com/echocat/caretakerd/errors"
    "encoding/json"
)

type Type int
const (
    FromFile Type = 0
    FromEnvironment Type = 1
    Generated Type = 2
)

var AllTypes []Type = []Type{
    FromFile,
    FromEnvironment,
    Generated,
}

func (i Type) String() string {
    result, err := i.CheckedString()
    if err != nil {
        panic(err)
    }
    return result
}

func (i Type) CheckedString() (string, error) {
    switch i {
    case FromFile:
        return "fromFile", nil
    case FromEnvironment:
        return "fromEnvironment", nil
    case Generated:
        return "generated", nil
    }
    return "", errors.New("Illegal security type: %d", i)
}

func (i *Type) Set(value string) error {
    if valueAsInt, err := strconv.Atoi(value); err == nil {
        for _, candidate := range AllTypes {
            if int(candidate) == valueAsInt {
                (*i) = candidate
                return nil
            }
        }
        return errors.New("Illegal security type: " + value)
    } else {
        lowerValue := strings.ToLower(value)
        for _, candidate := range AllTypes {
            if strings.ToLower(candidate.String()) == lowerValue {
                (*i) = candidate
                return nil
            }
        }
        return errors.New("Illegal security type: " + value)
    }
}

func (i Type) MarshalYAML() (interface{}, error) {
    return i.String(), nil
}

func (i *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var value string
    if err := unmarshal(&value); err != nil {
        return err
    }
    return i.Set(value)
}

func (i Type) MarshalJSON() ([]byte, error) {
    s, err := i.CheckedString()
    if err != nil {
        return []byte{}, err
    }
    return json.Marshal(s)
}

func (i *Type) UnmarshalJSON(b []byte) error {
    var value string
    if err := json.Unmarshal(b, &value); err != nil {
        return err
    }
    return i.Set(value)
}

func (i Type) IsTakingFilename() bool {
    return i == FromFile
}

func (i Type) IsGenerating() bool {
    return i == Generated
}

func (i Type) IsConsumingCaFile () bool {
    return i == FromFile || i == FromEnvironment
}

func (i Type) Validate() error {
    _, err := i.CheckedString()
    return err
}

