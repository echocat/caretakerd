package config

import (
    loggerConfig "github.com/echocat/caretakerd/logger/config"
    "github.com/echocat/caretakerd/service"
    "github.com/echocat/caretakerd/rpc"
    "github.com/echocat/caretakerd/control"
)

type Config struct {
    Rpc      rpc.Config `json:"rpc" yaml:"rpc,omitempty"`
    Control  control.Config `json:"control" yaml:"control,omitempty"`
    Logger   loggerConfig.Config `json:"logger" yaml:"logger,omitempty"`
    Services service.Configs `json:"services" yaml:"services,omitempty"`
}

func NewConfig() Config {
    result := Config{}
    result.init()
    return result
}

func (i Config) Validate() error {
    err := i.Rpc.Validate()
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
    (*i).Rpc = rpc.NewConfig()
    (*i).Control = control.NewConfig()
    (*i).Logger = loggerConfig.NewLoggerConfig()
    (*i).Services = service.NewServiceConfigs()
}

func (i *Config) BeforeUnmarshalYAML() error {
    i.init()
    return nil
}

func (s *Config) Configure(value string, with func(conf *Config, value string) error) error {
    return with(s, value)
}
