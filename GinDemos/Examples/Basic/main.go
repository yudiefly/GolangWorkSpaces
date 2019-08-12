package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

/*
	gin获取从Post请求中的JSON参数
	示例：http://127.0.0.1:8080/login
	参数（Body）：
		{
		"user":"yudiefly",
		"password":"1234",
		"data":{
				"number":12,
				"rows":1
			   }
		}
	返回：
		{
		    "data": {
		        "Number": 12,
		        "rows": 1
		    },
		    "status": "you are logged in"
		}

*/

type LoginForm struct {
	User     string         `form:"user" binding:"required"`
	Password string         `form:"password" binding:"required"`
	Data     UserCenterData `form:"data" json:"data"`
}

type UserCenterData struct {
	Number int `josn:"number"`
	Rows   int `json:"rows"`
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	//Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/login", func(c *gin.Context) {
		var form LoginForm

		if c.ShouldBindJSON(&form) == nil {
			//c.ShouldBind(&form)
			if form.User == "yudiefly" && form.Password == "1234" {
				c.JSON(200, gin.H{
					"status": "you are logged in",
					"data":   form.Data,
				})
			} else {
				c.JSON(401, gin.H{"status": "unauthorized"})
			}
		}
	})

	//Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})
	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", //user:foo password：bar
		"manu": "123", //user:manu password:123
	}))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		//Parse JOSN
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})
	return r

}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":9090")
}
