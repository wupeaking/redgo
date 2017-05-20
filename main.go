package main

import (
	"fmt"
	"math"
	"runtime/debug"
	"strconv"
	"unsafe"

	"time"

	"github.com/wupeaking/redgo/datastruct"
)

type hmap struct {
	count int // # live cells == size of map.  Must be first (used by len() builtin)
	flags uint8
	B     uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	hash0 uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)
	//overflow   *[2]*[]*runtime.bmap
}

func main() {
	m := make(map[string]string)
	c, b := getInfo(m)
	fmt.Println("count: ", c, "b: ", b)
	for i := 0; i < 10000; i++ {
		m[strconv.Itoa(i)] = strconv.Itoa(i)
		if i%200 == 0 {
			c, b := getInfo(m)
			cap := math.Pow(float64(2), float64(b))
			fmt.Printf("count: %d, b: %d, load: %f\n", c, b, float64(c)/cap)
		}
	}
	println("开始删除------")
	for i := 0; i < 10000; i++ {
		delete(m, strconv.Itoa(i))
		if i%200 == 0 {
			c, b := getInfo(m)
			cap := math.Pow(float64(2), float64(b))
			fmt.Println("count: ", c, "b:", b, "load: ", float64(c)/cap)
		}
	}

	debug.FreeOSMemory()
	c, b = getInfo(m)
	fmt.Println("释放后: ", "count: ", c, "b:", b)
}

func getInfo(m map[string]string) (int, int) {
	point := (**hmap)(unsafe.Pointer(&m))
	value := *point
	return value.count, int(value.B)
}

func mapdebug() {
	d := datastruct.NewDict(&datastruct.DemoDictFuncs{})
	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond)
		d.Set(strconv.Itoa(i), strconv.Itoa(i))
	}
}
