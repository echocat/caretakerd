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

func (this FilePermission) String() string {
    return fmt.Sprintf("%04o", this)
}

func (this *FilePermission) Set(value string) error {
    if filePermissionOctPattern.MatchString(value) {
        _, err := fmt.Sscanf(value, "%o", this)
        if err != nil {
            return err
        }
    } else {
        return errors.New("Illegal file permission format: %v", value)
    }
    return nil
}

func (this FilePermission) MarshalYAML() (interface{}, error) {
    return this.String(), nil
}

func (this *FilePermission) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var value string
    if err := unmarshal(&value); err != nil {
        return err
    }
    return this.Set(value)
}

func (this FilePermission) MarshalJSON() ([]byte, error) {
    return json.Marshal(this.String())
}

func (this *FilePermission) UnmarshalJSON(b []byte) error {
    var value string
    if err := json.Unmarshal(b, &value); err != nil {
        return err
    }
    return this.Set(value)
}

func (this FilePermission) Validate() {
    this.String()
}

func (this FilePermission) ThisOrDefault() FilePermission {
    if uint32(this) == 0 {
        return DefaultFilePermission()
    } else {
        return this
    }
}

func (this FilePermission) AsFileMode() os.FileMode {
    return os.FileMode(this)
}

func DefaultFilePermission() FilePermission {
    return defaultFilePermission
}