package app

import (
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/defaults"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/values"
)

// ConfigWrapper wraps the config of caretakerd and triggers the loading of this config file when
// calling Set(string).
type ConfigWrapper struct {
	instance    *caretakerd.Config
	filename    values.String
	explicitSet bool
}

// NewConfigWrapper creates a new instance of ConfigWrapper.
func NewConfigWrapper() *ConfigWrapper {
	instance := caretakerd.NewConfig()
	return &ConfigWrapper{
		instance:    &instance,
		filename:    defaults.ConfigFilename(),
		explicitSet: false,
	}
}

func (instance ConfigWrapper) String() string {
	return instance.filename.String()
}

// Sets the given string to the current object from a string.
// Returns an error object if there are problems while transforming the string.
func (instance *ConfigWrapper) Set(value string) error {
	if len(value) == 0 {
		return errors.New("There is an empty filename for configuration provided.")
	}
	filename := values.String(value)
	conf, err := caretakerd.LoadFromYamlFile(filename)
	if err != nil {
		return err
	}
	conf = conf.EnrichFromEnvironment()
	instance.filename = filename
	instance.instance = &conf
	instance.explicitSet = true
	return nil
}

// IsExplicitSet returns true if the Set(string) method was called, i.e. someone explicitly set the config file.
func (instance ConfigWrapper) IsExplicitSet() bool {
	return instance.explicitSet
}

// ConfigureAndValidate configures itself with the given parameters and
// also calls every validation method. If an error occurs, the corresponding error will be returned.
// If no errors occur, nil is returned.
func (instance *ConfigWrapper) ConfigureAndValidate(listenAddress *FlagWrapper, pemFile *FlagWrapper, validateAlsoMaster bool) error {
	listenAddress.AssignIfExplicitSet(&instance.instance.RPC.Listen)
	pemFile.AssignIfExplicitSet(&instance.instance.Control.Access.PemFile)
	err := instance.instance.Validate()
	if err != nil {
		return err
	}
	if validateAlsoMaster {
		return instance.instance.ValidateMaster()
	}
	return nil
}

// Instance returns the final caretakerd Config instance.
func (instance ConfigWrapper) Instance() *caretakerd.Config {
	return instance.instance
}
