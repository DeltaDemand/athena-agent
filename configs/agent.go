package appConfigs

import (
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	"github.com/DeltaDemand/athena-agent/internal/sampler"
	"sync"
)

//本agent的一些参数
type AgentConfs struct {
	Pause bool `json:"pause"`
	Exit  bool `json:"exit"`
}

//agent配置改变执行函数
func (agent *AgentConfs) Execute(wg *sync.WaitGroup) error {
	if agent.Exit {
		client.DelAgent()
		//检测到退出为真，直接Fatal退出
		global.Logger.Fatal("Agent退出...")
	}
	global.SetPause(agent.Pause)
	//如果读到非暂停，将暂停的采样器执行起来
	if !global.GetPause() {
		//如果注册成功
		if global.GetRegisterSuccess() {
			for _, metric := range global.Metrics {
				metric.(sampler.Sampler).Execute(wg)
			}
		} else {
			//将云端状态重设为false
			client.RefreshAgentState(true)
			global.Logger.Println("Agent未注册成功，请检查ReportServer设置..")
		}
	} else {
		//agent暂停
		global.Logger.Println("配置服务器修改Agent为暂停状态")
	}
	return nil
}
