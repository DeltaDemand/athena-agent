package client

import (
	"context"
	"encoding/json"
	"github.com/DeltaDemand/athena-agent/global"
	"go.etcd.io/etcd/client/v3"
	"strconv"
	"sync"
	"time"
)

type Etcd struct {
	ConfigServer string   `json:"config_server"`
	AgentGroup   string   `json:"agent_group"`
	AgentName    string   `json:"agent_name"`
	Apply        bool     `json:"apply"`
	EndPoints    []string `json:"endPoints"`
	DialTimeout  int      `json:"dialTimeout"`
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
func (e *Etcd) WatchConfig(configName string, configs interface{}, obj ConfigChangeExecuter, wg *sync.WaitGroup) {
	wg.Add(1)
	//要监听的配置作为key
	key := e.ConfigServer + global.Split + e.AgentGroup + global.Split + e.AgentName + global.Split + configName
	//先上传一份配置到etcd服务器，相当于注册
	value, _ := json.Marshal(configs)
	//etcd没有就阻塞，不会往下运行，节约资源
	cli.Put(context.TODO(), key, string(value))
	go func() {
		watchCh := cli.Watch(context.TODO(), key)
		for res := range watchCh {
			value := res.Events[0].Kv.Value
			if err := json.Unmarshal(value, configs); err != nil {
				global.Logger.Println(configName, " watchConfig err", err)
				continue
			}
			//该配置改变需要执行事件
			if obj != nil {
				err := obj.Execute(wg)
				//执行失败
				if err != nil {
					global.Logger.Printf("%s Configs change to %#v, fail to Execute", configName, configs)
				} else {
					//执行成功
					global.Logger.Printf("%s Configs change to %#v,change Execute", configName, configs)
				}
			} else {
				//不需要执行事件直接打印结构体
				global.Logger.Printf("%s Configs change to %#v", configName, configs)
			}
		}
		wg.Done()
	}()
}

// 更新配置后执行：重新连接etcd
func (e *Etcd) Execute(wg *sync.WaitGroup) error {
	//防止关闭连接后goroutine为零退出程序
	wg.Add(1)
	defer wg.Done()
	//释放旧连接
	e.CloseConn()
	global.AgentGroup = e.AgentGroup
	global.AgentName = e.AgentName
	//Apply==true就重连
	if e.Apply {
		//设置etcd变化和没成功处理此次变化，再次连接
		global.HandleChangeSuccess = false
		err := e.Connect()
		if err != nil {
			global.Logger.Println("etcd重连失败---", err)
			return err
		}
	}
	//通知此次etcd变化
	global.EtcdChange <- true
	//未成功就一直等到成功
	for !global.HandleChangeSuccess {
		time.Sleep(time.Second)
	}
	return nil
}

//把Agent状态更新到etcd服务器
func RefreshAgentState(pause bool) {
	key := global.ConfigServer + global.Split + global.AgentGroup + global.Split + global.AgentName + global.Split + global.Agent
	//根据Agent状态更新configServer上的状态
	if pause {
		//etcd没有开启就阻塞，不会往下运行，节约资源
		cli.Put(context.TODO(), key, "{\"checkAlive\":"+strconv.Itoa(global.CheckAlive)+",\"pause\":true,\"exit\":false}")
	} else {
		cli.Put(context.TODO(), key, "{\"checkAlive\":"+strconv.Itoa(global.CheckAlive)+",\"pause\":false,\"exit\":false}")
	}
	global.Logger.Println(key)
}

func (e *Etcd) CloseConn() {
	cli.Close()
}
