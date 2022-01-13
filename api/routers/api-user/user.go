package api_user

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/api/methods"
	"github.com/gin-gonic/gin"
)

func HandlerUpdateUserProfile(c *gin.Context) {
	success, status, data := methods.MethodUpdateUserProfile(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerGetUserProfile(c *gin.Context) {
	success, status, data := methods.MethodGetUserProfile(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerVerifyCode(c *gin.Context) {
	success, status, data := methods.MethodVerifyCode(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerForgotPassword(c *gin.Context) {
	success, status, data := methods.MethodForgotPassword(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerUpdatePassword(c *gin.Context) {
	success, status, data := methods.MethodUpdatePassword(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}
