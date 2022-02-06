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

	confs := appConfigs.LoadingConfigs() //加载configs/config.json下配置
	inputArgs.Parse(&confs)              //检测用户输入配置

	sampler.Init(confs) //采样器初始化，需要运行的添加到global.RunMetrics中去

	client.ConnectGRPC(confs) //连接服务器
	defer client.CloseConn()
	global.InitVar()  //初始化本机的全局变量（IP）
	client.Register() //注册本机到服务器

	wg.Add(global.RunMetricsNum)
	for _, metric := range global.RunMetrics {
		go metric.(sampler.Sampler).Execute(&wg)
	}
	wg.Wait()
}
