package methods

import (
	"advanced-web.hcmus/api/base"
	req_res "advanced-web.hcmus/api/req_res_struct"
	export_excel "advanced-web.hcmus/biz/export-excel"
	import_excel "advanced-web.hcmus/biz/import-excel"
	"advanced-web.hcmus/biz/upload"
	"advanced-web.hcmus/config"
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/services/smtp"
	"advanced-web.hcmus/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"path/filepath"
	"time"
)

func MethodGetClassroomList(c *gin.Context) (bool, string, interface{}) {
	var classroomArray = make([]model.Classroom, 0)
	model.DBInstance.
		Order("id ASC").
		Offset(base.GetIntQuery(c, "page") * base.PageSizeLimit).
		Limit(base.PageSizeLimit).
		Preload("Owner").
		Find(&classroomArray)

	var classroomResArray = make([]model.ClassroomRes, 0)
	for _, classroom := range classroomArray {
		classroomResArray = append(classroomResArray, classroom.ToRes())
	}

	return true, base.CodeSuccess, classroomResArray
}

func GetListClassroomByJWTType(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	var mappingArray = make([]model.UserClassroomMapping, 0)
	JWTType := util.ToUint(c.Query("jwt_type"))
	userRole := model.UserRole{}.GetRoleByJWTType(uint(JWTType))

	if userRole.ID == 0 {
		model.DBInstance.
			Order("id ASC").
			Offset(base.GetIntQuery(c, "page")*base.PageSizeLimit).
			Limit(base.PageSizeLimit).
			Preload("Classroom").
			Preload("Classroom.Owner").
			Where("user_id = ?", user.ID).
			Find(&mappingArray)
	}

	if userRole.ID > 0 {
		model.DBInstance.
			Order("id ASC").
			Offset(base.GetIntQuery(c, "page")*base.PageSizeLimit).
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
		Offset(base.GetIntQuery(c, "page")*base.PageSizeLimit).
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

func MethodGetClassroomByID(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))
	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))

	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	// Check user already in class
	ok, mapping := MiddlewareImplementUserInClassroom(user.ID, classroom.ID)
	if !ok {
		return false, base.CodeBadRequest, nil
	}

	classroom.GetListStudent()
	classroom.TeacherArray = classroom.GetListUserByJWTType(model.JWT_TYPE_TEACHER)

	classroomRes := classroom.ToRes()
	classroomRes.JWTType = mapping.UserRole.JWTType

	return true, base.CodeSuccess, classroomRes
}

func MethodCreateClassroom(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	_ = c.Request.ParseMultipartForm(20971520)

	var classroomInfo req_res.PostCreateClassroom
	if err := c.ShouldBind(&classroomInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	if util.EmptyOrBlankString(classroomInfo.Name) {
		return false, base.CodeEmptyClassroomName, nil
	}

	if util.EmptyOrBlankString(classroomInfo.Code) {
		return false, base.CodeEmptyClassroomCode, nil
	}

	var newClassroom = model.Classroom{
		OwnerID:       user.ID,
		Name:          classroomInfo.Name,
		CoverImageURL: "",
		Code:          classroomInfo.Code,
		Description:   classroomInfo.Description,
	}
	newClassroom.ClassroomGenerateInviteCode()

	err := model.DBInstance.Create(&newClassroom).Error

	if err != nil {
		return false, base.CodeCreateClassroomFail, nil
	}

	_, header, errFile := c.Request.FormFile("coverImage")
	if errFile == nil {
		newFileName := fmt.Sprintf("%v%v", time.Now().Unix(), filepath.Ext(header.Filename))
		folderDst := fmt.Sprintf("/system/classrooms/%v", newClassroom.ID)

		util.CreateFolder(folderDst)

		fileDst := fmt.Sprintf("%v/%v", folderDst, newFileName)
		if err := util.SaveUploadedFile(header, folderDst, newFileName); err == nil {
			model.DBInstance.
				Model(&newClassroom).
				Update("cover_image_url", fileDst)
		}
	}

	// Temp: Set owner to teacher of class
	var newMapping = model.UserClassroomMapping{
		ClassroomID: newClassroom.ID,
		UserID:      user.ID,
		UserRoleID:  model.UserRole{}.GetRoleByJWTType(model.JWT_TYPE_TEACHER).ID,
	}
	model.DBInstance.Create(&newMapping)

	return true, base.CodeSuccess, newClassroom.ToRes()
}

func MethodJoinClassroom(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	inviteCode := c.Query("code")

	// Validate invite code
	if util.EmptyOrBlankString(inviteCode) {
		return false, base.CodeInvalidClassroomInviteCode, nil
	}

	existed, classroom, jwtType := model.Classroom{}.GetClassroomByInviteCode(inviteCode)
	if !existed {
		return false, base.CodeInvalidClassroomInviteCode, nil
	}

	// Check user already be an owner of class
	if classroom.OwnerID == user.ID {
		return false, base.CodeUserAlreadyOwnerOfClass, nil
	}

	// Check user existed in class
	ok, _ := MiddlewareImplementUserInClassroom(user.ID, classroom.ID)
	if ok {
		return false, base.CodeUserAlreadyInClassroom, nil
	}

	// Create new mapping
	var newMapping = model.UserClassroomMapping{
		ClassroomID: classroom.ID,
		UserID:      user.ID,
		UserRoleID:  model.UserRole{}.GetRoleByJWTType(jwtType).ID,
	}
	model.DBInstance.Create(&newMapping)

	return true, base.CodeSuccess, nil
}

func MethodInviteToClassroom(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	var inviteToClassroomInfo req_res.PostInviteToClassroom
	if err := c.ShouldBindJSON(&inviteToClassroomInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	var classroom = model.Classroom{}.FindClassroomByID(inviteToClassroomInfo.ClassroomID)
	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	if classroom.OwnerID != user.ID {
		return false, base.CodeOnlyOwnerCanInviteMemberToClassroom, false
	}

	// If empty invite code => create new and save to database
	createNewCode := false
	if util.EmptyOrBlankString(classroom.InviteTeacherCode) {
		createNewCode = true
		classroom.InviteTeacherCode = model.GenerateInviteCode(classroom, model.JWT_TYPE_TEACHER)
	}
	if util.EmptyOrBlankString(classroom.InviteStudentCode) {
		createNewCode = true
		classroom.InviteStudentCode = model.GenerateInviteCode(classroom, model.JWT_TYPE_STUDENT)
	}

	if createNewCode {
		model.DBInstance.Save(&classroom)
	}

	// Generate invite link
	inviteTeacherLink := fmt.Sprintf("%v/join?code=%v", config.Config.FeDomain, classroom.InviteTeacherCode)
	inviteStudentLink := fmt.Sprintf("%v/join?code=%v", config.Config.FeDomain, classroom.InviteStudentCode)

	type TemplateData struct {
		URL string
	}
	teacherTemplateData := TemplateData{
		URL: inviteTeacherLink,
	}
	studentTemplateData := TemplateData{
		URL: inviteStudentLink,
	}

	// Send invite teacher
	r1 := smtp.NewRequest(inviteToClassroomInfo.TeacherEmailArray, "JOIN MY CLASS AS A TEACHER", "JOIN MY CLASS AS A TEACHER")
	if err1 := r1.ParseTemplate("./public/assets/email-template/invite-template.html", teacherTemplateData); err1 == nil {
		r1.SendEmail()
	}

	// Send invite student
	r2 := smtp.NewRequest(inviteToClassroomInfo.StudentEmailArray, "JOIN MY CLASS AS A STUDENT", "JOIN MY CLASS AS A STUDENT")
	if err2 := r2.ParseTemplate("./public/assets/email-template/invite-template.html", studentTemplateData); err2 == nil {
		r2.SendEmail()
	}

	return true, base.CodeSuccess, nil
}

func MethodExportStudentListByClassroomID(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))
	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))

	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	// Check user is a teacher in class
	ok, _ := MiddlewareImplementUserIsATeacherInClassroom(user.ID, classroom.ID)
	if !ok {
		return false, base.CodeBadRequest, nil
	}

	classroom.GetListStudent()

	fileUrl := export_excel.ProcessExportStudent(classroom.StudentArray)

	return true, base.CodeSuccess, fileUrl
}

func MethodImportStudentListByClassroomID(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	classroomID := util.ToUint(c.Param("id"))
	var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))

	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	// Check user is a teacher in class
	ok := MiddlewareImplementUserIsAnOwnerOfClassroom(user.ID, classroom)
	if !ok {
		return false, base.CodeBadRequest, nil
	}

	file, header, errFile := c.Request.FormFile("import-student-file")
	if errFile != nil {
		return false, base.CodeImportStudentFail, nil
	}

	fileBytes, err := ioutil.ReadAll(file)
	util.CheckErr(err)

	filePath := upload.WithTemporary().Save(header.Filename, fileBytes)

	biz := import_excel.SheetStudentStruct{}.Initialize(filePath, classroom.ID)

	ok, importResponseArray := biz.Importing()

	if !ok {
		return false, base.CodeImportStudentFail, nil
	}

	return true, base.CodeSuccess, importResponseArray
}
