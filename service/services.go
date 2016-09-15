package service

import (
	"github.com/echocat/caretakerd/keyStore"
	"github.com/echocat/caretakerd/panics"
	usync "github.com/echocat/caretakerd/sync"
)

// Services is a couple of services with their names.
type Services map[string]*Service

// NewServices creates a new instance of Services from the given Configs.
func NewServices(configs Configs, syncGroup *usync.Group, sec *keyStore.KeyStore) (*Services, error) {
	err := configs.Validate()
	if err != nil {
		return nil, err
	}
	result := Services{}
	for name, conf := range configs {
		newService, err := NewService(conf, name, syncGroup.NewGroup(), sec)
		if err != nil {
			return nil, err
		}
		result[name] = newService
	}
	return &result, nil
}

// Get returns a service for the given name if a service exists.
// If no service for the given name could be found, "nil" is returned.
func (instance Services) Get(name string) *Service {
	return instance[name]
}

// GetMaster returns the master of this instance.
// If there is no master, "nil" is returned.
func (instance Services) GetMaster() *Service {
	for _, service := range instance {
		if service.config.Type != Master {
			return service
		}
	}
	return nil
}

// GetMasterOrFail returns the master, if a master exists.
// Otherwise it will fail with a panic.
func (instance Services) GetMasterOrFail() *Service {
	master := instance.GetMaster()
	if master == nil {
		panics.New("There is no master service defined.").Throw()
	}
	return master
}

// GetAllButMaster returns every service but not the master.
func (instance Services) GetAllButMaster() Services {
	result := Services{}
	for name, service := range instance {
		if service.config.Type != Master {
			result[name] = service
		}
	}
	return result
}

// GetAllAutoStartable returns every service that is autostartable.
func (instance Services) GetAllAutoStartable() Services {
	result := Services{}
	for name, service := range instance {
		if service.config.Type.IsAutoStartable() {
			result[name] = service
		}
	}
	return result
}

// Close closes this instance and all of the services.
func (instance Services) Close() {
	for _, service := range instance {
		service.Close()
	}
}
