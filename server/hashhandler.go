package server

import (
	"errors"

	"github.com/wupeaking/redgo/datastruct"
)

// hash类型的相关操作

// 设置hash相关的操作......

// HSet 设置某个hash的某个字段的内容 如果不存在则创建
func (myself *SrvHandler) HSet(key string, field []byte, value []byte) (int, error) {
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != HASH {
			return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		hashValue.Set(string(field), value)
		return 1, nil
	}
	// 执行到此处说明不存在
	newHash := datastruct.NewMap()
	newHash.Set(string(field), value)
	newV := &Value{value: newHash, valueType: HASH}
	myself.db.data.Set(key, newV)
	return 1, nil
}

// HSetNx 设置一个hash的某个字段的内容 如果存在则不设置 如果key不存在 则肯定设置 如果key存在 field不存在也设置
func (myself *SrvHandler) HSetNx(key string, field []byte, value []byte) (int, error) {
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != HASH {
			return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		_, ok := hashValue.Get(string(field))
		if ok {
			return 0, nil
		}
		hashValue.Set(string(field), value)
		return 1, nil
	}
	// 如果key不存在 则放心大胆的设置
	newHash := datastruct.NewMap()
	newHash.Set(string(field), value)
	newV := &Value{value: newHash, valueType: HASH}
	myself.db.data.Set(key, newV)
	return 1, nil
}

// HMSet 设置某个hash的读个字段值
func (myself *SrvHandler) HMSet(key string, fieldValues [][]byte) error {
	// 检验参数
	//log.Debug("intser hmset.....", "key: ", key)
	length := len(fieldValues)
	if length%2 != 0 {
		return errors.New("ERR wrong number of arguments for 'hmset' command")
	}
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != HASH {
			return errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		for i := 0; i < length; i += 2 {
			field := fieldValues[i]
			value := fieldValues[i+1]
			//log.Error("hmset: ", "field: ", string(field), " value: ", string(value))
			hashValue.Set(string(field), value)
		}
		return nil
	}
	newHash := datastruct.NewMap()
	for i := 0; i < length; i += 2 {
		field := fieldValues[i]
		value := fieldValues[i+1]
		//log.Error("hmset: ", "field: ", string(field), " value: ", string(value))
		newHash.Set(string(field), value)
	}
	newV := &Value{value: newHash, valueType: HASH}
	myself.db.data.Set(key, newV)
	return nil
}

// 获取hash相关操作....

// HGet 获取hash的指定字段的内容
func (myself *SrvHandler) HGet(key string, field []byte) ([]byte, error) {
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != HASH {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		value, ok := hashValue.Get(string(field))
		if ok {
			return value.([]byte), nil
		}
		return nil, nil
	}
	return nil, nil
}

//HMGet 获取hash多个字段的内容
func (myself *SrvHandler) HMGet(key string, fields []byte) ([][]byte, error) {
	result := make([][]byte, 0, 10)
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != HASH {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		for _, field := range fields {
			value, ok := hashValue.Get(string(field))
			if ok {
				result = append(result, value.([]byte))
			}
		}

		return result, nil
	}
	return nil, nil
}

// HExists 检查hash的某个字段是否存在
func (myself *SrvHandler) HExists(key string, field []byte) (int, error) {

	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != HASH {
			return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		_, ok := hashValue.Get(string(field))
		if ok {
			return 1, nil
		}
		return 0, nil
	}
	return 0, nil
}

// HKeys 返回hash的所有字段
func (myself *SrvHandler) HKeys(key string) ([][]byte, error) {
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != HASH {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		return hashValue.Keys(".")
	}
	return nil, nil
}

// HVals 返回hash的所有字段
func (myself *SrvHandler) HVals(key string) ([][]byte, error) {
	v, ok := myself.db.data.Get(key)
	values := make([][]byte, 0, 10)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != HASH {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		fields, e := hashValue.Keys(".")
		if e != nil {
			return nil, e
		}

		for _, filed := range fields {
			value, _ := hashValue.Get(string(filed))
			values = append(values, value.([]byte))
		}
		return values, nil
	}
	return nil, nil
}

// HGetAll 返回hash的所有字段和值
func (myself *SrvHandler) HGetAll(key string) ([][]byte, error) {
	v, ok := myself.db.data.Get(key)
	filedValues := make([][]byte, 0, 10)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != HASH {
			return nil, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		fields, e := hashValue.Keys(".")
		if e != nil {
			return nil, e
		}

		for _, filed := range fields {
			filedValues = append(filedValues, filed)
			value, _ := hashValue.Get(string(filed))
			filedValues = append(filedValues, value.([]byte))
		}
		return filedValues, nil
	}
	return nil, nil
}

// HLen 返回hash字段的数量
func (myself *SrvHandler) HLen(key string) (int, error) {
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != HASH {
			return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		return int(hashValue.Len()), nil
	}
	return 0, nil
}

// 删除相关操作

// HDel 删除指定的字段
func (myself *SrvHandler) HDel(key string, fields [][]byte) (int, error) {
	v, ok := myself.db.data.Get(key)
	// 判断是否存在
	if ok {
		comValue := v.(*Value)
		if comValue.valueType != HASH {
			return 0, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hashValue := comValue.value.(*datastruct.Map)
		for _, field := range fields {
			hashValue.Delete(string(field))
		}
		return len(fields), nil
	}
	return 0, nil
}
