package api_admin_classroom

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/api/methods"
	"github.com/gin-gonic/gin"
)

func HandlerAdminGetListClassroom(c *gin.Context) {
	success, status, data := methods.MethodAdminGetListClassroom(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerAdminGetClassroomByID(c *gin.Context) {
	success, status, data := methods.MethodAdminGetClassroomByID(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}
