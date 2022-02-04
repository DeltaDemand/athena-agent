package appConfigs

type CpuConfs struct {
	Run              bool  `json:"run"`
	SamplingInterval int64 `json:"samplingInterval"`
	ReportInterval   int64 `json:"reportInterval"`
}
