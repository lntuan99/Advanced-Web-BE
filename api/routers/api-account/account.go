package api_account

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/api/methods"
	"github.com/gin-gonic/gin"
)

func HandlerRegisterAccount(c *gin.Context) {
	success, status, data := methods.MethodRegisterAccount(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerLoginAccount(c *gin.Context) {
	success, status, data := methods.MethodLoginAccount(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerGoogleLogin(c *gin.Context) {
	success, status, data := methods.MethodGoogleLogin(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}