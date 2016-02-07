package control

import (
	"github.com/echocat/caretakerd/access"
	"github.com/echocat/caretakerd/defaults"
)

// # Description
//
// Defines the access rights of caretakerctl to caretakerd.
type Config struct {
	// Configures the permission of caretakerctl to control caretakerd remotely
	// and how to obtain the credentials for it.
	//
	// For details see {@ref github.com/echocat/caretakerd/access.Config}.
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
