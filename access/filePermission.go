package access

import (
	"os"
	"encoding/json"
	"regexp"
	"fmt"
	"github.com/echocat/caretakerd/errors"
)

var filePermissionOctPattern = regexp.MustCompile("^(\\d?)(\\d)(\\d)(\\d)$")

type FilePermission os.FileMode

func (instance FilePermission) String() string {
	return fmt.Sprintf("%04o", instance)
}

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

func (instance FilePermission) MarshalYAML() (interface{}, error) {
	return instance.String(), nil
}

func (instance *FilePermission) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

func (instance FilePermission) MarshalJSON() ([]byte, error) {
	return json.Marshal(instance.String())
}

func (instance *FilePermission) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

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