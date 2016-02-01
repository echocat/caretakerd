package service

import (
	"github.com/echocat/caretakerd/errors"
)

// @inline
type Configs map[string]Config

func NewConfigs() Configs {
	return Configs{}
}

func (s *Configs) GetMasterName() (string, bool) {
	for name, service := range *s {
		if service.Type == Master {
			return name, true
		}
	}
	return "", false
}

func (s *Configs) Configure(name string, value string, with func(conf *Config, value string) error) error {
	conf, ok := (*s)[name]
	if !ok {
		return errors.New("There does no service with name '%s' exist.", name)
	}
	err := with(&conf, value)
	(*s)[name] = conf
	return err
}

func (s *Configs) ConfigureSub(name string, key string, value string, with func(conf *Config, key string, value string) error) error {
	conf, ok := (*s)[name]
	if !ok {
		return errors.New("There does no service with name '%s' exist.", name)
	}
	err := with(&conf, key, value)
	(*s)[name] = conf
	return err
}
