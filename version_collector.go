package main

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"xenbits.xenproject.org/git-http/xen.git/tools/golang/xenlight"
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

type VersionCollector struct {
	xenCtx *xenlight.Context
}

func NewVersionCollector(ctx *xenlight.Context) XenPromCollector {
	return &VersionCollector{
		xenCtx: ctx,
	}
}

func (VersionCollector) Name() string {
	return "version"
}

func (VersionCollector) DefaultEnabled() bool {
	return true
}

func (c *VersionCollector) PromCollector() prometheus.Collector {
	return c
}

func (c VersionCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c VersionCollector) Collect(ch chan<- prometheus.Metric) {
	versinfo, err := c.xenCtx.GetVersionInfo()
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
