package access

import (
    . "github.com/echocat/caretakerd/values"
    "github.com/echocat/caretakerd/errors"
)

type Config struct {
    Type              Type `json:"type" yaml:"type"`
    Permission        Permission `json:"permission" yaml:"permission"`
    PemFile           String `json:"pemFile,omitempty" yaml:"pemFile"`
    PemFilePermission FilePermission `json:"pemFilePermission,omitempty" yaml:"pemFilePermission"`
    PemFileUser       String `json:"pemFileUser,omitempty" yaml:"pemFileUser"`
}

func NewNoneConfig() Config {
    return Config{
        Type: None,
    }
}

func NewTrustedConfig(permission Permission) Config {
    return Config{
        Type: Trusted,
        Permission: permission,
    }
}

func NewGenerateToEnvironmentConfig(permission Permission) Config {
    return Config{
        Type: GenerateToEnvironment,
        Permission: permission,
    }
}

func NewGenerateToFileConfig(permission Permission, pemFile String) Config {
    return Config{
        Type: GenerateToFile,
        Permission: permission,
        PemFile: String(pemFile),
        PemFilePermission: DefaultFilePermission(),
        PemFileUser: String(""),
    }
}

func (this Config) Validate() error {
    err := this.Type.Validate()
    if err == nil {
        err = this.Permission.Validate()
    }
    if err == nil {
        err = this.validateRequireStringValue(this.PemFile, "pemFile", this.Type.IsTakingFilename)
    }
    if err == nil {
        err = this.validateStringOnlyAllowedValue(this.PemFileUser, "pemFileUser", this.Type.IsTakingFileUser)
    }
    // TODO!    i.validateUint32OnlyAllowedValue(uint32(i.PemFilePermission), "pemFilePermission", i.Auth.IsTakingPermission)
    // TODO!    i.validateStringOnlyAllowedValue(i.KeyArguments, "keyArguments", i.Auth.IsGeneratingCertificate)
    return err
}

func (this Config) validateRequireStringValue(value String, fieldName string, isAllowedMethod func() bool) error {
    if isAllowedMethod() {
        if value.IsEmpty() {
            return errors.New("There is no %s set for type %v.", fieldName, this.Type)
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

func (this Config) validateUint32OnlyAllowedValue(value uint32, fieldName string, isAllowedMethod func() bool) error {
    if ! isAllowedMethod() && value != 0 {
        return errors.New("There is no %s allowed for type %v.", fieldName, this.Type)
    }
    return nil
}
