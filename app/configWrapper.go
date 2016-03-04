package app

import (
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/defaults"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/values"
)

// ConfigWrapper wraps the config of caretakerd and trigger on
// call of Set(string) the loading of this config file.
type ConfigWrapper struct {
	instance    *caretakerd.Config
	filename    values.String
	explicitSet bool
}

// NewConfigWrapper creates a new instance of ConfigWrapper
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

// Set the given string to current object from a string.
// Return an error object if there are some problems while transforming the string.
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

// IsExplicitSet returns true if the Set(string) method was called.
// This means: Someone sets the config file explicit.
func (instance ConfigWrapper) IsExplicitSet() bool {
	return instance.explicitSet
}

// ConfigureAndValidate will configure itself with the given parameters and will
// also call every validation method. If an error occurs this will be returned.
// If there is no problem nil is retured.
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
