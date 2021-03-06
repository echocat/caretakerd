package access

import (
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/values"
)

// Config to access caretakerd.
type Config struct {
	// @default "generateToFile" (for control/caretakerctl) "none" (for services)
	//
	// Defines how this access will be ensured.
	//
	// For details see possible values {@ref github.com/echocat/caretakerd/access.Type}.
	Type Type `json:"type" yaml:"type"`

	// @default "readWrite" (for control/caretakerctl) "forbidden" (for services)
	//
	// Defines what the control/service can do with caretakerd.
	//
	// For details see possible values {@ref github.com/echocat/caretakerd/access.Permission}.
	Permission Permission `json:"permission" yaml:"permission"`

	// @default ""
	//
	// If the property {@ref #Type type} = {@ref github.com/echocat/caretakerd/access.Type#Trusted trusted},
	// the certificates specified in this file are used to trust remote connections. Not matching remote connections will be
	// rejected.
	//
	// If the property {@ref #Type type} = {@ref github.com/echocat/caretakerd/access.Type#GenerateToFile generateToFile},
	// caretakerd generates this file that must be used by remote connections.
	//
	// > **Important:** If the property {@ref #Type type} = {@ref github.com/echocat/caretakerd/access.Type#GenerateToFile generateToFile},
	// > this property is required.
	PemFile values.String `json:"pemFile,omitempty" yaml:"pemFile"`

	// @default 0600
	//
	// Permission in filesystem of the generated {@ref #PemFile pem file}.
	PemFilePermission FilePermission `json:"pemFilePermission,omitempty" yaml:"pemFilePermission"`

	// @default ""
	//
	// If set, this user owns the generated {@ref #PemFile pem file}.
	// Otherwise it is owned by the user caretakerd is running with.
	PemFileUser values.String `json:"pemFileUser,omitempty" yaml:"pemFileUser"`
}

// NewNoneConfig creates a new Config that denies access to anything.
func NewNoneConfig() Config {
	return Config{
		Type:       None,
		Permission: Forbidden,
	}
}

// NewTrustedConfig creates a new Config with the given Permission based on Trusted rules.
func NewTrustedConfig(permission Permission) Config {
	return Config{
		Type:       Trusted,
		Permission: permission,
	}
}

// NewGenerateToEnvironmentConfig creates a new Config with the given permission
// and will force a creation of certificates to environment variables.
func NewGenerateToEnvironmentConfig(permission Permission) Config {
	return Config{
		Type:       GenerateToEnvironment,
		Permission: permission,
	}
}

// NewGenerateToFileConfig creates a new Config with the given permission
// and will force a creation of certificates to the given pemFile.
func NewGenerateToFileConfig(permission Permission, pemFile values.String) Config {
	return Config{
		Type:              GenerateToFile,
		Permission:        permission,
		PemFile:           values.String(pemFile),
		PemFileUser:       values.String(""),
		PemFilePermission: DefaultFilePermission(),
	}
}

// Validate validates an action on this object and returns an error object if there is any.
func (instance Config) Validate() error {
	err := instance.Type.Validate()
	if err == nil {
		err = instance.Permission.Validate()
	}
	if err == nil {
		err = instance.validateRequireStringValue(instance.PemFile, "pemFile", instance.Type.IsTakingFilename)
	}
	if err == nil {
		err = instance.validateStringOnlyAllowedValue(instance.PemFileUser, "pemFileUser", instance.Type.IsTakingFileUser, values.String(""))
	}
	if err == nil {
		err = instance.validateUint32OnlyAllowedValue(uint32(instance.PemFilePermission), "pemFilePermission", instance.Type.IsTakingFilePermission, uint32(DefaultFilePermission()))
	}
	return err
}

func (instance Config) validateRequireStringValue(value values.String, fieldName string, isAllowedMethod func() bool) error {
	if isAllowedMethod() {
		if value.IsEmpty() {
			return errors.New("There is no %s set for type %v.", fieldName, instance.Type)
		}
	}
	return nil
}

func (instance Config) validateStringOnlyAllowedValue(value values.String, fieldName string, isAllowedMethod func() bool, defaultValue values.String) error {
	if !isAllowedMethod() && value != defaultValue && !value.IsEmpty() {
		return errors.New("There is no %s allowed for type %v.", fieldName, instance.Type)
	}
	return nil
}

func (instance Config) validateUint32OnlyAllowedValue(value uint32, fieldName string, isAllowedMethod func() bool, defaultValue uint32) error {
	if !isAllowedMethod() && value != defaultValue && value != 0 {
		return errors.New("There is no %s allowed for type %v.", fieldName, instance.Type)
	}
	return nil
}
