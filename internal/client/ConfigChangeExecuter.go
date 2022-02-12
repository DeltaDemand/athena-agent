package client

import (
	"sync"
)

type ConfigChangeExecuter interface {
	Execute(wg *sync.WaitGroup) error //配置变化时，执行的处理接口
}
