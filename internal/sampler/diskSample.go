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
	if config.Run && config.SamplingInterval > config.ReportInterval {
		global.Logger.Println("disk的采样间隔大于上报间隔，自动修改采样间隔等于上报间隔")
		config.SamplingInterval = config.ReportInterval
	}
	samplingInterval := int64(0)
	for config.Run {
		timeNow := time.Now().Unix()
		if samplingInterval == 0 {
			parts, _ := disk.Partitions(true)
			diskInfo, err := disk.Usage(parts[0].Mountpoint)
			if err == nil {
				everySamplingMem += diskInfo.UsedPercent
				samplingMemTimes++
			}
			samplingInterval = config.SamplingInterval
		}

		if timeNow%config.ReportInterval == 0 && samplingMemTimes != 0 {
			client.RequestToServer(pb.ReportReq{
				UId:        global.GetUId(),
				Timestamp:  timeNow,
				Metric:     global.DiskUsed,
				Dimensions: map[string]string{LocalIp: global.GetIP()},
				Value:      everySamplingMem / float64(samplingMemTimes),
			})
			everySamplingMem = 0.0
			samplingMemTimes = 0
		}
		samplingInterval--
		time.Sleep(time.Second)
	}
}
