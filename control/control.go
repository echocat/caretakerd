package control

import (
    "github.com/echocat/caretakerd/access"
    "github.com/echocat/caretakerd/rpc/security"
    "github.com/echocat/caretakerd/errors"
)

type Control struct {
    config Config
    access *access.Access
}

func NewControl(conf Config, sec *security.Security) (*Control, error) {
    err := conf.Validate()
    if err != nil {
        return nil, err
    }
    a, err := access.NewAccess(conf.Access, "caretakerctl", sec)
    if err != nil {
        return nil, errors.New("Could not create caretakerctl interface.").CausedBy(err)
    }
    return &Control{
        config: conf,
        access: a,
    }, nil
}

func (this *Control) Access() *access.Access {
    return this.access
}

func (this *Control) ConfigObject() interface{} {
    return this.config
}
