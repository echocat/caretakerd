package control

import (
	"github.com/echocat/caretakerd/access"
	"github.com/echocat/caretakerd/defaults"
	"runtime"
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

// NewConfigFor creates a new instance of Config.
func NewConfigFor(platform string) Config {
	result := Config{}
	result.init(platform)
	return result
}

func (instance *Config) init(platform string) {
	(*instance).Access = access.NewGenerateToFileConfig(access.ReadWrite, defaults.AuthFileKeyFilenameFor(platform))
}

// BeforeUnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *Config) BeforeUnmarshalYAML() error {
	instance.init(runtime.GOOS)
	return nil
}

// Validate validates the action on this object and returns an error object if there are any.
func (instance Config) Validate() error {
	return instance.Access.Validate()
}
