package server

import (
	"testing"

	"time"

	redis "github.com/garyburd/redigo/redis"
	redis_server "github.com/wupeaking/go-redis-server"
)

// 增加测试文件
var start bool

func stratServer() {
	if !start {
		handler := new(SrvHandler)
		server, _ := redis_server.NewServer(redis_server.DefaultConfig().Host("0.0.0.0").Port(6380).Handler(handler))
		go server.ListenAndServe()
		println("start server...")
		time.Sleep(10 * time.Second)
		start = true
	}
}

func TestString(t *testing.T) {
	con, e := redis.DialTimeout("tcp", "127.0.0.1:6379", 10*time.Second, 1*time.Second, 1*time.Second)
	if e != nil {
		t.Fatal("连接redis出错: ", e)
	}
	con.Flush()

	// set 一个值
	_, e = con.Do("SET", "sringdemo", "helloworld")
	if e != nil {
		t.Fatal("set command err : ", e)
	}

	// 获取设置的值
	v, e := redis.String(con.Do("GET", "sringdemo"))
	if v != "helloworld" || e != nil {
		t.Fatal("get command err")
	}

	// Mset
	_, e = con.Do("MSET", "sringdemo1", "helloworld1", "sringdemo2", "helloworld2")
	if e != nil {
		t.Fatal("set command err : ", e)
	}
	vs, e := redis.Strings(con.Do("MGET", "sringdemo1", "sringdemo2"))

	if (e != nil) || (vs[0] != "helloworld1") || (vs[1] != "helloworld2") {
		t.Fatal("mset command err : ", e, vs[0], vs[1])
	}

	// GETRANGE
	vb, e := redis.String(con.Do("GETRANGE", "sringdemo", 0, 5))
	if e != nil || string(vb) != "hellow" {
		t.Fatal("getrange command err : ", e, "v: ", vb)
	}

	// strlen
	l, e := redis.Int(con.Do("STRLEN", "sringdemo"))
	if e != nil || l != len("helloworld") {
		t.Fatal("getrange command err : ", e, v)
	}

	// append
	con.Do("APPEND", "sringdemo", " redis")
	v, e = redis.String(con.Do("GET", "sringdemo"))
	if e != nil || l != len("helloworld") {
		t.Fatal("append command err : ", e, v)
	}
}
