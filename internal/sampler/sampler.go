package sampler

import (
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"github.com/DeltaDemand/athena-agent/global"
	"sync"
)

type Sampler interface {
	Execute(wg *sync.WaitGroup)
	GetName() string
}

var (
	cpuSampler  = cpuSample{}
	memSampler  = memSample{}
	diskSampler = diskSample{}
)

func Init(config appConfigs.Config) {
	if config.CpuConfi.Run {
		cpuSampler.name = global.CpuRate
		cpuSampler.config = config.CpuConfi
		global.RunMetrics = append(global.RunMetrics, &cpuSampler)
	}
	if config.MemConfi.Run {
		memSampler.name = global.MemUsed
		memSampler.config = config.MemConfi
		global.RunMetrics = append(global.RunMetrics, &memSampler)
	}
	if config.DiskConfi.Run {
		diskSampler.name = global.DiskUsed
		diskSampler.config = config.DiskConfi
		global.RunMetrics = append(global.RunMetrics, &diskSampler)
	}

	setGlobalRunMetrics()
}
func setGlobalRunMetrics() {
	global.RunMetricsNum = len(global.RunMetrics)
	global.RunMetricsName = make([]string, global.RunMetricsNum)
	for i, metric := range global.RunMetrics {
		global.RunMetricsName[i] = metric.(Sampler).GetName()
	}
}