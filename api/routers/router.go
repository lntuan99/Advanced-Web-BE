package routers

import (
    "advanced-web.hcmus/api/base"
    api_account "advanced-web.hcmus/api/routers/api-account"
    api_classroom "advanced-web.hcmus/api/routers/api-classroom"
    api_status "advanced-web.hcmus/api/routers/api-status"
    "advanced-web.hcmus/config"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/module/apmgin"
    "os"
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

    routeVersion01 := r.Group("api/v1")

    // Multipart quota
    r.MaxMultipartMemory = 20971520 // Exactly 20MB
    r.Static("/media", "./public")

    accountRoute := routeVersion01.Group("account")
    {
        accountRoute.POST("/register", api_account.HandlerRegisterAccount)
    }

    classroomRoute := routeVersion01.Group("classroom")
    {
        classroomRoute.GET("/", api_classroom.HandlerGetClassroomList)
        classroomRoute.POST("/", api_classroom.HandlerCreateClassroom)
    }

    return r
}
