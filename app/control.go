package app

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/echocat/caretakerd/client"
	"github.com/echocat/caretakerd/stack"
	"github.com/echocat/caretakerd/values"
	"os"
)

func actionWrapper(clientFactory *client.Factory, command func(client *client.Client) error) func(context *kingpin.ParseContext) error {
	return func(*kingpin.ParseContext) error {
		cli, err := clientFactory.NewClient()
		if err == nil {
			err = command(cli)
		}
		switch err.(type) {
		case client.ConflictError, client.AccessDeniedError, client.ServiceNotFoundError:
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
		default:
			stack.Print(err, os.Stderr, 0)
		}
		os.Exit(1)
		return nil
	}
}

func handleJSONResponse(response interface{}, err error) error {
	if err != nil {
		return err
	}
	if s, ok := response.(string); ok {
		_, err := fmt.Fprintln(os.Stdout, s)
		return err
	} else if i, ok := response.(values.Integer); ok {
		_, err := fmt.Fprintln(os.Stdout, i)
		return err
	} else {
		jConf, err := json.MarshalIndent(response, "", "   ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(os.Stdout, string(jConf))
		return err
	}
}

func getActionWrapper(clientFactory *client.Factory, action func(client *client.Client) (interface{}, error)) func(context *kingpin.ParseContext) error {
	return actionWrapper(clientFactory, func(client *client.Client) error {
		return handleJSONResponse(action(client))
	})
}

func registerConfigCommand(app *kingpin.Application, clientFactory *client.Factory) {
	cmd := app.Command("config", "Returns configurations for defined service, '!daemon' or '!control'.")

	target := cmd.Arg("target", "Could be either '!daemon' for the daemon itself, '!control' for the control of the daemon or each name of a configuration service.").
		Required().
		String()

	cmd.Action(getActionWrapper(clientFactory, func(client *client.Client) (interface{}, error) {
		if *target == "!daemon" {
			return client.GetConfig()
		}
		if *target == "!control" {
			return client.GetControlConfig()
		}
		return client.GetServiceConfig(*target)
	}))
}

func registerGetCommand(at *kingpin.Application, clientFactory *client.Factory) {
	cmd := at.Command("get", "Query states for given service or if nothing specified for all services.")

	target := cmd.Arg("target", "If specified this service will be queried otherwise all services will be queried.").
		String()

	cmd.Action(getActionWrapper(clientFactory, func(client *client.Client) (interface{}, error) {
		if target != nil && len(*target) > 0 {
			return client.GetService(*target)
		}
		return client.GetServices()
	}))
}

func registerStatusCommand(at *kingpin.Application, clientFactory *client.Factory) {
	cmd := at.Command("status", "Query status of a service.")

	target := cmd.Arg("service", "Service to be queried.").
		Required().
		String()

	cmd.Action(getActionWrapper(clientFactory, func(client *client.Client) (interface{}, error) {
		return client.GetServiceStatus(*target)
	}))
}

func registerPidCommand(at *kingpin.Application, clientFactory *client.Factory) {
	cmd, serviceName := registerServiceNameEnabledCommand(at, "pid", "Query pid of a service.")

	cmd.Action(getActionWrapper(clientFactory, func(client *client.Client) (interface{}, error) {
		return client.GetServicePid(*serviceName)
	}))
}

func registerStartCommand(at *kingpin.Application, clientFactory *client.Factory) {
	cmd, serviceName := registerServiceNameEnabledCommand(at, "start", "Starts a service.")

	cmd.Action(actionWrapper(clientFactory, func(client *client.Client) error {
		return client.StartService(*serviceName)
	}))
}

func registerRestartCommand(at *kingpin.Application, clientFactory *client.Factory) {
	cmd, serviceName := registerServiceNameEnabledCommand(at, "restart", "Restarts a service.")

	cmd.Action(actionWrapper(clientFactory, func(client *client.Client) error {
		return client.RestartService(*serviceName)
	}))
}

func registerStopCommand(at *kingpin.Application, clientFactory *client.Factory) {
	cmd, serviceName := registerServiceNameEnabledCommand(at, "stop", "Stops a service.")

	cmd.Action(actionWrapper(clientFactory, func(client *client.Client) error {
		return client.StopService(*serviceName)
	}))
}

func registerKillCommand(at *kingpin.Application, clientFactory *client.Factory) {
	cmd, serviceName := registerServiceNameEnabledCommand(at, "kill", "Kills a service.")

	cmd.Action(actionWrapper(clientFactory, func(client *client.Client) error {
		return client.KillService(*serviceName)
	}))
}

func registerSignalCommand(at *kingpin.Application, clientFactory *client.Factory) {
	cmd, serviceName := registerServiceNameEnabledCommand(at, "signal", "Send a signal to service.")

	var signal values.Signal
	cmd.Arg("signal", "Signal to be send").
		Required().
		SetValue(&signal)

	cmd.Action(actionWrapper(clientFactory, func(client *client.Client) error {
		return client.SignalService(*serviceName, signal)
	}))
}

func registerServiceNameEnabledCommand(at *kingpin.Application, name, description string) (cmd *kingpin.CmdClause, serviceName *string) {
	cmd = at.Command(name, description)

	serviceName = cmd.Arg("service", "Service to execute the action on.").
		Required().
		String()

	return
}

func registerControlCommands(config *ConfigWrapper, at *kingpin.Application) {
	clientFactory := client.NewFactory(config)

	registerConfigCommand(at, clientFactory)
	registerGetCommand(at, clientFactory)
	registerStatusCommand(at, clientFactory)
	registerPidCommand(at, clientFactory)
	registerStartCommand(at, clientFactory)
	registerRestartCommand(at, clientFactory)
	registerStopCommand(at, clientFactory)
	registerKillCommand(at, clientFactory)
	registerSignalCommand(at, clientFactory)
}
