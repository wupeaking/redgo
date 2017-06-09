## 一个golang版本的redis的单机数据库实现 (仅仅为了学习)

[![Travis](https://travis-ci.org/wupeaking/redgo.svg?branch=master)](https://travis-ci.org/wupeaking/redgo)
[![GoDoc](https://godoc.org/github.com/wupeaking/kafkainfo?status.svg)](https://godoc.org/github.com/wupeaking/redgo)
[![codecov.io](https://codecov.io/gh/wupeaking/redgo/coverage.svg?branch=master)](https://codecov.io/gh/wupeaking/redgo?branch=master)

## 安装
> 源码build

```shell

# clone
export GOPATH=`pwd`
mkdir -p  $GOPATH/src/github.com/wupeaking && cd $GOPATH/src/github.com/wupeaking
git clone git@github.com:wupeaking/redgo.git && cd redgo

# build
go build -o redgo main.go

```

# 使用 

* 增加一个配置文件 config.yaml
```shell
host: 0.0.0.0
port: 6379
```
* 启动服务
```shell
# ./redgo
```

* 使用redis客户端连接 
```shell
# redis-cli
> set demo helloworld
> get demo
> helloworld
```