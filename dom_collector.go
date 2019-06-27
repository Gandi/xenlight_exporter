package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xen-project/xen/tools/golang/xenlight"
	"gopkg.in/alecthomas/kingpin.v2"
	"strconv"
)

//version_info

var (
	domCpuTimeDesc = prometheus.NewDesc(
		"xen_domain_cpu_time_total",
		"CPU time used by the domain",
		[]string{"domain_name"}, nil,
	)
	domVcpuTimeDesc = prometheus.NewDesc(
		"xen_domain_vcpu_time_total",
		"CPU time per vCPU for the domain",
		[]string{"domain_name", "cpu_id"}, nil,
	)
	domCpuCountDesc = prometheus.NewDesc(
		"xen_domain_cpu_count",
		"Number of available CPU for domain",
		[]string{"domain_name"}, nil,
	)
	domCpuOnlineDesc = prometheus.NewDesc(
		"xen_domain_cpu_online_count",
		"Number of online CPU for domain",
		[]string{"domain_name"}, nil,
	)
	domMemoryMaxDesc = prometheus.NewDesc(
		"xen_domain_memory_max_bytes",
		"Total ammount of RAM on the domain",
		[]string{"domain_name"}, nil,
	)
	domMemoryCurrentDesc = prometheus.NewDesc(
		"xen_domain_memory_current_bytes",
		"Current ammount of RAM used by the domain",
		[]string{"domain_name"}, nil,
	)
	domMemoryOutstandingDesc = prometheus.NewDesc(
		"xen_domain_memory_outstanding_bytes",
		"Total ammount of outstanding RAM for the domain",
		[]string{"domain_name"}, nil,
	)
	domVcpuShowDetails = kingpin.Flag(
		"collector.domain.show-vcpus-details",
		"Enable the collection of per-vcpu time",
	).Default("false").Bool()
)

type DomainCollector struct{}

func init() {
	registerCollector("domain", defaultEnabled, NewDomainCollector)
}

func NewDomainCollector() prometheus.Collector {
	return &DomainCollector{}
}

func (collector DomainCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(collector, ch)
}

func (collector DomainCollector) Collect(ch chan<- prometheus.Metric) {
	xenlight.Ctx.Open()

	dominfos := xenlight.Ctx.ListDomain()
	for _, dominfo := range dominfos {
		domName := xenlight.Ctx.DomidToName(dominfo.Domid)
		ch <- prometheus.MustNewConstMetric(
			domCpuCountDesc,
			prometheus.GaugeValue,
			float64(dominfo.VcpuMaxId+1),
			domName,
		)
		ch <- prometheus.MustNewConstMetric(
			domCpuOnlineDesc,
			prometheus.GaugeValue,
			float64(dominfo.VcpuOnline),
			domName,
		)
		ch <- prometheus.MustNewConstMetric(
			domCpuTimeDesc,
			prometheus.CounterValue,
			dominfo.CpuTime.Seconds(),
			domName,
		)
		if *domVcpuShowDetails {
			vcpus := xenlight.Ctx.ListVcpu(dominfo.Domid)
			for _, vcpu := range vcpus {
				ch <- prometheus.MustNewConstMetric(
					domVcpuTimeDesc,
					prometheus.CounterValue,
					vcpu.VCpuTime.Seconds(),
					domName, strconv.FormatUint(uint64(vcpu.Vcpuid), 10),
				)
			}
		}
		ch <- prometheus.MustNewConstMetric(
			domMemoryMaxDesc,
			prometheus.GaugeValue,
			float64(dominfo.MaxMemkb)*1024,
			domName,
		)
		ch <- prometheus.MustNewConstMetric(
			domMemoryCurrentDesc,
			prometheus.GaugeValue,
			float64(dominfo.CurrentMemkb)*1024,
			domName,
		)
		ch <- prometheus.MustNewConstMetric(
			domMemoryOutstandingDesc,
			prometheus.GaugeValue,
			float64(dominfo.OutstandingMemkb)*1024,
			domName,
		)
	}
}
