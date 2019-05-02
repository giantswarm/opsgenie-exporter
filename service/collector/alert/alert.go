package alert

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/opsgenie-exporter/service/collector/opsgenie"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

const (
	namespace = "opsgenie"
	subsystem = "alert"

	labelStatus = "status"
)

var (
	alertCount *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "count"),
		"Count of OpsGenie alerts.",
		[]string{
			labelStatus,
		},
		nil,
	)
)

type Config struct {
	Client *opsgenie.Client
}

type Alert struct {
	client *opsgenie.Client
}

func New(config Config) (*Alert, error) {
	if config.Client == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Client must not be empty", config)
	}

	a := &Alert{
		client: config.Client,
	}

	return a, nil
}

func (a *Alert) Collect(ch chan<- prometheus.Metric) error {
	var g errgroup.Group

	g.Go(func() error {
		numAlerts, err := a.client.CountAlerts()
		if err != nil {
			return microerror.Mask(err)
		}

		ch <- prometheus.MustNewConstMetric(
			alertCount,
			prometheus.GaugeValue,
			float64(numAlerts),
			"",
		)

		return nil
	})

	g.Go(func() error {
		numOpenAlerts, err := a.client.CountOpenAlerts()
		if err != nil {
			return microerror.Mask(err)
		}

		ch <- prometheus.MustNewConstMetric(
			alertCount,
			prometheus.GaugeValue,
			float64(numOpenAlerts),
			"open",
		)

		return nil
	})

	g.Go(func() error {
		numClosedAlerts, err := a.client.CountClosedAlerts()
		if err != nil {
			return microerror.Mask(err)
		}

		ch <- prometheus.MustNewConstMetric(
			alertCount,
			prometheus.GaugeValue,
			float64(numClosedAlerts),
			"closed",
		)

		return nil
	})

	if err := g.Wait(); err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (a *Alert) Describe(ch chan<- *prometheus.Desc) error {
	ch <- alertCount

	return nil
}
