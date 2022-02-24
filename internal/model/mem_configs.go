package model

//memory_used配置接收结构体
type MemConfs struct {
	Run              bool  `json:"run"`
	SamplingInterval int64 `json:"sampling_interval"`
	ReportInterval   int64 `json:"report_interval"`
}
