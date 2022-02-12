package sampler

import (
	"sync"
)

// Sampler 指标(名词)--->采样器（执行者）.采样器接口定义
type Sampler interface {
	Execute(wg *sync.WaitGroup) error //执行采样并发送
	GetMetricName() string            //获取该采样器的指标名字
	GetConfigPtr() interface{}        //用于更新配置
}
