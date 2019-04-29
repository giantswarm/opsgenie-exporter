package main

import (
	"flag"
	"fmt"

	"github.com/giantswarm/exporterkit"
	"github.com/giantswarm/exporterkit/collector"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/opsgenie-exporter/alert"
	"github.com/giantswarm/opsgenie-exporter/opsgenie"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	var err error

	opsgenieApiKey := flag.String("api-key", "", "Opsgenie API key")
	flag.Parse()

	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var opsgenieClient *opsgenie.Client
	{
		c := opsgenie.Config{
			Key: *opsgenieApiKey,
		}

		opsgenieClient, err = opsgenie.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var alertCollector collector.Interface
	{
		c := alert.Config{
			Client: opsgenieClient,
		}

		alertCollector, err = alert.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var collectorSet *collector.Set
	{
		c := collector.SetConfig{
			Collectors: []collector.Interface{
				alertCollector,
			},
			Logger: logger,
		}

		collectorSet, err = collector.NewSet(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	var exporter *exporterkit.Exporter
	{
		c := exporterkit.Config{
			Collectors: []prometheus.Collector{
				collectorSet,
			},
			Logger: logger,
		}

		exporter, err = exporterkit.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v\n", err))
		}
	}

	exporter.Run()
}
