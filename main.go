package main

import (
	"runtime/debug"
	"strconv"

	"time"

	"github.com/wupeaking/redgo/datastruct"
)

func main() {
	// 不定时的创建和删除map
	mapdebug()

	// 查看gc状态
	for {
		debug.FreeOSMemory()
		time.Sleep(1 * time.Second)
	}

}

func mapdebug() {
	d := datastruct.NewDict(&datastruct.DemoDictFuncs{})
	for i := 0; i < 1000000; i++ {
		//time.Sleep(50 * time.Millisecond)
		d.Set(strconv.Itoa(i), "this is test hahhhahh "+strconv.Itoa(i))
	}
	println("start delete........")
	for i := 0; i < 1000000; i++ {
		//time.Sleep(50 * time.Millisecond)
		d.Delete(strconv.Itoa(i))
	}
}

func mapdebug1() {
	d := make(map[string]string)
	for i := 0; i < 1000000; i++ {
		//time.Sleep(50 * time.Millisecond)
		d[strconv.Itoa(i)] = "this is test hahhhahh " + strconv.Itoa(i)
	}
	println("start delete........")
	for i := 0; i < 1000000; i++ {
		//time.Sleep(50 * time.Millisecond)
		delete(d, strconv.Itoa(i))
	}
}
