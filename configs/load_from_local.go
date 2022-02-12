package appConfigs

import (
	"encoding/json"
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	"github.com/DeltaDemand/athena-agent/internal/model"
	"io/ioutil"
)

type Config struct {
	AgentConfi   AgentConfs          `json:"Agent"`
	Etcd         client.Etcd         `json:"Etcd"`
	ReportServer client.ReportServer `json:"ReportServer"`
	CpuConfi     model.CpuConfs      `json:"cpu_rate"`
	MemConfi     model.MemConfs      `json:"memory_used"`
	DiskConfi    model.DiskConfs     `json:"disk_used"`
	CpuMemConfi  model.CpuMemConfs   `json:"cpu_mem"`
}

var config = Config{}
var configFile = "configs/config.json"

func LoadingConfigs() Config {
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		global.Logger.Println("配置文件读取失败\n", err)
		//读取失败，返回初始化零值的Config
		return config
	}
	configErr := json.Unmarshal(configBytes, &config)
	if configErr != nil {
		//解异常，再尝试解析etcd的配置
		err = json.Unmarshal(configBytes, &config.Etcd)
		if err != nil {
			//etcd也无解析失败
			global.Logger.Println("解析json配置文件读取失败")
		}
	}
	return config
}
