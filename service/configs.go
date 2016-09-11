package service

import (
	"github.com/echocat/caretakerd/errors"
)

// Configs represents a couple of named service configs.
// @inline
type Configs map[string]Config

// NewConfigs creates a new instance of Configs.
func NewConfigs() Configs {
	return Configs{}
}

// GetMasterName returns the name of the configured master.
// If there is no valid master configured false is returned.
func (s *Configs) GetMasterName() (string, bool) {
	for name, service := range *s {
		if service.Type == Master {
			return name, true
		}
	}
	return "", false
}

// Configure executes a configuring action for a service with the given name.
func (s *Configs) Configure(serviceName string, value string, with func(conf *Config, value string) error) error {
	conf, ok := (*s)[serviceName]
	if !ok {
		return errors.New("There does no service with name '%s' exist.", serviceName)
	}
	err := with(&conf, value)
	(*s)[serviceName] = conf
	return err
}

// ConfigureSub executes a configuring action for a service with the given name.
func (s *Configs) ConfigureSub(serviceName string, key string, value string, with func(conf *Config, key string, value string) error) error {
	conf, ok := (*s)[serviceName]
	if !ok {
		return errors.New("There does no service with name '%s' exist.", serviceName)
	}
	err := with(&conf, key, value)
	(*s)[serviceName] = conf
	return err
}
