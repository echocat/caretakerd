package app

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/echocat/caretakerd/client"
	"github.com/echocat/caretakerd/stack"
	. "github.com/echocat/caretakerd/values"
	"os"
	"strings"
)

type DirectError struct {
	error string
}

func (instance DirectError) Error() string {
	return instance.error
}

func actionWrapper(clientFactory *client.ClientFactory, command func(context *cli.Context, client *client.Client) error) func(context *cli.Context) {
	return func(context *cli.Context) {
		cli, err := clientFactory.NewClient()
		if err != nil {
			stack.Print(err, os.Stderr, 0)
			os.Exit(1)
		}
		err = command(context, cli)
		if de, ok := err.(DirectError); ok {
			fmt.Fprintln(os.Stderr, de.Error())
			os.Exit(1)
		} else if _, ok := err.(client.ConflictError); ok {
			os.Exit(1)
		} else if ade, ok := err.(client.AccessDeniedError); ok {
			fmt.Fprintln(os.Stderr, ade.Error())
			os.Exit(1)
		} else if snfe, ok := err.(client.ServiceNotFoundError); ok {
			fmt.Fprintln(os.Stderr, snfe.Error())
			os.Exit(1)
		} else if err != nil {
			stack.Print(err, os.Stderr, 0)
			os.Exit(1)
		}
	}
}

func handleJsonResponse(response interface{}, err error) error {
	if err != nil {
		return err
	}
	if s, ok := response.(string); ok {
		fmt.Fprintln(os.Stdout, s)
	} else if i, ok := response.(Integer); ok {
		fmt.Fprintln(os.Stdout, i)
	} else {
		jConf, err := json.MarshalIndent(response, "", "   ")
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, string(jConf))
	}
	return nil
}

func globalSpecificGetActionWrapper(clientFactory *client.ClientFactory, action func(client *client.Client) (interface{}, error)) func(context *cli.Context) {
	return actionWrapper(clientFactory, func(context *cli.Context, client *client.Client) error {
		return globalSpecificGetAction(context, client, action)
	})
}

func globalSpecificGetAction(context *cli.Context, client *client.Client, action func(client *client.Client) (interface{}, error)) error {
	response, err := action(client)
	return handleJsonResponse(response, err)
}

func serviceSpecificGetActionWrapper(clientFactory *client.ClientFactory, action func(name string, client *client.Client) (interface{}, error)) func(context *cli.Context) {
	return actionWrapper(clientFactory, func(context *cli.Context, client *client.Client) error {
		return serviceSpecificGetAction(context, client, action)
	})
}

func serviceSpecificGetAction(context *cli.Context, client *client.Client, action func(name string, client *client.Client) (interface{}, error)) error {
	args := context.Args()
	if len(args) != 1 {
		return DirectError{error: fmt.Sprintf("Illegal number of arguments (%d) for command %v", len(args), context.Command.Name)}
	}
	response, err := action(args[0], client)
	return handleJsonResponse(response, err)
}

func serviceSpecificTriggerActionWrapper(clientFactory *client.ClientFactory, action func(name string, client *client.Client) error) func(context *cli.Context) {
	return actionWrapper(clientFactory, func(context *cli.Context, client *client.Client) error {
		return serviceSpecificTriggerAction(context, client, action)
	})
}

func serviceSpecificTriggerAction(context *cli.Context, client *client.Client, action func(name string, client *client.Client) error) error {
	args := context.Args()
	return action(args[0], client)
}

func baseControlEnsure(context *cli.Context) error {
	return ensureConfig(false, conf)
}

func ensureNoControlArgument(context *cli.Context) error {
	err := baseControlEnsure(context)
	if err != nil {
		return err
	}
	if len(context.Args()) != 0 {
		return DirectError{error: "There is only no argument allowed."}
	}
	return nil
}

func ensureServiceNameArgument(context *cli.Context) error {
	err := baseControlEnsure(context)
	if err != nil {
		return err
	}
	args := context.Args()
	if (len(args) <= 0) || (len(strings.TrimSpace(args[0])) == 0) {
		return DirectError{error: "There is no service name provided."}
	}
	if len(args) > 1 {
		return DirectError{error: "There is only one argument allowed."}
	}
	return nil
}

func ensureServiceNameAndSignalArgument(context *cli.Context) error {
	err := baseControlEnsure(context)
	if err != nil {
		return err
	}
	args := context.Args()
	if (len(args) <= 0) || (len(strings.TrimSpace(args[0])) == 0) {
		return DirectError{error: "There is no service name provided."}
	}
	if (len(args) <= 1) || (len(strings.TrimSpace(args[1])) == 0) {
		return DirectError{error: "There is no signal provided."}
	}
	var sig Signal
	err = sig.Set(args[1])
	if err != nil {
		return DirectError{error: err.Error()}
	}
	if len(args) > 2 {
		return DirectError{error: "There are only two arguments allowed."}
	}
	return nil
}

func createConfigCommand(commonClientFlags []cli.Flag, clientFactory *client.ClientFactory) cli.Command {
	return cli.Command{
		Name:      "config",
		Flags:     commonClientFlags,
		ArgsUsage: " ",
		Usage:     "Query whole daemon configuration.",
		Before:    ensureNoControlArgument,
		Action: globalSpecificGetActionWrapper(clientFactory, func(client *client.Client) (interface{}, error) {
			return client.GetConfig()
		}),
		OnUsageError: onUsageErrorFor("config"),
	}
}

func createControlConfigCommand(commonClientFlags []cli.Flag, clientFactory *client.ClientFactory) cli.Command {
	return cli.Command{
		Name:      "controlConfig",
		Flags:     commonClientFlags,
		ArgsUsage: " ",
		Usage:     "Query control configuration.",
		Before:    ensureNoControlArgument,
		Action: globalSpecificGetActionWrapper(clientFactory, func(client *client.Client) (interface{}, error) {
			return client.GetControlConfig()
		}),
		OnUsageError: onUsageErrorFor("controlConfig"),
	}
}

func createServicesCommand(commonClientFlags []cli.Flag, clientFactory *client.ClientFactory) cli.Command {
	return cli.Command{
		Name:      "services",
		Flags:     commonClientFlags,
		ArgsUsage: " ",
		Usage:     "Query whole daemon configuration with all its actual service stats.",
		Before:    ensureNoControlArgument,
		Action: globalSpecificGetActionWrapper(clientFactory, func(client *client.Client) (interface{}, error) {
			return client.GetServices()
		}),
		OnUsageError: onUsageErrorFor("services"),
	}
}

func createServiceCommand(commonClientFlags []cli.Flag, clientFactory *client.ClientFactory) cli.Command {
	return cli.Command{
		Name:      "service",
		Flags:     commonClientFlags,
		ArgsUsage: "<service name>",
		Usage:     "Query service configuration and its actual stats.",
		Before:    ensureServiceNameArgument,
		Action: serviceSpecificGetActionWrapper(clientFactory, func(name string, client *client.Client) (interface{}, error) {
			return client.GetService(name)
		}),
		OnUsageError: onUsageErrorFor("service"),
	}
}

func createServiceConfigCommand(commonClientFlags []cli.Flag, clientFactory *client.ClientFactory) cli.Command {
	return cli.Command{
		Name:      "serviceConfig",
		Flags:     commonClientFlags,
		ArgsUsage: "<service name>",
		Usage:     "Query service configuration.",
		Before:    ensureServiceNameArgument,
		Action: serviceSpecificGetActionWrapper(clientFactory, func(name string, client *client.Client) (interface{}, error) {
			return client.GetServiceConfig(name)
		}),
		OnUsageError: onUsageErrorFor("serviceConfig"),
	}
}

func createServiceStatusCommand(commonClientFlags []cli.Flag, clientFactory *client.ClientFactory) cli.Command {
	return cli.Command{
		Name:      "serviceStatus",
		Flags:     commonClientFlags,
		ArgsUsage: "<service name>",
		Aliases:   []string{"status"},
		Usage:     "Query service status.",
		Before:    ensureServiceNameArgument,
		Action: serviceSpecificGetActionWrapper(clientFactory, func(name string, client *client.Client) (interface{}, error) {
			return client.GetServiceStatus(name)
		}),
		OnUsageError: onUsageErrorFor("serviceStatus"),
	}
}

func createServicePidCommand(commonClientFlags []cli.Flag, clientFactory *client.ClientFactory) cli.Command {
	return cli.Command{
		Name:      "servicePid",
		Flags:     commonClientFlags,
		ArgsUsage: "<service name>",
		Aliases:   []string{"pid"},
		Usage:     "Query service pid.",
		Before:    ensureServiceNameArgument,
		Action: serviceSpecificGetActionWrapper(clientFactory, func(name string, client *client.Client) (interface{}, error) {
			return client.GetServicePid(name)
		}),
		OnUsageError: onUsageErrorFor("servicePid"),
	}
}

func createStartServiceCommand(commonClientFlags []cli.Flag, clientFactory *client.ClientFactory) cli.Command {
	return cli.Command{
		Name:      "serviceStart",
		Flags:     commonClientFlags,
		ArgsUsage: "<service name>",
		Aliases:   []string{"start"},
		Usage:     "Start a service.",
		Before:    ensureServiceNameArgument,
		Action: serviceSpecificTriggerActionWrapper(clientFactory, func(name string, client *client.Client) error {
			return client.StartService(name)
		}),
		OnUsageError: onUsageErrorFor("serviceStart"),
	}
}

func createRestartServiceCommand(commonClientFlags []cli.Flag, clientFactory *client.ClientFactory) cli.Command {
	return cli.Command{
		Name:      "serviceRestart",
		Flags:     commonClientFlags,
		ArgsUsage: "<service name>",
		Aliases:   []string{"restart"},
		Usage:     "Restart a service.",
		Before:    ensureServiceNameArgument,
		Action: serviceSpecificTriggerActionWrapper(clientFactory, func(name string, client *client.Client) error {
			return client.RestartService(name)
		}),
		OnUsageError: onUsageErrorFor("serviceRestart"),
	}
}

func createStopServiceCommand(commonClientFlags []cli.Flag, clientFactory *client.ClientFactory) cli.Command {
	return cli.Command{
		Name:      "serviceStop",
		Flags:     commonClientFlags,
		ArgsUsage: "<service name>",
		Aliases:   []string{"stop"},
		Usage:     "Stop a service.",
		Before:    ensureServiceNameArgument,
		Action: serviceSpecificTriggerActionWrapper(clientFactory, func(name string, client *client.Client) error {
			return client.StopService(name)
		}),
		OnUsageError: onUsageErrorFor("serviceStop"),
	}
}

func createKillServiceCommand(commonClientFlags []cli.Flag, clientFactory *client.ClientFactory) cli.Command {
	return cli.Command{
		Name:      "serviceKill",
		Flags:     commonClientFlags,
		ArgsUsage: "<service name>",
		Aliases:   []string{"kill"},
		Usage:     "Kill a service.",
		Before:    ensureServiceNameArgument,
		Action: serviceSpecificTriggerActionWrapper(clientFactory, func(name string, client *client.Client) error {
			return client.KillService(name)
		}),
		OnUsageError: onUsageErrorFor("serviceKill"),
	}
}

func createSignalServiceCommand(commonClientFlags []cli.Flag, clientFactory *client.ClientFactory) cli.Command {
	return cli.Command{
		Name:      "serviceSignal",
		Flags:     commonClientFlags,
		ArgsUsage: "<service name> <signal>",
		Aliases:   []string{"signal"},
		Usage:     "Send a signal to service.",
		Before:    ensureServiceNameAndSignalArgument,
		Action: actionWrapper(clientFactory, func(context *cli.Context, client *client.Client) error {
			args := context.Args()
			name := args[0]
			var sig Signal
			err := sig.Set(args[1])
			if err != nil {
				return DirectError{error: err.Error()}
			}
			return client.SignalService(name, sig)
		}),
		OnUsageError: onUsageErrorFor("serviceSignal"),
	}
}

func registerControlCommandsAt(app *cli.App) {
	clientFactory := client.NewClientFactory(conf.Instance())

	commonClientFlags := []cli.Flag{}

	app.Commands = append(app.Commands,
		createConfigCommand(commonClientFlags, clientFactory),
		createControlConfigCommand(commonClientFlags, clientFactory),
		createServicesCommand(commonClientFlags, clientFactory),
		createServiceCommand(commonClientFlags, clientFactory),
		createServiceConfigCommand(commonClientFlags, clientFactory),
		createServiceStatusCommand(commonClientFlags, clientFactory),
		createServicePidCommand(commonClientFlags, clientFactory),
		createStartServiceCommand(commonClientFlags, clientFactory),
		createRestartServiceCommand(commonClientFlags, clientFactory),
		createStopServiceCommand(commonClientFlags, clientFactory),
		createKillServiceCommand(commonClientFlags, clientFactory),
		createSignalServiceCommand(commonClientFlags, clientFactory),
	)
}
