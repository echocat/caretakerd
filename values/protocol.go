package values

import (
    "strconv"
    "strings"
    "github.com/echocat/caretakerd/errors"
    "encoding/json"
)

type Protocol int
const (
    Tcp Protocol = 0
    Unix Protocol = 1
)

var AllProtocols []Protocol = []Protocol{
    Tcp,
    Unix,
}

func (i Protocol) String() string {
    result, err := i.CheckedString()
    if err != nil {
        panic(err)
    }
    return result
}

func (i Protocol) CheckedString() (string, error) {
    switch i {
    case Tcp:
        return "tcp", nil
    case Unix:
        return "unix", nil
    }
    return "", errors.New("Illegal protocol: %d", i)
}

func (i *Protocol) Set(value string) error {
    if valueAsInt, err := strconv.Atoi(value); err == nil {
        for _, candidate := range AllProtocols {
            if int(candidate) == valueAsInt {
                (*i) = candidate
                return nil
            }
        }
        return errors.New("Illegal protocol: " + value)
    } else {
        lowerValue := strings.ToLower(value)
        for _, candidate := range AllProtocols {
            if strings.ToLower(candidate.String()) == lowerValue {
                (*i) = candidate
                return nil
            }
        }
        return errors.New("Illegal protocol: " + value)
    }
}

func (i Protocol) MarshalYAML() (interface{}, error) {
    return i.String(), nil
}

func (i *Protocol) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var value string
    if err := unmarshal(&value); err != nil {
        return err
    }
    return i.Set(value)
}

func (i Protocol) MarshalJSON() ([]byte, error) {
    s, err := i.CheckedString()
    if err != nil {
        return []byte{}, err
    }
    return json.Marshal(s)
}

func (i *Protocol) UnmarshalJSON(b []byte) error {
    var value string
    if err := json.Unmarshal(b, &value); err != nil {
        return err
    }
    return i.Set(value)
}

func (i Protocol) Validate() error {
    _, err := i.CheckedString()
    return err
}

