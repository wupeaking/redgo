package server

import (
	"errors"

	"strconv"

	"github.com/wupeaking/redgo/datastruct"
)

// 符串类型的命令相关处理

func (myself *SrvHandler) setValue(key string, value []byte) {
	str := datastruct.NewSds(value)
	v := &Value{valueType: "string", value: str}
	myself.db.data.Set(key, v)
}

//Set 设置一个键的值
func (myself *SrvHandler) Set(key string, value []byte) error {
	str := datastruct.NewSds(value)
	v := &Value{valueType: "string", value: str}
	myself.db.data.Set(key, v)
	return nil
}

// SetNx 如果键值不存在 则设置
func (myself *SrvHandler) SetNx(key string, value []byte) (int, error) {
	_, ok := myself.db.data.Get(key)
	// 如果存在键 则不需要设置
	if ok {
		return 0, nil
	}
	str := datastruct.NewSds(value)
	myself.db.data.Set(key, &Value{valueType: "string", value: str})
	return 1, nil
}

// SetRange 设置某个键值的范围的内容
func (myself *SrvHandler) SetRange(key string, index int, value []byte) (int, error) {
	v, ok := myself.db.data.Get(key)
	// 如果不存在键 则不需要设置
	if !ok {
		str := datastruct.NewSds(value)
		myself.db.data.Set(key, &Value{valueType: "string", value: str})
		return str.SdsLen(), nil
	}
	// 如果存在 则更改
	comValue := v.(*Value)
	sdsValue := comValue.value.(*datastruct.Sds)
	sdsValue.SdsRange(0, index)
	sdsValue.SdsCat(value)
	return sdsValue.SdsLen(), nil
}

// Mset 多个键值设置
func (myself *SrvHandler) Mset(keyValues [][]byte) error {
	// 键值对应该成对出现
	if len(keyValues)%2 != 0 {
		return errors.New("key value should be pair")
	}
	for i, kv := range keyValues {
		if i%2 != 0 {
			continue
		}
		key := kv
		value := keyValues[i+1]
		myself.setValue(string(key), value)
	}
	return nil
}

// Msetnx 设置多个不存在的值
func (myself *SrvHandler) Msetnx(keyValues [][]byte) error {
	// 键值对应该成对出现
	if len(keyValues)%2 != 0 {
		return errors.New("key value should be pair")
	}
	for i, kv := range keyValues {
		if i%2 != 0 {
			continue
		}
		key := kv
		value := keyValues[i+1]
		_, e := myself.SetNx(string(key), value)
		if e != nil {
			return e
		}
	}
	return nil
}

// Get 获取键的值
func (myself *SrvHandler) Get(key string) ([]byte, error) {
	v, ok := myself.db.data.Get(key)
	// 如果存在键 则不需要设置
	if !ok {
		return nil, nil
	}
	// 如果存在 判定值类型
	comValue := v.(*Value)
	if comValue.valueType != "string" {
		return nil, errors.New("当前key不是string类型")
	}
	sdsValue := comValue.value.(*datastruct.Sds)
	return sdsValue.Buffer(), nil
}

// GetSet 设置值 并返回之前的旧值
func (myself *SrvHandler) GetSet(key string, value []byte) ([]byte, error) {
	result, e1 := myself.Get(key)
	e2 := myself.Set(key, value)
	if e2 != nil {
		return nil, e2
	}
	return result, e1
}

// Mget 获取多个键值
func (myself *SrvHandler) Mget(keys []string) ([][]byte, error) {
	var result [][]byte
	for _, key := range keys {
		v, e := myself.Get(key)
		if e != nil {
			return nil, e
		}
		result = append(result, v)
	}
	return result, nil
}

// GetRange 返回指定返回的子串
func (myself *SrvHandler) GetRange(key string, left int, right int) ([]byte, error) {
	if left > right || left < 0 {
		return nil, errors.New("start end must be a positive numberand start should Less than end")
	}
	v, ok := myself.db.data.Get(key)
	// 如果存在键 则不需要设置
	if !ok {
		return []byte(""), nil
	}
	// 如果存在 判定值类型
	comValue := v.(*Value)
	if comValue.valueType != "string" {
		return nil, errors.New("当前key不是string类型")
	}
	sdsValue := comValue.value.(*datastruct.Sds)
	if left > sdsValue.SdsLen() {
		return []byte(""), nil
	}

	if right > sdsValue.SdsLen() {
		right = sdsValue.SdsLen()
	}
	return sdsValue.Buffer()[left:right], nil
}

// Incr 自增
func (myself *SrvHandler) Incr(key string) (int, error) {
	//
	v, e := myself.Get(key)
	if e != nil {
		return 0, e
	}
	if v == nil {
		e = myself.Set(key, []byte("1"))
		return 1, e
	}
	// 执行到此处说明 存在key
	vi, e := strconv.Atoi(string(v))
	if e != nil {
		return 0, e
	}
	return vi + 1, myself.Set(key, []byte(strconv.Itoa(vi+1)))
}

// Decr 自减
func (myself *SrvHandler) Decr(key string) (int, error) {
	v, e := myself.Get(key)
	if e != nil {
		return 0, e
	}
	if v == nil {
		e = myself.Set(key, []byte("-1"))
		return -1, e
	}
	// 执行到此处说明 存在key
	vi, e := strconv.Atoi(string(v))
	if e != nil {
		return 0, errors.New("value canot transform int ")
	}
	return vi - 1, myself.Set(key, []byte(strconv.Itoa(vi-1)))
}

// Append 追加内容 如果不存在 则设置
func (myself *SrvHandler) Append(key string, value []byte) (int, error) {
	v, ok := myself.db.data.Get(key)
	// 不存在 直接设置
	if !ok {
		return len(value), myself.Set(key, value)
	}
	// 如果存在 判定值类型
	comValue := v.(*Value)
	if comValue.valueType != "string" {
		return 0, errors.New("当前key不是string类型")
	}
	sdsValue := comValue.value.(*datastruct.Sds)
	sdsValue.SdsCat(value)
	return sdsValue.SdsLen(), nil
}

// StrLen 返回key内容的长度 如果不存在 返回0
func (myself *SrvHandler) StrLen(key string) (int, error) {
	v, ok := myself.db.data.Get(key)
	// 不存在 直接设置
	if !ok {
		return 0, nil
	}
	// 如果存在 判定值类型
	comValue := v.(*Value)
	if comValue.valueType != "string" {
		return 0, errors.New("当前key不是string类型")
	}
	sdsValue := comValue.value.(*datastruct.Sds)
	return sdsValue.SdsLen(), nil
}
