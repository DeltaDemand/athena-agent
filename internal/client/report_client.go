package client

import (
	"context"
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"github.com/DeltaDemand/athena-agent/global"
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

var (
	conn       *grpc.ClientConn
	clientPool = &sync.Pool{
		New: func() interface{} {
			// 创建 gRPC 客户端
			return pb.NewReportServerClient(conn)
		},
	}
)

type RegisterError struct {
	Msg string
}

func ConnectGRPC(confs appConfigs.Config) {
	var err error
	conn, err = grpc.Dial(confs.ReportServer.GetAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("连接 gPRC 服务失败,", err)
	}
}

//注册并获取Uid,设置成全局变量
func Register() {
	//获取一个client结构体
	client := clientPool.Get().(pb.ReportServerClient)

	resp, err := client.Register(context.Background(), &pb.RegisterReq{
		Timestamp:   time.Now().Unix(),
		Metrics:     global.Metrics,
		Description: global.GetIP(),
	})
	if err != nil {
		log.Println("Register失败，或服务器返回UID失败")
	}
	log.Printf("client.Register resp{code: %d, Uid:%s, message: %s}\n", resp.Code, resp.UId, resp.Msg)
	global.SetUId(resp.UId)
	clientPool.Put(client)
}

func RequestToServer(req pb.ReportReq) *pb.ReportRsp {
	client := clientPool.Get().(pb.ReportServerClient)
	rep, err := client.Report(context.TODO(), &req)
	if err != nil {
		log.Fatal("gPRC服务发送信息失败\n", err)
	}
	clientPool.Put(client)
	return rep
}
func CloseConn() {
	clientPool.New()
	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
