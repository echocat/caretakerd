package app

import (
	"github.com/codegangsta/cli"
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/defaults"
	"github.com/echocat/caretakerd/errors"
	"os"
	"fmt"
)

var conf = NewConfigWrapper()
var defaultListenAddress = defaults.ListenAddress()
var defaultPemFile = defaults.AuthFileKeyFilename()
var listenAddress = NewFlagWrapper(&defaultListenAddress)
var pemFile = NewFlagWrapper(&defaultPemFile)

func init() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .Flags}}[global options]{{end}} command ...
{{if .Commands}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`
}

func NewApps() map[ExecutableType]*cli.App {
	result := map[ExecutableType]*cli.App{}
	for _, executableType := range AllExecutableTypes {
		result[executableType] = NewAppFor(executableType)
	}
	return result
}

func NewAppFor(executableType ExecutableType) *cli.App{
	app := newAppFor(executableType)
	registerCommandsFor(executableType, app)

	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(os.Stderr, "Command does not exist: %v\n\n", command)
		cli.HelpPrinter(os.Stderr, cli.AppHelpTemplate, app)
	}

	return app
}

func newAppFor(executableType ExecutableType) *cli.App {
	var configDescription string
	var configEnvVar string
	switch executableType {
	case Daemon:
		configDescription = "Configuration file for daemon."
		configEnvVar = "CTD_CONFIG"
	case Control:
		configDescription = "Configuration file for control."
		configEnvVar = "CTCTL_CONFIG"
	default:
		configDescription = "Configuration file for daemon and control."
		configEnvVar = "CT_CONFIG"
	}

	app := cli.NewApp()
	app.Version = caretakerd.VERSION
	app.Commands = []cli.Command{}
	app.OnUsageError =  func(context *cli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(app.Writer, "Error: %v\n\n", err)
		if isSubcommand {
			cli.ShowSubcommandHelp(context)
		} else {
			cli.ShowAppHelp(context)
		}
		return err
	}
	app.Flags = []cli.Flag{
		cli.GenericFlag{
			Name:   "config,c",
			Value:  conf,
			Usage:  configDescription,
			EnvVar: configEnvVar,
		},
		cli.GenericFlag{
			Name:  "address,a",
			Value: listenAddress,
			Usage: "Listen address of the daemon.",
		},
	}

	if executableType != Daemon {
		app.Flags = append(app.Flags, cli.GenericFlag{
			Name:  "pem,p",
			Value: pemFile,
			Usage: "Location of PEM file which contains the private public key pair for access to the daemon.",
		})
	}

	switch executableType {
	case Daemon:
		app.Name = caretakerd.DAEMON_NAME
		app.Usage = "Simple control daemon for processes."
	case Control:
		app.Name = caretakerd.CONTROL_NAME
		app.Usage = "Remote control for " + caretakerd.DAEMON_NAME
	default:
		app.Name = caretakerd.BASE_NAME
		app.Usage = "Simple control daemon for processes including remote control for itself."
	}
	return app
}

func registerCommandsFor(executableType ExecutableType, at *cli.App) {
	switch executableType {
	case Daemon:
		registerDaemonCommandsAt(executableType, at)
	case Control:
		registerControlCommandsAt(at)
	default:
		registerDaemonCommandsAt(executableType, at)
		registerControlCommandsAt(at)
	}
}

type ExecutableType int

const (
	Daemon  ExecutableType = 0
	Control ExecutableType = 1
	Generic ExecutableType = 2
)

var AllExecutableTypes = []ExecutableType{
	Daemon,
	Control,
	Generic,
}

func (instance ExecutableType) String() string {
	switch instance {
	case Daemon:
		return caretakerd.DAEMON_NAME
	case Control:
		return caretakerd.CONTROL_NAME
	}
	return caretakerd.BASE_NAME
}

func ensureConfig(daemonChecks bool) error {
	if conf.explicitSet {
		return conf.ConfigureAndValidate(listenAddress, pemFile, daemonChecks)
	}
	newConf := NewConfigWrapper()
	err := newConf.Set(newConf.String())
	if err != nil {
		if _, ok := err.(caretakerd.ConfigDoesNotExistError); ok {
			if daemonChecks {
				return errors.New("There is neither the --config flag set nor does a configuration file under default position (%v) exist.", newConf.String())
			} else {
				return conf.ConfigureAndValidate(listenAddress, pemFile, daemonChecks)
			}
		} else {
			return err
		}
	}
	err = newConf.ConfigureAndValidate(listenAddress, pemFile, daemonChecks)
	if err != nil {
		return err
	}
	conf = newConf
	return nil
}

func onUsageErrorFor(commandName string) func(context *cli.Context, err error) error {
	return func(context *cli.Context, err error) error {
		fmt.Fprintf(context.App.Writer, "Error: %v\n\n", err)
		cli.ShowCommandHelp(context, commandName)
		return err
	}
}
