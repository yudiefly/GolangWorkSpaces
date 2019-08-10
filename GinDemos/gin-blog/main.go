package main

import (
	"fmt"
	"gin-blog/pkg/setting"
	"gin-blog/routers"
	"log"
	"net/http"

	"gin-blog/pkg/gredis"
)

func init() {
	gredis.Setup()
}

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://github.com/EDDYCJY/go-gin-example
// @license.name MIT
// @license.url https://github.com/EDDYCJY/go-gin-example/blob/master/LICENSE
func main() {

	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	//注册路由
	router := routers.InitRouter()
	//注册http服务器
	s := &http.Server{
		Addr:           endPoint, //fmt.Sprintf(":%d", setting.ServerSetting.HttpPort), //setting.HTTPPort
		Handler:        router,
		ReadTimeout:    readTimeout,    //setting.ServerSetting.ReadTimeout,  //setting.ReadTimeout
		WriteTimeout:   writeTimeout,   //setting.ServerSetting.WriteTimeout, //setting.WriteTimeout
		MaxHeaderBytes: maxHeaderBytes, //1 << 20,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	//启动http server
	s.ListenAndServe()

}
