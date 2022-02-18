package model

//cpu_rate配置接收结构体
type CpuConfs struct {
	Run              bool  `json:"run"`
	SamplingInterval int64 `json:"samplingInterval"`
	ReportInterval   int64 `json:"reportInterval"`
}
