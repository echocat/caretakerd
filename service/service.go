package service

import (
	"github.com/echocat/caretakerd/access"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/keyStore"
	"github.com/echocat/caretakerd/logger"
	usync "github.com/echocat/caretakerd/sync"
	"runtime"
)

// Service represents a service instance in caretakerd that was created from a Config object.
type Service struct {
	config    Config
	logger    *logger.Logger
	name      string
	syncGroup *usync.Group
	access    *access.Access
}

func finalize(what *Service) {
	what.Close()
}

// NewService creates a new service instance from given Config.
func NewService(conf Config, name string, syncGroup *usync.Group, sec *keyStore.KeyStore) (*Service, error) {
	err := conf.Validate()
	if err != nil {
		return nil, errors.New("Config of service '%v' is not valid.", name).CausedBy(err)
	}
	acc, err := access.NewAccess(conf.Access, name, sec)
	if err != nil {
		return nil, errors.New("Could not create access for service '%v'.", name).CausedBy(err)
	}
	log, err := logger.NewLogger(conf.Logger, name, syncGroup)
	if err != nil {
		return nil, errors.New("Could not create logger for service '%v'.", name).CausedBy(err)
	}
	result := &Service{
		config:    conf,
		logger:    log,
		name:      name,
		syncGroup: syncGroup,
		access:    acc,
	}
	runtime.SetFinalizer(result, finalize)
	return result, nil
}

// Close closes all resources of this service.
func (instance *Service) Close() {
	instance.logger.Close()
}

func (instance Service) String() string {
	return instance.Name()
}

// Name returns the name of this service.
func (instance Service) Name() string {
	return instance.name
}

// Config returns the config this services was created from.
func (instance Service) Config() Config {
	return instance.config
}

// Logger returns the logger this service uses.
func (instance Service) Logger() *logger.Logger {
	return instance.logger
}

// Access returns the access instance this service uses.
func (instance *Service) Access() *access.Access {
	return instance.access
}
