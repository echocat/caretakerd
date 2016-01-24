package securityStore

import (
    . "github.com/echocat/caretakerd/values"
    "strconv"
    "strings"
    "github.com/echocat/caretakerd/errors"
)

var defaults = map[string]interface{} {
    "Type": Generated,
    "PemFile": String(""),
    "Hints": String("algorithm:`rsa` bits:`1024`"),
    "CaFile": String(""),
}

type Config struct {
    Type    Type `json:"type" yaml:"type"`
    PemFile String `json:"pemFile,omitempty" yaml:"pemFile"`
    Hints   String `json:"hints,omitempty" yaml:"hints"`
    CaFile  String `json:"caFile,omitempty" yaml:"caFile"`
}

func NewConfig() Config {
    result := Config{}
    SetDefaultsTo(defaults, &result)
    return result
}

func (this Config) Validate() error {
    err := this.Type.Validate()
    if err == nil {
        err = this.validateRequireStringOrNotValue(this.PemFile, "pemFile", this.Type.IsTakingFilename)
    }
    if err == nil {
        err = this.validateStringOnlyAllowedValue(this.CaFile, "caFile", this.Type.IsConsumingCaFile)
    }
    if err == nil {
        algorithm := this.GetKeyArgument("algorithm")
        if len(algorithm) > 0 && strings.ToLower(algorithm) != "rsa" {
            err = errors.New("Unsupported algorithm: %s", algorithm)
        }
    }
    return err
}

func (this Config) validateRequireStringOrNotValue(value String, fieldName string, isAllowedMethod func() bool) error {
    if isAllowedMethod() {
        if value.IsEmpty() {
            return errors.New("There is no %s set for type %v.", fieldName, this.Type)
        }
    } else {
        if !value.IsEmpty() {
            return errors.New("There is no %s allowed for type %v.", fieldName, this.Type)
        }
    }
    return nil
}

func (this Config) validateStringOnlyAllowedValue(value String, fieldName string, isAllowedMethod func() bool) error {
    if ! isAllowedMethod() && !value.IsEmpty() {
        return errors.New("There is no %s allowed for type %v.", fieldName, this.Type)
    }
    return nil
}

func (this Config) GetKeyArgument(key string) string {
    arguments := this.Hints
    for arguments != "" {
        i := 0
        for i < len(arguments) && arguments[i] == ' ' {
            i++
        }
        arguments = arguments[i:]
        if arguments == "" {
            break
        }
        i = 0
        for i < len(arguments) && arguments[i] > ' ' && arguments[i] != ':' && arguments[i] != '`' && arguments[i] != 0x7f {
            i++
        }
        if i == 0 || i + 1 >= len(arguments) || arguments[i] != ':' || arguments[i + 1] != '`' {
            break
        }
        name := string(arguments[:i])
        arguments = arguments[i + 1:]

        i = 1
        for i < len(arguments) && arguments[i] != '`' {
            if arguments[i] == '\\' {
                i++
            }
            i++
        }
        if i >= len(arguments) {
            break
        }
        qvalue := string(arguments[:i + 1])
        arguments = arguments[i + 1:]

        if key == name {
            value, err := strconv.Unquote(qvalue)
            if err != nil {
                break
            }
            return value
        }
    }
    return ""
}
