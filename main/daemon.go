package main

import (
    "os"
    "github.com/codegangsta/cli"
    "github.com/echocat/caretakerd"
    "github.com/echocat/caretakerd/config"
    "github.com/echocat/caretakerd/logger"
    "github.com/echocat/caretakerd/values"
    "github.com/echocat/caretakerd/sync"
    "github.com/echocat/caretakerd/stack"
)

func attachArgsToMasterIfPossible(args []string, to *config.Config) {
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
    attachArgsToMasterIfPossible(context.Args(), conf.Instance())
    return ensureConfig(true, conf)
}

func runDaemon(conf config.Config, args []string) {
    instance, err := caretakerd.NewCaretakerd(conf, sync.NewSyncGroup())
    if err != nil {
        stack.Print(err, os.Stderr, 0)
        os.Exit(1)
    }

    instance.Logger().Log(logger.Debug, caretakerd.DAEMON_NAME  + " successful loaded. Starting now services...")
    exitCode, _ := instance.Run()
    instance.Logger().Log(logger.Debug, caretakerd.DAEMON_NAME  + " done.")

    instance.Close()

    os.Exit(int(exitCode))
}

func registerDaemonCommandsAt(executableType ExecutableType, app *cli.App) {
    var name string
    switch executableType {
    case daemon:
        name = "run"
    default:
        name = "daemon"
    }

    app.Commands = append(app.Commands, cli.Command{
        Name: name,
        SkipFlagParsing: true,
        ArgsUsage: "[<args pass to master service>...]",
        Usage: "Run " + caretakerd.DAEMON_NAME + " in forground.",
        Before: ensureDaemonConfig,
        Action: func(context *cli.Context) {
            runDaemon(*conf.instance, context.Args())
        },
    })

}
