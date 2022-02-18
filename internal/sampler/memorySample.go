package sampler

import (
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	"github.com/DeltaDemand/athena-agent/internal/model"
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"sync"
	"time"
)

//memory_used采样器结构，实现Sampler接口
type MemSample struct {
	Config  model.MemConfs //对应配置
	running bool           //运行状态
	name    string         //指标名字
}

func (s *MemSample) GetMetricName() string {
	return s.name
}

func (s *MemSample) SetName(str string) {
	s.name = str
}
func (s *MemSample) GetConfigPtr() interface{} {
	return &s.Config
}

//采样执行的动作，配置改变时也进行进入判断，
func (s *MemSample) Execute(wg *sync.WaitGroup) error {
	//agent暂停或已经有实例在跑就不执行
	if !global.GetPause() && s.running == false {
		wg.Add(1)
		//设置该采样器在运行的状态
		s.running = true
		go func() {
			s.sendMemPercent()
			s.running = false
			wg.Done()
		}()
	}
	return nil
}

//采集多个时间间隔内存百分比数据，求平均后上报
func (s *MemSample) sendMemPercent() {
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
			memInfo, err := mem.VirtualMemory()
			if err == nil {
				samplingResult += memInfo.UsedPercent
				samplingTimes++
			}
			//采集完数据，重置采集时间间隔
			samplingInterval = s.Config.SamplingInterval
		}
		timeNow := time.Now().Unix()
		//整点发送数据
		if timeNow%s.Config.ReportInterval == 0 && samplingTimes != 0 {
			memInfo, _ := mem.VirtualMemory()
			client.RequestToServer(pb.ReportReq{
				UId:        global.GetUId(),
				Timestamp:  timeNow,
				Metric:     s.name,
				Dimensions: map[string]string{LocalIp: global.GetIP(), "memTotal": strconv.FormatUint(memInfo.Total, 10)},
				Value:      samplingResult / float64(samplingTimes),
			})
			samplingResult = 0.0
			samplingTimes = 0
		}
		samplingInterval-- //采集时间间隔减一个单位
		time.Sleep(time.Second)
	}

}
