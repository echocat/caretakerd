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

// NewConfig create a new config instance.
func NewConfig() Config {
	result := Config{}
	result.init()
	return result
}

// Validate return an error instance on every validation problem of the
// config instance.
func (instance Config) Validate() error {
	err := instance.KeyStore.Validate()
	if err == nil {
		err = instance.Rpc.Validate()
	}
	if err == nil {
		err = instance.Control.Validate()
	}
	if err == nil {
		err = instance.Logger.Validate()
	}
	if err == nil {
		err = instance.Services.Validate()
	}
	return err
}

// ValidateMaster return an error instance on every validation problem of the
// master service config instance.
func (instance Config) ValidateMaster() error {
	return instance.Services.ValidateMaster()
}

func (instance *Config) init() {
	(*instance).KeyStore = keyStore.NewConfig()
	(*instance).Rpc = rpc.NewConfig()
	(*instance).Control = control.NewConfig()
	(*instance).Logger = logger.NewConfig()
	(*instance).Services = service.NewConfigs()
}

// BeforeUnmarshalYAML is used by yaml unmarcshalling. Do not call direct.
func (instance *Config) BeforeUnmarshalYAML() error {
	instance.init()
	return nil
}
