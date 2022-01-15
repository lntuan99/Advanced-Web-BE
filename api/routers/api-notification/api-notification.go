package api_notification

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/api/methods"
	"github.com/gin-gonic/gin"
)

func HandlerGetListNotification(c *gin.Context) {
	success, status, data := methods.MethodGetNotificationList(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerMarkReadNotification(c *gin.Context) {
	success, status, data := methods.MethodMarkReadNotification(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}
