package collector

import (
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	config *Config
	logger log.Logger

	metricSpecs map[string]MetricSpec
}

func NewCollector(logger log.Logger, config *Config) *Collector {

	c := &Collector{
		config: config,
		logger: logger,

		metricSpecs: make(map[string]MetricSpec),
	}

	RegisterQMgrMetricSpec(c.metricSpecs)

	return c
}

// Describe Redis metric descriptions.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, spec := range c.metricSpecs {
		ch <- spec.Desc
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	level.Debug(c.logger).Log("msg", "Collecting metrics.")

	// Create a WaitGroup to block closing the database until all goroutines are done
	var wg sync.WaitGroup

	// var up float64 = 1

	wg.Add(1)
	go func() {
		c.registerMetric(ch, c.mustFindMetricSpec(SERVER_CPU_USAGE), 0, "server_id")

		wg.Done()
	}()

	wg.Wait()
	// ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, up)
	level.Debug(c.logger).Log("msg", "Finished collecting metrics.")
}
