package global

import (
	"log"
	"os"
)

const (
	MetricsNum = 3 //总指标数
	CpuRate    = "cpu_rate"
	MemUsed    = "memory_used"
	DiskUsed   = "disk_used"
)

var (
	RunMetrics     = make([]interface{}, 0, MetricsNum)
	RunMetricsName []string
	RunMetricsNum  int
	Logger         = log.New(os.Stdout, "<Agent>", log.Lshortfile|log.Ldate|log.Ltime)
)

func InitVar() {
	err := initIP()
	if err != nil {
		Logger.Println("获取本机IP失败")
	}
}
