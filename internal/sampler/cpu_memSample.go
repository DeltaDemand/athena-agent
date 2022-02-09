package sampler

import (
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"sync"
	"time"
)

type cpuMemSample struct {
	name   string
	config appConfigs.CpuMemConfs
}

func (s *cpuMemSample) GetName() string {
	return s.name
}
func (s *cpuMemSample) Execute(wg *sync.WaitGroup) {
	sendCpuMemLogic(s.config)
	wg.Done()
}

func sendCpuMemLogic(config appConfigs.CpuMemConfs) {
	if config.Run && config.SamplingInterval > config.ReportInterval {
		global.Logger.Println("mem的采样间隔大于上报间隔，自动修改采样间隔等于上报间隔")
		config.SamplingInterval = config.ReportInterval
		cpu.Percent(0, false)
	}

	value := 0.0
	times := 0
	warnStr := "cpu percent over" + strconv.FormatFloat(config.CpuPercentThreshold, 'f', 2, 64) + "% and mem percent over" + strconv.FormatFloat(config.MemPercentThreshold, 'f', 2, 64) + "% times"
	samplingInterval := int64(0)
	for config.Run {
		if samplingInterval == 0 {
			percent, cpuErr := cpu.Percent(0, false)
			memInfo, memErr := mem.VirtualMemory()
			if memErr == nil && cpuErr == nil {
				if percent[0] > config.CpuPercentThreshold && memInfo.UsedPercent > config.MemPercentThreshold {
					times++
				}
			}
			samplingInterval = config.SamplingInterval
		}

		timeNow := time.Now().Unix()
		if timeNow%config.ReportInterval == 0 {
			if times > 0 {
				value = 1.0
			}
			client.RequestToServer(pb.ReportReq{
				UId:        global.GetUId(),
				Timestamp:  timeNow,
				Metric:     global.CpuMem,
				Dimensions: map[string]string{LocalIp: global.GetIP(), warnStr: strconv.Itoa(times)},
				Value:      value,
			})
			times = 0
			value = 0.0
		}
		samplingInterval--
		time.Sleep(time.Second)
	}
}
