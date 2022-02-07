package inputArgs

import (
	"flag"
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
)

func Parse(confs *appConfigs.Config) {
	flag.StringVar(&confs.ReportServer.Ip, "ip", confs.ReportServer.Ip, "监控服务器ip地址")
	flag.StringVar(&confs.ReportServer.Port, "p", confs.ReportServer.Port, "监控服务器监听端口号")
	flag.BoolVar(&confs.CpuConfi.Run, "cpu", confs.CpuConfi.Run, "是否开启cpu采样")
	flag.BoolVar(&confs.MemConfi.Run, "mem", confs.MemConfi.Run, "是否开启mem采样")
	flag.BoolVar(&confs.DiskConfi.Run, "disk", confs.DiskConfi.Run, "是否开启disk采样")
	flag.Int64Var(&confs.CpuConfi.ReportInterval, "cpuR", confs.CpuConfi.ReportInterval, "cpu上报时间间隔")
	flag.Int64Var(&confs.MemConfi.ReportInterval, "memR", confs.MemConfi.ReportInterval, "mem上报时间间隔")
	flag.Int64Var(&confs.DiskConfi.ReportInterval, "diskR", confs.DiskConfi.ReportInterval, "disk上报时间间隔")
	flag.Int64Var(&confs.MemConfi.SamplingInterval, "memS", confs.MemConfi.SamplingInterval, "mem采样时间间隔")
	flag.Parse()
}
