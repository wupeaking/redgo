package main

import (
	"strconv"

	"unsafe"

	"github.com/wupeaking/redgo/datastruct"
)

func main() {
	// 不定时的创建和删除map
	mapdebug()

	// 查看gc状态
	// for {
	// 	debug.FreeOSMemory()
	// 	time.Sleep(1 * time.Second)
	// }

	d := &demo{c: "hello world hehahhhhhhhhhhhhhhhhh", d: []string{"a", "dadada", "adhiwhihdiwh"}}
	ptr := unsafe.Pointer(d)
	println(*(*string)(unsafe.Pointer(uintptr(ptr) + unsafe.Offsetof(d.c))))
	println(unsafe.Sizeof(d.d))

}

type demo struct {
	a int
	b int64
	c string
	d []string
	e map[string]string
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
