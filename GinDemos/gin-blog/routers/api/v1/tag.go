package v1

import (
	"gin-blog/models"
	"gin-blog/pkg/e"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"log"
	"net/http"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
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

// @Summary 获取多个文章标签
// @Produce  json
// @Parameter  name query string false "Name"
// @Parameter  state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"系统错误"}"
// @Router /api/v1/tags [get]
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

	//list, err := models.GetTags(util.GetPage(c), setting.PageSize, maps)
	list, err := models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	if err != nil {
		log.Fatal(err)
	} else {
		data["list"] = list
	}

	count, err := models.GetTagTotal(maps)
	if err != nil {
		log.Fatal(err)
	} else {
		data["total"] = count
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

/*
 用Postman用POST访问http://127.0.0.1:8000/api/v1/tags?name=1&state=1&created_by=test，查看code是否返回200及blog_tag表中是否有值，有值则正确。
*/

// @Summary 新增文章标签
// @Produce  json
// @Parameter  name query string true "Name"
// @Parameter  state query int false "State"
// @Parameter  created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state, _ := com.StrTo(c.DefaultQuery("state", "0")).Int()
	createdBy := c.Query("created_by")

	//开始对参数进行验证
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.Required(name, "created_by").Message("创建人不能为空")
	valid.MaxSize(name, 100, "created_by").Message("创建人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

//修改文章标签
// @Summary 修改文章标签
// @Produce  json
// @Parameter  id param int true "ID"
// @Parameter  name query string true "ID"
// @Parameter  state query int false "State"
// @Parameter  modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	id, _ := com.StrTo(c.Param("id")).Int()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state, _ = com.StrTo(arg).Int()
		valid.Range(state, 0, 1, "state").Message("状态只允许是0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagById(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}
			models.EditTag(id, data)
		} else {
			for _, err := range valid.Errors {
				log.Println(err.Key, err.Message)
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// @Summary 删除文章标签
// @Produce  json
// @Parameter  id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"系统错误"}"
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	id, _ := com.StrTo(c.Param("id")).Int()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagById(id) {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}
