package api_admin_user

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/api/methods"
	"github.com/gin-gonic/gin"
)

func HandlerGetListUser(c *gin.Context) {
	success, status, data := methods.MethodGetListUser(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerAdminGetUserByID(c *gin.Context) {
	success, status, data := methods.MethodAdminGetUserByID(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerAdminBanUserByID(c *gin.Context) {
	success, status, data := methods.MethodAdminBanUserByID(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerMapStudentCode(c *gin.Context) {
	success, status, data := methods.MethodMapStudentCode(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}
