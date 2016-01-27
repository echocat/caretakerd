package rpc

import (
	. "github.com/echocat/caretakerd/defaults"
	. "github.com/echocat/caretakerd/values"
)

var defaults = map[string]interface{}{
	"Enabled": Boolean(false),
	"Listen":  ListenAddress(),
}

type Config struct {
	Enabled Boolean       `json:"enabled" yaml:"enabled"`
	Listen  SocketAddress `json:"listen" yaml:"listen"`
}

func NewConfig() Config {
	result := Config{}
	result.init()
	return result
}

func (instance *Config) init() {
	SetDefaultsTo(defaults, instance)
}

func (instance *Config) BeforeUnmarshalYAML() error {
	instance.init()
	return nil
}

func (instance Config) Validate() error {
	err := instance.Enabled.Validate()
	if err == nil {
		err = instance.Listen.Validate()
	}
	return err
}
