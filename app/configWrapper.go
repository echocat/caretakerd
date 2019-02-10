package app

import (
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/defaults"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/values"
	"runtime"
	"sync"
)

// ConfigWrapper wraps the config of caretakerd and triggers the loading of this config file when
// calling Set(string).
type ConfigWrapper struct {
	config        *caretakerd.Config
	explicitSet   bool
	listenAddress *FlagWrapper
	pemFile       *FlagWrapper
	platform      string

	loaded bool
	mutex  *sync.Mutex
}

// NewConfigWrapperFor creates a new instance of ConfigWrapper.
func NewConfigWrapperFor(platform string) *ConfigWrapper {
	instance := caretakerd.NewConfigFor(platform)
	defaultListenAddress := defaults.ListenAddressFor(platform)
	defaultPemFile := defaults.AuthFileKeyFilenameFor(platform)
	return &ConfigWrapper{
		config:        &instance,
		explicitSet:   false,
		listenAddress: NewFlagWrapper(&defaultListenAddress),
		pemFile:       NewFlagWrapper(&defaultPemFile),
		platform:      platform,
		mutex:         new(sync.Mutex),
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
		instance.explicitSet = true
		instance.config = config
		return nil
	}
}

func (instance ConfigWrapper) loadConfigFrom(fileName values.String) (*caretakerd.Config, error) {
	if len(fileName) == 0 {
		return nil, errors.New("There is an empty filename for configuration provided.")
	}
	conf, err := caretakerd.LoadFromYamlFile(instance.platform, fileName)
	if err != nil {
		return nil, err
	}
	conf = conf.EnrichFromEnvironment()
	return &conf, nil
}

func (instance *ConfigWrapper) populateAndValidate(forDaemon bool, conf *caretakerd.Config) error {
	instance.listenAddress.AssignIfExplicitSet(&conf.RPC.Listen)
	instance.pemFile.AssignIfExplicitSet(&conf.Control.Access.PemFile)
	err := conf.Validate()
	if err != nil {
		return err
	}
	if forDaemon {
		return conf.ValidateMaster()
	}
	return nil
}

// Instance returns the final caretakerd Config instance.
func (instance ConfigWrapper) Instance() *caretakerd.Config {
	return instance.config
}

// ListenAddress returns the listenAddress FlagWrapper
func (instance ConfigWrapper) ListenAddress() *FlagWrapper {
	return instance.listenAddress
}

// PemFile returns the pemFile FlagWrapper
func (instance ConfigWrapper) PemFile() *FlagWrapper {
	return instance.pemFile
}

// ProvideConfig will either return the already loaded configuration or will load it
func (instance *ConfigWrapper) ProvideConfig(forDaemon bool) (*caretakerd.Config, error) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	if instance.loaded {
		return instance.config, nil
	}

	var config *caretakerd.Config
	if instance.explicitSet {
		config = instance.config
	} else {
		filename := defaults.ConfigFilenameFor(instance.platform)
		if lConfig, err := instance.loadConfigFrom(filename); caretakerd.IsConfigNotExists(err) {
			if forDaemon {
				return nil, errors.New("There is neither the --config flag set nor does a configuration file under default position (%v) exist.", filename)
			}
			if lConfig == nil {
				nConfig := caretakerd.NewConfigFor(runtime.GOOS)
				lConfig = &nConfig
			}
			config = lConfig
		} else if err != nil {
			return nil, err
		} else {
			config = lConfig
		}
	}

	if err := instance.populateAndValidate(forDaemon, config); err != nil {
		return nil, err
	}
	instance.config = config
	instance.loaded = true
	return config, nil
}
