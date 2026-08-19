package monitor

import (
	"errors"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var metricsRegistryMu sync.Mutex

func registerMetric(collector prometheus.Collector) prometheus.Collector {
	metricsRegistryMu.Lock()
	defer metricsRegistryMu.Unlock()

	if err := prometheus.Register(collector); err != nil {
		var alreadyRegistered prometheus.AlreadyRegisteredError
		if errors.As(err, &alreadyRegistered) {
			return alreadyRegistered.ExistingCollector
		}
		panic(err)
	}
	return collector
}
