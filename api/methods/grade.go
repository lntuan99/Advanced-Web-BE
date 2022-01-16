package methods

import (
	"advanced-web.hcmus/api/base"
	req_res "advanced-web.hcmus/api/req_res_struct"
	export_excel "advanced-web.hcmus/biz/export-excel"
	import_excel "advanced-web.hcmus/biz/import-excel"
	"advanced-web.hcmus/biz/upload"
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"sort"
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

	// Check existed grade in class
	var dbGrade model.Grade
	model.DBInstance.
		Preload("Classroom").
		First(&dbGrade, gradeInfo.ID) // => model.DBInstance.First(&existedGrade, gradeInfo.GradeID)

	if dbGrade.ID == 0 {
		return false, base.CodeGradeNotExisted, nil
	}

	// Check classroom ID existed
	if dbGrade.Classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	ok, _ := MiddlewareImplementUserIsATeacherInClassroom(user.ID, dbGrade.ClassroomID)
	if !ok {
		return false, base.CodeGradeUserInvalid, nil
	}

	// check name of grade existed in classroom
	if dbGrade.Name != gradeInfo.Name {
		var existedNameGrade model.Grade
		model.DBInstance.First(&existedNameGrade, "classroom_id = ? AND name = ?", dbGrade.ClassroomID, gradeInfo.Name)

		if existedNameGrade.ID > 0 {
			return false, base.CodeGradeAlreadyInClassroom, nil
		}
	}

	// If Mark Finalized
	var createNotification = false
	if dbGrade.IsFinalized == false && gradeInfo.IsFinalized == true {
		createNotification = true
	}

	dbGrade.Name = gradeInfo.Name
	dbGrade.MaxPoint = gradeInfo.MaxPoint
	dbGrade.OrdinalNumber = gradeInfo.OrdinalNumber
	dbGrade.IsFinalized = gradeInfo.IsFinalized

	model.DBInstance.Save(&dbGrade)

	if createNotification {
		var newNotification = model.Notification{
			Title:   fmt.Sprintf("Đã có thể xem điểm '%v' lớp '%v'", dbGrade.Name, dbGrade.Classroom.Name),
			Message: fmt.Sprintf("Điểm '%v' lớp '%v' đã được công bố. Vui lòng nhấn vào thông báo hoặc truy cập vào lớp để xem điểm", dbGrade.Name, dbGrade.Classroom.Name),
			Payload: fmt.Sprintf("%v", dbGrade.ID),
			Users:   nil,
		}

		dbGrade.Classroom.GetListStudent()

		for _, student := range dbGrade.Classroom.StudentArray {
			if student.User.ID > 0 {
				newNotification.Users = append(newNotification.Users, student.User)
			}
		}

		model.DBInstance.Save(&newNotification)
	}

	return true, base.CodeSuccess, dbGrade.ToRes()
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
	//ok, mapping := MiddlewareImplementUserInClassroom(gradeInfo.StudentID, classroom.ID)
	//if !ok && mapping.UserRole.JWTType != model.JWT_TYPE_STUDENT {
	//	return false, base.CodeUserIsNotAStudentInClass, nil
	//}

	var dbStudentGradeMapping model.StudentGradeMapping
	model.DBInstance.First(&dbStudentGradeMapping, "student_id = ? AND grade_id = ?", gradeInfo.StudentID, gradeInfo.GradeID)

	dbStudentGradeMapping.StudentID = gradeInfo.StudentID
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

	classroom.GetListStudent()

	// Find all user grade mapping in classroom
	var dataResponse = make([]model.ResponseStudentGradeInClassroom, 0)
	for _, student := range classroom.StudentArray {
		var studentGradeResponse = student.MappedStudentInformationToResponseStudentGradeInClassroom(classroom.ID, nil)
		dataResponse = append(dataResponse, studentGradeResponse)
	}

	return true, base.CodeSuccess, dataResponse
}

func MethodExportGradeBoardByClassroomID(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))
	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))

	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	// validate user is a teacher in class
	ok, _ := MiddlewareImplementUserIsATeacherInClassroom(user.ID, classroom.ID)
	if !ok {
		return false, base.CodeBadRequest, nil
	}

	var gradeBoardInfo req_res.PostExportGradeBoard
	if err := c.ShouldBindJSON(&gradeBoardInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	// Validate grade already belong to classroom
	var okeGradeArray = make([]model.Grade, 0)
	for _, gradeID := range gradeBoardInfo.GradeIDArray {
		var dbGrade model.Grade

		model.DBInstance.First(&dbGrade, gradeID)
		if dbGrade.ClassroomID == classroom.ID {
			okeGradeArray = append(okeGradeArray, dbGrade)
		}
	}

	sort.Slice(okeGradeArray, func(i, j int) bool {
		return okeGradeArray[i].OrdinalNumber < okeGradeArray[j].OrdinalNumber
	})

	classroom.GetListStudent()

	var responseStudentGradeInClassroomArray = make([]model.ResponseStudentGradeInClassroom, len(classroom.StudentArray))
	for i, student := range classroom.StudentArray {
		var studentGradeResponse = student.MappedStudentInformationToResponseStudentGradeInClassroom(classroom.ID, nil)
		responseStudentGradeInClassroomArray[i] = studentGradeResponse
	}

	fileUrl := export_excel.ProcessExportGradeBoard(responseStudentGradeInClassroomArray, okeGradeArray)

	return true, base.CodeSuccess, fileUrl
}

func MethodImportGradeBoardByClassroomID(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))
	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))

	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	//Check user is a teacher in class
	ok := MiddlewareImplementUserIsAnOwnerOfClassroom(user.ID, classroom)
	if !ok {
		return false, base.CodeBadRequest, nil
	}

	file, header, errFile := c.Request.FormFile("import-grade-board-file")
	if errFile != nil {
		return false, base.CodeImportStudentFail, nil
	}

	fileBytes, err := ioutil.ReadAll(file)
	util.CheckErr(err)

	filePath := upload.WithTemporary().Save(header.Filename, fileBytes)

	biz := import_excel.SheetGradeBoardStruct{}.Initialize(filePath, classroom.ID)

	ok, importResponseArray := biz.Importing()

	if !ok {
		return false, base.CodeImportGradeBoardFail, nil
	}

	return true, base.CodeSuccess, importResponseArray
}

func MethodGetGradeBoardForStudentInClassroom(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))

	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	// Validate user is a student in classroom
	ok, mapping := MiddlewareImplementUserInClassroom(user.ID, classroom.ID)
	if !ok || mapping.UserRole.JWTType != model.JWT_TYPE_STUDENT {
		return false, base.CodeUserIsNotAStudentInClass, nil
	}

	// Find student mapping this user in classroom
	var mappingStudent model.Student
	model.DBInstance.
		Preload("User").
		Where("code = ? and classroom_id = ?", user.Code, classroom.ID).
		First(&mappingStudent)

	var isFinalized = true

	return true, base.CodeSuccess, mappingStudent.MappedStudentInformationToResponseStudentGradeInClassroom(classroom.ID, &isFinalized)
}

func MethodGetGradeReviewRequested(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))
	gradeID := util.ToUint(c.Param("grade-id"))
	reviewRequestedID := util.ToUint(c.Query("review-id"))

	// Validate classroom
	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	// Validate grade already belong to classroom
	var dbGrade model.Grade
	model.DBInstance.First(&dbGrade, gradeID)
	if dbGrade.ClassroomID != classroom.ID {
		return false, base.CodeGradeNotBelongToClassroom, nil
	}

	// Validate user already in classroom
	ok, mapping := MiddlewareImplementUserInClassroom(user.ID, classroom.ID)
	if !ok {
		return false, base.CodeUserNotInClassroom, nil
	}

	// Validate grade review requested in classroom
	var dbGradeReviewRequested model.GradeReviewRequested
	model.DBInstance.
		Preload("Comments").
		Preload("StudentGradeMapping.Student").
		Preload("StudentGradeMapping.Grade").
		First(&dbGradeReviewRequested, reviewRequestedID)

	if dbGradeReviewRequested.ID == 0 ||
		dbGradeReviewRequested.StudentGradeMapping.Grade.ClassroomID != classroom.ID ||
		dbGradeReviewRequested.StudentGradeMapping.GradeID != uint(gradeID) {
		return false, base.CodeReviewRequestedNotInClassroom, nil
	}

	// If user is a student
	// validate user is an owner of requested
	if mapping.UserRole.JWTType == model.JWT_TYPE_STUDENT {
		// Find student mapping this user in classroom
		var mappingStudent model.Student
		model.DBInstance.
			Preload("User").
			Where("code = ? and classroom_id = ?", user.Code, classroom.ID).
			First(&mappingStudent)

		if dbGradeReviewRequested.StudentGradeMapping.StudentID != mappingStudent.ID {
			return false, base.CodeUserNotAnOwnerOfRequested, nil
		}
	}

	return true, base.CodeSuccess, dbGradeReviewRequested.ToRes()
}

func MethodCreateGradeReviewRequested(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))
	gradeID := util.ToUint(c.Param("grade-id"))

	// Validate classroom
	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	// Validate grade already belong to classroom
	var dbGrade model.Grade
	model.DBInstance.First(&dbGrade, gradeID)
	if dbGrade.ClassroomID != classroom.ID {
		return false, base.CodeGradeNotBelongToClassroom, nil
	}

	// Validate user already in classroom
	ok, mapping := MiddlewareImplementUserInClassroom(user.ID, classroom.ID)
	if !ok || mapping.UserRole.JWTType != model.JWT_TYPE_STUDENT {
		return false, base.CodeUserIsNotAStudentInClass, nil
	}

	var gradeReviewRequestedInfo req_res.PostCreateGradeReviewRequested
	if err := c.ShouldBindJSON(&gradeReviewRequestedInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	// Find student mapping this user in classroom
	var mappingStudent model.Student
	model.DBInstance.
		Preload("User").
		Where("code = ? and classroom_id = ?", user.Code, classroom.ID).
		First(&mappingStudent)

	// Find student grade mapping
	var dbStudentGradeMapping model.StudentGradeMapping
	model.DBInstance.
		Preload("Student").
		Preload("Grade").
		Where("grade_id = ? AND student_id = ?", gradeID, mappingStudent.ID).
		First(&dbStudentGradeMapping)

	var dbGradeReviewRequested model.GradeReviewRequested
	model.DBInstance.
		Preload("Comments").
		Where("student_grade_mapping_id = ?", dbStudentGradeMapping.ID).
		First(&dbGradeReviewRequested)

	if dbGradeReviewRequested.IsProcessed {
		return false, base.CodeGradeReviewRequestedHasBeenProcessed, nil
	}

	// Mapping review requested data
	dbGradeReviewRequested.StudentGradeMappingID = dbStudentGradeMapping.ID
	dbGradeReviewRequested.StudentGradeMapping = dbStudentGradeMapping
	dbGradeReviewRequested.StudentExplanation = gradeReviewRequestedInfo.StudentExplanation
	dbGradeReviewRequested.StudentExpectation = gradeReviewRequestedInfo.StudentExpectation

	var createNewNotification = false
	if dbGradeReviewRequested.ID == 0 {
		createNewNotification = true
	}

	model.DBInstance.Save(&dbGradeReviewRequested)

	// Create notification
	if createNewNotification {
		var newNotification = model.Notification{
			Title:   fmt.Sprintf("Có một đơn phúc khảo mới"),
			Message: fmt.Sprintf("Học sinh '%v' lớp '%v' đã yêu cầu phúc khảo điểm '%v'", dbStudentGradeMapping.Student.Name, classroom.Name, dbGrade.Name),
			Payload: fmt.Sprintf("%v", dbGradeReviewRequested.ID),
			Users:   nil,
		}

		newNotification.Users = classroom.GetListUserByJWTType(model.JWT_TYPE_TEACHER)
		model.DBInstance.Save(&newNotification)
	}

	return true, base.CodeSuccess, dbGradeReviewRequested.ToRes()
}

func MethodMakeFinalDecisionGradeReviewRequested(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))
	gradeID := util.ToUint(c.Param("grade-id"))

	// Validate classroom
	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	// Validate grade already belong to classroom
	var dbGrade model.Grade
	model.DBInstance.First(&dbGrade, gradeID)
	if dbGrade.ClassroomID != classroom.ID {
		return false, base.CodeGradeNotBelongToClassroom, nil
	}

	//validate user is a teacher in class
	ok, _ := MiddlewareImplementUserIsATeacherInClassroom(user.ID, classroom.ID)
	if !ok {
		return false, base.CodeBadRequest, nil
	}

	var gradeReviewRequestedInfo req_res.PostMakeFinalDecisionGradeReviewRequested
	if err := c.ShouldBindJSON(&gradeReviewRequestedInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	// Validate grade review requested in classroom
	var dbGradeReviewRequested model.GradeReviewRequested
	model.DBInstance.
		Preload("Comments").
		Preload("StudentGradeMapping.Student.User").
		Preload("StudentGradeMapping.Grade").
		First(&dbGradeReviewRequested, gradeReviewRequestedInfo.GradeReviewRequestedID)

	if dbGradeReviewRequested.ID == 0 ||
		dbGradeReviewRequested.StudentGradeMapping.Grade.ClassroomID != classroom.ID ||
		dbGradeReviewRequested.StudentGradeMapping.GradeID != uint(gradeID) {
		return false, base.CodeReviewRequestedNotInClassroom, nil
	}

	if dbGradeReviewRequested.IsProcessed {
		return false, base.CodeGradeReviewRequestedHasBeenProcessed, nil
	}

	// Mapping review requested data
	dbGradeReviewRequested.FinalPoint = &gradeReviewRequestedInfo.FinalPoint
	dbGradeReviewRequested.StudentGradeMapping.Point = gradeReviewRequestedInfo.FinalPoint
	dbGradeReviewRequested.IsProcessed = true

	model.DBInstance.Save(&dbGradeReviewRequested)

	var newNotification = model.Notification{
		Title:   fmt.Sprintf("Đã có kết quả phúc khảo"),
		Message: fmt.Sprintf("Đơn phúc khảo điểm '%v' lớp '%v' đã có kết quả", dbGradeReviewRequested.StudentGradeMapping.Grade.Name, classroom.Name),
		Payload: fmt.Sprintf("%v", dbGradeReviewRequested.ID),
		Users:   []model.User{dbGradeReviewRequested.StudentGradeMapping.Student.User},
	}
	model.DBInstance.Save(&newNotification)

	return true, base.CodeSuccess, dbGradeReviewRequested.ToRes()
}

func MethodCreateCommentInGradeReviewRequested(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))
	gradeID := util.ToUint(c.Param("grade-id"))

	// Validate classroom
	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	// Validate grade already belong to classroom
	var dbGrade model.Grade
	model.DBInstance.First(&dbGrade, gradeID)
	if dbGrade.ClassroomID != classroom.ID {
		return false, base.CodeGradeNotBelongToClassroom, nil
	}

	// Validate user is a student in classroom
	ok, mapping := MiddlewareImplementUserInClassroom(user.ID, classroom.ID)
	if !ok {
		return false, base.CodeUserNotInClassroom, nil
	}

	var commentInfo req_res.PostCreateCommentInGradeReviewRequested
	if err := c.ShouldBindJSON(&commentInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	// Validate grade review requested in classroom
	var dbGradeReviewRequested model.GradeReviewRequested
	model.DBInstance.
		Preload("StudentGradeMapping.Student.User").
		Preload("StudentGradeMapping.Grade").
		First(&dbGradeReviewRequested, commentInfo.GradeReviewRequestedID)

	if dbGradeReviewRequested.ID == 0 ||
		dbGradeReviewRequested.StudentGradeMapping.Grade.ClassroomID != classroom.ID ||
		dbGradeReviewRequested.StudentGradeMapping.GradeID != uint(gradeID) {
		return false, base.CodeReviewRequestedNotInClassroom, nil
	}

	// If user is a student
	// validate user is an owner of requested
	if mapping.UserRole.JWTType == model.JWT_TYPE_STUDENT {
		// Find student mapping this user in classroom
		var mappingStudent model.Student
		model.DBInstance.
			Preload("User").
			Where("code = ? and classroom_id = ?", user.Code, classroom.ID).
			First(&mappingStudent)

		if dbGradeReviewRequested.StudentGradeMapping.StudentID != mappingStudent.ID {
			return false, base.CodeUserNotAnOwnerOfRequested, nil
		}
	}

	var newComment model.GradeReviewRequestedComment
	newComment.GradeReviewRequestedID = commentInfo.GradeReviewRequestedID
	newComment.Comment = commentInfo.Comment
	newComment.UserID = user.ID

	err := model.DBInstance.Create(&newComment).Error
	if err != nil {
		return false, base.CodeCreateCommentFail, nil
	}

	// Create notification
	// Only push notification for student if owner of comment is teacher
	if mapping.UserRole.JWTType == model.JWT_TYPE_TEACHER {
		var newNotification = model.Notification{
			Title:   fmt.Sprintf("Có một bình luận mới"),
			Message: fmt.Sprintf("Giáo viên '%v' lớp '%v' đã bình luận vào đơn phúc khảo của bạn", user.Name, classroom.Name),
			Payload: fmt.Sprintf("%v", newComment.GradeReviewRequestedID),
			Users:   []model.User{dbGradeReviewRequested.StudentGradeMapping.Student.User},
		}

		model.DBInstance.Save(&newNotification)
	}

	return true, base.CodeSuccess, newComment.ToRes()
}

func MethodGetListGradeReviewRequestedByClassroomId(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))

	// Validate classroom
	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	// validate user is a teacher in classroom
	ok, _ := MiddlewareImplementUserIsATeacherInClassroom(user.ID, classroom.ID)
	if !ok {
		return false, base.CodeBadRequest, nil
	}

	// Validate grade review requested in classroom
	var gradeReviewRequestedArray = make([]model.GradeReviewRequested, 0)
	model.DBInstance.
		Preload("Comments").
		Preload("StudentGradeMapping.Student").
		Preload("StudentGradeMapping.Grade").
		Find(&gradeReviewRequestedArray)

	var gradeReviewRequestedResArray = make([]model.GradeReviewRequestedRes, 0)

	for _, review := range gradeReviewRequestedArray {
		gradeReviewRequestedResArray = append(gradeReviewRequestedResArray, review.ToRes())
	}

	return true, base.CodeSuccess, gradeReviewRequestedResArray
}
