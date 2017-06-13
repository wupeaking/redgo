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
	con.Do("FLUSHALL")

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

func TestList(t *testing.T) {
	con, e := redis.DialTimeout("tcp", "127.0.0.1:6379", 10*time.Second, 1*time.Second, 1*time.Second)
	if e != nil {
		t.Fatal("连接redis出错: ", e)
	}
	con.Do("FLUSHALL")
	// LPUSH
	_, e = con.Do("LPUSH", "listdemo", "aa", "bb", "cc", "dd")
	if e != nil {
		t.Fatal("lpush command err: ", e)
	}

	// LRANGE
	lists, e := redis.Strings(con.Do("LRANGE", "listdemo", 0, -1))
	if e != nil {
		t.Fatal("lrange command err: ", e)
	}
	if len(lists) != 4 || lists[0] != "dd" || lists[1] != "cc" || lists[2] != "bb" || lists[3] != "aa" {
		t.Fatal("lrange command err: ", lists)
	}
	con.Do("DEL", "listdemo")

	// RPUSH
	_, e = con.Do("RPUSH", "listdemo", "aa", "bb", "cc", "dd")
	if e != nil {
		t.Fatal("rpush command err: ", e)
	}
	lists, e = redis.Strings(con.Do("LRANGE", "listdemo", 0, -1))
	if e != nil {
		t.Fatal("lrange command err: ", e)
	}
	if len(lists) != 4 || lists[0] != "aa" || lists[1] != "bb" || lists[2] != "cc" || lists[3] != "dd" {
		t.Fatal("lrange command err: ", lists)
	}

	// LSET
	_, e = con.Do("LSET", "listdemo", 1, "xxoo")
	if e != nil {
		t.Fatal("rpush command err: ", e)
	}
	lists, e = redis.Strings(con.Do("LRANGE", "listdemo", 0, -1))
	if e != nil {
		t.Fatal("lrange command err: ", e)
	}
	if len(lists) != 4 || lists[0] != "aa" || lists[1] != "xxoo" || lists[2] != "cc" || lists[3] != "dd" {
		t.Fatal("lrange command err: ", lists)
	}

	// LLEN
	listlen, e := redis.Int(con.Do("LLEN", "listdemo"))
	if e != nil {
		t.Fatal("lrange command err: ", e)
	}
	if listlen != 4 {
		t.Fatal("lrange command err: ", lists)
	}

	// LINSERT
	_, e = con.Do("LINSERT", "listdemo", "AFTER", "aa", "aaxxoo")
	if e != nil {
		t.Fatal("rpush command err: ", e)
	}
	lists, e = redis.Strings(con.Do("LRANGE", "listdemo", 0, -1))
	if e != nil {
		t.Fatal("lrange command err: ", e)
	}
	if len(lists) != 5 || lists[1] != "aaxxoo" {
		t.Fatal("lrange command err: ", lists)
	}

	// LPOP
	lpop, e := redis.String(con.Do("LPOP", "listdemo"))
	if e != nil || lpop != "aa" {
		t.Fatal("rpush command err: ", e)
	}

	con.Do("DEL", "listdemo")
	_, e = con.Do("RPUSH", "listdemo", "aa", "bb", "cc", "dd")
	// RPOP
	rpop, e := redis.String(con.Do("RPOP", "listdemo"))
	if e != nil || rpop != "dd" {
		t.Fatal("rpush command err: ", e)
	}

	//LREM
	_, e = con.Do("LREM", "listdemo", 1, "aa")
	lists, e = redis.Strings(con.Do("LRANGE", "listdemo", 0, -1))

	if len(lists) != 2 || lists[0] != "bb" || lists[1] != "cc" || e != nil {
		t.Fatal("lrange command err: ", e, lists)
	}

}

func TestHash(t *testing.T) {
	con, e := redis.DialTimeout("tcp", "127.0.0.1:6379", 10*time.Second, 1*time.Second, 1*time.Second)
	if e != nil {
		t.Fatal("连接redis出错: ", e)
	}
	con.Do("FLUSHALL")
	// HSET
	_, e = con.Do("HSET", "hsetdemo", "name", "wupeaking")
	if e != nil {
		t.Fatal("hset command err: ", e)
	}

	// HGET
	getv, e := redis.String(con.Do("HGET", "hsetdemo", "name"))
	if e != nil || getv != "wupeaking" {
		t.Fatal("hset command err: ", e)
	}

	// HMSET
	_, e = con.Do("HMSET", "hmsetdemo", "name", "wupeaking", "age", 27)
	if e != nil {
		t.Fatal("hmset command err: ", e)
	}
	// HMGET
	mgetv, e := redis.Strings(con.Do("HMGET", "hmsetdemo", "name", "age"))
	if e != nil || len(mgetv) != 2 || mgetv[0] != "wupeaking" || mgetv[1] != "27" {
		t.Fatal("hmget command err: ", e, mgetv)
	}

	// HExists
	is, e := redis.Int(con.Do("HEXISTS", "hmsetdemo", "name"))
	if e != nil || is != 1 {
		t.Fatal("HExists command err: ", e, mgetv)
	}

	// HKeys
	keys, e := redis.Strings(con.Do("HKEYS", "hmsetdemo"))
	if e != nil || len(keys) != 2 {
		t.Fatal("HKeys command err: ", e, mgetv)
	}

	// HVals
	vals, e := redis.Strings(con.Do("HVALS", "hmsetdemo"))
	if e != nil || len(vals) != 2 {
		t.Fatal("HKeys command err: ", e, mgetv)
	}

	// HGetAll
	kv, e := redis.Strings(con.Do("HGETALL", "hmsetdemo"))
	if e != nil || len(kv) != 4 {
		t.Fatal("HGetAll command err: ", e, mgetv)
	}

	// HLen
	length, e := redis.Int(con.Do("HLEN", "hmsetdemo"))
	if e != nil || length != 2 {
		t.Fatal("HLen command err: ", e, mgetv)
	}

	// HDel
	con.Do("HDEL", "hmsetdemo", "age")

	ret, e := con.Do("HGET", "hmsetdemo", "age")
	if ret != nil {
		t.Fatal("HDel command err: ", e, ret)
	}

}

func TestSet(t *testing.T) {
	con, e := redis.DialTimeout("tcp", "127.0.0.1:6379", 10*time.Second, 1*time.Second, 1*time.Second)
	if e != nil {
		t.Fatal("连接redis出错: ", e)
	}
	con.Do("FLUSHALL")

	// SAdd
	_, e = con.Do("SADD", "setdemo", "golang", "python", "java")
	if e != nil {
		t.Fatal("sadd command err: ", e)
	}
	// SMembers
	members, e := redis.Strings(con.Do("SMEMBERS", "setdemo"))
	if len(members) != 3 || e != nil {
		t.Fatal("SMembers command err: ", e)
	}
	// SIsMember
	is, e := redis.Bool(con.Do("SISMEMBER", "setdemo", "python"))
	if !is || e != nil {
		t.Fatal("SIsMember command err: ", e, is)
	}

	// SCard
	length, e := redis.Int(con.Do("SCARD", "setdemo"))
	if length != 3 || e != nil {
		t.Fatal("SCard command err: ", e, is)
	}

	// SInter
	con.Do("SADD", "setdemo1", "golang", "php", "ruby", "c++")

	inter, e := redis.Strings(con.Do("SINTER", "setdemo", "setdemo1"))
	if len(inter) != 1 || e != nil {
		t.Fatal("SInter command err: ", e, inter)
	}

	// SInterStore
	con.Do("SINTERSTORE", "setdemointer", "setdemo", "setdemo1")
	is, e = redis.Bool(con.Do("SISMEMBER", "setdemointer", "golang"))
	if !is || e != nil {
		t.Fatal("SInterStore command err: ", e, is)
	}

	// SDiff
	diff, e := redis.Strings(con.Do("SDIFF", "setdemo", "setdemo1"))
	if len(diff) != 2 || e != nil {
		t.Fatal("SDiff command err: ", e, diff)
	}

	//SDiffStore
	con.Do("SDIFFSTORE", "setdemodiff", "setdemo", "setdemo1")
	is, e = redis.Bool(con.Do("SISMEMBER", "setdemodiff", "python"))
	if !is || e != nil {
		t.Fatal("SDiffStore command err: ", e, is)
	}

	//SUnion
	union, e := redis.Strings(con.Do("SUNION", "setdemo", "setdemo1"))
	if len(union) != 6 || e != nil {
		t.Fatal("SDiff command err: ", e, union)
	}

	// SUnionStore
	con.Do("SUNIONSTORE", "setdemounion", "setdemo", "setdemo1")
	is, e = redis.Bool(con.Do("SISMEMBER", "setdemounion", "php"))
	if !is || e != nil {
		t.Fatal("SUnionStore command err: ", e, is)
	}

	// SRem
	con.Do("SREM", "setdemodiff", "golang")
	is, e = redis.Bool(con.Do("SISMEMBER", "setdemodiff", "golang"))
	if is || e != nil {
		t.Fatal("SUnionStore command err: ", e, is)
	}

	// SPop
	pop, e := redis.String(con.Do("SPOP", "setdemounion"))
	if pop == "" || e != nil {
		t.Fatal("spop  command err: ", e, pop)
	}

	// SMove

	con.Do("SMOVE", "setdemo", "newsetdemo")
	is, e = redis.Bool(con.Do("SISMEMBER", "newsetdemo", "golang"))
	if !is || e != nil {
		t.Fatal("SUnionStore command err: ", e, is)
	}

}
