package global

//不变的三个配置项
const (
	Agent        = "Agent"
	ReportServer = "ReportServer"
	Etcd         = "Etcd"
)
const (
	MetricsNum = 4 //总指标数，只是用于初始化空间
	CpuRate    = "cpu_rate"
	MemUsed    = "memory_used"
	DiskUsed   = "disk_used"
	CpuMem     = "cpu_mem"
)

var (
	Metrics     = make(map[string]interface{}, MetricsNum) //K:指标名称 V:采样器（Sampler）
	MetricsName []string
)
