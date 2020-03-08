package main

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

type AppConf struct {
	KafkaConf   `ini:"kafka"`
	TaillogConf `ini:"taillog"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}
type TaillogConf struct {
	FileName string `ini:"filename"`
}

var (
	appConfig = new(AppConf)
)

func main() {
	// cfg, err := ini.Load("./config/conf.ini")
	// if err != nil {
	// 	fmt.Println("Failed to read File:%v", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(cfg.Section("kafka").Key("address").String())
	// fmt.Println(cfg.Section("kafka").Key("topic").String())
	// fmt.Println(cfg.Section("taillog").Key("filename").String())
	//---------------------------------------------------------------
	err := ini.MapTo(appConfig, "./config/conf.ini")
	if err != nil {
		fmt.Println("Failed to read File:%v", err)
		os.Exit(1)
	}
	fmt.Println(appConfig.KafkaConf.Address)
	fmt.Println(appConfig.KafkaConf.Topic)
	fmt.Println(appConfig.TaillogConf.FileName)

}
