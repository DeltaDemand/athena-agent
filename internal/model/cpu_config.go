package model

//cpu_rate配置接收结构体
type CpuConfs struct {
	Run            bool  `json:"run"`
	ReportInterval int64 `json:"report_interval"`
}
