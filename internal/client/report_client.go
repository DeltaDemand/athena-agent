//实现pb协议的客户端函数

package client

import (
	"context"
	"github.com/DeltaDemand/athena-agent/global"
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"sync"
	"time"
)

var (
	// 保证正在调用client时,不会突然改变连接参数
	// 保证改变参数重连时，不会有goroutine调用client
	connSafe sync.RWMutex
)

const (
	notFoundCode = 10000001
	sqlErr       = 10000002
)

// Register 注册到ReportServer并获取Uid,设置成全局变量
func Register() error {
	//加读锁，防止获取client对象时被重新连接（ConnectGRPC函数）修改clientPool
	connSafe.RLock()
	//获取一个pb.ReportServerClient结构体
	client := clientPool.Get().(pb.ReportServerClient)
	defer clientPool.Put(client)
	//调用pb的注册函数
	resp, err := client.Register(context.Background(), &pb.RegisterReq{
		Timestamp:       time.Now().Unix(),
		Metrics:         global.MetricsName,
		CheckAliveTime:  int32(global.CheckAlive),
		AggregationTime: int32(global.AggregationTime),
		Description:     "agent group: " + global.AgentGroup + "; name: " + global.AgentName + "; ip: " + global.GetIP(),
	})
	connSafe.RUnlock()
	if err != nil {
		global.Logger.Printf("Register失败，服务器返回UID失败,Agent暂停...\n", err)
		//注册失败，注册状态设为失败
		global.SetRegisterSuccess(false)
		//注册失败，先暂停Agent
		global.SetPause(true)
		//把Agent状态更新configServer上的状态
		RefreshAgentState(true)
		return err
	}
	//注册失败，注册状态设为成功
	global.SetRegisterSuccess(true)
	global.Logger.Printf("client.Register resp{code: %d, Uid:%s, message: %s}\n", resp.Code, resp.UId, resp.Msg)
	global.SetUId(resp.UId)
	return nil
}

// RequestToServer 封装处理pb的Report函数
func RequestToServer(req pb.ReportReq) (*pb.ReportRsp, error) {
	//加读锁，防止获取client对象时被重新连接（ConnectGRPC函数）修改clientPool
	connSafe.RLock()
	client := clientPool.Get().(pb.ReportServerClient)
	defer clientPool.Put(client)
	//调用pb的报告函数
	rep, err := client.Report(context.Background(), &req)
	connSafe.RUnlock()
	if err != nil {
		global.Logger.Printf("gPRC服务发送信息失败\n", err)
		//处理发送直接返回
		return nil, err
	}
	global.Logger.Printf("client.Request resp{code: %d, message: %s}\n", rep.Code, rep.Msg)
	//ReportServer找不到本机uid，重新注册
	if rep.Code == notFoundCode {
		//再次注册
		Register()
	} else if rep.Code == sqlErr {
		//返回数据库错误，暂停Agent
		global.SetPause(true)
		//把Agent状态更新configServer上的状态
		RefreshAgentState(true)
	}
	return rep, nil
}
