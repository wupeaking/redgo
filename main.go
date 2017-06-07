package main

import (
	"github.com/wupeaking/redgo/server"
)

type demo struct {
	a int
}

// Jon公共方法
func (myself *demo) Jon() int {
	return myself.a
}

func main() {
	// d := &demo{a: 3}
	// // newfunc := func(in []reflect.Value) []reflect.Value {
	// // 	return in
	// // }
	// //fn, _ := reflect.TypeOf(d).MethodByName("Jon")
	// // fn := reflect.ValueOf(d.Jon).Elem()
	// // v := reflect.MakeFunc(reflect.TypeOf(d.Jon), newfunc)
	// // println(fn.CanSet())
	// // fn.Set(v)
	// // nfn := reflect.ValueOf(d).Elem()
	// // println(nfn.CanSet())
	// // nfn.Set(v)
	// println(d.Jon())

	server.StartServer("./config.yaml")
}
