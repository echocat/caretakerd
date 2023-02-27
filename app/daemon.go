package app

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/stack"
	"github.com/echocat/caretakerd/sync"
	"github.com/echocat/caretakerd/values"
	"os"
)

func attachArgsToMasterIfPossible(args []string, to *caretakerd.Config) {
	if len(args) > 0 {
		if masterName, ok := to.Services.GetMasterName(); ok {
			master := to.Services[masterName]
			for i := 0; i < len(args); i++ {
				master.Command = append(master.Command, values.String(args[i]))
			}
			to.Services[masterName] = master
		}
	}
}

func runDaemon(conf *caretakerd.Config, args []string) {
	attachArgsToMasterIfPossible(args, conf)
	instance, err := caretakerd.NewCaretakerd(conf, sync.NewGroup())
	if err != nil {
		stack.Print(err, os.Stderr, 0)
		os.Exit(1)
	}

	instance.Logger().Log(logger.Debug, caretakerd.DaemonName+" successful loaded. Starting now services...")
	exitCode, _ := instance.Run()
	instance.Logger().Log(logger.Debug, caretakerd.DaemonName+" done.")

	instance.Close()

	os.Exit(int(exitCode))
}

func registerDaemonCommandsAt(config *ConfigWrapper, executableType ExecutableType, app *kingpin.Application) {
	var name string
	switch executableType {
	case Daemon:
		name = "run"
	default:
		name = "daemon"
	}

	cmd := app.Command(name, "Run "+caretakerd.DaemonName+" in foreground.")

	if name == "daemon" {
		cmd.Alias("run")
	}

	arguments := cmd.Arg("args", "argument to be passed to master service").
		Strings()

	cmd.Action(func(*kingpin.ParseContext) error {
		config, err := config.ProvideConfig(true)
		if err != nil {
			return err
		}
		runDaemon(config, *arguments)
		return nil
	})
}
