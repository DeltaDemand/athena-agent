package sampler

import (
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/shirou/gopsutil/cpu"
	"strconv"
	"sync"
	"time"
)

type cpuSample struct {
	name   string
	config appConfigs.CpuConfs
}

func (s *cpuSample) GetName() string {
	return s.name
}

func (s *cpuSample) Execute(wg *sync.WaitGroup) {
	sendCpuPercent(s.config)
	wg.Done()
}

func sendCpuPercent(config appConfigs.CpuConfs) {
	if config.Run {
		cpu.Percent(0, false)
	}
	for config.Run {
		timeNow := time.Now().Unix()

		if timeNow%config.ReportInterval == 0 {
			//如果interval=0或者None时，比较自上次调用或模块导入后经过的系统CPU时间，立即返回。所以第一次的返回的数据是个无意义的数据。
			//当percpu是True返回表示利用率的浮点数列表，以每个逻辑CPU的百分比表示。
			percent, _ := cpu.Percent(0, false)
			physicalCores, _ := cpu.Counts(false)
			client.RequestToServer(pb.ReportReq{
				UId:        global.GetUId(),
				Timestamp:  timeNow,
				Metric:     global.CpuRate,
				Dimensions: map[string]string{LocalIp: global.GetIP(), "physicalCores": strconv.Itoa(physicalCores)},
				Value:      percent[0],
			})
		}
		time.Sleep(time.Second)
	}
}
