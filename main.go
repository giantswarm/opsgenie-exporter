package main

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/microkit/command"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/viper"

	"github.com/giantswarm/opsgenie-exporter/flag"
	"github.com/giantswarm/opsgenie-exporter/server"
	"github.com/giantswarm/opsgenie-exporter/service"
)

var (
	description string     = "The opsgenie-exporter exports Opsgenie data as Prometheus metrics."
	f           *flag.Flag = flag.New()
	gitCommit   string     = "n/a"
	name        string     = "opsgenie-exporter"
	source      string     = "https://github.com/giantswarm/opsgenie-exporter"
)

func main() {
	err := mainError()
	if err != nil {
		panic(fmt.Sprintf("%#v\n", err))
	}
}

func mainError() error {
	var err error

	// Create a new logger which is used by all packages.
	var newLogger micrologger.Logger
	{
		c := micrologger.Config{}

		newLogger, err = micrologger.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is storted out.
	newServerFactory := func(v *viper.Viper) microserver.Server {
		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			c := service.Config{
				Logger: newLogger,

				Description: description,
				Flag:        f,
				GitCommit:   gitCommit,
				ProjectName: name,
				Source:      source,
				Viper:       v,
			}

			newService, err = service.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v", err))
			}
			go newService.Boot(context.Background())
		}

		// Create a new custom server which bundles our endpoints.
		var newServer microserver.Server
		{
			c := server.Config{
				Logger:  newLogger,
				Service: newService,
				Viper:   v,

				ProjectName: name,
			}

			newServer, err = server.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v", err))
			}
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		c := command.Config{
			Logger:        newLogger,
			ServerFactory: newServerFactory,

			Description: description,
			GitCommit:   gitCommit,
			Name:        name,
			Source:      source,
		}

		newCommand, err = command.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	daemonCommand := newCommand.DaemonCommand().CobraCommand()

	daemonCommand.PersistentFlags().String(f.Service.Opsgenie.API.Token, "", "Auth token to access the Opsgenie API.")

	newCommand.CobraCommand().Execute()

	return nil
}
