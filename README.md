### Agent部署

拉取[Agent端](https://github.com/DeltaDemand/athena-agent)

```bash
git clone https://github.com/DeltaDemand/athena-agent.git
```

进入athena-agent目录，执行以下docker命令即可启动Agent端：

```bash
#进入athena-agent目录
cd athena-agent
#构建docker镜像
docker build -t athena-agent .

#本机测试：使用docker内网
docker run -d -i --name host01 --network athena_frontend athena-agent -aggregationTime=5 -checkAlive=30 -cpuR=10 -memR=10 -diskR=10 -cpu_memR=10

#云服务器测试
#阿里云
docker run -d --name host01 athena-agent -ip="112.74.60.132" -aggregationTime=5 -checkAlive=30 -cpuR=10 -memR=10 -diskR=10 -cpu_memR=10 -group=group01 -name=agent01
#腾讯云
docker run -d --name host01 athena-agent -ip="1.12.242.39" -aggregationTime=5 -checkAlive=30 -cpuR=10 -memR=10 -diskR=10 -cpu_memR=10 -group=group01 -name=agent01
# -group string
#        etcd上Agent分组 (default "g01")
# -name string
#        etcd上Agent名字 (default "A01")
# -cpuR int
#        cpu上报时间间隔 (default 60)
````
告警测试：
```bash
#进入agent终端
docker exec -it host01 /bin/sh
#运行程序，使cpu、内存等数据有所波动
./test/poseidon -n=12 -t=10
# -append int                                 
#        每个goroutine内append字符串的次数     
#  -goN int                                    
#        创造goroutine跑死循环个数 (default 10)
#  -sleep int                                  
#        每次循环睡眠时间(ns)                  
#  -times int                                  
#        死循环时间(s) (default 90) 
```
