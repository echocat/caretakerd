package main

import (
    "os"
    "regexp"
    "strings"
    "fmt"
    "github.com/codegangsta/cli"
    "github.com/echocat/caretakerd"
    "github.com/echocat/caretakerd/panics"
    "github.com/echocat/caretakerd/defaults"
    "github.com/echocat/caretakerd/errors"
    "github.com/echocat/caretakerd/config"
)

var executableNamePattern = regexp.MustCompile("(?:^|" + regexp.QuoteMeta(string(os.PathSeparator)) + ")" + caretakerd.BASE_NAME + "(d|ctl)(?:$|[\\.\\-\\_].*$)")
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

func main() {
    defer panics.DefaultPanicHandler()

    executableType := getExecutableType()

    app := newAppFor(executableType)
    registerCommandsFor(executableType, app)

    app.CommandNotFound = func(c *cli.Context, command string) {
        fmt.Fprintf(os.Stderr, "Command does not exist: %v\n\n", command)
        cli.HelpPrinter(os.Stderr, cli.AppHelpTemplate, app)
    }

    err := app.Run(os.Args)
    if err != nil {
        os.Exit(1)
    }
}

func newAppFor(executableType ExecutableType) *cli.App {
    var configDescription string
    var configEnvVar string
    switch executableType {
    case daemon:
        configDescription = "Configuration file for daemon."
        configEnvVar = "CTD_CONFIG"
    case control:
        configDescription = "Configuration file for control."
        configEnvVar = "CTCTL_CONFIG"
    default:
        configDescription = "Configuration file for daemon and control."
        configEnvVar = "CT_CONFIG"
    }

    app := cli.NewApp()
    app.Version = caretakerd.VERSION
    app.Commands = []cli.Command{}
    app.Flags = []cli.Flag{
        cli.GenericFlag{
            Name: "config,c",
            Value: conf,
            Usage: configDescription,
            EnvVar: configEnvVar,
        },
        cli.GenericFlag{
            Name: "address,a",
            Value: listenAddress,
            Usage: "Listen address of the daemon.",
        },
    }

    if executableType != daemon {
        app.Flags = append(app.Flags, cli.GenericFlag{
            Name: "pem,p",
            Value: pemFile,
            Usage: "Location of PEM file which contains the private public key pair for access to the daemon.",
        })
    }

    switch executableType {
    case daemon:
        app.Name = caretakerd.DAEMON_NAME
        app.Usage = "Simple control daemon for processes."
    case control:
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
    case daemon:
        registerDaemonCommandsAt(executableType, at)
    case control:
        registerControlCommandsAt(at)
    default:
        registerDaemonCommandsAt(executableType, at)
        registerControlCommandsAt(at)
    }
}

type ExecutableType int

const (
    generic ExecutableType = 0
    daemon ExecutableType = 1
    control ExecutableType = 2
)

func getExecutableType() ExecutableType {
    executable := strings.ToLower(os.Args[0])
    match := executableNamePattern.FindStringSubmatch(executable)
    if match != nil && len(match) == 2 {
        switch match[1] {
        case "d":
            return daemon
        case "ctl":
            return control
        }
    }
    return generic
}

func ensureConfig(daemonChecks bool, target *ConfigWrapper) error {
    if target.explicitSet {
        return target.ConfigureAndValidate(listenAddress, pemFile, daemonChecks)
    }
    newConf := NewConfigWrapper()
    err := newConf.Set(defaults.ConfigFilename().String())
    if err != nil {
        if _, ok := err.(config.ConfigDoesNotExistError); ok {
            if daemonChecks {
                return errors.New("There is neither the --config flag set nor does a configuration file under default position (%v) exist.", defaults.ConfigFilename())
            } else {
                return target.ConfigureAndValidate(listenAddress, pemFile, daemonChecks)
            }
        } else {
            return err
        }
    }
    err = newConf.ConfigureAndValidate(listenAddress, pemFile, daemonChecks)
    if err != nil {
        return err
    }
    target = newConf
    return nil
}
