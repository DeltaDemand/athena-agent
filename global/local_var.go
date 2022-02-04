package global

import (
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"log"
	"net"
	"strings"
)

var (
	ip      = "0.0.0.0"
	uId     string
	Metrics = make([]string, 0, 3)
)

const (
	CPU_RATE  = "cpu_rate"
	MEM_RATE  = "memory_used"
	DISK_RATE = "disk_used"
)

func SetUId(id string) {
	uId = id
}
func GetIP() string {
	return ip
}
func GetUId() string {
	return uId
}

func initIP() (err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	ip = localAddr[0:idx]
	return nil
}

func InitVar(config appConfigs.Config) {
	err := initIP()
	if err != nil {
		log.Println("获取本机Addr失败")
	}
	if config.MemConfi.Run {
		Metrics = append(Metrics, CPU_RATE)
	}
	if config.MemConfi.Run {
		Metrics = append(Metrics, MEM_RATE)
	}
	if config.MemConfi.Run {
		Metrics = append(Metrics, DISK_RATE)
	}
}
