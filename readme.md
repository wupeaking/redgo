## 一个golang版本的redis的单机数据库实现 (仅仅为了学习)

[![Travis](https://travis-ci.org/wupeaking/redgo.svg?branch=master)](https://travis-ci.org/wupeaking/redgo)
[![GoDoc](https://godoc.org/github.com/wupeaking/redgo?status.svg)](https://godoc.org/github.com/wupeaking/redgo)
[![codecov.io](https://codecov.io/gh/wupeaking/redgo/coverage.svg?branch=master)](https://codecov.io/gh/wupeaking/redgo?branch=master)

## 安装
> 源码build

```shell
# 1. 要求安装go
# 2. clone

export GOPATH=`pwd`
mkdir -p  $GOPATH/src/github.com/wupeaking && cd $GOPATH/src/github.com/wupeaking
git clone git@github.com:wupeaking/redgo.git && cd redgo

# build
make
```

> 生成docker镜像

```shell

# 1. 要求安装docker
# 2. clone

mkdir -p  `pwd`/src/github.com/wupeaking && cd $GOPATH/src/github.com/wupeaking
git clone git@github.com:wupeaking/redgo.git && cd redgo
export GOPATH=`pwd`
# 执行make
make docker

```

## 在宿主机上使用 

* 增加一个配置文件 config.yaml
```shell
host: 0.0.0.0
port: 6379
```

* 启动服务
```shell
# ./redgo
```
## 在docker环境运行
```shell
# 使用默认配置文件
docker run -p 6379:6379 --name=redgo -d redgo:v0.0.1 
# 如果需要更改配置文件 挂载配置文件路径 配置文件名为 config.yaml
docker run --net=host -v `pwd`:/etc/redgo --name=redgo -d redgo:v0.0.1 
```

* 使用redis客户端连接 
```shell
# redis-cli
> set demo helloworld
> get demo
> helloworld

> lpush listdemo aa bb cc dd
> lrange listdemo -1
aa
bb
cc
dd
> ....
```