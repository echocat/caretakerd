package rpc

import (
	"github.com/echocat/caretakerd/defaults"
	"github.com/echocat/caretakerd/values"
	"runtime"
)

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

// NewConfigFor creates a new instance of Config.
func NewConfigFor(platform string) Config {
	result := Config{}
	result.init(platform)
	return result
}

func (instance *Config) init(platform string) {
	values.SetDefaultsTo(map[string]interface{}{
		"Enabled": values.Boolean(false),
		"Listen":  defaults.ListenAddressFor(platform),
	}, instance)
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	instance.init(runtime.GOOS)

	type noMethods Config
	return unmarshal((*noMethods)(instance))
}

// Validate validates actions on this object and returns an error object there are any.
func (instance Config) Validate() error {
	err := instance.Enabled.Validate()
	if err == nil {
		err = instance.Listen.Validate()
	}
	return err
}
