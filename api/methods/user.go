package methods

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"github.com/gin-gonic/gin"
)

func GetListClassroomByJWTType(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	var mappingArray = make([]model.UserClassroomMapping, 0)
	JWTType := util.ToUint(c.Query("jwt_type"))
	userRole :=  model.UserRole{}.GetRoleByJWTType(uint(JWTType))

	if userRole.ID == 0 {
		model.DBInstance.
			Order("id ASC").
			Offset(base.GetIntQuery(c, "page") * base.PageSizeLimit).
			Limit(base.PageSizeLimit).
			Preload("Classroom").
			Preload("Classroom.Owner").
			Where("user_id = ?", user.ID).
			Find(&mappingArray)
	}

	if userRole.ID > 0 {
		model.DBInstance.
			Order("id ASC").
			Offset(base.GetIntQuery(c, "page") * base.PageSizeLimit).
			Limit(base.PageSizeLimit).
			Preload("Classroom").
			Preload("Classroom.Owner").
			Where("user_id = ?", user.ID).
			Where("user_role_id = ?", userRole.ID).
			Find(&mappingArray)

	}

	var classroomResArray = make([]model.ClassroomRes, 0)
	for _, mapping := range mappingArray {
		classroomResArray = append(classroomResArray, mapping.Classroom.ToRes())
	}

	return true, base.CodeSuccess, classroomResArray
}

func GetListClassroomOwnedByUser(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	var classroomArray = make([]model.Classroom, 0)
	model.DBInstance.
		Order("id ASC").
		Offset(base.GetIntQuery(c, "page") * base.PageSizeLimit).
		Limit(base.PageSizeLimit).
		Where("owner_id = ?", user.ID).
		Find(&classroomArray)

	var classroomResArray = make([]model.ClassroomRes, 0)
	for _, classroom := range classroomArray {
		classroom.Owner = user
		classroomResArray = append(classroomResArray, classroom.ToRes())
	}

	return true, base.CodeSuccess, classroomResArray
}
