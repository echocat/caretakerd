package access

import (
	"encoding/json"
	"fmt"
	"github.com/echocat/caretakerd/errors"
	"os"
	"regexp"
)

var filePermissionOctPattern = regexp.MustCompile("^(\\d?)(\\d)(\\d)(\\d)$")

// FilePermission represents a operating system file permission.
// @inline
type FilePermission os.FileMode

func (instance FilePermission) String() string {
	return fmt.Sprintf("%04o", instance)
}

// Set sets the given string to current object from a string.
// Returns an error object if there are problems while transforming the string.
func (instance *FilePermission) Set(value string) error {
	if filePermissionOctPattern.MatchString(value) {
		_, err := fmt.Sscanf(value, "%o", instance)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Illegal file permission format: %v", value)
}

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (instance FilePermission) MarshalYAML() (interface{}, error) {
	return instance.String(), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *FilePermission) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call this method directly.
func (instance FilePermission) MarshalJSON() ([]byte, error) {
	return json.Marshal(instance.String())
}

// UnmarshalJSON is used until json unmarshalling. Do not call this method directly.
func (instance *FilePermission) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate validates action on this object and returns an error object if errors occur.
func (instance FilePermission) Validate() error { return nil }

// ThisOrDefault returns this instance if not empty. Otherwise the default FilePermission will be returned.
func (instance FilePermission) ThisOrDefault() FilePermission {
	if uint32(instance) == 0 {
		return DefaultFilePermission()
	}
	return instance
}

// AsFileMode returns this instance as os.FileMode instance.
func (instance FilePermission) AsFileMode() os.FileMode {
	return os.FileMode(instance)
}

// DefaultFilePermission returns the default FilePermission instance.
func DefaultFilePermission() FilePermission {
	return defaultFilePermission
}
