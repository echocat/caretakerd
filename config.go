package caretakerd

import (
	"github.com/echocat/caretakerd/control"
	"github.com/echocat/caretakerd/keyStore"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/rpc"
	"github.com/echocat/caretakerd/service"
)

// Root configuration of caretakerd.
type Config struct {
	// Defines how the encryption of caretakerd works.
	// This is especially important if {@ref #Rpc RPC} is used.
	//
	// For details see {@ref github.com/echocat/caretakerd/keyStore.Config}.
	KeyStore keyStore.Config `json:"keyStore" yaml:"keyStore,omitempty"`

	// Defines how caretaker can controlled remotely.
	//
	// For details see {@ref github.com/echocat/caretakerd/rpc.Config}.
	Rpc rpc.Config `json:"rpc" yaml:"rpc,omitempty"`

	// Defines the access rights of caretakerctl to caretakerd.
	// This requires {@ref #Rpc RPC} enabled.
	//
	// For details see {@ref github.com/echocat/caretakerd/control.Config}.
	Control control.Config `json:"control" yaml:"control,omitempty"`

	// Configures the logger for caretakerd itself.
	// This does not include output of services.
	//
	// For details see {@ref github.com/echocat/caretakerd/logger.Config}.
	Logger logger.Config `json:"logger" yaml:"logger,omitempty"`

	// Services configuration to run with caretakerd.
	//
	// > **Important**: This is a map and requires exact one service
	// > configured as {@ref github.com/echocat/caretakerd/service.Config#Type type} = {@ref github.com/echocat/caretakerd/service.Type#Master master}.
	//
	// For details see {@ref github.com/echocat/caretakerd/service.Config}.
	Services service.Configs `json:"services" yaml:"services,omitempty"`
}

func NewConfig() Config {
	result := Config{}
	result.init()
	return result
}

func (i Config) Validate() error {
	err := i.KeyStore.Validate()
	if err == nil {
		err = i.Rpc.Validate()
	}
	if err == nil {
		err = i.Control.Validate()
	}
	if err == nil {
		err = i.Logger.Validate()
	}
	if err == nil {
		err = i.Services.Validate()
	}
	return err
}

func (i Config) ValidateMaster() error {
	return i.Services.ValidateMaster()
}

func (i *Config) init() {
	(*i).KeyStore = keyStore.NewConfig()
	(*i).Rpc = rpc.NewConfig()
	(*i).Control = control.NewConfig()
	(*i).Logger = logger.NewConfig()
	(*i).Services = service.NewConfigs()
}

func (i *Config) BeforeUnmarshalYAML() error {
	i.init()
	return nil
}

func (s *Config) Configure(value string, with func(conf *Config, value string) error) error {
	return with(s, value)
}
