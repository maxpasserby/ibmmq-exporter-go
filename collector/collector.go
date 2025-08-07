package collector

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

const (
	namespace = "ibmmq"
	labelName = "name"
	labelID   = "id"
)

type Collector struct {
	config     *Config
	logger     log.Logger
	testMetric *prometheus.Desc
	up         *prometheus.Desc
}

func NewCollector(logger log.Logger, c *Config) *Collector {
	return &Collector{
		config: c,
		logger: logger,
		testMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "_test_metric"),
			"这里是测试指标说明",
			[]string{labelName, labelID}, //这里是需要放的标签，送值时必须一致
			nil,
		),
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Metric indicating the status of the exporter collection. 1 indicates that the connection IBM MQ was successful, and all available metrics were collected. "+
				"0 indicates that the exporter failed to collect 1 or more metrics, due to an inability to connect to IBM MQ.",
			nil,
			nil,
		),
	}
}

func (c *Collector) Describe(descs chan<- *prometheus.Desc) {
	descs <- c.testMetric
	descs <- c.up
}
func (c *Collector) Collect(metrics chan<- prometheus.Metric) {
	level.Debug(c.logger).Log("msg", "Collecting metrics.")

	// Create a WaitGroup to block closing the database until all goroutines are done
	var wg sync.WaitGroup

	var up float64 = 1

	wg.Add(1)
	go func() {
		if err := c.collectTestMetrics(metrics); err != nil {
			level.Error(c.logger).Log("msg", "Failed to collect storage metrics.", "err", err)
			up = 0
		}
		wg.Done()
	}()

	wg.Wait()
	metrics <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, up)
	level.Debug(c.logger).Log("msg", "Finished collecting metrics.")
}

func (c *Collector) collectTestMetrics(metrics chan<- prometheus.Metric) error {
	level.Debug(c.logger).Log("msg", "测试指标.")
	//开始写业务逻辑
	metrics <- prometheus.MustNewConstMetric(c.testMetric, prometheus.GaugeValue, 0, "设置的标签1值", "设置的标签2值")

	return nil
}
