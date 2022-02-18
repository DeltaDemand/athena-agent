package model

//cpu_mem配置接收结构体
type CpuMemConfs struct {
	Run                 bool    `json:"run"`
	SamplingInterval    int64   `json:"samplingInterval"`
	ReportInterval      int64   `json:"reportInterval"`
	CpuPercentThreshold float64 `json:"cpuPercentThreshold"`
	MemPercentThreshold float64 `json:"memPercentThreshold"`
}