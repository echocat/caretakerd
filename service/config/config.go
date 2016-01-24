package config

import (
    loggerConfig "github.com/echocat/caretakerd/logger/config"
    "github.com/echocat/caretakerd/service/kind"
    "github.com/echocat/caretakerd/service/signal"
    "github.com/echocat/caretakerd/service/restartType"
    "github.com/echocat/caretakerd/service/environment"
    . "github.com/echocat/caretakerd/values"
    "github.com/echocat/caretakerd/service/exitCode"
    "github.com/echocat/caretakerd/access"
    "github.com/echocat/caretakerd/service/cron"
)

type Config struct {
    Logger                loggerConfig.Config `json:"logger" yaml:"logger,omitempty"`
    Command               []String `json:"command" yaml:"command,flow"`
    Kind                  kind.Kind `json:"kind" yaml:"kind"`
    CronExpression        cron.Expression `json:"cronExpression" yaml:"cronExpression"`
    StartDelayInSeconds   NonNegativeInteger `json:"startDelayInSeconds" yaml:"startDelayInSeconds"`
    RestartDelayInSeconds NonNegativeInteger `json:"restartDelayInSeconds" yaml:"restartDelayInSeconds"`
    SuccessExitCodes      exitCode.ExitCodes `json:"successExitCodes" yaml:"successExitCodes,flow"`
    StopSignal            signal.Signal `json:"stopSignal" yaml:"stopSignal"`
    StopWaitInSeconds     NonNegativeInteger `json:"stopWaitInSeconds" yaml:"stopWaitInSeconds"`
    User                  String `json:"user" yaml:"user"`
    Environment           environment.Environments `json:"environment" yaml:"environment"`
    Directory             String `json:"directory" yaml:"directory"`
    AutoRestart           restartType.RestartType `json:"autoRestart" yaml:"autoRestart"`
    InheritEnvironment    Boolean `json:"inheritEnvironment" yaml:"inheritEnvironment"`
    Access                access.Config `json:"access" yaml:"access,omitempty"`
}

func NewServiceConfig() Config {
    result := Config{}
    result.init()
    return result
}

func (i Config) WithCommand(command ...String) Config {
    i.Command = command
    return i
}

func (i *Config) init() {
    (*i).Logger = loggerConfig.NewLoggerConfig()
    (*i).Command = []String{}
    (*i).Kind = kind.AutoStart
    (*i).CronExpression = cron.NewCronExpression()
    (*i).StartDelayInSeconds = NonNegativeInteger(0)
    (*i).RestartDelayInSeconds = NonNegativeInteger(5)
    (*i).SuccessExitCodes = exitCode.ExitCodes{exitCode.ExitCode(0)}
    (*i).StopSignal = defaultStopSignal()
    (*i).StopWaitInSeconds = NonNegativeInteger(30)
    (*i).User = String("")
    (*i).Environment = environment.Environments{}
    (*i).Directory = String("")
    (*i).AutoRestart = restartType.OnFailures
    (*i).InheritEnvironment = Boolean(true)
    (*i).Access = access.NewNoneConfig()
}

func (i *Config) BeforeUnmarshalYAML() error {
    i.init()
    return nil
}
