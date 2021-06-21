package collect

import (
	"github.com/prometheus/client_golang/prometheus"
	"xenbits.xenproject.org/git-http/xen.git/tools/golang/xenlight"
)

var (
	physTopologyNodesDesc = prometheus.NewDesc(
		"xen_physical_topology_nodes_number",
		"Number of socket on the host",
		nil, nil,
	)
	physTopologyCoresDesc = prometheus.NewDesc(
		"xen_physical_topology_cores_per_socket",
		"Number of cores per socket on the host",
		nil, nil,
	)
	physTopologyThreadsDesc = prometheus.NewDesc(
		"xen_physical_topology_threads_per_core",
		"Number of threads per core on the host",
		nil, nil,
	)
	physMemoryTotalDesc = prometheus.NewDesc(
		"xen_physical_memory_total_bytes",
		"Total ammount of RAM on the host",
		nil, nil,
	)
	physMemoryFreeDesc = prometheus.NewDesc(
		"xen_physical_memory_free_bytes",
		"Total ammount of free RAM on the host",
		nil, nil,
	)
	physMemoryScrubDesc = prometheus.NewDesc(
		"xen_physical_memory_scrub_bytes",
		"Total ammount of scrub RAM on the host",
		nil, nil,
	)
	physMemoryOutstandingDesc = prometheus.NewDesc(
		"xen_physical_memory_outstanding_bytes",
		"Total ammount of outstanding RAM on the host",
		nil, nil,
	)
)

type PhysicalCollector struct {
	xenCtx *xenlight.Context
}

func NewPhysicalCollector(ctx *xenlight.Context) XenPromCollector {
	return &PhysicalCollector{
		xenCtx: ctx,
	}
}

func (PhysicalCollector) Name() string {
	return "physical"
}

func (PhysicalCollector) DefaultEnabled() bool {
	return true
}

func (c *PhysicalCollector) PromCollector() prometheus.Collector {
	return c
}

func (c PhysicalCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c PhysicalCollector) Collect(ch chan<- prometheus.Metric) {
	physinfo, err := c.xenCtx.GetPhysinfo()
	if err != nil {
		return
	}
	versinfo, err := c.xenCtx.GetVersionInfo()
	if err != nil {
		return
	}
	pageSize := versinfo.Pagesize
	ch <- prometheus.MustNewConstMetric(
		physTopologyNodesDesc,
		prometheus.GaugeValue,
		float64(physinfo.NrNodes),
	)
	ch <- prometheus.MustNewConstMetric(
		physTopologyCoresDesc,
		prometheus.GaugeValue,
		float64(physinfo.CoresPerSocket),
	)
	ch <- prometheus.MustNewConstMetric(
		physTopologyThreadsDesc,
		prometheus.GaugeValue,
		float64(physinfo.ThreadsPerCore),
	)
	ch <- prometheus.MustNewConstMetric(
		physMemoryTotalDesc,
		prometheus.GaugeValue,
		float64(physinfo.TotalPages*uint64(pageSize)),
	)
	ch <- prometheus.MustNewConstMetric(
		physMemoryFreeDesc,
		prometheus.GaugeValue,
		float64(physinfo.FreePages*uint64(pageSize)),
	)
	ch <- prometheus.MustNewConstMetric(
		physMemoryScrubDesc,
		prometheus.GaugeValue,
		float64(physinfo.ScrubPages*uint64(pageSize)),
	)
	ch <- prometheus.MustNewConstMetric(
		physMemoryOutstandingDesc,
		prometheus.GaugeValue,
		float64(physinfo.OutstandingPages*uint64(pageSize)),
	)
}
