package sampler

import (
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	"github.com/DeltaDemand/athena-agent/internal/model"
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/shirou/gopsutil/cpu"
	"strconv"
	"sync"
	"time"
)

type CpuSample struct {
	Config  model.CpuConfs //对应配置
	running bool           //运行状态
	name    string         //指标名字
}

func (s *CpuSample) GetMetricName() string {
	return s.name
}
func (s *CpuSample) SetName(str string) {
	s.name = str
}
func (s *CpuSample) GetConfigPtr() interface{} {
	return &s.Config
}

func (s *CpuSample) Execute(wg *sync.WaitGroup) error {
	//agent暂停或注册没成功或已经有实例在跑就不执行
	if !global.GetPause() && global.GetRegisterSuccess() && s.running == false {
		wg.Add(1)
		//设置该采样器在运行的状态
		s.running = true
		go func() {
			s.sendCpuPercent()
			s.running = false
			wg.Done()
		}()
	}
	return nil
}

func (s *CpuSample) sendCpuPercent() {
	if s.Config.Run {
		cpu.Percent(0, false)
	}
	//处理采样间隔大于上报间隔情况
	if s.Config.Run && s.Config.SamplingInterval > s.Config.ReportInterval {
		global.Logger.Println("mem的采样间隔大于上报间隔，自动修改采样间隔等于上报间隔")
		s.Config.SamplingInterval = s.Config.ReportInterval
	}
	samplingResult := 0.0
	samplingTimes := 0
	samplingInterval := int64(0)
	//循环判断agent是否停止，或该指标是否启动
	for !global.GetPause() && s.Config.Run {
		if samplingInterval == 0 {
			//采集时间间隔到了，采集数据
			//如果interval=0或者None时，比较自上次调用或模块导入后经过的系统CPU时间，立即返回。所以第一次的返回的数据是个无意义的数据。
			//当percpu是True返回表示利用率的浮点数列表，以每个逻辑CPU的百分比表示。
			percent, err := cpu.Percent(0, false)
			if err == nil {
				samplingResult += percent[0]
				samplingTimes++
			}
			//采集完数据，重置采集时间间隔
			samplingInterval = s.Config.SamplingInterval
		}
		timeNow := time.Now().Unix()
		//整点发送数据
		if timeNow%s.Config.ReportInterval == 0 && samplingTimes != 0 {
			logicalCores, _ := cpu.Counts(true)
			client.RequestToServer(pb.ReportReq{
				UId:        global.GetUId(),
				Timestamp:  timeNow,
				Metric:     s.name,
				Dimensions: map[string]string{LocalIp: global.GetIP(), "physicalCores": strconv.Itoa(logicalCores)},
				Value:      samplingResult / float64(samplingTimes),
			})
			samplingResult = 0.0
			samplingTimes = 0
		}
		samplingInterval-- //采集时间间隔减一个单位
		time.Sleep(time.Second)
	}

}
