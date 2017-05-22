package datastruct

/*
* 当前文档 是用于示例如何实现获取一个golang原生的map的容量
* 在golang的map设计中，当负载因子达到6.25时会进行再次的扩容
* 操作，但是没有实现缩容操作，而且当我们删除map的键值时，内存是
* 不会释放的。
 */
// 下面是一个示例 验证在不停gc的过程中map是不会释放删除的键值的内存，同时展示如何获取map的容量

/*
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

func demo() {
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
*/
