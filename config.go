package caretakerd

import (
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/service"
	"github.com/echocat/caretakerd/rpc"
	"github.com/echocat/caretakerd/control"
	"github.com/echocat/caretakerd/keyStore"
)

type Config struct {
	KeyStore keyStore.Config `json:"keyStore" yaml:"keyStore,omitempty"`
	Rpc      rpc.Config `json:"rpc" yaml:"rpc,omitempty"`
	Control  control.Config `json:"control" yaml:"control,omitempty"`
	Logger   logger.Config `json:"logger" yaml:"logger,omitempty"`
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
