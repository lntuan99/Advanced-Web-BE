package api_grade

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/api/methods"
	"github.com/gin-gonic/gin"
)

func HandlerGetListGradeByClassroomId(c *gin.Context) {
	success, status, data := methods.MethodGetListGradeByClassroomId(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerCreateGrade(c *gin.Context) {
	success, status, data := methods.MethodCreateGrade(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerUpdateGrade(c *gin.Context) {
	success, status, data := methods.MethodUpdateGrade(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerDeleteGrade(c *gin.Context) {
	success, status, data := methods.MethodDeleteGrade(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerInputGradeForAStudent(c *gin.Context) {
	success, status, data := methods.MethodInputGradeForAStudent(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerGetGradeBoardByClassroomID(c *gin.Context) {
	success, status, data := methods.MethodGetGradeBoardByClassroomID(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerExportGradeBoardByClassroomID(c *gin.Context) {
	success, status, data := methods.MethodExportGradeBoardByClassroomID(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerImportGradeBoardByClassroomID(c *gin.Context) {
	success, status, data := methods.MethodImportGradeBoardByClassroomID(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerGetGradeBoardForStudentInClassroom(c *gin.Context) {
	success, status, data := methods.MethodGetGradeBoardForStudentInClassroom(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerGetGradeReviewRequested(c *gin.Context) {
	success, status, data := methods.MethodGetGradeReviewRequested(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerCreateGradeReviewRequested(c *gin.Context) {
	success, status, data := methods.MethodCreateGradeReviewRequested(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerMakeFinalDecisionGradeReviewRequested(c *gin.Context) {
	success, status, data := methods.MethodMakeFinalDecisionGradeReviewRequested(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerCreateCommentInGradeReviewRequested(c *gin.Context) {
	success, status, data := methods.MethodCreateCommentInGradeReviewRequested(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerGetListGradeReviewRequestedByClassroomId(c *gin.Context) {
	success, status, data := methods.MethodGetListGradeReviewRequestedByClassroomId(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}
