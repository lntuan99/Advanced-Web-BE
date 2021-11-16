package methods

import (
    "advanced-web.hcmus/api/base"
    req_res "advanced-web.hcmus/api/req_res_struct"
    "advanced-web.hcmus/model"
    "advanced-web.hcmus/util"
    "fmt"
    "github.com/gin-gonic/gin"
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


func MethodGetClassroomByID(c *gin.Context) (bool, string, interface{}) {
    classroomID := util.ToUint(c.Param("id"))
    var classroom = model.Classroom{}.FindClassroomByID(uint(classroomID))

    if classroom.ID == 0 {
        return false, base.CodeClassroomIDNotExisted, nil
    }

    classroom.StudentArray = classroom.GetListUserByJWTType(model.JWT_TYPE_STUDENT)
    classroom.TeacherArray = classroom.GetListUserByJWTType(model.JWT_TYPE_TEACHER)

    return true, base.CodeSuccess, classroom.ToRes()
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

    existedClassroomCode := model.Classroom{}.FindClassroomByCodeBelongToOwner(classroomInfo.Code, user.ID)

    if existedClassroomCode.ID > 0 {
        return false, base.CodeClassroomCodeExisted, nil
    }

    var newClassroom = model.Classroom {
        OwnerID:           user.ID,
        Name:              classroomInfo.Name,
        CoverImageURL:     "",
        Code:              classroomInfo.Code,
        Description:       classroomInfo.Description,
    }
    newClassroom.GenerateInviteLink()

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

    return true, base.CodeSuccess, nil
}