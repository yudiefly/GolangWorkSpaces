package v1

import (
	"gin-blog/models"
	"gin-blog/pkg/e"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"net/http"

	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
)

/*
GetTags的相关注释：
1、c.Query可用于获取?name=test&state=1这类URL参数，而c.DefaultQuery则支持设置一个默认值
2、code变量使用了e模块的错误编码，这正是先前规划好的错误码，方便排错和识别记录
3、util.GetPage保证了各接口的page处理是一致的
4、c *gin.Context是Gin很重要的组成部分，可以理解为上下文，它允许我们在中间件之间传递变量、管理流、验证请求的JSON和呈现JSON响应
【在本机执行curl 127.0.0.1:8000/api/v1/tags，正确的返回值为{"code":200,"data":{"lists":[],"total":0},"msg":"ok"}，若存在问题请结合gin结果进行拍错。】
【在获取标签列表接口中，我们可以根据name、state、page来筛选查询条件，分页的步长可通过app.ini进行配置，以lists、total的组合返回达到分页效果。】
*/

//获取多个文章标签
func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state, _ = com.StrTo(arg).Int()
		maps["state"] = state
	}
	code := e.SUCCESS

	data["list"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//新增文章标签
func AddTag(c *gin.Context) {

}

//修改文章标签
func EditTag(c *gin.Context) {

}

//删除文章标签
func DeleteTag(c *gin.Context) {

}
