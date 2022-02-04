package sampler

import (
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

var (
	everySamplingCpu = 0.0
	samplingCpuTimes = 0
)

func SendCpuPercent(config appConfigs.CpuConfs) {
	for config.Run {
		timeNow := time.Now().Unix()
		percent, err := cpu.Percent(time.Duration(config.SamplingInterval)*time.Second, false)
		if err == nil {
			everySamplingCpu += percent[0]
			samplingCpuTimes++
		}
		if timeNow%config.ReportInterval == 0 && samplingCpuTimes != 0 {
			client.RequestToServer(pb.ReportReq{
				UId:        global.GetUId(),
				Timestamp:  timeNow,
				Metric:     global.CPU_RATE,
				Dimensions: map[string]string{LOCAL_IP: global.GetIP()},
				Value:      everySamplingCpu / float64(samplingCpuTimes),
			})
			everySamplingCpu = 0.0
			samplingCpuTimes = 0
		}
	}

}
