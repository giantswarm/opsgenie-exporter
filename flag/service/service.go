package service

import (
	"github.com/giantswarm/opsgenie-exporter/flag/service/opsgenie"
)

type Service struct {
	Opsgenie opsgenie.Opsgenie
}
