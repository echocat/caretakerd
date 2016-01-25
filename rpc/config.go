package rpc

import (
    . "github.com/echocat/caretakerd/values"
    . "github.com/echocat/caretakerd/defaults"
)

var defaults = map[string]interface{}{
    "Enabled": Boolean(false),
    "Listen": ListenAddress(),
}

type Config struct {
    Enabled       Boolean `json:"enabled" yaml:"enabled"`
    Listen        SocketAddress `json:"listen" yaml:"listen"`
}

func NewConfig() Config {
    result := Config{}
    result.init()
    return result
}

func (this *Config) init() {
    SetDefaultsTo(defaults, this)
}

func (this *Config) BeforeUnmarshalYAML() error {
    this.init()
    return nil
}

func (this Config) Validate() error {
    err := this.Enabled.Validate()
    if err == nil {
        err = this.Listen.Validate()
    }
    return err
}
