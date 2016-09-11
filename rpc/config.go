package rpc

import (
	"github.com/echocat/caretakerd/defaults"
	"github.com/echocat/caretakerd/values"
)

var defaultValues = map[string]interface{}{
	"Enabled": values.Boolean(false),
	"Listen":  defaults.ListenAddress(),
}

// # Description
//
// Defines the remote access to caretakerd.
type Config struct {
	// @default false
	//
	// If this is set to ``true`` it is possible to control caretakerd remotely.
	// This includes the [``caretakerctl``](#commands.caretakerctl) command and also
	// by the services itself.
	//
	// > **Hint:** This does **NOT** automatically grants each of it caretakerd access rights.
	// > This is separately handled by the following access properties:
	// >
	// > * {@ref github.com/echocat/caretakerd/control.Config#Access Control.access} for caretakerctl
	// > * {@ref github.com/echocat/caretakerd/service.Config#Access Services.access} for services
	Enabled values.Boolean `json:"enabled" yaml:"enabled"`

	// @default "tcp://localhost:57955"
	//
	// Address where caretakerd RPC interface is listened to.
	//
	// For details of possible values see {@ref github.com/echocat/caretakerd/values.SocketAddress}.
	Listen values.SocketAddress `json:"listen" yaml:"listen"`
}

// NewConfig creates a new instance of Config.
func NewConfig() Config {
	result := Config{}
	result.init()
	return result
}

func (instance *Config) init() {
	values.SetDefaultsTo(defaultValues, instance)
}

// BeforeUnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *Config) BeforeUnmarshalYAML() error {
	instance.init()
	return nil
}

// Validate validates actions on this object and returns an error object there are any.
func (instance Config) Validate() error {
	err := instance.Enabled.Validate()
	if err == nil {
		err = instance.Listen.Validate()
	}
	return err
}
