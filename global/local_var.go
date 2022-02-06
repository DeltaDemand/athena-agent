package global

import (
	"net"
	"strings"
)

var (
	ip  = "0.0.0.0"
	uId string
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
		Logger.Println(err.Error())
		return err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	ip = localAddr[0:idx]
	return nil
}
