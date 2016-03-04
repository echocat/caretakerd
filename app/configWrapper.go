package app

import (
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/defaults"
	"github.com/echocat/caretakerd/errors"
	. "github.com/echocat/caretakerd/values"
)

type ConfigWrapper struct {
	instance    *caretakerd.Config
	filename    String
	explicitSet bool
}

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
	filename := String(value)
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

func (instance ConfigWrapper) IsExplicitSet() bool {
	return instance.explicitSet
}

func (instance *ConfigWrapper) ConfigureAndValidate(listenAddress *FlagWrapper, pemFile *FlagWrapper, validateAlsoMaster bool) error {
	listenAddress.AssignIfExplicitSet(&instance.instance.Rpc.Listen)
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

func (instance ConfigWrapper) Instance() *caretakerd.Config {
	return instance.instance
}
