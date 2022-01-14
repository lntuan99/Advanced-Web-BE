package api_jwt

import (
	"advanced-web.hcmus/model"
	"github.com/gin-gonic/gin"
	"time"
)

func NewGinJWTMiddleware() *GinJWTMiddleware {
	return &GinJWTMiddleware{
		Realm:            "production zone",
		SigningAlgorithm: "RS512",
		PrivKeyFile:      "zzz/key",
		PubKeyFile:       "zzz/key.pub",
		Timeout:          time.Hour * 24 * 7,
		MaxRefresh:       time.Hour,
		Authorizator: func(user interface{}, c *gin.Context) bool {
			if user == nil {
				return false
			}
			_, ok := user.(model.User)
			if ok {
				return true
			}

			_, ok = user.(model.AdminUser)
			if ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc:   time.Now,
		IsRequired: true,
	}
}

func DefaultGinJWTMiddleware() *GinJWTMiddleware {
	return NewGinJWTMiddleware()
}

var (
	GwtAuthMiddleware = DefaultGinJWTMiddleware()
)
