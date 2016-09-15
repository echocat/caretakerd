package app

import (
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/stack"
	"github.com/echocat/caretakerd/sync"
	"github.com/echocat/caretakerd/values"
	"github.com/urfave/cli"
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

func ensureDaemonConfig(context *cli.Context) error {
	err := ensureConfig(true)
	if err != nil {
		return err
	}
	attachArgsToMasterIfPossible(context.Args(), conf.Instance())
	return nil
}

func runDaemon(conf caretakerd.Config, args []string) {
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

func registerDaemonCommandsAt(executableType ExecutableType, app *cli.App) {
	var name string
	switch executableType {
	case Daemon:
		name = "run"
	default:
		name = "daemon"
	}

	app.Commands = append(app.Commands, cli.Command{
		Name:            name,
		SkipFlagParsing: true,
		ArgsUsage:       "[<args pass to master service>...]",
		Usage:           "Run " + caretakerd.DaemonName + " in forground.",
		Before:          ensureDaemonConfig,
		Action: func(context *cli.Context) {
			runDaemon(*conf.instance, context.Args())
		},
		OnUsageError: onUsageErrorFor(name),
	})

}
