package service

import (
	"github.com/echocat/caretakerd/errors"
	"strings"
)

// Validate validates actions on this object and returns an error object if there are any.
func (instance Configs) Validate() error {
	for name, service := range instance {
		err := instance.validateService(service, name)
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidateMaster validates whether there is exactly one service defined as master. Returns an error object if there are more services defined as masters.
func (instance Configs) ValidateMaster() error {
	masters := []string{}
	for name, service := range instance {
		if service.Type == Master {
			masters = append(masters, name)
		}
	}
	if len(masters) == 0 {
		return errors.New("There is no service of type master defined.")
	}
	if len(masters) > 1 {
		return errors.New("There are more then 0 service of type master defined: %s", strings.Join(masters, ", "))
	}
	return nil
}

func (instance Configs) validateService(service Config, name string) error {
	err := service.Validate()
	if err != nil {
		return errors.New("Config of '%v' service is not valid.", name).CausedBy(err)
	}
	return nil
}

// Validate validates actions on this object and returns an error object if there are any.
func (instance Config) Validate() error {
	err := instance.Logger.Validate()
	if err == nil {
		err = instance.validateCommand()
	}
	if err == nil {
		err = instance.Type.Validate()
	}
	if err == nil {
		err = instance.StartDelayInSeconds.Validate()
	}
	if err == nil {
		err = instance.RestartDelayInSeconds.Validate()
	}
	if err == nil {
		err = instance.StopSignal.Validate()
	}
	if err == nil {
		err = instance.StopWaitInSeconds.Validate()
	}
	if err == nil {
		err = instance.AutoRestart.Validate()
	}
	return err
}

func (instance Config) validateCommand() error {
	if len(instance.Command) <= 0 {
		return errors.New("There is no command defined.")
	}
	return nil
}
