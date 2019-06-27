package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
)

const (
	defaultEnabled  = true
	defaultDisabled = false
)

type registeredCollector struct {
	enabled *bool
	factory func() prometheus.Collector
}

var availableCollectors = make([]registeredCollector, 0)

func registerCollector(collector string, enabled bool, factory func() prometheus.Collector) {
	var helpDefaultState string
	if enabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}

	flagName := fmt.Sprintf("collector.%s", collector)
	flagHelp := fmt.Sprintf("Enable the %s collector (default: %s).", collector, helpDefaultState)
	defaultValue := fmt.Sprintf("%v", enabled)

	flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Bool()
	availableCollectors = append(availableCollectors, registeredCollector{
		enabled: flag,
		factory: factory,
	})
}

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9603").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("xenlight_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)
	for _, col := range availableCollectors {
		if *col.enabled {
			reg.MustRegister(col.factory())
		}
	}

	http.Handle(*metricsPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
			<head><title>ctld Exporter</title></head>
			<body>
			<h1>ctld Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Errorln(err)
		}
	})

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
