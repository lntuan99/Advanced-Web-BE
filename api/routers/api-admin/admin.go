package api_admin

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/api/methods"
	"github.com/gin-gonic/gin"
)

func HandlerLoginAdminAccount(c *gin.Context) {
	success, status, data := methods.MethodLoginAdminUser(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerGetListAdminUser(c *gin.Context) {
	success, status, data := methods.MethodGetListAdminUser(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerGetAdminUserByID(c *gin.Context) {
	success, status, data := methods.MethodGetAdminUserByID(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerCreateAdminUser(c *gin.Context) {
	success, status, data := methods.MethodCreateAdminUser(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}
