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
	// @default false
	//
	// If this is set to ``true`` it is possible to control caretakerd remotely.
	// This includes the [``caretakerctl``](#commands.caretakerctl) command and also
	// by the services itself.
	//
	// > **Hint:** This does **NOT** give automatically access each of it access rights to caretakerd.
	// > This is separately handled by access properties:
	// >
	// > * {@ref github.com/echocat/caretakerd/control.Config#Access Control.access} for caretakerctl
	// > * {@ref github.com/echocat/caretakerd/service.Config#Access Services.access} for services
	Enabled Boolean `json:"enabled" yaml:"enabled"`

	// @default "tcp://localhost:57955"
	//
	// Address where caretakerd RPC interface is listen to.
	//
	// For details of possible values see {@ref github.com/echocat/caretakerd/values.SocketAddress}.
	Listen SocketAddress `json:"listen" yaml:"listen"`
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
