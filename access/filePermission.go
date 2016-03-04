package access

import (
	"encoding/json"
	"fmt"
	"github.com/echocat/caretakerd/errors"
	"os"
	"regexp"
)

var filePermissionOctPattern = regexp.MustCompile("^(\\d?)(\\d)(\\d)(\\d)$")

// @inline
type FilePermission os.FileMode

func (instance FilePermission) String() string {
	return fmt.Sprintf("%04o", instance)
}

// Set the given string to current object from a string.
// Return an error object if there are some problems while transforming the string.
func (instance *FilePermission) Set(value string) error {
	if filePermissionOctPattern.MatchString(value) {
		_, err := fmt.Sscanf(value, "%o", instance)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Illegal file permission format: %v", value)
	}
	return nil
}

// MarshalYAML is used until yaml marshalling. Do not call directly.
func (instance FilePermission) MarshalYAML() (interface{}, error) {
	return instance.String(), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call directly.
func (instance *FilePermission) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call directly.
func (instance FilePermission) MarshalJSON() ([]byte, error) {
	return json.Marshal(instance.String())
}

// UnmarshalJSON is used until json unmarshalling. Do not call directly.
func (instance *FilePermission) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate do validate action on this object and return an error object if any.
func (instance FilePermission) Validate() {
	instance.String()
}

func (instance FilePermission) ThisOrDefault() FilePermission {
	if uint32(instance) == 0 {
		return DefaultFilePermission()
	} else {
		return instance
	}
}

func (instance FilePermission) AsFileMode() os.FileMode {
	return os.FileMode(instance)
}

func DefaultFilePermission() FilePermission {
	return defaultFilePermission
}
