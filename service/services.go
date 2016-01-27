package service

import (
	usync "github.com/echocat/caretakerd/sync"
	"github.com/echocat/caretakerd/panics"
	"github.com/echocat/caretakerd/keyStore"
)

type Services map[string]*Service

func NewServices(confs Configs, syncGroup *usync.SyncGroup, sec *keyStore.KeyStore) (*Services, error) {
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
		if service.config.Type != Master {
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
		if service.config.Type != Master {
			result[name] = service
		}
	}
	return result
}

func (s Services) GetAllAutoStartable() Services {
	result := Services{}
	for name, service := range s {
		if service.config.Type.IsAutoStartable() {
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
