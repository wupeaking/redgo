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
)

//Map 构造map实现
type Map struct {
	refrash int // 当为-1 时 没有进行值的刷新 当为0的时候正在刷新 当为1时表示所有的值已经刷新到value2中 还没转到value1
	value1  map[string]interface{}
	value2  map[string]interface{}
	sync.RWMutex
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

// Set 设置一个键值
func (m *Map) Set(key string, value interface{}) {
	m.RLock()
	defer m.RUnlock()
	if m.refrash == -1 {
		m.value1[key] = value
	}
	m.value2[key] = value
}

// Get 获取一个键的内容
func (m *Map) Get(key string) (interface{}, bool) {
	// 	首先从value1获取内容
	var value interface{} = struct{}{}
	ok := false

	if m.value1 != nil {
		value, ok = m.value1[key]
	}
	if ok {
		return value, ok
	}
	// 执行此处 说明在value1中没有找到
	if m.value2 != nil {
		value, ok = m.value2[key]
	}
	return value, ok
}

// Len 获取当前键值个数
func (m *Map) Len() int {
	if m.refrash <= 0 {
		return len(m.value1)
	}
	return len(m.value2)
}

//Cap 获取当前的map的cap
func (m *Map) Cap() int {
	var value *hmap
	if m.refrash <= 0 {
		point := (**hmap)(unsafe.Pointer(&(m.value1)))
		value = *point
	} else {
		point := (**hmap)(unsafe.Pointer(&(m.value2)))
		value = *point
	}
	return int(math.Pow(float64(2), float64(value.B)))
}

//StartReduce 开始进行缩容操作
func (m *Map) StartReduce() {
	m.RLock()
	defer m.RUnlock()
	if m.refrash != -1 {
		// 说明已经在进行缩容  直接返回
		return
	}
	m.refrash = 0
	// 获取value1的值
	m.value2 = make(map[string]interface{}, len(m.value1))
	for k := range m.value1 {
		m.value2[k] = m.value1[k]
	}
	m.value1 = m.value2
	m.refrash = 1

}
