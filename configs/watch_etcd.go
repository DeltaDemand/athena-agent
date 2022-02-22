package appConfigs

import (
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/sampler"
	"sync"
)

func (c *Config) RunEtcd(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		//监听Etcd，配置变化会值，刚启动也属于变化
		for range global.EtcdChange {
			c.StartWatchEtcd(wg)
			global.HandleChangeSuccess = true
			//不使用etcd直接打破循环，释放goroutine
			if c.Etcd.Apply == false {
				global.EtcdOnline = false
				break
			}
		}
		wg.Done()
	}()
}

func (c *Config) StartWatchEtcd(wg *sync.WaitGroup) {
	//该agent配置值使用etcd
	if c.Etcd.Apply {
		//初始化连接etcd的参数
		_ = c.Etcd.Connect()
		//检查etcd上是否正常连接且是否存在同名Agent
		c.Etcd.CheckConfigServer()
		global.Logger.Println("已连接（", global.ConfigServer, "),开始监听配置<", global.AgentGroup, "|", global.AgentName, ">")
		//每个配置项监内部实现都用一个goroutine,所以要传wg
		//监听Agent
		c.Etcd.WatchConfig(global.Agent, &c.AgentConfi, &c.AgentConfi, wg)

		//监听ReportServer
		c.Etcd.WatchConfig(global.ReportServer, &c.ReportServer, &c.ReportServer, wg)

		//监听Etcd
		c.Etcd.WatchConfig(global.Etcd, &c.Etcd, &c.Etcd, wg)

		//监听全局的指标
		for _, metric := range global.Metrics {
			//每个指标的采样器
			sam := metric.(sampler.Sampler)
			c.Etcd.WatchConfig(sam.GetMetricName(), sam.GetConfigPtr(), sam, wg)
		}
	}

}
