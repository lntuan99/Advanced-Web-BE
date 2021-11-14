package api_user

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/api/methods"
	"github.com/gin-gonic/gin"
)

func HandlerGetListClassroomByJWTType(c *gin.Context) {
	success, status, data := methods.GetListClassroomByJWTType(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerGetListClassroomOwnedByUser(c *gin.Context) {
	success, status, data := methods.GetListClassroomOwnedByUser(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}