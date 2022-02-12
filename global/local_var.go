package global

import (
	"net"
	"strings"
	"sync"
)

var (
	ip              string
	uId             string
	CheckAlive      int
	exit            bool
	pause           bool
	AgentGroup      string
	AgentName       string
	AggregationTime int64
	agentStatLock   sync.RWMutex //只允许一个goroutine修改Agent状态(暂停)
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
}
