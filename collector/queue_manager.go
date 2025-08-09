package collector

const (
	SERVER_CPU_USAGE  = "server_cpu_usage"
	SERVER_MEM_USAGE  = "server_mem_usage"
	SERVER_DISK_USAGE = "server_disk_usage"
)

func RegisterQMgrMetricSpec(metricSpecs map[string]MetricSpec) {
	metricSpecs[SERVER_CPU_USAGE] = NewMetricSpec(SERVER_CPU_USAGE, Counter, "当前服务器CPU使用率", nil, nil, []string{"server_id"})
	metricSpecs[SERVER_MEM_USAGE] = NewMetricSpec(SERVER_MEM_USAGE, Gauge, "当前服务器内存使用率", nil, nil, []string{"serverbeijing"})
	metricSpecs[SERVER_DISK_USAGE] = NewMetricSpec(SERVER_DISK_USAGE, Gauge, "当前服务器磁盘使用率", nil, nil, []string{"serverzone"})
}
