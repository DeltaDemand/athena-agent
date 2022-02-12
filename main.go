package main

import (
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	"github.com/DeltaDemand/athena-agent/internal/inputArgs"
	"github.com/DeltaDemand/athena-agent/internal/sampler"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {
	//加载configs/config.json下配置
	confs := appConfigs.LoadingConfigs()
	//检测用户输入配置
	inputArgs.Parse(&confs)

	//指标和对应的采样器初始化，添加到global.RunMetrics中去
	confs.InitGlobalMetrics()

	//监听Etcd，查看是否变化
	confs.StartWatchEtcd(&wg)

	//连接ReportServer
	confs.ReportServer.ConnectGRPC()
	defer confs.ReportServer.CloseConn()

	//初始化本机的一些全局变量（IP）
	global.InitVar()
	//注册本机到服务器
	client.Register()

	//每个指标都启动它的采样器
	for _, metric := range global.Metrics {
		metric.(sampler.Sampler).Execute(&wg)
	}

	wg.Wait()
}
