package model

//disk_used配置接收结构体
type DiskConfs struct {
	Run            bool  `json:"run"`
	ReportInterval int64 `json:"reportInterval"`
}
