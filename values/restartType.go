package values

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"strings"
)

// # Description
//
// RestartType tells caretakerd what to do if a process ends.
type RestartType int

const (
	// @id never
	// Never restart the process.
	Never RestartType = 0
	// @id onFailures
	// Only restart the process on failures.
	OnFailures RestartType = 1
	// @id always
	// Always restart the process. This means on success and on failures.
	Always RestartType = 2
)

// AllRestartTypes contains all possible variants of RestartType.
var AllRestartTypes = []RestartType{
	Never,
	OnFailures,
	Always,
}

func (instance RestartType) String() string {
	result, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return result
}

// CheckedString is like String but also returns an optional error if there are any
// validation errors.
func (instance RestartType) CheckedString() (string, error) {
	switch instance {
	case Never:
		return "never", nil
	case OnFailures:
		return "onFailures", nil
	case Always:
		return "always", nil
	}
	return "", errors.New("Illegal restart type: %v", instance)
}

// Sets the given string to current object from a string.
// Returns an error object if there are any problems while transforming the string.
func (instance *RestartType) Set(value string) error {
	lowerValue := strings.ToLower(value)
	for _, candidate := range AllRestartTypes {
		if strings.ToLower(candidate.String()) == lowerValue {
			(*instance) = candidate
			return nil
		}
	}
	return errors.New("Illegal restart type: %s", value)
}

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (instance RestartType) MarshalYAML() (interface{}, error) {
	return instance.CheckedString()
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *RestartType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call this method directly.
func (instance RestartType) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

// UnmarshalJSON is used until json unmarshalling. Do not call this method directly.
func (instance *RestartType) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// OnSuccess returns "true" of the current RestartType tells to restart on success.
func (instance RestartType) OnSuccess() bool {
	return instance == Always
}

// OnFailures returns "true" of the current RestartType tells to restart on failures.
func (instance RestartType) OnFailures() bool {
	return instance == OnFailures || instance == Always
}

// Validate validates actions on this object and returns an error object if there are any.
func (instance RestartType) Validate() error {
	_, err := instance.CheckedString()
	return err
}
