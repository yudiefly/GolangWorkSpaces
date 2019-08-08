package util

import (
	"gin-blog/pkg/setting"

	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
)

//分页的方法
func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize //setting.PageSize
	}
	return result
}
