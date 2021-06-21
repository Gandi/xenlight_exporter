package main

import (
	"fmt"
	"net/http"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"xenbits.xenproject.org/git-http/xen.git/tools/golang/xenlight"
)

func init() {
	ctx, err := xenlight.NewContext()
	if err != nil {
		panic(err.Error())
	}

	registerCollector(NewDomainCollector(ctx))
	registerCollector(NewPhysicalCollector(ctx))
	registerCollector(NewVersionCollector(ctx))
}

type registeredCollector struct {
	enabled *bool
	factory func() prometheus.Collector
}

var availableCollectors = make([]registeredCollector, 0)

func registerCollector(collector XenPromCollector) {
	var helpDefaultState string
	if collector.DefaultEnabled() {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}

	flagName := fmt.Sprintf("collector.%s", collector)
	flagHelp := fmt.Sprintf("Enable the %s collector (default: %s).", collector, helpDefaultState)
	defaultValue := fmt.Sprintf("%t", collector.DefaultEnabled())

	flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Bool()
	availableCollectors = append(availableCollectors, registeredCollector{
		enabled: flag,
		factory: collector.PromCollector,
	})
}

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9603").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	)

	logCfg := new(promlog.Config)
	flag.AddFlags(kingpin.CommandLine, logCfg)

	kingpin.Version(version.Print("xenlight_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(logCfg)

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
			level.Error(logger).Log(err)
		}
	})

	level.Info(logger).Log("Listening on", *listenAddress)
	level.Error(logger).Log(http.ListenAndServe(*listenAddress, nil))

}
