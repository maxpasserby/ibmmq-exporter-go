package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"gopkg.in/alecthomas/kingpin.v2"
	"ibmmq-exporter-go/collector"
)

var (
	webConfig  = webflag.AddFlags(kingpin.CommandLine, ":9975")
	metricPath = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").Envar("IBM_MQ_EXPORTER_WEB_TELEMETRY_PATH").String()
	username   = kingpin.Flag("username", "The username for the user used when querying metrics.").Envar("IBM_USERNAME").Required().String()
	password   = kingpin.Flag("password", "The password for the user used when querying metrics.").Envar("IBM_PASSWORD").String()
)

const (
	// The name of the exporter.
	exporterName    = "ibmmq_exporter"
	landingPageHtml = `<html>
<head><title>ibmmq exporter</title></head>
	<body>
		<h1>ibmmq exporter</h1>
		<p><a href='%s'>Metrics</a></p>
	</body>
</html>`
)

func main() {
	kingpin.Version(version.Print(exporterName))

	promlogConfig := &promlog.Config{}

	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(promlogConfig)

	// Construct the collector, using the flags for configuration
	c := &collector.Config{
		Username: *username,
		Password: *password,
	}

	if err := c.Validate(); err != nil {
		level.Error(logger).Log("msg", "Configuration is invalid.", "err", err)
		os.Exit(1)
	}

	col := collector.NewCollector(logger, c)

	// Register collector with prometheus client library
	prometheus.MustRegister(version.NewCollector(exporterName))
	prometheus.MustRegister(col)

	serveMetrics(logger)
}

func serveMetrics(logger log.Logger) {
	landingPage := []byte(fmt.Sprintf(landingPageHtml, *metricPath))

	http.Handle(*metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8") // nolint: errcheck
		w.Write(landingPage)                                       // nolint: errcheck
	})

	srv := &http.Server{}
	if err := web.ListenAndServe(srv, webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error running HTTP server", "err", err)
		os.Exit(1)
	}
}
