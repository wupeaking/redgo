package server

// 实现redis的服务端协议
import (
	"io/ioutil"
	"os"

	"github.com/wupeaking/redgo/datastruct"

	redis "github.com/wupeaking/go-redis-server"
	yaml "gopkg.in/yaml.v2"
)

//SrvHandler 处理客户端请求
type SrvHandler struct {
	Host string // 监听的主机地址
	Port int    // 端口号
	db   *DataBase
	// 其他配置待续 比如保存策略 是否守护进程运行。。。
}

// StartServer 启动redis服务
func StartServer(configFile string) error {
	handler := new(SrvHandler)
	e := loadConfigFile(configFile, handler)
	if e != nil {
		return e
	}
	server, _ := redis.NewServer(redis.DefaultConfig().Host(handler.Host).Port(handler.Port).Handler(handler))
	return server.ListenAndServe()
}

// loadConfigFile 加载配置文件
func loadConfigFile(configFile string, handler *SrvHandler) error {
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	configure := make(map[interface{}]interface{})

	err = yaml.Unmarshal(content, &configure)
	if err != nil {
		return err
	}
	// host
	hosti, ok := configure["host"]
	if !ok {
		handler.Host = "0.0.0.0"
	} else {
		handler.Host = hosti.(string)
	}
	// port
	porti, ok := configure["port"]
	if !ok {
		handler.Port = 6379
	} else {
		handler.Port = porti.(int)
	}
	handler.db = CreateDataBase()
	return nil
}

// DataBase 数据类型
type DataBase struct {
	data *datastruct.Map
}

// Value 数据库的值
type Value struct {
	value     interface{}
	valueType string
}

// 定义数据类型
const (
	STRING = "string"
	LIST   = "list"
	HASH   = "hash"
	SET    = "set"
	ZSET   = "zset"
)

// CreateDataBase 创建一个数据库
func CreateDataBase() *DataBase {
	db := new(DataBase)
	db.data = datastruct.NewMap()
	return db
}
