package app

import (
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/defaults"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/values"
	"sync"
)

// ConfigWrapper wraps the config of caretakerd and triggers the loading of this config file when
// calling Set(string).
type ConfigWrapper struct {
	config      *caretakerd.Config
	explicitSet bool
	forDaemon   bool

	loaded bool
	mutex  *sync.Mutex
}

// NewConfigWrapper creates a new instance of ConfigWrapper.
func NewConfigWrapper() *ConfigWrapper {
	instance := caretakerd.NewConfig()
	return &ConfigWrapper{
		config:      &instance,
		explicitSet: false,
	}
}

func (instance ConfigWrapper) String() string {
	return instance.config.Source.String()
}

// Set sets the given string to the current object from a string.
// Returns an error object if there are problems while transforming the string.
func (instance *ConfigWrapper) Set(value string) error {
	if config, err := instance.loadConfigFrom(values.String(value)); err != nil {
		return err
	} else {
		instance.config = config
		return nil
	}
}

func (instance ConfigWrapper) loadConfigFrom(fileName values.String) (*caretakerd.Config, error) {
	if len(fileName) == 0 {
		return nil, errors.New("There is an empty filename for configuration provided.")
	}
	conf, err := caretakerd.LoadFromYamlFile(fileName)
	if err != nil {
		return nil, err
	}
	conf = conf.EnrichFromEnvironment()
	return &conf, nil
}

func (instance *ConfigWrapper) populateAndValidate(conf *caretakerd.Config) error {
	listenAddress.AssignIfExplicitSet(&conf.RPC.Listen)
	pemFile.AssignIfExplicitSet(&conf.Control.Access.PemFile)
	err := conf.Validate()
	if err != nil {
		return err
	}
	if instance.forDaemon {
		return conf.ValidateMaster()
	}
	return nil
}

// Instance returns the final caretakerd Config instance.
func (instance ConfigWrapper) Instance() *caretakerd.Config {
	return instance.config
}

// ProvideConfig will either return the already loaded configuration or will load it
func (instance *ConfigWrapper) ProvideConfig() (*caretakerd.Config, error) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	if instance.loaded {
		return instance.config, nil
	}

	var config *caretakerd.Config
	if instance.explicitSet {
		config = instance.config
	} else {
		filename := defaults.ConfigFilename()
		if lConfig, err := instance.loadConfigFrom(filename); caretakerd.IsConfigNotExists(err) {
			if instance.forDaemon {
				return nil, errors.New("There is neither the --config flag set nor does a configuration file under default position (%v) exist.", filename)
			}
			config = lConfig
		} else if err != nil {
			return nil, err
		} else {
			config = lConfig
		}
	}

	if err := instance.populateAndValidate(config); err != nil {
		return nil, err
	}
	instance.config = config
	instance.loaded = true
	return config, nil
}
