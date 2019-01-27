package app

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/defaults"
	"os"
	"runtime"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05Z"
)

var (
	version  = "development"
	revision = "development"
	compiled = ""

	config               = NewConfigWrapper()
	defaultListenAddress = defaults.ListenAddress()
	defaultPemFile       = defaults.AuthFileKeyFilename()
	listenAddress        = NewFlagWrapper(&defaultListenAddress)
	pemFile              = NewFlagWrapper(&defaultPemFile)
)

func init() {
	if compiled == "" {
		compiled = time.Now().Format(timeFormat)
	}
}

func handleVersion(name string) func(*kingpin.ParseContext) error {
	return func(*kingpin.ParseContext) error {
		_, err := fmt.Fprintf(os.Stderr, `%s
 Version:      %s
 Git revision: %s
 Built:        %s
 Go version:   %s
 OS/Arch:      %s/%s
`,
			name, version, revision, compiled, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		if err == nil {
			os.Exit(0)
		}
		return err
	}
}

// NewApps creates new instances of the command line parser (cli.App) for every ExecutableType.
func NewApps() map[ExecutableType]*kingpin.Application {
	result := map[ExecutableType]*kingpin.Application{}
	for _, executableType := range AllExecutableTypes {
		result[executableType] = NewAppFor(executableType)
	}
	return result
}

// NewAppFor creates a new instance of the command line parser (cli.App) for the given executableType.
func NewAppFor(executableType ExecutableType) *kingpin.Application {
	app := newAppFor(executableType)
	registerCommandsFor(executableType, app)

	return app
}

func newAppFor(executableType ExecutableType) *kingpin.Application {
	var app *kingpin.Application
	switch executableType {
	case Daemon:
		app = kingpin.New(caretakerd.DaemonName, "Simple control daemon for processes.")
		app.Flag("config", "Configuration file for daemon.").
			Short('c').
			Envar("CTD_CONFIG").
			PlaceHolder(defaults.ConfigFilename().String()).
			SetValue(config)
	case Control:
		app = kingpin.New(caretakerd.ControlName, "Remote control for "+caretakerd.DaemonName)
		app.Flag("config", "Configuration file for control.").
			Short('c').
			Envar("CTCTL_CONFIG").
			PlaceHolder(defaults.ConfigFilename().String()).
			SetValue(config)
	default:
		app = kingpin.New(caretakerd.BaseName, "Simple control daemon for processes including remote control for itself.")
		app.Flag("config", "Configuration file for daemon and control.").
			Short('c').
			Envar("CT_CONFIG").
			PlaceHolder(defaults.ConfigFilename().String()).
			SetValue(config)
	}

	app.Flag("address", "Listen address of the daemon.").
		Short('a').
		PlaceHolder(listenAddress.String()).
		SetValue(listenAddress)

	if executableType == Daemon {
		config.forDaemon = true
	} else {
		app.Flag("pem", "Location of PEM file which contains the private public key pair for access to the daemon.").
			Short('p').
			PlaceHolder(pemFile.String()).
			SetValue(pemFile)
	}

	app.Command("version", "Print the actual version and other useful information.").
		Action(handleVersion(app.Name))

	return app
}

func registerCommandsFor(executableType ExecutableType, at *kingpin.Application) {
	switch executableType {
	case Daemon:
		registerDaemonCommandsAt(executableType, at)
	case Control:
		registerControlCommands(at)
	default:
		registerDaemonCommandsAt(executableType, at)
		registerControlCommands(at)
	}
}

// ExecutableType represents a type of the caretakerd executable.
type ExecutableType int

const (
	// Daemon indicates that this executable is the caretaker daemon itself.
	Daemon ExecutableType = 0
	// Control indicates that this executable is the caretaker control binary.
	Control ExecutableType = 1
	// Generic indicates that this executable is the caretaker binary which combines
	// daemon and control binary together.
	Generic ExecutableType = 2
)

// AllExecutableTypes contains all possible variants of ExecutableType.
var AllExecutableTypes = []ExecutableType{
	Daemon,
	Control,
	Generic,
}

func (instance ExecutableType) String() string {
	switch instance {
	case Daemon:
		return caretakerd.DaemonName
	case Control:
		return caretakerd.ControlName
	}
	return caretakerd.BaseName
}
