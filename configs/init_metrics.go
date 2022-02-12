package appConfigs

import (
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/sampler"
)

var (
	//存入全局变量Metrics map[string]interface{}的元素
	cpuSampler    = sampler.CpuSample{}
	memSampler    = sampler.MemSample{}
	diskSampler   = sampler.DiskSample{}
	cpuMemSampler = sampler.CpuMemSample{}
)

//初始化全局变量Metrics和MetricsName
func (c *Config) InitGlobalMetrics() {
	c.setGlobalMetrics()
	setGlobalMetricsName()
}

//将配置中的参数放入到每个sampler（采样器）中,初始化全局变量Metrics
func (config *Config) setGlobalMetrics() {

	cpuSampler.SetName(global.CpuRate)
	cpuSampler.Config = config.CpuConfi
	global.Metrics[global.CpuRate] = &cpuSampler

	memSampler.SetName(global.MemUsed)
	memSampler.Config = config.MemConfi
	global.Metrics[global.MemUsed] = &memSampler

	diskSampler.SetName(global.DiskUsed)
	diskSampler.Config = config.DiskConfi
	global.Metrics[global.DiskUsed] = &diskSampler

	cpuMemSampler.SetName(global.CpuMem)
	cpuMemSampler.Config = config.CpuMemConfi
	global.Metrics[global.CpuMem] = &cpuMemSampler
}

//初始化MetricsName(指标的名字数组，用于发送到ReportServer)
func setGlobalMetricsName() {
	global.MetricsName = make([]string, 0, global.MetricsNum)
	for _, metric := range global.Metrics {
		global.MetricsName = append(global.MetricsName, metric.(sampler.Sampler).GetMetricName())
	}
}
