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

func (instance Config) Validate() error {
	err := instance.Type.Validate()
	if err == nil {
		err = instance.Permission.Validate()
	}
	if err == nil {
		err = instance.validateRequireStringValue(instance.PemFile, "pemFile", instance.Type.IsTakingFilename)
	}
	if err == nil {
		err = instance.validateStringOnlyAllowedValue(instance.PemFileUser, "pemFileUser", instance.Type.IsTakingFileUser)
	}
	// TODO!    i.validateUint32OnlyAllowedValue(uint32(i.PemFilePermission), "pemFilePermission", i.Auth.IsTakingPermission)
	// TODO!    i.validateStringOnlyAllowedValue(i.KeyArguments, "keyArguments", i.Auth.IsGeneratingCertificate)
	return err
}

func (instance Config) validateRequireStringValue(value String, fieldName string, isAllowedMethod func() bool) error {
	if isAllowedMethod() {
		if value.IsEmpty() {
			return errors.New("There is no %s set for type %v.", fieldName, instance.Type)
		}
	}
	return nil
}

func (instance Config) validateStringOnlyAllowedValue(value String, fieldName string, isAllowedMethod func() bool) error {
	if ! isAllowedMethod() && !value.IsEmpty() {
		return errors.New("There is no %s allowed for type %v.", fieldName, instance.Type)
	}
	return nil
}

func (instance Config) validateUint32OnlyAllowedValue(value uint32, fieldName string, isAllowedMethod func() bool) error {
	if ! isAllowedMethod() && value != 0 {
		return errors.New("There is no %s allowed for type %v.", fieldName, instance.Type)
	}
	return nil
}
