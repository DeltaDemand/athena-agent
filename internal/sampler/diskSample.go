package sampler

import (
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/shirou/gopsutil/disk"
	"sync"
	"time"
)

const LocalIp = "local_ip"

type diskSample struct {
	name   string
	config appConfigs.DiskConfs
}

func (s *diskSample) GetName() string {
	return s.name
}
func (s *diskSample) Execute(wg *sync.WaitGroup) {
	sendDiskPercent(s.config)
	wg.Done()
}

func sendDiskPercent(config appConfigs.DiskConfs) {

	for config.Run {
		timeNow := time.Now().Unix()

		if timeNow%config.ReportInterval == 0 {
			parts, _ := disk.Partitions(true)
			diskInfo, err := disk.Usage(parts[0].Mountpoint)
			if err == nil {
				client.RequestToServer(pb.ReportReq{
					UId:        global.GetUId(),
					Timestamp:  timeNow,
					Metric:     global.DiskUsed,
					Dimensions: map[string]string{LocalIp: global.GetIP()},
					Value:      diskInfo.UsedPercent,
				})
			} else {
				global.Logger.Println("获取硬盘数据失败")
			}
		}
		time.Sleep(time.Second)
	}
}
