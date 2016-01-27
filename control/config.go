package control

import (
	"github.com/echocat/caretakerd/defaults"
	"github.com/echocat/caretakerd/access"
)

type Config struct {
	Access access.Config `json:"access" yaml:"access,omitempty"`
}

func NewConfig() Config {
	result := Config{}
	result.init()
	return result
}

func (instance *Config) init() {
	(*instance).Access = access.NewGenerateToFileConfig(access.ReadWrite, defaults.AuthFileKeyFilename())
}

func (instance *Config) BeforeUnmarshalYAML() error {
	instance.init()
	return nil
}

func (instance Config) Validate() error {
	return instance.Access.Validate()
}
