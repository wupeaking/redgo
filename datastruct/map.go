/*Package datastruct 这个包是用于构造map结构，在golang里面已经实现了map这个基本类型。
基本的map实现还是比较高效的。在1.8的版本中map的负载因子是6.5.当超过这个负载因子时，map就是进行
扩容，但是目前遇到的问题是，golang没有对删除的键进行缩容。也即是就算是当前的map没有键也不会被gc掉
因此这一块内容需要自己扩展
现在的思路是 底层依然使用go的map，当负载因子减小到一定程度时 新建map将值拷贝过来 然后释放之前的map
*/
package datastruct

import (
	"math"
	"sync"
	"unsafe"

	log "github.com/wupeaking/logrus"
)

// 最小负载因子的阈值
var minMapFactor = 0.3

// 每次转移的键值个数
var countPerRerash = 100

//Map 构造map实现
type Map struct {
	refrash int64 // 当为-1 时 没有进行值的刷新 当为0的时候准备刷新
	value1  map[string]interface{}
	value2  map[string]interface{}
	sync.RWMutex
}

// NewMap 新建map对象
func NewMap() *Map {
	myself := new(Map)
	myself.refrash = -1
	myself.value1 = make(map[string]interface{}, 10)
	myself.value2 = nil
	return myself
}

// Set 设置一个键值
func (myself *Map) Set(key string, value interface{}) {
	myself.RLock()
	defer myself.RUnlock()

	if myself.needRefrash() {
		myself.startRefrash()
		myself.value2[key] = value
		return
	}
	myself.value1[key] = value
}

// Get 获取一个键的内容
func (myself *Map) Get(key string) (interface{}, bool) {
	// 	首先从value1获取内容
	var value interface{} = struct{}{}
	ok := false
	// 判断是否需要needrefrash
	if !myself.needRefrash() {
		value, ok = myself.value1[key]
		return value, ok
	}

	myself.startRefrash()
	if myself.value1 != nil {
		value, ok = myself.value1[key]
	}
	if ok {
		return value, ok
	}
	// 执行此处 说明在value1中没有找到
	if myself.value2 != nil {
		value, ok = myself.value2[key]
	}
	return value, ok
}

// Delete 删除一个key
func (myself *Map) Delete(key string) {
	myself.RLock()
	defer myself.RUnlock()
	if !myself.needRefrash() {
		delete(myself.value1, key)
		return
	}
	myself.startRefrash()
	delete(myself.value1, key)
	if myself.value2 != nil {
		delete(myself.value2, key)
	}
}

// Len 获取当前键值个数
func (myself *Map) Len() int64 {
	return int64(len(myself.value1) + len(myself.value2))
}

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

//Cap 获取当前的map的cap
func (myself *Map) Cap() int64 {
	var value *hmap
	point := (**hmap)(unsafe.Pointer(&(myself.value1)))
	value = *point
	return int64(math.Pow(float64(2), float64(value.B)))
}

// needRefrash 确认是否需要refrash
func (myself *Map) needRefrash() bool {
	// 1. 如果cap小于100 不需要进行refrash
	// 2. 如果负载率小于0.2 转移
	// 3. 如果refrash值不为负 所说明已经开始refrash
	myself.RLock()
	defer myself.RUnlock()
	if myself.refrash != -1 {
		return true
	}
	caption := myself.Cap()
	if caption < 100 {
		return false
	}
	// 计算负载率
	loadFac := float64(len(myself.value1)) * 1.0 / float64(caption)
	if loadFac < minMapFactor {
		myself.refrash = 0
		return true
	}
	return false
}

// startRefrash 开始rehash
func (myself *Map) startRefrash() {
	myself.RLock()
	defer myself.RUnlock()

	// 判断是否是否是刚刚开始refrash
	if myself.refrash == 0 {
		myself.value2 = make(map[string]interface{}, len(myself.value1))
		log.Debug("需要开始刷新，第一次刷新, length:", len(myself.value1))
	}
	length := len(myself.value1)
	// 判断是否rehash完成
	if length == 0 {
		log.Debug("已经刷新完成")
		// 如果完成 将value2 转到value1
		myself.value1 = myself.value2
		myself.value2 = nil
		myself.refrash = -1
		return
	}
	// 执行到此处 说明没有rehash完 每次转移100个键值
	sum := 0
	if length > countPerRerash {
		sum = countPerRerash
	} else {
		sum = length
	}
	i := 0
	for key := range myself.value1 {
		myself.value2[key] = myself.value1[key]
		delete(myself.value1, key)
		i++
		if i > sum {
			break
		}
	}
	log.Debug("进行了一次刷新，刷新个数: ", sum)
	myself.refrash++
	return
}
