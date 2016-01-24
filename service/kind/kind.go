package kind

import (
    "strconv"
    "strings"
    "github.com/echocat/caretakerd/errors"
    "encoding/json"
)

type Kind int
const (
    OnDemand Kind = 0
    AutoStart Kind = 1
    Master Kind = 2
)

var All []Kind = []Kind{
    OnDemand,
    AutoStart,
    Master,
}

func (i Kind) String() string {
    result, err := i.CheckedString()
    if err != nil {
        panic(err)
    }
    return result
}

func (i Kind) CheckedString() (string, error) {
    switch i {
    case OnDemand:
        return "onDemand", nil
    case AutoStart:
        return "autoStart", nil
    case Master:
        return "master", nil
    }
    return "", errors.New("Illegal kind: %d", i)
}

func (i *Kind) Set(value string) error {
    if valueAsInt, err := strconv.Atoi(value); err == nil {
        for _, candidate := range All {
            if int(candidate) == valueAsInt {
                (*i) = candidate
                return nil
            }
        }
        return errors.New("Illegal kind: " + value)
    } else {
        lowerValue := strings.ToLower(value)
        for _, candidate := range All {
            if strings.ToLower(candidate.String()) == lowerValue {
                (*i) = candidate
                return nil
            }
        }
        return errors.New("Illegal kind: " + value)
    }
}

func (i Kind) MarshalYAML() (interface{}, error) {
    return i.CheckedString()
}

func (i *Kind) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var value string
    if err := unmarshal(&value); err != nil {
        return err
    }
    return i.Set(value)
}

func (i Kind) MarshalJSON() ([]byte, error) {
    s, err := i.CheckedString()
    if err != nil {
        return []byte{}, err
    }
    return json.Marshal(s)
}

func (i *Kind) UnmarshalJSON(b []byte) error {
    var value string
    if err := json.Unmarshal(b, &value); err != nil {
        return err
    }
    return i.Set(value)
}

func (i Kind) IsAutoStartable() bool {
    switch i {
    case Master:
        fallthrough
    case AutoStart:
        return true
    }
    return false
}

func (i Kind) Validate() error {
    _, err := i.CheckedString()
    return err
}
