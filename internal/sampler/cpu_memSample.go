package sampler

import (
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	"github.com/DeltaDemand/athena-agent/internal/model"
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"sync"
	"time"
)

type CpuMemSample struct {
	Config  model.CpuMemConfs //对应配置
	running bool              //运行状态
	name    string            //指标名字
}

func (s *CpuMemSample) GetMetricName() string {
	return s.name
}
func (s *CpuMemSample) SetName(str string) {
	s.name = str
}
func (s *CpuMemSample) GetConfigPtr() interface{} {
	return &s.Config
}

func (s *CpuMemSample) Execute(wg *sync.WaitGroup) error {
	//agent暂停或注册没成功或已经有实例在跑就不执行
	if !global.GetPause() && global.GetRegisterSuccess() && s.running == false {
		wg.Add(1)
		//设置该采样器在运行的状态
		s.running = true
		go func() {
			s.sendCpuMemLogic()
			s.running = false
			wg.Done()
		}()
	}
	return nil
}

//采集cpu和内存数据，并进行逻辑判断，上报逻辑判断结果
func (s *CpuMemSample) sendCpuMemLogic() {
	//处理采样间隔大于上报间隔情况
	if s.Config.Run && s.Config.SamplingInterval > s.Config.ReportInterval {
		global.Logger.Println("mem的采样间隔大于上报间隔，自动修改采样间隔等于上报间隔")
		s.Config.SamplingInterval = s.Config.ReportInterval
	}
	cpu.Percent(0, false)
	value := 0.0 //上报值
	times := 0   //超过预警条件次数
	warnStr := "cpu percent over" + strconv.FormatFloat(s.Config.CpuPercentThreshold, 'f', 2, 64) + "% and mem percent over" + strconv.FormatFloat(s.Config.MemPercentThreshold, 'f', 2, 64) + "% times"
	samplingInterval := int64(0)
	//循环判断agent是否停止，或该指标是否启动
	for !global.GetPause() && s.Config.Run {
		if samplingInterval == 0 {
			//采集时间间隔到了，采集数据
			percent, cpuErr := cpu.Percent(0, false)
			memInfo, memErr := mem.VirtualMemory()
			if memErr == nil && cpuErr == nil {
				if percent[0] > s.Config.CpuPercentThreshold && memInfo.UsedPercent > s.Config.MemPercentThreshold {
					//采集到数据，超过预警条件，则次数加一
					times++
				}
			}
			//采集完数据，重置采集时间间隔
			samplingInterval = s.Config.SamplingInterval
		}

		timeNow := time.Now().Unix()
		//整点采集数据并发送数据
		if timeNow%s.Config.ReportInterval == 0 {
			if times > 0 {
				//超过预警条件次数不为零，上报value = 1.0
				value = 1.0
			}
			////超过预警条件次数放在Dimensions中上报
			client.RequestToServer(pb.ReportReq{
				UId:        global.GetUId(),
				Timestamp:  timeNow,
				Metric:     s.name,
				Dimensions: map[string]string{LocalIp: global.GetIP(), warnStr: strconv.Itoa(times)},
				Value:      value,
			})
			//重置变量
			times = 0
			value = 0.0
		}
		samplingInterval--
		time.Sleep(time.Second)
	}
}
