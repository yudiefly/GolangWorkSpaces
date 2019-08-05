package main

import (
	"fmt"

	"gin-blog/pkg/setting"
	"gin-blog/routers"
	"net/http"
)

func main() {

	//注册路由
	//router:=gin.Default()
	//router.GET("/test",func(c *gin.Context){
	//	c.JSON(200,gin.H{
	//		"message":"test"
	//	})
	//})

	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
