package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"xenbits.xen.org/git-http/xen.git/tools/golang/xenlight"
	"strconv"
)

var (
	versionInfoDesc = prometheus.NewDesc(
		"xen_version_info",
		"Version information of the host",
		[]string{
			"major",
			"minor",
			"extra",
			"changeset",
			"compile_date",
			"build_id",
			"commandline",
		}, nil,
	)
)

type VersionCollector struct{}

func init() {
	registerCollector("version", defaultEnabled, NewVersionCollector)
}

func NewVersionCollector() prometheus.Collector {
	return &VersionCollector{}
}

func (collector VersionCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(collector, ch)
}

func (collector VersionCollector) Collect(ch chan<- prometheus.Metric) {
	xenlight.Ctx.Open()
	versinfo, err := xenlight.Ctx.GetVersionInfo()
	if err != nil {
		return
	}
	ch <- prometheus.MustNewConstMetric(
		versionInfoDesc,
		prometheus.GaugeValue,
		float64(1),
		strconv.Itoa(versinfo.XenVersionMajor),
		strconv.Itoa(versinfo.XenVersionMinor),
		versinfo.XenVersionExtra,
		versinfo.Changeset,
		versinfo.CompileDate,
		versinfo.BuildId,
		versinfo.Commandline,
	)
}
