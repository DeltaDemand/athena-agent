package sampler

import (
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/shirou/gopsutil/disk"
	"time"
)

const LOCAL_IP = "ip"

func SendDiskPercent(config appConfigs.DiskConfs) {

	for config.Run {
		timeNow := time.Now().Unix()
		parts, _ := disk.Partitions(true)
		diskInfo, err := disk.Usage(parts[0].Mountpoint)
		if err == nil {
			everySamplingMem += diskInfo.UsedPercent
			samplingMemTimes++
		}
		if timeNow%config.ReportInterval == 0 && samplingMemTimes != 0 {
			client.RequestToServer(pb.ReportReq{
				UId:        global.GetUId(),
				Timestamp:  timeNow,
				Metric:     global.DISK_RATE,
				Dimensions: map[string]string{LOCAL_IP: global.GetIP()},
				Value:      everySamplingMem / float64(samplingMemTimes),
			})
			everySamplingMem = 0.0
			samplingMemTimes = 0
		}
		time.Sleep(time.Second * time.Duration(config.SamplingInterval))
	}

}
