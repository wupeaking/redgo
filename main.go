package main

import (
	"flag"

	"fmt"
	"time"

	"runtime"

	log "github.com/wupeaking/logrus"
	"github.com/wupeaking/redgo/server"
)

var (
	version   = flag.Bool("version", false, "show version")
	configire = flag.String("configure", "./config.yaml", "configure file path")
	loglevel  = flag.String("loglevel", "error", "set log level")
)

var (
	// GitCommit git版本号
	GitCommit string
	// Branch 分支名称
	Branch string
)

func main() {
	flag.Parse()
	if *version {
		println("commit: ", GitCommit, " branch: ", Branch)
		return
	}
	setLogLevel()

	log.SetFormatter(&LogFormat{})

	err := server.StartServer(*configire)
	log.Error("start server faild: ", err)
}

// LogFormat 自定义日志格式
type LogFormat struct {
}

// Format 实现Formatter接口
func (format *LogFormat) Format(entry *log.Entry) ([]byte, error) {
	_, file, line, _ := runtime.Caller(5)
	info := fmt.Sprintf("[%s] (%s@%d level=%s) %s\n", time.Now().Format("2006-01-02 15:04:05"),
		file, line, entry.Level.String(), entry.Message)

	return []byte(info), nil
}

func setLogLevel() {
	switch *loglevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}
}
