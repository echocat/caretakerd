package config

import (
    "strings"
    "github.com/echocat/caretakerd/service/kind"
    "github.com/echocat/caretakerd/errors"
)

func (this Configs) Validate() error {
    for name, service := range this {
        err := this.validateService(service, name)
        if err != nil {
            return err
        }
    }
    return nil
}

func (this Configs) ValidateMaster() error {
    masters := []string{}
    for name, service := range this {
        if service.Kind == kind.Master {
            masters = append(masters, name)
        }
    }
    if len(masters) == 0 {
        return errors.New("There is no service of kind master defined.")
    }
    if len(masters) > 1 {
        return errors.New("There are more then 0 service of kind master defined: %s", strings.Join(masters, ", "))
    }
    return nil
}

func (this Configs) validateService(service Config, name string) error {
    err := service.Validate()
    if err != nil {
        return errors.New("Config of '%v' service is not valid.", name).CausedBy(err)
    }
    return nil
}

func (this Config) Validate() error {
    err := this.Logger.Validate()
    if err == nil {
        err = this.validateCommand()
    }
    if err == nil {
        err = this.Kind.Validate()
    }
    if err == nil {
        err = this.StartDelayInSeconds.Validate()
    }
    if err == nil {
        err = this.RestartDelayInSeconds.Validate()
    }
    if err == nil {
        err = this.StopSignal.Validate()
    }
    if err == nil {
        err = this.StopWaitInSeconds.Validate()
    }
    if err == nil {
        err = this.AutoRestart.Validate()
    }
    return err
}

func (i Config) validateCommand() error {
    if len (i.Command) <= 0 {
        return errors.New("There is no command defined.")
    }
    return nil
}
