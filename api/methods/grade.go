package methods

import (
	"advanced-web.hcmus/api/base"
	req_res "advanced-web.hcmus/api/req_res_struct"
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"github.com/gin-gonic/gin"
)

func MethodCreateGrade(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	var gradeInfo req_res.PostCreateGrade
	if err := c.ShouldBind(&gradeInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	if util.EmptyOrBlankString(gradeInfo.Name) {
		return false, base.CodeEmptyGradeName, nil
	}

	var classroom = model.Classroom{}.FindClassroomByID(gradeInfo.ClassroomID)
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	//Check user is a teacher
	var mapping model.UserClassroomMapping
	model.DBInstance.
		Preload("UserRole").
		First(&mapping, "user_id = ? AND classroom_id = ?", user.ID, classroom.ID)
	if mapping.ID == 0 {
		return false, base.CodeBadRequest, nil
	}
	if mapping.UserRole.JWTType != model.JWT_TYPE_TEACHER {
		return false, base.CodeGradeUserInvalid, nil
	}

	// Check existed grade in class
	var existedGrade model.Grade
	model.DBInstance.First(&existedGrade, "classroom_id = ? AND name = ?", classroom.ID, gradeInfo.Name)

	if existedGrade.ID > 0 {
		return false, base.CodeGradeAlreadyInClassroom, nil
	}

	newGrade := model.Grade{
		ClassroomID:   gradeInfo.ClassroomID,
		Name:          gradeInfo.Name,
		MaxPoint:      gradeInfo.MaxPoint,
		OrdinalNumber: gradeInfo.OrdinalNumber,
	}
	err := model.DBInstance.Create(&newGrade).Error

	if err != nil {
		return false, base.CodeCreateGradeFail, nil
	}

	return true, base.CodeSuccess, newGrade.ToRes()
}
func MethodUpdateGrade(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	var gradeInfo req_res.PostUpdateGrade
	if err := c.ShouldBind(&gradeInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	if util.EmptyOrBlankString(gradeInfo.Name) {
		return false, base.CodeEmptyGradeName, nil
	}

	// Check existed grade in class
	var existedGrade model.Grade
	model.DBInstance.First(&existedGrade, "id = ?", gradeInfo.GradeID)

	if existedGrade.ID == 0 {
		return false, base.CodeGradeNotExisted, nil
	}

	//check nameExisted grade
	if existedGrade.Name != gradeInfo.Name {
		var existedNameGrade model.Grade
		model.DBInstance.First(&existedNameGrade, "classroom_id = ? AND name = ?", existedGrade.ClassroomID, gradeInfo.Name)
		if existedNameGrade.ID > 0 {
			return false, base.CodeGradeAlreadyInClassroom, nil
		}
	}

	var classroom = model.Classroom{}.FindClassroomByID(existedGrade.ClassroomID)
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	//Check user is a teacher
	var mapping model.UserClassroomMapping
	model.DBInstance.
		Preload("UserRole").
		First(&mapping, "user_id = ? AND classroom_id = ?", user.ID, classroom.ID)
	if mapping.ID == 0 {
		return false, base.CodeBadRequest, nil
	}
	if mapping.UserRole.JWTType != model.JWT_TYPE_TEACHER {
		return false, base.CodeGradeUserInvalid, nil
	}

	existedGrade.Name = gradeInfo.Name
	existedGrade.MaxPoint = gradeInfo.MaxPoint
	existedGrade.OrdinalNumber = gradeInfo.OrdinalNumber

	model.DBInstance.Save(&existedGrade)

	return true, base.CodeSuccess, existedGrade.ToRes()
}
func MethodDeleteGrade(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	gradeID := util.ToUint(c.Param("id"))

	// Check existed grade in class
	var existedGrade model.Grade
	model.DBInstance.First(&existedGrade, "id = ?", gradeID)

	if existedGrade.ID == 0 {
		return false, base.CodeGradeNotExisted, nil
	}

	var classroom = model.Classroom{}.FindClassroomByID(existedGrade.ClassroomID)
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	//Check user is a teacher
	var mapping model.UserClassroomMapping
	model.DBInstance.
		Preload("UserRole").
		First(&mapping, "user_id = ? AND classroom_id = ?", user.ID, classroom.ID)
	if mapping.ID == 0 {
		return false, base.CodeBadRequest, nil
	}
	if mapping.UserRole.JWTType != model.JWT_TYPE_TEACHER {
		return false, base.CodeGradeUserInvalid, nil
	}

	model.DBInstance.Delete(&existedGrade)

	return true, base.CodeSuccess, nil
}

func MethodGetListGradeByClassroomId(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))

	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	//Check user in classroom
	var mapping model.UserClassroomMapping
	model.DBInstance.
		Preload("UserRole").
		First(&mapping, "user_id = ? AND classroom_id = ?", user.ID, classroom.ID)
	if mapping.ID == 0 {
		return false, base.CodeBadRequest, nil
	}

	var gradeArray = make([]model.Grade, 0)

	model.DBInstance.
		Order("ordinal_number ASC").
		Where("classroom_id = ?", classroomID).
		Find(&gradeArray)

	var gradeResArray = make([]model.GradeRes, 0)
	for _, grade := range gradeArray {
		gradeResArray = append(gradeResArray, grade.ToRes())
	}
	return true, base.CodeSuccess, gradeResArray
}
