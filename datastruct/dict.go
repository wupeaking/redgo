package datastruct

// 注意这不是一个并发安全的结构

import (
	"fmt"

	murmur3 "github.com/spaolacci/murmur3"
	log "github.com/wupeaking/logrus"
)

var minFactor = float64(0.5)
var maxFactor = float64(3.0)

// Dict 字典的定义
type Dict struct {
	funcs  DictFuncType
	ht     [2]dictht // 保存2个hash表 用于rehash
	rehash uint64
}

// dictht 字典哈希表
type dictht struct {
	table    []*dictEntry
	size     uint64
	sizeMask uint64
	used     uint64
}

//dictEntry 字典的入口
type dictEntry struct {
	next  *dictEntry
	key   string
	value interface{}
}

//DictFuncType 字典处理相关的函数接口
type DictFuncType interface {
	// 计算hash值的函数
	HashCalc(key string) uint64
	// 复制键的函数
	KeyCopy(key string) string
	// 复制值的函数
	ValueCopy(v interface{}) interface{}
	// 对比键的函数
	KeyCompare(k1, k2 string) int
	// 销毁键的函数
	KeyDestructor(key string)
	// 销毁值的函数
	ValueDestructor(value interface{})
}

// NewDict 创建一个新的字典
func NewDict(funcs DictFuncType) *Dict {
	dict := Dict{rehash: 0}
	var size uint64 = 4
	dict.ht[0] = dictht{size: size, sizeMask: size - 1, used: 0, table: make([]*dictEntry, size, size)}
	dict.ht[1] = dictht{size: 4, sizeMask: 3, used: 0, table: make([]*dictEntry, 4, 4)}
	dict.funcs = funcs
	return &dict
}

// Set 增加一对键值
func (d *Dict) Set(key string, value interface{}) {
	// 首先计算hash值
	insert := func(i int, k string, v interface{}) {
		hash := d.funcs.HashCalc(k)
		offset := hash & d.ht[i].sizeMask
		entry := dictEntry{key: k, value: v}

		if d.ht[i].table[offset] == nil {
			d.ht[i].table[offset] = &entry
			entry.next = nil
		} else {
			// 说明已经在此表下有内容 需要判断是否有相同的key
			curEntry := d.ht[i].table[offset]
			headEntry := curEntry
			isExist := false
			for curEntry != nil {
				if curEntry.key == k {
					curEntry.value = value
					isExist = true
					break
				}
				curEntry = curEntry.next
				continue
			}
			if !isExist {
				entry.next = headEntry
				d.ht[i].table[offset] = &entry
			}
		}
		d.ht[i].used++
		return
	}

	// 判断是否需要rehash
	need := d.needRehash()
	if !need {
		log.Debug("不需要rehash ", "size: ", d.ht[0].size, " used: ", d.ht[0].used)
		// 计算偏移 todo:: 先不考虑rehash
		insert(0, key, value)
		return
	}
	// 执行此处 说明需要rehash
	// 首先判断是否已经在hash了
	if d.rehash > 0 {
		// 说明已经在hash了 判断是否已经hash完成
		if d.ht[0].used == 0 {
			log.Debug("全部已经rehash完成，将ht[1]的内容移到ht[0]", "h2.size: ", d.ht[1].size, "h1.used: ", d.ht[1].used)
			// 将ht[0]指向ht[1]
			d.ht[0].table = d.ht[1].table
			d.ht[0].size = d.ht[1].size
			d.ht[0].sizeMask = d.ht[1].sizeMask
			d.ht[0].used = d.ht[1].used
			d.rehash = 0
			d.ht[1].table = nil
			insert(0, key, value)
			return
		}
		// 执行到这里说明还没有hash完成 从table中取出一个值 放到ht[1]中
		log.Debug("开始rehash, rehash: ", d.rehash, "len(d.ht[0].table): ", len(d.ht[0].table), " used: ", d.ht[0].used)
		entry := d.ht[0].table[d.rehash-1]
		if entry == nil {
			// 说明表中已经没有值了 那就增加rehash值 等下回操作下一个table
			log.Debug("当前table没有内容了", d.rehash)
			d.rehash++
			//d.ht[0].used--
			insert(1, key, value)
			return
		}
		// 说明里面有内容 将内容挂到ht[1]里
		for entry != nil {
			log.Debug("读取table的内容", d.rehash)
			insert(1, entry.key, entry.value)
			d.ht[0].used--
			entry = entry.next
		}
		d.ht[0].table[d.rehash-1] = nil
		insert(1, key, value)
		d.rehash++
		return
	}
	// 执行到此处 说明需要hash 但是还没有进行 需要新建内容 rehash肯定是为0
	d.ht[1].table = make([]*dictEntry, d.ht[0].used, d.ht[0].used)
	d.ht[1].size = d.ht[0].used
	d.ht[1].sizeMask = d.ht[0].used - 1
	d.ht[1].used = 0
	log.Debug("需要rehash,但是还没有开始扩缩容 ", "sizeMask: ", d.ht[0].used-1, " size: ", d.ht[0].used, " rehash:", d.rehash)
	entry := d.ht[0].table[d.rehash]
	if entry == nil {
		// 说明表中已经没有值了 那就增加rehash值 等下回操作下一个table
		d.rehash++
		//d.ht[0].used--
		insert(1, key, value)
		return
	}
	// 说明里面有内容 将内容挂到ht[1]里
	for entry != nil {
		insert(1, entry.key, entry.value)
		d.ht[0].used--
		entry = entry.next
	}
	d.ht[0].table[d.rehash] = nil
	insert(1, key, value)
	d.rehash++
	return
}

// Get 返回一个值类型
func (d *Dict) Get(key string) (interface{}, bool) {
	hash := d.funcs.HashCalc(key)

	get := func(i int) (interface{}, bool) {
		offset := hash & d.ht[i].sizeMask
		log.Debug("offset:", offset, " key:", key, " sizeMask:", d.ht[i].sizeMask, "hash: ", hash)
		entry := d.ht[i].table[offset]
		for entry != nil {
			if entry.key != key {
				entry = entry.next
				continue
			}
			return entry.value, true
		}
		return nil, false
	}
	// 首先尝试从ht[0]获取
	v, ok := get(0)
	if ok {
		return v, ok
	}

	// 从ht[1]获取
	if d.ht[1].table == nil {
		return nil, false
	}
	v, ok = get(1)
	if ok {
		return v, ok
	}
	return nil, false
}

// Delete 删除一个键值
func (d *Dict) Delete(key string) {
	hash := d.funcs.HashCalc(key)
	// 计算偏移 todo:: 先不考虑rehash
	offset := hash & d.ht[0].sizeMask
	entry := d.ht[0].table[offset]
	prev := entry
	i := 0
	for entry != nil {
		if entry.key != key {
			prev = entry
			entry = entry.next
			i++
			continue
		}
		// 执行到此处 说明存在此key 直接不在引用 由gc释放
		if i == 0 {
			// 说明就是首节点
			d.ht[0].table[offset] = entry.next
		} else {
			prev.next = entry.next
		}
		d.ht[0].used--
		break
	}
}

// Free 释放字典
func (d *Dict) Free() {
	d.ht[0].table = nil
	d.ht[1].table = nil
}

// Print 打印key和value
func (d *Dict) Print() {
	for _, entry := range d.ht[0].table {
		for entry != nil {
			fmt.Println("key:", entry.key, " v:", entry.value)
			entry = entry.next
		}
	}
	if d.ht[1].table == nil {
		return
	}
	for _, entry := range d.ht[1].table {
		for entry != nil {
			fmt.Println("key:", entry.key, " v:", entry.value)
			entry = entry.next
		}
	}

}

// 判断是否需要重新hash
func (d *Dict) needRehash() bool {
	if d.rehash > 0 {
		// 说明已经开始了刷新 返回true
		log.Debug("开始rehash")
		return true
	}
	// 计算负载因子
	if d.ht[0].used < 10 {
		return false
	}
	loadFactor := float64(d.ht[0].used) / float64(d.ht[0].size)
	if loadFactor > maxFactor || loadFactor < minFactor {
		log.Debug("负载因子过大或过小,需要rehash", loadFactor)
		return true
	}
	return false
}

// GetRandomKey 获取随机的键值

// DemoDictFuncs 函数簇示例
type DemoDictFuncs struct {
}

// HashCalc 计算hash值
func (*DemoDictFuncs) HashCalc(key string) uint64 {
	return murmur3.Sum64([]byte(key))
}

// KeyCopy 返回键的复制
func (*DemoDictFuncs) KeyCopy(key string) string {
	return key
}

//ValueCopy 值拷贝
func (*DemoDictFuncs) ValueCopy(v interface{}) interface{} {
	return v
}

// KeyCompare 值比较
func (*DemoDictFuncs) KeyCompare(k1, k2 string) int {
	if k1 == k2 {
		return 0
	}
	if len(k1) > len(k2) {
		return 1
	}
	return -1
}

// KeyDestructor 键销毁
func (*DemoDictFuncs) KeyDestructor(key string) {
	return
}

// ValueDestructor 值销毁
func (*DemoDictFuncs) ValueDestructor(value interface{}) {
	return
}
