package collector

import (
	"github.com/giantswarm/exporterkit/collector"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/viper"

	"github.com/giantswarm/opsgenie-exporter/flag"
	"github.com/giantswarm/opsgenie-exporter/service/collector/alert"
	"github.com/giantswarm/opsgenie-exporter/service/collector/opsgenie"
)

type SetConfig struct {
	Flag   *flag.Flag
	Logger micrologger.Logger
	Viper  *viper.Viper
}

// Set is basically only a wrapper for the operator's collector implementations.
// It eases the iniitialization and prevents some weird import mess so we do not
// have to alias packages.
type Set struct {
	*collector.Set
}

func NewSet(config SetConfig) (*Set, error) {
	var err error

	var opsgenieClient *opsgenie.Client
	{
		c := opsgenie.Config{
			Key: config.Viper.GetString(config.Flag.Service.Opsgenie.API.Token),
		}

		opsgenieClient, err = opsgenie.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var alertCollector collector.Interface
	{
		c := alert.Config{
			Client: opsgenieClient,
		}

		alertCollector, err = alert.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var collectorSet *collector.Set
	{
		c := collector.SetConfig{
			Collectors: []collector.Interface{
				alertCollector,
			},
			Logger: config.Logger,
		}

		collectorSet, err = collector.NewSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Set{
		Set: collectorSet,
	}

	return s, nil
}
