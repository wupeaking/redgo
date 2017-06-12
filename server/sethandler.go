package server

import (
	"errors"

	log "github.com/wupeaking/logrus"
	"github.com/wupeaking/redgo/datastruct"
)

// set类型相关操作

// 增加相关操作

// SAdd 增加一个元素
func (myself *SrvHandler) SAdd(key string, members [][]byte) (int, error) {
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != SET {
			return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		count := 0
		for _, member := range members {
			if _, ok := hashValue.Get(string(member)); ok {
				continue
			}
			hashValue.Set(string(member), nil)
			count++
		}
		return count, nil
	}
	// 执行到此处说明不存在
	newHash := datastruct.NewMap()
	for _, member := range members {
		newHash.Set(string(member), nil)
	}

	newV := &Value{value: newHash, valueType: SET}
	myself.db.data.Set(key, newV)
	return len(members), nil
}

// 读取相关操作....

// SMembers 返回集合的所有元素
func (myself *SrvHandler) SMembers(key string) ([][]byte, error) {
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != SET {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		return hashValue.Keys(".")
	}
	// 执行到此处说明不存在

	return nil, nil
}

// SIsMember 测试是否存在该元素
func (myself *SrvHandler) SIsMember(key string, member []byte) (int, error) {
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != SET {
			return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		if _, ok := hashValue.Get(string(member)); ok {
			return 1, nil
		}
		return 0, nil
	}
	// 执行到此处说明不存在

	return 0, nil
}

// SRandMember 随机返回一个
func (myself *SrvHandler) SRandMember(key string) ([]byte, error) {
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != SET {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		return hashValue.RandomKeys(key)
	}
	// 执行到此处说明不存在

	return nil, nil
}

// SCard 返回集合中总的元素数量
func (myself *SrvHandler) SCard(key string) (int, error) {
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != SET {
			return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		return int(hashValue.Len()), nil
	}
	// 执行到此处说明不存在
	return 0, nil
}

// 集合的运算操作...

// SInter 获取集合的交集
func (myself *SrvHandler) SInter(key string, keys [][]byte) ([][]byte, error) {
	var minLenValue *datastruct.Map
	var minLen int64
	values := make([]*datastruct.Map, 0, 10)

	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if !ok {
		return nil, nil
	}
	comValue := v.(*Value)
	if comValue.valueType != SET {
		return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	hashValue := comValue.value.(*datastruct.Map)
	minLenValue = hashValue
	minLen = hashValue.Len()
	values = append(values, hashValue)

	// 迭代剩下的key 求出成员最小的那个key
	for _, k := range keys {
		v, ok := myself.db.data.Get(string(k))
		// 判断是否存在
		if !ok {
			return nil, nil
		}
		comValue := v.(*Value)
		if comValue.valueType != SET {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		length := hashValue.Len()
		if length < minLen {
			minLen = length
			minLenValue = hashValue
		}
		values = append(values, hashValue)
	}

	result := make([][]byte, 0, 10)
	// 根据最小值 来进行求交集
	members, e := minLenValue.Keys(".")
	if e != nil {
		return nil, e
	}
	log.Debug("len(values): ", len(values))
	for _, member := range members {
		var ok bool
		log.Debug("member: ", string(member))
		for _, value := range values {
			_, ok = value.Get(string(member))
			if !ok {
				break
			}
		}
		if ok {
			ok = false
			result = append(result, member)
		}
	}
	return result, nil
}

// SInterStore 将交集存储到新的key
func (myself *SrvHandler) SInterStore(newkey string, key1 []byte, keys [][]byte) (int, error) {
	members, e := myself.SInter(string(key1), keys)
	if e != nil {
		return 0, e
	}
	if len(members) == 0 {
		return 0, nil
	}

	newHash := datastruct.NewMap()
	for _, member := range members {
		newHash.Set(string(member), nil)
	}

	newV := &Value{value: newHash, valueType: SET}
	myself.db.data.Set(newkey, newV)
	return len(members), nil
}

// SDiff 比较多个集合的差集
func (myself *SrvHandler) SDiff(key string, keys [][]byte) ([][]byte, error) {

	values := make([]*datastruct.Map, 0, 10)

	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if !ok {
		return nil, nil
	}
	comValue := v.(*Value)
	if comValue.valueType != SET {
		return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	hashValue := comValue.value.(*datastruct.Map)

	// 迭代剩下的key
	for _, k := range keys {
		v, ok := myself.db.data.Get(string(k))
		// 判断是否存在
		if !ok {
			continue
		}
		comValue := v.(*Value)
		if comValue.valueType != SET {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hash := comValue.value.(*datastruct.Map)
		values = append(values, hash)
	}

	result := make([][]byte, 0, 10)
	members, e := hashValue.Keys(".")
	if e != nil {
		return nil, e
	}

	for _, member := range members {
		for _, value := range values {
			_, ok := value.Get(string(member))
			// 属于A且不属于B的才能是差集
			if ok {
				continue
			}
			// 执行到此处 说明当前member 在索引的集合中都不存在
			result = append(result, member)
		}
	}
	return result, nil
}

// SDiffStore 多个集合的差集保存为新值
func (myself *SrvHandler) SDiffStore(newkey string, key []byte, keys [][]byte) (int, error) {
	members, e := myself.SDiff(string(key), keys)
	if e != nil {
		return 0, e
	}
	// 如果成员为0 则直接不保存  不浪费空间
	if len(members) == 0 {
		return 0, nil
	}
	newHash := datastruct.NewMap()
	for _, member := range members {
		newHash.Set(string(member), nil)
	}

	newV := &Value{value: newHash, valueType: SET}
	myself.db.data.Set(newkey, newV)
	return len(members), nil
}

// SUnion 返回集合的并集
func (myself *SrvHandler) SUnion(key string, keys [][]byte) ([][]byte, error) {

	values := make([]*datastruct.Map, 0, 10)

	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if !ok {
		return nil, nil
	}
	comValue := v.(*Value)
	if comValue.valueType != SET {
		return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	hashValue := comValue.value.(*datastruct.Map)

	// 迭代剩下的key
	for _, k := range keys {
		v, ok := myself.db.data.Get(string(k))
		// 判断是否存在
		if !ok {
			continue
		}
		comValue := v.(*Value)
		if comValue.valueType != SET {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hash := comValue.value.(*datastruct.Map)
		values = append(values, hash)
	}

	// 创建一个新的map
	newMap := make(map[string]bool, 10)
	members, e := hashValue.Keys(".")
	if e != nil {
		return nil, e
	}

	for _, member := range members {
		newMap[string(member)] = true
	}

	for _, value := range values {
		members, e := value.Keys(".")
		if e != nil {
			return nil, e
		}
		for _, member := range members {
			newMap[string(member)] = true
		}
	}
	result := make([][]byte, 0, 100)
	for key := range newMap {
		result = append(result, []byte(key))
	}
	return result, nil
}

// SUnionStore 多个集合合并成一个新的
func (myself *SrvHandler) SUnionStore(newkey string, key []byte, keys [][]byte) (int, error) {
	members, e := myself.SUnion(string(key), keys)
	if e != nil {
		return 0, nil
	}
	// 如果成员为0 则直接不保存  不浪费空间
	if len(members) == 0 {
		return 0, nil
	}
	newHash := datastruct.NewMap()
	for _, member := range members {
		newHash.Set(string(member), nil)
	}

	newV := &Value{value: newHash, valueType: SET}
	myself.db.data.Set(newkey, newV)
	return len(members), nil
}

// 删除相关操作

// SRem 删除集合中的某些元素 可以一次删除多个 返回删除成功的个数 如果该集合不存在 返回0
func (myself *SrvHandler) SRem(key string, members [][]byte) (int, error) {
	//
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if !ok {
		return 0, nil
	}
	comValue := v.(*Value)
	if comValue.valueType != SET {
		return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	hashValue := comValue.value.(*datastruct.Map)
	count := 0
	for _, member := range members {
		mk := string(member)
		_, ok := hashValue.Get(mk)
		if ok {
			hashValue.Delete(mk)
			count++
		}
	}
	return count, nil
}

// SPop 随机删除一个元素 并返回删除的元素 如果集合已经没有元素或者该集合不存在 返回nil
func (myself *SrvHandler) SPop(key string) ([]byte, error) {
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if !ok {
		return nil, nil
	}
	comValue := v.(*Value)
	if comValue.valueType != SET {
		return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	hashValue := comValue.value.(*datastruct.Map)
	keyb, e := hashValue.RandomKeys(".")
	if e != nil {
		return nil, e
	}
	if keyb != nil {
		hashValue.Delete(string(keyb))
	}
	return keyb, nil
}

// SMove 将第一个集合中的元素移到另一个集合中去 如果操作成功返回1 失败返回0 如果目的集合不存在则创建
func (myself *SrvHandler) SMove(keySrc string, keyDst []byte) (int, error) {
	v, ok := myself.db.data.Get(keySrc)
	// 判断是否存在
	if !ok {
		return 0, nil
	}
	comValue := v.(*Value)
	if comValue.valueType != SET {
		return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	hashValue := comValue.value.(*datastruct.Map)
	memb, e := hashValue.RandomKeys(".")
	if e != nil {
		return 0, e
	}
	if memb == nil {
		return 0, nil
	}
	hashValue.Delete(string(memb))

	vdst, ok := myself.db.data.Get(string(keyDst))
	if ok {
		comValue := vdst.(*Value)
		if comValue.valueType != SET {
			return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hash := comValue.value.(*datastruct.Map)
		hash.Set(string(memb), nil)
		return 1, nil
	}
	// 需要新建set
	newHash := datastruct.NewMap()
	newHash.Set(string(memb), nil)
	newV := &Value{value: newHash, valueType: SET}
	myself.db.data.Set(string(keyDst), newV)
	return 1, nil
}
