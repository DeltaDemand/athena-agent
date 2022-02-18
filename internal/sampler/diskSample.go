package sampler

import (
	"github.com/DeltaDemand/athena-agent/global"
	"github.com/DeltaDemand/athena-agent/internal/client"
	"github.com/DeltaDemand/athena-agent/internal/model"
	pb "github.com/EZ4BRUCE/athena-proto/proto"
	"github.com/shirou/gopsutil/disk"
	"strconv"
	"sync"
	"time"
)

const LocalIp = "local_ip"

type DiskSample struct {
	Config  model.DiskConfs //对应配置
	running bool            //运行状态
	name    string          //指标名字
}

func (s *DiskSample) GetMetricName() string {
	return s.name
}
func (s *DiskSample) SetName(str string) {
	s.name = str
}

func (s *DiskSample) GetConfigPtr() interface{} {
	return &s.Config
}
func (s *DiskSample) Execute(wg *sync.WaitGroup) error {
	//agent暂停或注册没成功或已经有实例在跑就不执行
	if !global.GetPause() && global.GetRegisterSuccess() && s.running == false {
		wg.Add(1)
		//设置该采样器在运行的状态
		s.running = true
		go func() {
			s.sendDiskPercent()
			s.running = false
			wg.Done()
		}()
	}
	return nil
}

//采集并发送磁盘百分比数据
func (s *DiskSample) sendDiskPercent() {
	//循环判断agent是否停止，或该指标是否启动
	for !global.GetPause() && s.Config.Run {
		timeNow := time.Now().Unix()

		if timeNow%s.Config.ReportInterval == 0 {
			//整点采集数据并发送数据
			parts, _ := disk.Partitions(true)
			diskInfo, err := disk.Usage(parts[0].Mountpoint)
			if err == nil {
				client.RequestToServer(pb.ReportReq{
					UId:        global.GetUId(),
					Timestamp:  timeNow,
					Metric:     s.name,
					Dimensions: map[string]string{LocalIp: global.GetIP(), "diskTotal": strconv.FormatUint(diskInfo.Total, 10)},
					Value:      diskInfo.UsedPercent,
				})
			} else {
				global.Logger.Println("获取硬盘数据失败")
			}
		}
		time.Sleep(time.Second)
	}
}
