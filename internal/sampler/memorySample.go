package sampler

import (
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"sync"
	"time"
)

type memSample struct {
	name   string
	config appConfigs.MemConfs
}

func (s *memSample) GetName() string {
	return s.name
}
func (s *memSample) Execute(wg *sync.WaitGroup) {
	sendMemPercent(s.config)
	wg.Done()
}

func sendMemPercent(config appConfigs.MemConfs) {
	if config.Run && config.SamplingInterval > config.ReportInterval {
		global.Logger.Println("mem的采样间隔大于上报间隔，自动修改采样间隔等于上报间隔")
		config.SamplingInterval = config.ReportInterval
	}
	samplingResult := 0.0
	samplingTimes := 0
	samplingInterval := int64(0)

	for config.Run {
		timeNow := time.Now().Unix()
		if samplingInterval == 0 {
			memInfo, err := mem.VirtualMemory()
			if err == nil {
				samplingResult += memInfo.UsedPercent
				samplingTimes++
			}
			samplingInterval = config.SamplingInterval
		}

		if timeNow%config.ReportInterval == 0 && samplingTimes != 0 {
			memInfo, _ := mem.VirtualMemory()
			client.RequestToServer(pb.ReportReq{
				UId:        global.GetUId(),
				Timestamp:  timeNow,
				Metric:     global.MemUsed,
				Dimensions: map[string]string{LocalIp: global.GetIP(), "memTotal": strconv.FormatUint(memInfo.Total, 10)},
				Value:      samplingResult / float64(samplingTimes),
			})
			samplingResult = 0.0
			samplingTimes = 0
		}
		samplingInterval--
		time.Sleep(time.Second)
	}
}
