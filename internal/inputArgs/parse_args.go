package inputArgs

import (
	"flag"
	appConfigs "github.com/DeltaDemand/athena-agent/configs"
	"github.com/DeltaDemand/athena-agent/global"
	"strings"
)

// Parse :Agent命令行参数解析
func Parse(confs *appConfigs.Config) {
	flag.StringVar(&confs.AgentConfi.Group, "group", confs.AgentConfi.Group, "Agent分组")
	flag.StringVar(&confs.AgentConfi.Name, "name", confs.AgentConfi.Name, "Agent名字")

	flag.IntVar(&confs.AgentConfi.CheckAlive, "checkAlive", confs.AgentConfi.CheckAlive, "检测是否存活时间间隔")
	flag.Int64Var(&global.AggregationTime, "aggregationTime", 0, "建议聚合时间，默认0(由server端决定)")
	flag.StringVar(&confs.ReportServer.Ip, "ip", confs.ReportServer.Ip, "监控服务器ip地址")
	flag.StringVar(&confs.ReportServer.Port, "p", confs.ReportServer.Port, "监控服务器监听端口号")

	flag.BoolVar(&confs.CpuConfi.Run, "cpu", confs.CpuConfi.Run, "是否开启cpu采样")
	flag.BoolVar(&confs.MemConfi.Run, "mem", confs.MemConfi.Run, "是否开启mem采样")
	flag.BoolVar(&confs.DiskConfi.Run, "disk", confs.DiskConfi.Run, "是否开启disk采样")
	flag.BoolVar(&confs.CpuMemConfi.Run, "cpu_mem", confs.CpuMemConfi.Run, "是否开启cpu_mem采样")

	flag.Int64Var(&confs.CpuConfi.ReportInterval, "cpuR", confs.CpuConfi.ReportInterval, "cpu上报时间间隔")
	flag.Int64Var(&confs.MemConfi.ReportInterval, "memR", confs.MemConfi.ReportInterval, "mem上报时间间隔")
	flag.Int64Var(&confs.DiskConfi.ReportInterval, "diskR", confs.DiskConfi.ReportInterval, "disk上报时间间隔")
	flag.Int64Var(&confs.CpuMemConfi.ReportInterval, "cpu_memR", confs.CpuMemConfi.ReportInterval, "cpu_mem上报时间间隔")

	flag.Int64Var(&confs.MemConfi.SamplingInterval, "memS", confs.MemConfi.SamplingInterval, "mem采样时间间隔")
	flag.Int64Var(&confs.CpuMemConfi.SamplingInterval, "cpu_memS", confs.CpuMemConfi.SamplingInterval, "cpu_mem采样时间间隔")

	var ends string
	flag.BoolVar(&confs.Etcd.Apply, "etcd", true, "是否连接etcd")
	flag.StringVar(&ends, "endPoints", "112.74.60.132:2379,localhost:2379", "etcd节点的地址可以多个，用逗号(,)隔开")

	flag.Parse()
	//初始化全局变量值
	confs.Etcd.EndPoints = strings.Split(ends, ",")
	global.CheckAlive = confs.AgentConfi.CheckAlive
	global.AgentGroup = confs.AgentConfi.Group
	global.AgentName = confs.AgentConfi.Name
}
