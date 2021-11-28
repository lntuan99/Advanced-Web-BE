package api_grade

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/api/methods"
	"github.com/gin-gonic/gin"
)

func HandlerGetListGradeByClassroomId(c *gin.Context) {
	success, status, data := methods.MethodGetListGradeByClassroomId(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}
func HandlerCreateGrade(c *gin.Context) {
	success, status, data := methods.MethodCreateGrade(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}
func HandlerUpdateGrade(c *gin.Context) {
	success, status, data := methods.MethodUpdateGrade(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}
func HandlerDeleteGrade(c *gin.Context) {
	success, status, data := methods.MethodDeleteGrade(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}
