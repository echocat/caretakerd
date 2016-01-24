package control

import (
    "github.com/echocat/caretakerd/defaults"
    "github.com/echocat/caretakerd/access"
)

type Config struct {
    Access  access.Config `json:"access" yaml:"access,omitempty"`
}

func NewConfig() Config {
    result := Config{}
    result.init()
    return result
}

func (this *Config) init() {
    (*this).Access = access.NewGenerateToFileConfig(access.ReadWrite, defaults.AuthFileKeyFilename())
}

func (this *Config) BeforeUnmarshalYAML() error {
    this.init()
    return nil
}

func (this Config) Validate() error {
    return this.Access.Validate()
}
