package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	ReconcilesApisTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "reconciles_apidefinition_total",
			Help: "Number of total API definition reconciliation attempts",
		},
	)
)

func init() {
	metrics.Registry.MustRegister(ReconcilesApisTotal)
}
