package server

import (
	"errors"

	"fmt"

	log "github.com/wupeaking/logrus"
)

// 用于处理其他的非类型操作的相关命令

// Keys 返回所有的满足指定模式的键
func (myself *SrvHandler) Keys(pattern string) ([][]byte, error) {
	return myself.db.data.Keys(pattern)
}

// Exists 判断键值是否存在
func (myself *SrvHandler) Exists(key string, keys ...[]byte) (int, error) {
	count := 0

	_, ok := myself.db.data.Get(key)
	if ok {
		count++
	}
	for _, k := range keys {
		_, ok := myself.db.data.Get(string(k))
		if ok {
			count++
		}
	}
	return count, nil
}

// Type 判断一个键的类型 不存在返回none
func (myself *SrvHandler) Type(key string) ([]byte, error) {
	v, ok := myself.db.data.Get(key)
	if ok {
		return nil, nil
	}
	comValue := v.(*Value)

	return []byte(comValue.valueType), nil
}

// Rename 重命名一个键
func (myself *SrvHandler) Rename(key string, newkey []byte) error {
	v, ok := myself.db.data.Get(key)
	if ok {
		return errors.New("ERR no such key")
	}
	comValue := v.(*Value)
	myself.db.data.Delete(key)
	myself.db.data.Set(string(newkey), comValue)
	return nil
}

// Del 删除一个键
func (myself *SrvHandler) Del(key string, keys ...[]byte) (int, error) {
	count := 0
	_, ok := myself.db.data.Get(key)
	if ok {
		myself.db.data.Delete(key)
		count++
	}
	for _, k := range keys {
		_, ok := myself.db.data.Get(string(k))
		if ok {
			myself.db.data.Delete(string(k))
			count++
		}
	}
	return count, nil
}

// Move 将当前数据库的一个键移到另一个库中
func (myself *SrvHandler) Move(key string, db int) error {
	// todo::
	log.Error("暂时还不支持多库操作")
	return errors.New("暂时还不支持多库操作")
}

// Expire 设置一个key的过期时间 如果key不存在 返回0
func (myself *SrvHandler) Expire(key string) (int, error) {
	// todo:
	log.Error("暂时还不支持键ttl功能")
	return 0, errors.New("暂时还不支持键ttl功能")
}

// Persist 移除过期时间
func (myself *SrvHandler) Persist(key string) (int, error) {
	// todo:
	log.Error("暂时还不支持键ttl功能")
	return 0, errors.New("暂时还不支持键ttl功能")
}

// TTL 查看键值时间
func (myself *SrvHandler) TTL(key string) (int, error) {
	// todo:
	log.Error("暂时还不支持键ttl功能")
	return 0, errors.New("暂时还不支持键ttl功能")
}

// 服务器相关命令----

// Ping 测试连接
func (myself *SrvHandler) Ping() ([]byte, error) {
	return []byte("PONG"), nil
}

// Dbsize 当前库的所有key数量
func (myself *SrvHandler) Dbsize() (int, error) {
	return int(myself.db.data.Len()), nil
}

// Info 返回服务信息
func (myself *SrvHandler) Info() ([]byte, error) {
	info := fmt.Sprintf(`info:
	host: %s,
	port: %d,
	dbseize: %d,
	其他信息待续...
`, myself.Host, myself.Port, myself.db.data.Len())
	return []byte(info), nil
}

//FlushDb 删除数据库的所有key
func (myself *SrvHandler) FlushDb() error {
	log.Warn("start flush db...")
	myself.db.data.Free()
	return nil
}

// FlushAll 删除所有库的key
func (myself *SrvHandler) FlushAll() error {
	//todo:: 暂时和删除单库一样功能
	myself.db.data.Free()
	return nil
}
