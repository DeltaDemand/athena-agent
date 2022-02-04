FROM alpine:latest
ADD test /usr/local/bin/

# Alpine Linux doesn't use pam, which means that there is no /etc/nsswitch.conf,
# but Golang relies on /etc/nsswitch.conf to check the order of DNS resolving
# (see https://github.com/golang/go/commit/9dee7771f561cf6aee081c0af6658cc81fac3918)
# To fix this we just create /etc/nsswitch.conf and add the following line:
RUN echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf
WORKDIR /usr/local/bin/

CMD ["test"]

MAINTAINER="2390647320@qq.com"

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct


COPY sensors/sensor /tmp/

RUN chmod +x /tmp/sensor

ENTRYPOINT ["/tmp/sensor"]



# 移动到工作目录：/build
WORKDIR
# 将代码复制到容器中
COPY . .
RUN go mod tidy -compat=1.17

# go generate 编译前自动执行代码
# go env 查看go的环境变量
# go build -o athena-server . 打包项目生成文件名为athena-server的二进制文件
RUN go generate && go env && go build -o agent


FROM alpine:latest

LABEL MAINTAINER="EZ4BRUCE@lhy122786302@gmail.com"
WORKDIR /go/src/athena-server

# 把/go/src/gin-vue-admin整个文件夹的文件到当前工作目录
COPY --from=0 /go/src/athena-server ./

EXPOSE 8888

ENTRYPOINT ./athena-server