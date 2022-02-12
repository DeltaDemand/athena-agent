package client

import (
	"context"
	"encoding/json"
	"github.com/DeltaDemand/athena-agent/global"
	"go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

type Etcd struct {
	Apply       bool     `json:"apply"`
	EndPoints   []string `json:"endPoints"`
	DialTimeout int      `json:"dialTimeout"`
}

var (
	cli *clientv3.Client
)

func (e *Etcd) Connect() error {
	if e.Apply {
		var err error
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   e.EndPoints,
			DialTimeout: time.Duration(e.DialTimeout) * time.Second,
		})
		if err != nil {
			global.Logger.Printf("connect to etcd failed, err:%v\n", err)
			return err
		}
	}
	return nil
}

// WatchConfig
//参数  ConfigChangeExecuter表示该配置更新是否要执行事件，nil表示不用
func (e *Etcd) WatchConfig(key string, configs interface{}, obj ConfigChangeExecuter, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		watchCh := cli.Watch(context.TODO(), global.AgentGroup+global.AgentName+key)
		for res := range watchCh {
			//Etcd的配置改变了，不需监听变化了，打断循环
			if e.Apply == false {
				break
			}
			value := res.Events[0].Kv.Value
			if err := json.Unmarshal(value, configs); err != nil {
				global.Logger.Println(key, " watchConfig err", err)
				continue
			}
			//该配置改变需要执行事件
			if obj != nil {
				err := obj.Execute(wg)
				//执行失败
				if err != nil {
					global.Logger.Printf("%s Configs change to %#v, fail to Execute", key, configs)
				}
				//重连成功
				global.Logger.Printf("%s Configs change to %#v,change Execute", key, configs)
			} else {
				//不需要执行事件
				global.Logger.Printf("%s Configs change to %#v", key, configs)
			}
		}
		wg.Done()
	}()
}

// 更新配置后执行：重新连接etcd
func (e *Etcd) Execute(wg *sync.WaitGroup) error {
	//先释放旧连接
	e.CloseConn()
	if e.Apply == true {
		//还使用etcd配置，再次连接
		err := e.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Etcd) CloseConn() {
	cli.Close()
}
