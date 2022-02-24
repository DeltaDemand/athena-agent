package model

//cpu_mem配置接收结构体
type CpuMemConfs struct {
	Run                 bool    `json:"run"`
	SamplingInterval    int64   `json:"sampling_interval"`
	ReportInterval      int64   `json:"report_interval"`
	CpuPercentThreshold float64 `json:"cpu_percent_threshold"`
	MemPercentThreshold float64 `json:"mem_percent_threshold"`
}
