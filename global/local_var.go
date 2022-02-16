package global

import (
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	ip         string
	uId        string //本Agent唯一id
	CheckAlive int    //本Agent可能最长多长时间不上报
	pause      bool   //本Agent暂停状态

	exit         bool   //退出，仅可通过etcd设置
	ConfigServer string //本Agent需连接的etcd上的修改配置的服务
	AgentGroup   string //本Agent所属etcd上的群组
	AgentName    string ////本Agent在etcd上的名称

	EtcdChange          = make(chan bool, 1) //etcd有变，开始监听etcd
	HandleChangeSuccess bool                 //etcd变化处理成功标识
	AggregationTime     int64                //上报几次进行聚合
	agentStatLock       sync.RWMutex         //只允许一个goroutine修改Agent状态(暂停)
	Split               = "(]"               //用于监控etcd连接 AgentGroup、AgentName、配置项的字符串

	Logger = log.New(os.Stdout, "<Agent>", log.Lshortfile|log.Ldate|log.Ltime)
)

// SetPause 设置Agent暂停状态，可能会有多个goroutine同时访问，加锁
func SetPause(p bool) {
	agentStatLock.Lock()
	pause = p
	agentStatLock.Unlock()
}

// GetPause 读取Agent暂停状态，可能会有多个goroutine同时访问，加锁
func GetPause() bool {
	agentStatLock.RLock()
	defer agentStatLock.RUnlock()
	return pause
}

//只在Register中使用
func SetUId(id string) {
	uId = id
}
func GetIP() string {
	return ip
}
func GetUId() string {
	return uId
}

//初始化本机ip
func initIP() (err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		Logger.Println("获取本机IP失败:", err.Error())
		return err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	ip = localAddr[0:idx]
	return nil
}

func InitVar() {
	initIP()
	EtcdChange <- true
}
