package service

import (
	"github.com/echocat/caretakerd/access"
	"github.com/echocat/caretakerd/logger"
	. "github.com/echocat/caretakerd/values"
)

type Config struct {
	Logger                logger.Config      `json:"logger" yaml:"logger,omitempty"`
	Command               []String           `json:"command" yaml:"command,flow"`
	Type                  Type               `json:"type" yaml:"type"`
	CronExpression        CronExpression     `json:"cronExpression" yaml:"cronExpression"`
	StartDelayInSeconds   NonNegativeInteger `json:"startDelayInSeconds" yaml:"startDelayInSeconds"`
	RestartDelayInSeconds NonNegativeInteger `json:"restartDelayInSeconds" yaml:"restartDelayInSeconds"`
	SuccessExitCodes      ExitCodes          `json:"successExitCodes" yaml:"successExitCodes,flow"`
	StopSignal            Signal             `json:"stopSignal" yaml:"stopSignal"`
	StopWaitInSeconds     NonNegativeInteger `json:"stopWaitInSeconds" yaml:"stopWaitInSeconds"`
	User                  String             `json:"user" yaml:"user"`
	Environment           Environments       `json:"environment" yaml:"environment"`
	Directory             String             `json:"directory" yaml:"directory"`
	AutoRestart           RestartType        `json:"autoRestart" yaml:"autoRestart"`
	InheritEnvironment    Boolean            `json:"inheritEnvironment" yaml:"inheritEnvironment"`
	Access                access.Config      `json:"access" yaml:"access,omitempty"`
}

func NewConfig() Config {
	result := Config{}
	result.init()
	return result
}

func (i Config) WithCommand(command ...String) Config {
	i.Command = command
	return i
}

func (i *Config) init() {
	(*i).Logger = logger.NewConfig()
	(*i).Command = []String{}
	(*i).Type = AutoStart
	(*i).CronExpression = NewCronExpression()
	(*i).StartDelayInSeconds = NonNegativeInteger(0)
	(*i).RestartDelayInSeconds = NonNegativeInteger(5)
	(*i).SuccessExitCodes = ExitCodes{ExitCode(0)}
	(*i).StopSignal = defaultStopSignal()
	(*i).StopWaitInSeconds = NonNegativeInteger(30)
	(*i).User = String("")
	(*i).Environment = Environments{}
	(*i).Directory = String("")
	(*i).AutoRestart = OnFailures
	(*i).InheritEnvironment = Boolean(true)
	(*i).Access = access.NewNoneConfig()
}

func (i *Config) BeforeUnmarshalYAML() error {
	i.init()
	return nil
}
