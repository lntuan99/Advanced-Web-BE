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
