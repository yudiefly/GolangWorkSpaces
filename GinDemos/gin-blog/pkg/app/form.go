package app

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"gin-blog/pkg/e"
)

func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, e.ERROR
	}

	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	return http.StatusOK, e.SUCCESS
}

/*
	获取Post方式提交的Body中的JSON对象，并验证之
*/
func BindPostJsonAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.ShouldBindJSON(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	// valid := validation.Validation{}
	// check, err := valid.Valid(form)
	// if err != nil {
	// 	return http.StatusInternalServerError, e.ERROR
	// }

	// if !check {
	// 	MarkErrors(valid.Errors)
	// 	return http.StatusBadRequest, e.INVALID_PARAMS
	// }

	return http.StatusOK, e.SUCCESS
}
