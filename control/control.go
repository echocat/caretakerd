package control

import (
	"github.com/echocat/caretakerd/access"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/keyStore"
)

// Control represents how a remote caretakerctl/control process is able to control
// the current caretakerd instance.
type Control struct {
	config Config
	access *access.Access
}

// NewControl creates a new instance of Control with the given config.
func NewControl(conf Config, ks *keyStore.KeyStore) (*Control, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}
	a, err := access.NewAccess(conf.Access, "caretakerctl", ks)
	if err != nil {
		return nil, errors.New("Could not create caretakerctl interface.").CausedBy(err)
	}
	return &Control{
		config: conf,
		access: a,
	}, nil
}

// Access returns the enclosed Access instance.
func (instance *Control) Access() *access.Access {
	return instance.access
}

// ConfigObject returns the Config object this control was created with.
func (instance *Control) ConfigObject() interface{} {
	return instance.config
}
