package main

import (
	"github.com/wupeaking/redgo/server"
)

// type demo struct {
// 	a int
// 	B time.Time
// }

// type NewTime struct {
// 	time.Time
// }

// func (t NewTime) String() string {
// 	return "xxxooo"
// }

// // Jon公共方法
// func (myself *demo) Jon() int {
// 	return myself.a
// }

func main() {
	// d := &demo{a: 3, B: time.Now()}
	// fmt.Println(d.B)
	// // // newfunc := func(in []reflect.Value) []reflect.Value {
	// // // 	return in
	// // // }
	// // //fn, _ := reflect.TypeOf(d).MethodByName("Jon")
	// // // fn := reflect.ValueOf(d.Jon).Elem()
	// // // v := reflect.MakeFunc(reflect.TypeOf(d.Jon), newfunc)
	// // // println(fn.CanSet())
	// // // fn.Set(v)
	// // // nfn := reflect.ValueOf(d).Elem()
	// // // println(nfn.CanSet())
	// // // nfn.Set(v)
	// // println(d.Jon())
	// a := NewTime{time.Now().Add(time.Hour)}
	// ta := (*time.Time)(unsafe.Pointer(&a))
	// v := reflect.ValueOf(*ta)

	// reflect.ValueOf(d).Elem().FieldByName("B").Set(v)

	// fmt.Println(d.B.String())
	server.StartServer("./config.yaml")
}
