package service

import (
	"github.com/echocat/caretakerd/logger"
	"runtime"
	usync "github.com/echocat/caretakerd/sync"
	"github.com/echocat/caretakerd/access"
	"github.com/echocat/caretakerd/keyStore"
	"github.com/echocat/caretakerd/errors"
)

type Service struct {
	config    Config
	logger    *logger.Logger
	name      string
	syncGroup *usync.SyncGroup
	access    *access.Access
}

func finalize(what *Service) {
	what.Close()
}

func NewService(conf Config, name string, syncGroup *usync.SyncGroup, sec *keyStore.KeyStore) (*Service, error) {
	err := conf.Validate()
	if err != nil {
		return nil, errors.New("Config of service '%v' is not valid.", name).CausedBy(err)
	}
	acc, err := access.NewAccess(conf.Access, name, sec)
	if err != nil {
		return nil, errors.New("Could not create access for service '%v'.", name).CausedBy(err)
	}
	log, err := logger.NewLogger(conf.Logger, name, syncGroup.NewSyncGroup())
	if err != nil {
		return nil, errors.New("Could not create logger for service '%v'.", name).CausedBy(err)
	}
	result := &Service{
		config: conf,
		logger: log,
		name: name,
		syncGroup: syncGroup,
		access: acc,
	}
	runtime.SetFinalizer(result, finalize)
	return result, nil
}

func (this *Service) Close() {
	this.logger.Close()
}

func (this Service) String() string {
	return this.Name()
}

func (this Service) Name() string {
	return this.name
}

func (this Service) Config() Config {
	return this.config
}

func (this Service) Logger() *logger.Logger {
	return this.logger
}

func (this *Service) Access() *access.Access {
	return this.access
}
