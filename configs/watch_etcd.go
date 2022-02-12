package appConfigs

import (
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/sampler"
	"sync"
)

func (c *Config) StartWatchEtcd(wg *sync.WaitGroup) {
	//连接etcd
	err := c.Etcd.Connect()
	//该agent使用etcd，并连接不报错
	if c.Etcd.Apply && err == nil {
		//每个配置项监控都用一个goroutine
		go func() {
			//监听Agent
			c.Etcd.WatchConfig(global.Agent, &c.AgentConfi, &c.AgentConfi, wg)
		}()
		go func() {
			//监听ReportServer
			c.Etcd.WatchConfig(global.ReportServer, &c.ReportServer, &c.ReportServer, wg)
			wg.Done()
		}()
		go func() {
			//监听Etcd
			c.Etcd.WatchConfig(global.Etcd, &c.Etcd, &c.Etcd, wg)
			wg.Done()
		}()
		//监听全局的指标
		for _, metric := range global.Metrics {
			//注意不能将sam值放到goroutine中解析，否则容易只取到最后一个值
			sam := metric.(sampler.Sampler)
			go func() {
				c.Etcd.WatchConfig(sam.GetMetricName(), sam.GetConfigPtr(), sam, wg)
				wg.Done()
			}()
		}
	}

}
