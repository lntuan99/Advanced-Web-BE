package api_classroom

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/api/methods"
	"github.com/gin-gonic/gin"
)

func HandlerGetClassroomList(c *gin.Context) {
	success, status, data := methods.MethodGetListClassroom(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

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

func HandlerGetClassroomByID(c *gin.Context) {
	success, status, data := methods.MethodGetClassroomByID(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerCreateClassroom(c *gin.Context) {
	success, status, data := methods.MethodCreateClassroom(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerJoinClassroom(c *gin.Context) {
	success, status, data := methods.MethodJoinClassroom(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerInviteToClassroom(c *gin.Context) {
	success, status, data := methods.MethodInviteToClassroom(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerExportStudentListByClassroomID(c *gin.Context) {
	success, status, data := methods.MethodExportStudentListByClassroomID(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}

func HandlerImportStudentListByClassroomID(c *gin.Context) {
	success, status, data := methods.MethodImportStudentListByClassroomID(c)

	if !success {
		base.ResponseError(c, status)
	} else {
		base.ResponseResult(c, data)
	}
}
