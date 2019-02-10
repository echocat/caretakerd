package values

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"strconv"
	"strings"
)

// Protocol represents a type of an protocol.
type Protocol int

const (
	// TCP represents the TCP protocol type.
	TCP Protocol = 0
	// Unix represents a Unix socket files based protocol type.
	Unix Protocol = 1
)

// AllProtocols contains all possible variants of Protocol.
var AllProtocols = []Protocol{
	TCP,
	Unix,
}

func (instance Protocol) String() string {
	result, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return result
}

// CheckedString is like String but also returns an optional error if there are any
// validation errors.
func (instance Protocol) CheckedString() (string, error) {
	switch instance {
	case TCP:
		return "tcp", nil
	case Unix:
		return "unix", nil
	}
	return "", errors.New("Illegal protocol: %d", instance)
}

// Set sets the given string to current object from a string.
// Returns an error object if there are any problems while transforming the string.
func (instance *Protocol) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllProtocols {
			if int(candidate) == valueAsInt {
				*instance = candidate
				return nil
			}
		}
		return errors.New("Illegal protocol: " + value)
	}
	lowerValue := strings.ToLower(value)
	for _, candidate := range AllProtocols {
		if strings.ToLower(candidate.String()) == lowerValue {
			*instance = candidate
			return nil
		}
	}
	return errors.New("Illegal protocol: " + value)
}

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (instance Protocol) MarshalYAML() (interface{}, error) {
	return instance.String(), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *Protocol) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call this method directly.
func (instance Protocol) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

// UnmarshalJSON is used until json unmarshalling. Do not call this method directly.
func (instance *Protocol) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate validates actions on this object and returns an error object if there are any.
func (instance Protocol) Validate() error {
	_, err := instance.CheckedString()
	return err
}
