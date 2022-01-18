package api_status

import (
	"advanced-web.hcmus/api/base"
	"github.com/gin-gonic/gin"
)

func HandlerStatus(c *gin.Context) {
	base.ResponseResult(c, gin.H{"status": "hello"})
}
