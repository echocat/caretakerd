package service

import (
    "github.com/echocat/caretakerd/service/kind"
    "github.com/echocat/caretakerd/service/config"
    usync "github.com/echocat/caretakerd/sync"
    "github.com/echocat/caretakerd/panics"
    "github.com/echocat/caretakerd/rpc/security"
)

type Services map[string]*Service

func NewServices(confs config.Configs, syncGroup *usync.SyncGroup, sec *security.Security) (*Services, error) {
    err := confs.Validate()
    if err != nil {
        return nil, err
    }
    result := Services{}
    for name, conf := range confs {
        newService, err := NewService(conf, name, syncGroup.NewSyncGroup(), sec)
        if err != nil {
            return nil, err
        }
        result[name] = newService
    }
    return &result, nil
}

func (s Services) Get(name string) (*Service, bool) {
    result, ok := s[name]
    return result, ok
}

func (s Services) GetMaster() *Service {
    for _, service := range s {
        if service.config.Kind != kind.Master {
            return service
        }
    }
    return nil
}

func (s Services) GetMasterOrFail() *Service {
    master := s.GetMaster()
    if master == nil {
        panics.New("There is no master service defined.").Throw()
    }
    return master
}

func (s Services) GetAllButMaster() Services {
    result := Services{}
    for name, service := range s {
        if service.config.Kind != kind.Master {
            result[name] = service
        }
    }
    return result
}

func (s Services) GetAllAutoStartable() Services {
    result := Services{}
    for name, service := range s {
        if service.config.Kind.IsAutoStartable() {
            result[name] = service
        }
    }
    return result
}

func (s Services) Close() {
    for _, service := range s {
        service.Close()
    }
}
