package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusType string

const (
	namespace = "ibmmq"

	Gauge     PrometheusType = "gauge"
	Counter   PrometheusType = "counter"
	Summary   PrometheusType = "summary"
	Histogram PrometheusType = "histogram"
)

type MetricSpec struct {
	Name      string
	Type      PrometheusType
	Quantiles map[float64]float64
	Buckets   map[float64]uint64
	Desc      *prometheus.Desc
}

func NewMetricSpec(name string, prometheusType PrometheusType, help string, quantiles map[float64]float64, buckets map[float64]uint64, labels []string) MetricSpec {
	return MetricSpec{
		Name:      name,
		Type:      prometheusType,
		Quantiles: quantiles,
		Buckets:   buckets,
		Desc:      prometheus.NewDesc(prometheus.BuildFQName(namespace, "", name), help, labels, nil),
	}
}

func (c *Collector) registerMetric(ch chan<- prometheus.Metric, spec MetricSpec, value float64, labelValues ...string) {
	switch spec.Type {
	case Gauge:
		ch <- c.registerConstMetric(spec, prometheus.GaugeValue, value, labelValues...)
	case Counter:
		ch <- c.registerConstMetric(spec, prometheus.CounterValue, value, labelValues...)
	case Summary:
		ch <- c.registerConstSummary(spec, 0, 0, spec.Quantiles, labelValues...)
	case Histogram:
		ch <- c.registerConstHistogram(spec, 0, 0, spec.Buckets, labelValues...)
	}
}
func (c *Collector) registerConstSummary(spec MetricSpec, count uint64, sum float64, quantiles map[float64]float64, labelValues ...string) prometheus.Metric {
	return prometheus.MustNewConstSummary(spec.Desc, count, sum, quantiles, labelValues...)
}

func (c *Collector) registerConstHistogram(spec MetricSpec, count uint64, sum float64, buckets map[float64]uint64, labelValues ...string) prometheus.Metric {
	return prometheus.MustNewConstHistogram(spec.Desc, count, sum, buckets, labelValues...)
}

func (c *Collector) registerConstMetric(spec MetricSpec, valueType prometheus.ValueType, value float64, labelValues ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(spec.Desc, valueType, value, labelValues...)
}

func (c *Collector) mustFindMetricSpec(name string) MetricSpec {
	spec, found := c.metricSpecs[name]
	if !found {
		panic(fmt.Sprintf("couldn't find metric description for %s", name))
	}
	return spec
}
