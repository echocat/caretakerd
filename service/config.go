package service

import (
	"github.com/echocat/caretakerd/access"
	"github.com/echocat/caretakerd/logger"
	. "github.com/echocat/caretakerd/values"
)

type Config struct {
	// @default autoStart
	//
	// Defines how this service will run by caretakerd.
	//
	// For details of possible values see {@ref github.com/echocat/caretakerd/values.RestartType}.
	//
	// > **Important**: Exact one of the services have to be configured as
	// > {@ref github.com/echocat/caretakerd/service.Type#Master master}.
	Type                  Type               `json:"type" yaml:"type"`

	// @default []
	//
	// The command the service process has to start with. The called command have to be run in foreground - or in other words: Should not daemonize.
	//
	// > **Hint**: If there is no command line provided this service cannot be started and caretakerd will
	// > fail.
	//
	// # PATH expansion
	//
	// The provided commands are resolved from the ``PATH`` environment provided to caretakerd.
	// This makes it possible to use the names of the binaries like ``sleep`` instead of ``/usr/bin/sleep``.
	//
	// # Parameter evaluation
	//
	// Environment variables could be included like:
	// ```yaml
	// command: ["echo", "${MESSAGE}"]
	// environment:
	//     MESSAGE: "Hello world!"
	// ```
	//
	// ```bash
	// $ caretakerd run
	// Hello world!
	// ```
	//
	// # Special master handling
	//
	// If the service is configured as {@ref #Type type} = {@ref github.com/echocat/caretakerd/service.Type#Master master}
	// every parameter which was passed to caretakerd itself will be enriched to the called command line of the service process.
	//
	// Config example:
	// ```yaml
	// command: ["echo", "Hello"]
	// ```
	//
	// Run examples:
	// ```bash
	// $ caretakerd run
	// Hello
	//
	// $ caretakerd run "world!"
	// Hello world!
	// ```
	Command               []String           `json:"command" yaml:"command,flow"`

	// @default []
	//
	// Commands to be executed before execution of the actual {@ref #Command command}.
	//
	// If one of these commands fails also the whole service is marked as failed. The actual
	// {@ref #Command command} will not be invoked and the {@ref #AutoRestart autoRestart} handling will be initiated.
	//
	// Only exit codes of value ``0`` will be accepted as success.
	//
	// If there is a minus (``-``) provided as first item of the command every error of this command will be ignored.
	//
	// Example:
	// ```yaml
	// preCommands:
	// - ["-", "program.sh", "prepare"]        # Ignore if fails
	// - ["program.sh", "prepareAndDoNotFail"] # Do not ignore if fails
	// command: ["program.sh", "run"]
	// ```
	PreCommands           [][]String         `json:"preCommands" yaml:"preCommands,flow"`

	// @default []
	//
	// Commands to be executed after execution of the actual {@ref #Command command}.
	//
	// Every result of these commands are ignored and will not force another beauvoir - except: an error log entry in log files.
	//
	// Only exit codes of value ``0`` will be accepted as success.
	//
	// If there is a minus (``-``) provided as first item of the command every error of this command will be ignored.
	//
	// Example:
	// ```yaml
	// command: ["program.sh", "run"]
	// postCommands:
	// - ["-", "program.sh", "cleanUp"]        # Ignore if fails
	// - ["program.sh", "cleanUpAndDoNotFail"] # Log if fails
	// ```
	PostCommands          [][]String         `json:"postCommands" yaml:"postCommands,flow"`

	// @default ""
	//
	// If configured this will trigger the service at this specific times. If not the service will
	// run as a normal process just once (except of the {@ref #AutoRestart autoRestart} handling).
	//
	// For details of possible values see {@ref github.com/echocat/caretakerd/service.CronExpression}.
	CronExpression        CronExpression     `json:"cronExpression" yaml:"cronExpression"`

	// @default 0
	//
	// Wait before the service process will start the first time.
	//
	// > **Hint:** Every run triggered by {@ref #CronExpression cronExpression} will also wait for this delay.
	StartDelayInSeconds   NonNegativeInteger `json:"startDelayInSeconds" yaml:"startDelayInSeconds"`

	// @default [0]
	//
	// Every of these values represents an expected success exit codes.
	// If a service ends with one of these values, the service will not be restarted.
	// Other values will trigger a auto restart if configured.
	//
	// See: {@ref #AutoRestart autoRestart}
	SuccessExitCodes      ExitCodes          `json:"successExitCodes" yaml:"successExitCodes,flow"`

	// @default "TERM" (on Unix like systems) - "KILL" (on Windows)
	//
	// Signal which will be send to the service when a stop is requested.
	// You can use the signal number here and also names like ``"TERM"`` or ``"KILL"``.
	StopSignal            Signal             `json:"stopSignal" yaml:"stopSignal"`

	// @default []
	//
	// Command to be execute to stop the service.
	//
	// From the moment on this command is called the {@ref #StopWaitInSeconds stopWaitInSeconds} are running.
	// It is not important when this stopCommand ends or what is the exit code.
	// If this command is executed and the service does not end within the configured {@ref #StopWaitInSeconds stopWaitInSeconds}
	// the service will be killed.
	//
	// If this property is configured {@ref #StopSignal stopSignal} will not be evaluated.
	StopCommand           []String           `json:"stopCommand" yaml:"stopCommand,flow"`

	// @default 30
	//
	// Timeout to wait before kill the service process after a stop is requested.
	StopWaitInSeconds     NonNegativeInteger `json:"stopWaitInSeconds" yaml:"stopWaitInSeconds"`

	// @default ""
	//
	// User under which the service process will be started.
	User                  String             `json:"user" yaml:"user"`

	// @default []
	//
	// Environment variables to pass to the process.
	Environment           Environments       `json:"environment" yaml:"environment"`

	// @default true
	//
	// Pass the environment variables started with caretakerd also to the service process.
	InheritEnvironment    Boolean            `json:"inheritEnvironment" yaml:"inheritEnvironment"`

	// @default ""
	//
	// Working directory to start the service process in.
	Directory             String             `json:"directory" yaml:"directory"`

	// @default onFailures
	//
	// Configure how caretakerd will handle the end of a process.
	// It depends mainly on the {@ref #SuccessExitCodes successExitCodes} property.
	//
	// For details of possible values see {@ref github.com/echocat/caretakerd/values.RestartType}.
	AutoRestart           RestartType        `json:"autoRestart" yaml:"autoRestart"`

	// @default 5
	//
	// Seconds to wait before restart of a process.
	//
	// If a process should be restarted (because of {@ref #AutoRestart autoRestart}) caretakerd will wait this seconds before restart is initiated.
	RestartDelayInSeconds NonNegativeInteger `json:"restartDelayInSeconds" yaml:"restartDelayInSeconds"`

	// Configures the permission of this service to control caretakerd remotely
	// and how to obtain the credentials for it.
	//
	// For details see {@ref github.com/echocat/caretakerd/access.Config}.
	Access                access.Config      `json:"access" yaml:"access,omitempty"`

	// Configures the logger for this specific service.
	//
	// For details see {@ref github.com/echocat/caretakerd/logger.Config}.
	Logger                logger.Config      `json:"logger" yaml:"logger,omitempty"`
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
	(*i).PreCommands = [][]String{}
	(*i).PostCommands = [][]String{}
	(*i).Type = AutoStart
	(*i).CronExpression = NewCronExpression()
	(*i).StartDelayInSeconds = NonNegativeInteger(0)
	(*i).RestartDelayInSeconds = NonNegativeInteger(5)
	(*i).SuccessExitCodes = ExitCodes{ExitCode(0)}
	(*i).StopSignal = defaultStopSignal()
	(*i).StopCommand = []String{}
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
