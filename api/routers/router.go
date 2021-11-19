package routers

import (
	api_jwt "advanced-web.hcmus/api/api-jwt"
	api_user "advanced-web.hcmus/api/routers/api-user"
	"os"

	"advanced-web.hcmus/api/base"
	api_account "advanced-web.hcmus/api/routers/api-account"
	api_classroom "advanced-web.hcmus/api/routers/api-classroom"
	api_status "advanced-web.hcmus/api/routers/api-status"
	"advanced-web.hcmus/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin"
)

func Initialize() *gin.Engine {
	r := gin.New()

	//Set GIN in RELEASE_MODE if ENV is "production"
	if os.Getenv("ENV") == config.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	corConfig := cors.DefaultConfig()
	corConfig.AllowAllOrigins = true
	corConfig.AllowHeaders = []string{
		"authorization", "Authorization",
		"content-type", "accept",
		"referer", "user-agent",
	}
	r.Use(cors.New(corConfig))

	r.Use(apmgin.Middleware(r))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(base.MiddlewareClientVersion())

	r.GET("/status", api_status.HandlerStatus)

	authMiddleware := api_jwt.GwtAuthMiddleware
	routeVersion01 := r.Group("api/v1")

	// Multipart quota
	r.MaxMultipartMemory = 20971520 // Exactly 20MB
	r.Static("/media", "./public")

	accountRoute := routeVersion01.Group("account")
	{
		accountRoute.POST("/register", api_account.HandlerRegisterAccount)
		accountRoute.POST("/login", api_account.HandlerLoginAccount)
		accountRoute.POST("/google-login", api_account.HandlerGoogleLogin)
	}

	userRoute := routeVersion01.Group("user")
	userRoute.Use(authMiddleware.MiddlewareFuncUser())
	{
		userRoute.POST("/classroom/get-list-classroom-by-jwt-type", api_user.HandlerGetListClassroomByJWTType)
		userRoute.GET("/classroom/get-list-classroom-by-jwt-type", api_user.HandlerGetListClassroomByJWTType)
		userRoute.GET("/classroom/get-list-classroom-owned-by-user", api_user.HandlerGetListClassroomOwnedByUser)
	}

	classroomRoute := routeVersion01.Group("classroom")
	classroomRoute.Use(authMiddleware.MiddlewareFuncUser())
	{
		classroomRoute.GET("/", api_classroom.HandlerGetClassroomList)
		classroomRoute.GET("/:id", api_classroom.HandlerGetClassroomByID)
		classroomRoute.POST("/", api_classroom.HandlerCreateClassroom)
	}

	return r
}
