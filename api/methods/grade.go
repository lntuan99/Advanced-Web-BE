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
	if err := c.ShouldBindJSON(&gradeInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	if util.EmptyOrBlankString(gradeInfo.Name) {
		return false, base.CodeEmptyGradeName, nil
	}

	var classroom = model.Classroom{}.FindClassroomByID(gradeInfo.ClassroomID)
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	ok, _ := MiddlewareImplementUserIsATeacherInClassroom(user.ID, classroom.ID)
	if !ok {
		return false, base.CodeGradeUserInvalid, nil
	}

	// Check existed grade in class
	var existedGrade model.Grade
	model.DBInstance.First(&existedGrade, "classroom_id = ? AND name = ?", classroom.ID, gradeInfo.Name)

	if existedGrade.ID > 0 {
		return false, base.CodeGradeAlreadyInClassroom, nil
	}

	var gradeMaxOrdinary model.Grade
	var ordinalNumber uint
	if errorOrdinal := model.DBInstance.
		Order("ordinal_number DESC").
		Where("classroom_id = ?", classroom.ID).
		First(&gradeMaxOrdinary).Error; errorOrdinal != nil {
		ordinalNumber = 1
	} else {
		ordinalNumber = gradeMaxOrdinary.OrdinalNumber + 1
	}

	newGrade := model.Grade{
		ClassroomID:   gradeInfo.ClassroomID,
		Name:          gradeInfo.Name,
		MaxPoint:      gradeInfo.MaxPoint,
		OrdinalNumber: ordinalNumber,
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
	if err := c.ShouldBindJSON(&gradeInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	if util.EmptyOrBlankString(gradeInfo.Name) {
		return false, base.CodeEmptyGradeName, nil
	}

	// Check existed grade in class
	var existedGrade model.Grade
	model.DBInstance.
		Preload("Classroom").
		First(&existedGrade, "id = ?", gradeInfo.ID) // => model.DBInstance.First(&existedGrade, gradeInfo.GradeID)

	if existedGrade.ID == 0 {
		return false, base.CodeGradeNotExisted, nil
	}

	// Check classroom ID existed
	if existedGrade.Classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	ok, _ := MiddlewareImplementUserIsATeacherInClassroom(user.ID, existedGrade.ClassroomID)
	if !ok {
		return false, base.CodeGradeUserInvalid, nil
	}

	// check name of grade existed in classroom
	if existedGrade.Name != gradeInfo.Name {
		var existedNameGrade model.Grade
		model.DBInstance.First(&existedNameGrade, "classroom_id = ? AND name = ?", existedGrade.ClassroomID, gradeInfo.Name)

		if existedNameGrade.ID > 0 {
			return false, base.CodeGradeAlreadyInClassroom, nil
		}
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
	model.DBInstance.
		Preload("Classroom").
		First(&existedGrade, "id = ?", gradeID)

	if existedGrade.ID == 0 {
		return false, base.CodeGradeNotExisted, nil
	}

	if existedGrade.Classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	ok, _ := MiddlewareImplementUserIsATeacherInClassroom(user.ID, existedGrade.ClassroomID)
	if !ok {
		return false, base.CodeGradeUserInvalid, nil
	}

	model.DBInstance.
		Where("id = ?", existedGrade.ID).
		Delete(&model.Grade{})

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

	// Check user in classroom
	ok, _ := MiddlewareImplementUserInClassroom(user.ID, classroom.ID)
	if !ok {
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

func MethodInputGradeForAStudent(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))

	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	// Validate user is a teacher in classroom
	ok, _ := MiddlewareImplementUserIsATeacherInClassroom(user.ID, classroom.ID)
	if !ok {
		return false, base.CodeGradeUserInvalid, nil
	}

	var gradeInfo req_res.PostInputGradeForAStudent
	if err := c.ShouldBindJSON(&gradeInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	// Validate grade already belong to classroom
	var dbGrade model.Grade
	model.DBInstance.First(&dbGrade, gradeInfo.GradeID)
	if dbGrade.ClassroomID != classroom.ID {
		return false, base.CodeGradeNotBelongToClassroom, nil
	}

	// Validate student already in classroom
	ok, mapping := MiddlewareImplementUserInClassroom(gradeInfo.StudentID, classroom.ID)
	if !ok && mapping.UserRole.JWTType != model.JWT_TYPE_STUDENT {
		return false, base.CodeUserIsNotAStudentInClass, nil
	}

	var dbStudentGradeMapping model.UserGradeMapping
	model.DBInstance.First(&dbStudentGradeMapping, "user_id = ? AND grade_id = ?", gradeInfo.StudentID, gradeInfo.GradeID)

	dbStudentGradeMapping.UserID = gradeInfo.StudentID
	dbStudentGradeMapping.GradeID = gradeInfo.GradeID
	dbStudentGradeMapping.Point = gradeInfo.Point
	model.DBInstance.Save(&dbStudentGradeMapping)

	return true, base.CodeSuccess, nil
}

func MethodGetGradeBoardByClassroomID(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))

	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	// Validate user is a teacher in classroom
	ok, _ := MiddlewareImplementUserIsATeacherInClassroom(user.ID, classroom.ID)
	if !ok {
		return false, base.CodeGradeUserInvalid, nil
	}

	classroom.StudentArray = classroom.GetListUserByJWTType(model.JWT_TYPE_STUDENT)

	// Find all user grade mapping in classroom
	var dataResponse = make([]model.ResponseStudentGradeInClassroom, 0)
	for _, student := range classroom.StudentArray {
		var studentGradeResponse = student.MappedUserInformationToResponseStudentGradeInClassroom(classroom.ID)
		dataResponse = append(dataResponse, studentGradeResponse)
	}

	return true, base.CodeSuccess, dataResponse
}
