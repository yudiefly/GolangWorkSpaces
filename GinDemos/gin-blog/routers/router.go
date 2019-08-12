package routers

import (

	// _ "gin-blog/docs"
	// "github.com/swaggo/gin-swagger"
	// "github.com/swaggo/gin-swagger/swaggerFiles"

	"gin-blog/middleware/jwt"
	"gin-blog/pkg/qrcode"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/upload"
	"gin-blog/routers/api"
	"gin-blog/routers/api/v1"

	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	// gin.SetMode(setting.RunMode)
	gin.SetMode(setting.ServerSetting.RunMode)

	// //注册swagger支持
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//注册路由（认证/获取Token）
	r.GET("/auth", api.GetAuth)
	//实现图片文件的访问
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	//上传图片
	r.POST("/upload", api.UploadImage)

	//注册路由（标签管理、文章管理）
	apiv1 := r.Group("/api/v1")

	//引入中间件
	apiv1.Use(jwt.JWT())
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
		//生成二维码图片
		apiv1.POST("/articles/poster/generate", v1.GetnerateArticlePoster)

	}
	return r
}
