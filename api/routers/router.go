package routers

import (
	api_jwt "advanced-web.hcmus/api/api-jwt"
	"advanced-web.hcmus/api/base"
	api_account "advanced-web.hcmus/api/routers/api-account"
	api_admin "advanced-web.hcmus/api/routers/api-admin"
	api_admin_classroom "advanced-web.hcmus/api/routers/api-admin/api-admin-classroom"
	api_admin_user "advanced-web.hcmus/api/routers/api-admin/api-admin-user"
	api_classroom "advanced-web.hcmus/api/routers/api-classroom"
	api_grade "advanced-web.hcmus/api/routers/api-grade"
	api_notification "advanced-web.hcmus/api/routers/api-notification"
	api_status "advanced-web.hcmus/api/routers/api-status"
	api_user "advanced-web.hcmus/api/routers/api-user"
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

	authMiddleware := api_jwt.GwtAuthMiddleware
	routeVersion01 := r.Group("api/v1")

	// Multipart quota
	r.MaxMultipartMemory = 20971520 // Exactly 20MB
	r.Static("/media", "./public")
	r.Static("/export-data", "./public/system")

	accountRoute := routeVersion01.Group("account")
	{
		accountRoute.POST("/register", api_account.HandlerRegisterAccount)
		accountRoute.POST("/login", api_account.HandlerLoginAccount)
		accountRoute.POST("/google-login", api_account.HandlerGoogleLogin)
	}

	userRoute := routeVersion01.Group("user")

	// these api no need authorization
	userRoute.GET("/verify", api_user.HandlerVerifyCode)
	userRoute.GET("/forgot-password", api_user.HandlerForgotPassword)

	userRoute.Use(authMiddleware.MiddlewareFuncUser())
	{
		userRoute.GET("/", api_user.HandlerGetUserProfile)

		// old: POST
		userRoute.POST("/", api_user.HandlerUpdateUserProfile)

		// new: PUT
		// userRoute.PUT("/", api_user.HandlerUpdateUserProfile)

		userRoute.PUT("/update-password", api_user.HandlerUpdatePassword)
	}

	classroomRoute := routeVersion01.Group("classroom")
	classroomRoute.Use(authMiddleware.MiddlewareFuncUser())
	{
		classroomRoute.GET("/", api_classroom.HandlerGetClassroomList)
		classroomRoute.GET("/get-list-classroom-by-jwt-type", api_classroom.HandlerGetListClassroomByJWTType)
		classroomRoute.GET("/get-list-classroom-owned-by-user", api_classroom.HandlerGetListClassroomOwnedByUser)
		classroomRoute.GET("/:id", api_classroom.HandlerGetClassroomByID)
		classroomRoute.POST("/", api_classroom.HandlerCreateClassroom)
		classroomRoute.GET("/join", api_classroom.HandlerJoinClassroom)
		classroomRoute.POST("/invite", api_classroom.HandlerInviteToClassroom)
		classroomRoute.GET("/:id/export-student", api_classroom.HandlerExportStudentListByClassroomID)
		classroomRoute.POST("/:id/import-student", api_classroom.HandlerImportStudentListByClassroomID)

		gradeStructureRoute := classroomRoute.Group("grade")
		{
			gradeStructureRoute.GET("/:id", api_grade.HandlerGetListGradeByClassroomId)
			gradeStructureRoute.GET("/review-requested/:id", api_grade.HandlerGetListGradeReviewRequestedByClassroomId)
			gradeStructureRoute.POST("/add", api_grade.HandlerCreateGrade)

			// old: POST, /update
			gradeStructureRoute.POST("/update", api_grade.HandlerUpdateGrade)

			// new: PUT
			// gradeStructureRoute.PUT("/", api_grade.HandlerUpdateGrade)

			// old: GET, /delete/:id
			gradeStructureRoute.GET("/delete/:id", api_grade.HandlerDeleteGrade)

			// new: DELETE, /:id
			// gradeStructureRoute.GET("/:id", api_grade.HandlerDeleteGrade)

			gradeStructureRoute.POST("/:id", api_grade.HandlerInputGradeForAStudent)

			gradeStructureBoardRoute := gradeStructureRoute.Group("/board")
			{
				gradeStructureBoardRoute.GET("/:id", api_grade.HandlerGetGradeBoardByClassroomID)
				gradeStructureBoardRoute.POST("/:id/export-grade-board", api_grade.HandlerExportGradeBoardByClassroomID)
				gradeStructureBoardRoute.POST("/:id/import-grade-board", api_grade.HandlerImportGradeBoardByClassroomID)
			}

			gradeStudentRoute := gradeStructureRoute.Group("/student")
			{
				// this id is classroom ID
				gradeStudentRoute.GET("/:id", api_grade.HandlerGetGradeBoardForStudentInClassroom)
				gradeStudentRoute.GET("/:id/:grade-id", api_grade.HandlerGetGradeReviewRequested)
				gradeStudentRoute.POST("/:id/:grade-id", api_grade.HandlerCreateGradeReviewRequested)
				gradeStudentRoute.PUT("/:id/:grade-id", api_grade.HandlerMakeFinalDecisionGradeReviewRequested)
				gradeStudentRoute.POST("/:id/:grade-id/comment", api_grade.HandlerCreateCommentInGradeReviewRequested)
			}
		}
	}

	notificationRoute := routeVersion01.Group("notification")
	notificationRoute.Use(authMiddleware.MiddlewareFuncUser())
	{
		notificationRoute.GET("/list", api_notification.HandlerGetListNotification)
		notificationRoute.POST("/mark-read", api_notification.HandlerMarkReadNotification)
	}

	adminRoute := routeVersion01.Group("admin")
	adminRoute.POST("/login", api_admin.HandlerLoginAdminAccount)
	adminRoute.Use(authMiddleware.MiddlewareFuncAdminUser())
	{
		adminRoute.GET("", api_admin.HandlerGetListAdminUser)
		adminRoute.GET("/:id", api_admin.HandlerGetAdminUserByID)
		adminRoute.POST("/", api_admin.HandlerCreateAdminUser)

		manageUserRoute := adminRoute.Group("user")
		{
			manageUserRoute.GET("", api_admin_user.HandlerGetListUser)
			manageUserRoute.GET("/:id", api_admin_user.HandlerAdminGetUserByID)
			manageUserRoute.POST("/ban/:id", api_admin_user.HandlerAdminBanUserByID)
			manageUserRoute.POST("/map-student-code", api_admin_user.HandlerMapStudentCode)
		}

		manageClassroomRoute := adminRoute.Group("classroom")
		{
			manageClassroomRoute.GET("", api_admin_classroom.HandlerAdminGetListClassroom)
			manageClassroomRoute.GET("/:id", api_admin_classroom.HandlerAdminGetClassroomByID)
		}
	}

	return r
}
