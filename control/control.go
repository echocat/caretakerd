package control

import (
	"github.com/echocat/caretakerd/access"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/keyStore"
)

type Control struct {
	config Config
	access *access.Access
}

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

func (instance *Control) Access() *access.Access {
	return instance.access
}

func (instance *Control) ConfigObject() interface{} {
	return instance.config
}
