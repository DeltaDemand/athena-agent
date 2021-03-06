### Agent部署

####  拉取并构建镜像[Agent端](https://github.com/DeltaDemand/athena-agent)
```bash
#构建docker镜像
docker build -t athena-agent https://github.com/DeltaDemand/athena-agent.git#main
```
#### 启动镜像

```bash
#本机测试：使用docker内网
docker run -d -i --name host-test --network athena_frontend athena-agent -aggregationTime=5 -checkAlive=30 -cpuR=10 -memR=10 -diskR=10 -cpu_memR=10
```
```bash
#agent运行连接云服务器
#阿里云
docker run -d --name host-test athena-agent -ip="112.74.60.132" -aggregationTime=5 -checkAlive=30 -cpuR=10 -memR=10 -diskR=10 -cpu_memR=10 -group=group01 -name=agent01
#腾讯云
docker run -d --name host-test athena-agent -ip="1.12.242.39" -aggregationTime=5 -checkAlive=30 -cpuR=10 -memR=10 -diskR=10 -cpu_memR=10 -group=group01 -name=agent01
# athena-agent 部分参数解释:
# -aggregationTime int                                                      
#        上报几次进行聚合，默认0(由server端决定)
# -checkAlive int                                                           
#        检测是否存活时间间隔 (default 120) 
# -ip string
#        监控服务器ip地址 (default "athena-server")
# -cpuR int
#        cpu上报时间间隔 (default 60)
# -group string
#        etcd上Agent分组 (default "g01")
# -name string
#        etcd上Agent名字 (default "A01")
````
### 告警测试：

```bash
#进入agent终端
docker exec -it host-test /bin/sh
```
#### 运行程序，使cpu、内存等数据有所波动
goN=12情况下改变append值能让内存占用大致如下，可根据机器内存大小来调整测试。

200000--->150M

2000000--->917M

3000000->1.4G

4000000->2.0G
```bash
./test/poseidon -goN=12 -append=2000000 -sleep=100000000
#以下函数测试能让cpu跑90%以上，如主机过热请增大sleep时间。
./test/poseidon   -sleep=0 -goN=12  -time=60
# poseidon参数解释:
# -append int                                 
#        每个goroutine内append字符串的次数     
#  -goN int                                    
#        创造goroutine跑死循环个数 (default 10)
#  -sleep int                                  
#        每次循环睡眠时间(ns)                  
#  -time int                                  
#        死循环时间(s) (default 90) 
```
