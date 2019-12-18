package main

import (
	"mylogger_test/mylogger"
)

var log mylogger.Logger

func main() {
	log = mylogger.NewConsoleLogger("info")
	log = mylogger.NewFileLogger("Info", "./", "zhoulinwan.log", 10*1024*1024) //文件日志

	for {
		log.Debug("这是一条Debug日志")
		log.Info("这是一条Info日志")
		id := 10010
		name := "理想"
		log.Error("这是一条Error日志,id:%d,name:%s", id, name)
		log.Warning("这是一条Warning日志")
		log.Fatal("这是一条Fatal日志")
		// time.Sleep(time.Second * 2)
	}
}
