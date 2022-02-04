package main

import (
	"github.com/DeltaDemand/athena-agent/Global"
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"github.com/DeltaDemand/athena-agent/internal/client"
	"github.com/DeltaDemand/athena-agent/internal/inputArgs"
	"github.com/DeltaDemand/athena-agent/internal/sampler"
	"sync"
)

func main() {

	confs := appConfigs.LoadingConfigs() //加载configs/config.json下配置
	inputArgs.Parse(&confs)              //检测用户输入配置

	client.ConnectGRPC(confs) //连接服务器
	defer client.CloseConn()
	global.InitVar(confs) //初始化本机的全局变量（IP）
	client.Register()     //注册本机到服务器

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		//开启cpu发送
		sampler.SendCpuPercent(confs.CpuConfi)
		wg.Done()
	}()
	go func() {
		//开启mem发送
		sampler.SendMemPercent(confs.MemConfi)
		wg.Done()
	}()
	go func() {
		//开disk发送
		sampler.SendDiskPercent(confs.DiskConfi)
		wg.Done()
	}()
	wg.Wait()
}
