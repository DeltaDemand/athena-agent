package main

import (
	"fmt"
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	"github.com/DeltaDemand/athena-agent/internal/inputArgs"
	"github.com/DeltaDemand/athena-agent/internal/sampler"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

func main() {
	go exitHandle()
	//加载configs/config.json下配置
	confs := appConfigs.LoadingConfigs()
	//检测用户输入配置
	inputArgs.Parse(&confs)
	//启动初始化本机的一些全局变量
	global.InitVar()

	//指标和对应的采样器初始化，添加到global.RunMetrics中去
	confs.InitGlobalMetrics()
	//执行Etcd，内部会判断是否需要连接监听
	confs.RunEtcd(&wg)

	//连接ReportServer
	confs.ReportServer.ConnectGRPC()
	defer confs.ReportServer.CloseConn()

	//注册本机到服务器
	client.Register()

	//每个指标都启动它的采样器
	for _, metric := range global.Metrics {
		metric.(sampler.Sampler).Execute(&wg)
	}
	wg.Wait()
	global.Logger.Println("无采样器在运行，Agent退出...")
}

func exitHandle() {
	exitChan := make(chan os.Signal)
	//监听退出得信号
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case sig := <-exitChan:
			fmt.Println("接受到来自系统的信号：", sig)
			if global.EtcdOnline {
				//退出把云端配置删了
				client.DelAgent()
			}
			os.Exit(1) //如果ctrl+c 关不掉程序，使用os.Exit强行关掉
		}
	}

}
