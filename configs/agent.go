package appConfigs

import (
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/sampler"
	"sync"
)

//本agent的一些参数
type AgentConfs struct {
	CheckAlive int  `json:"checkAlive"` //给ReportServer的检查本agent存活的时间阈值
	Pause      bool `json:"pause"`
	Exit       bool `json:"exit"`
}

//agent配置改变执行函数
func (agent *AgentConfs) Execute(wg *sync.WaitGroup) error {
	if agent.Exit {
		//检测到退出为真，直接Fatal退出
		global.Logger.Fatal("Agent退出...")
	}
	global.SetPause(agent.Pause)
	//如果读到非暂停，将暂停的采样器执行起来
	if !global.GetPause() {
		for _, metric := range global.Metrics {
			metric.(sampler.Sampler).Execute(wg)
		}
	} else {
		//agent暂停
		global.Logger.Println("Agent暂停...")
	}
	return nil
}
