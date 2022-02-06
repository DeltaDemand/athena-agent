package appConfigs

import (
	"encoding/json"
	"github.com/DeltaDemand/athena-agent/global"
	"io/ioutil"
)

type Config struct {
	ReportServer ReportServer `json:"ReportServer"`
	CpuConfi     CpuConfs     `json:"CpuSampling"`
	MemConfi     MemConfs     `json:"MemSampling"`
	DiskConfi    DiskConfs    `json:"DiskSampling"`
}

var config = Config{}
var configFile = "configs/config.json"

func LoadingConfigs() Config {
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		global.Logger.Fatal("配置文件读取失败", err)
		return config
	}
	configErr := json.Unmarshal(configBytes, &config)
	if configErr != nil {
		global.Logger.Fatal("解析json配置文件读取失败")
	}
	return config
}
