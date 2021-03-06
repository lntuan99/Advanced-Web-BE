package methods

import (
	api_jwt "advanced-web.hcmus/api/api-jwt"
	"advanced-web.hcmus/api/base"
	req_res "advanced-web.hcmus/api/req_res_struct"
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"path/filepath"
	"strings"
	"time"
)

//====================================================
//====================================================
//============ MANAGE ADMIN FUNCTIONS ================
//====================================================
//====================================================

func MethodLoginAdminUser(c *gin.Context) (bool, string, interface{}) {
	var loginAccountInfo req_res.PostLoginAccount
	if err := c.ShouldBindJSON(&loginAccountInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	existedAdminUsername := model.AdminUser{}.FindAccountByUsername(loginAccountInfo.Username)

	if existedAdminUsername.ID == 0 {
		return false, base.CodeUsernameNotExisted, nil
	}

	loginAccountInfo.Password = util.StandardizedString(loginAccountInfo.Password)
	if util.EmptyOrBlankString(loginAccountInfo.Password) {
		return false, base.CodeEmptyPassword, nil
	}

	if !util.CompareHashingPasswordAndPassword(existedAdminUsername.Password, loginAccountInfo.Password) {
		return false, base.CodeWrongPassword, nil
	}

	adminLogin := req_res.RespondUserLogin{
		Token:           "",
		ID:              0,
		Name:            "",
		AvatarURL:       "",
		IsEmailVerified: false,
	}

	adminLogin = generateAdminToken(existedAdminUsername)

	return true, base.CodeSuccess, adminLogin
}

func MethodGetListAdminUser(c *gin.Context) (bool, string, interface{}) {
	var dbInstance = model.DBInstance

	var adminArray = make([]model.AdminUser, 0)

	var orderBy = "created_at ASC"
	if strings.ToLower(c.Query("sort")) == "desc" {
		orderBy = "created_at DESC"
	}

	var key = c.Query("key")
	if !util.EmptyOrBlankString(key) {
		formattedKeyword := util.RemoveAccent(util.TrimSpace(key))
		formattedKeyword = "%" + strings.ToLower(formattedKeyword) + "%"
		dbInstance = dbInstance.Where("f_unaccent(name) LIKE ? OR f_unaccent(email) LIKE ?", formattedKeyword, formattedKeyword)
	}

	dbInstance.
		Order(orderBy).
		Offset(base.GetIntQuery(c, "page") * base.PageSizeLimit).
		Limit(base.PageSizeLimit).
		Find(&adminArray)

	var adminResArray = make([]model.AdminUserRes, 0)
	for _, admin := range adminArray {
		adminResArray = append(adminResArray, admin.ToRes())
	}

	return true, base.CodeSuccess, adminResArray
}

func MethodGetAdminUserByID(c *gin.Context) (bool, string, interface{}) {
	adminUserID := util.ToUint(c.Param("id"))
	var adminUser = model.AdminUser{}.FindAdminUserByID(uint(adminUserID))

	if adminUser.ID == 0 {
		return false, base.CodeAdminUserIDNotExisted, nil
	}

	return true, base.CodeSuccess, adminUser.ToRes()
}

func MethodCreateAdminUser(c *gin.Context) (bool, string, interface{}) {
	_ = c.Request.ParseMultipartForm(20971520)

	var createAdminUserInfo req_res.PostCreateAdminUser
	if err := c.ShouldBind(&createAdminUserInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	// Check admin username valid
	existedAdminUsername := model.AdminUser{}.FindAccountByUsername(createAdminUserInfo.Username)
	if existedAdminUsername.ID > 0 {
		return false, base.CodeUsernameExisted, nil
	}

	// Check password valid
	if util.EmptyOrBlankString(util.StandardizedString(createAdminUserInfo.Password)) {
		return false, base.CodeEmptyPassword, nil
	}

	if createAdminUserInfo.Password != createAdminUserInfo.RetypePassword {
		return false, base.CodePasswordAndRetypeDoesNotMatch, nil
	}

	if util.EmptyOrBlankString(createAdminUserInfo.Name) {
		return false, base.CodeNameOfUserEmpty, nil
	}

	// Check email of user valid
	if util.EmptyOrBlankString(createAdminUserInfo.Email) {
		return false, base.CodeEmptyEmail, nil
	}
	if !util.IsEmailValid(createAdminUserInfo.Email) {
		return false, base.CodeInvalidEmailFormat, nil
	}
	existedAdminEmailUser := model.AdminUser{}.FindAdminUserByEmail(createAdminUserInfo.Email)
	if existedAdminEmailUser.ID > 0 {
		return false, base.CodeEmailExisted, nil
	}

	// Check phone of user valid
	phone := util.FormatPhoneNumber(createAdminUserInfo.Phone)
	if util.EmptyOrBlankString(phone) && !util.EmptyOrBlankString(createAdminUserInfo.Phone) {
		return false, base.CodePhoneInvalid, nil
	}
	createAdminUserInfo.Phone = phone
	if !util.EmptyOrBlankString(createAdminUserInfo.Phone) {
		existedAdminPhoneUser := model.AdminUser{}.FindAdminUserByPhone(createAdminUserInfo.Phone)
		if existedAdminPhoneUser.ID > 0 {
			return false, base.CodePhoneExisted, nil
		}
	}

	newAdminUser, success := model.AdminUser{}.ConvertPostRegisterAccountToModelAccount(createAdminUserInfo)
	if !success {
		return false, base.CodeRegisterAccountFail, nil
	}

	code := base.CodeSuccess
	err := model.DBInstance.Transaction(func(tx *gorm.DB) error {
		if err1 := tx.Create(&newAdminUser).Error; err1 != nil {
			code = base.CodeRegisterAccountFail
			return err1
		}

		// FormFile returns the first file for the given key `avatar`
		_, header, errFile := c.Request.FormFile("avatar")
		if errFile == nil {
			newFileName := fmt.Sprintf("%v%v", time.Now().Unix(), filepath.Ext(header.Filename))
			folderDst := fmt.Sprintf("/system/admins/%v", newAdminUser.ID)

			util.CreateFolder(folderDst)

			fileDst := fmt.Sprintf("%v/%v", folderDst, newFileName)
			if err := util.SaveUploadedFile(header, folderDst, newFileName); err == nil {
				tx.Model(&newAdminUser).
					Update("avatar", fileDst)
			}
		}

		return nil
	})

	if err != nil {
		return false, code, nil
	}

	return true, base.CodeSuccess, newAdminUser.ToRes()
}

//====================================================
//====================================================
//============= MANAGE USER FUNCTIONS ================
//====================================================
//====================================================

func MethodGetListUser(c *gin.Context) (bool, string, interface{}) {
	var dbInstance = model.DBInstance

	var userArray = make([]model.User, 0)

	var orderBy = "created_at ASC"
	if strings.ToLower(c.Query("sort")) == "desc" {
		orderBy = "created_at DESC"
	}

	var key = c.Query("key")
	if !util.EmptyOrBlankString(key) {
		formattedKeyword := util.RemoveAccent(util.TrimSpace(key))
		formattedKeyword = "%" + strings.ToLower(formattedKeyword) + "%"
		dbInstance = dbInstance.Where("f_unaccent(name) LIKE ? OR f_unaccent(email) LIKE ?", formattedKeyword, formattedKeyword)
	}

	dbInstance.
		Order(orderBy).
		Offset(base.GetIntQuery(c, "page") * base.PageSizeLimit).
		Limit(base.PageSizeLimit).
		Find(&userArray)

	var userResArray = make([]model.UserRes, 0)
	for _, user := range userArray {
		userResArray = append(userResArray, user.ToRes())
	}

	return true, base.CodeSuccess, userResArray
}

func MethodAdminGetUserByID(c *gin.Context) (bool, string, interface{}) {
	userID := util.ToUint(c.Param("id"))
	existedUser, _, user := model.User{}.FindUserByID(uint(userID))

	if !existedUser {
		return false, base.CodeUserNotExisted, nil
	}

	return true, base.CodeSuccess, user.ToRes()
}

func MethodAdminBanUserByID(c *gin.Context) (bool, string, interface{}) {
	userID := util.ToUint(c.Param("id"))
	existedUser, _, user := model.User{}.FindUserByID(uint(userID))

	if !existedUser {
		return false, base.CodeUserNotExisted, nil
	}

	if user.Enabled == false {
		return false, base.CodeUserAlreadyBanned, nil
	}

	if err := model.DBInstance.Model(&user).
		Updates(map[string]interface{}{
			"enabled": false,
		}).Error; err != nil {
		return false, base.CodeBanUserFail, nil
	}

	return true, base.CodeSuccess, user.ToRes()
}

func MethodAdminUnBanUserByID(c *gin.Context) (bool, string, interface{}) {
	userID := util.ToUint(c.Param("id"))
	existedUser, _, user := model.User{}.FindUserByID(uint(userID))

	if !existedUser {
		return false, base.CodeUserNotExisted, nil
	}

	if user.Enabled == true {
		return false, base.CodeUserAlreadyEnabled, nil
	}

	if err := model.DBInstance.Model(&user).
		Updates(map[string]interface{}{
			"enabled": true,
		}).Error; err != nil {
		return false, base.CodeUnBanUserFail, nil
	}

	return true, base.CodeSuccess, user.ToRes()
}

func MethodMapStudentCode(c *gin.Context) (bool, string, interface{}) {
	var mapStudentCodeInfo req_res.PostMapStudentCode
	if err := c.ShouldBind(&mapStudentCodeInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	existedUser, _, user := model.User{}.FindUserByID(mapStudentCodeInfo.UserID)

	if !existedUser {
		return false, base.CodeUserNotExisted, nil
	}

	err := model.DBInstance.Transaction(func(tx *gorm.DB) error {
		_, _, existedUserCode := model.User{}.FindUserByCode(mapStudentCodeInfo.StudentCode)
		if existedUserCode.ID > 0 {
			err1 := tx.Model(&existedUserCode).
				Updates(map[string]interface{}{
					"code": nil,
				}).Error

			if err1 != nil {
				return err1
			}
		}

		err2 := tx.Model(&user).
			Updates(map[string]interface{}{
				"code":           mapStudentCodeInfo.StudentCode,
				"is_code_locked": mapStudentCodeInfo.IsCodeLocked,
			}).Error

		if err2 != nil {
			return err2
		}

		return nil
	})

	if err != nil {
		return false, base.CodeMapStudentCodeFail, nil
	}

	return true, base.CodeSuccess, user.ToRes()
}

//====================================================
//====================================================
//========== MANAGE CLASSROOM FUNCTIONS ==============
//====================================================
//====================================================
func MethodAdminGetListClassroom(c *gin.Context) (bool, string, interface{}) {
	var dbInstance = model.DBInstance

	var classroomArray = make([]model.Classroom, 0)

	var orderBy = "created_at ASC"
	if strings.ToLower(c.Query("sort")) == "desc" {
		orderBy = "created_at DESC"
	}

	var key = c.Query("key")
	if !util.EmptyOrBlankString(key) {
		formattedKeyword := util.RemoveAccent(util.TrimSpace(key))
		formattedKeyword = "%" + strings.ToLower(formattedKeyword) + "%"
		dbInstance = dbInstance.Where("f_unaccent(name) LIKE ?", formattedKeyword)
	}

	dbInstance.
		Order(orderBy).
		Offset(base.GetIntQuery(c, "page") * base.PageSizeLimit).
		Limit(base.PageSizeLimit).
		Preload("Owner").
		Find(&classroomArray)

	var classroomResArray = make([]model.ClassroomResLite, 0)
	for _, classroom := range classroomArray {
		classroomResArray = append(classroomResArray, classroom.ToResLite())
	}

	return true, base.CodeSuccess, classroomResArray
}

func MethodAdminGetClassroomByID(c *gin.Context) (bool, string, interface{}) {
	classroomID := util.ToUint(c.Param("id"))
	classroom := model.Classroom{}.FindClassroomByID(uint(classroomID))

	if classroom.ID == 0 {
		return false, base.CodeClassroomIDNotExisted, nil
	}

	return true, base.CodeSuccess, classroom.ToRes(model.UserClassroomMapping{})
}

//====================================================
//====================================================
//============== PRIVATE FUNCTION ====================
//====================================================
//====================================================
func generateAdminToken(admin model.AdminUser) req_res.RespondUserLogin {
	result := req_res.RespondUserLogin{
		Token:           "",
		ID:              0,
		Name:            "",
		AvatarURL:       "",
		IsEmailVerified: false,
	}

	mw := api_jwt.GwtAuthMiddleware
	_ = mw.MiddlewareInit()

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)
	if claims["id"] == nil {
		claims["id"] = admin.ID
	}
	expire := mw.TimeFunc().Add(mw.Timeout)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.SignedString(token)

	if err != nil {
		return result
	}

	result = req_res.RespondUserLogin{
		Token:           tokenString,
		ID:              admin.ID,
		Name:            admin.Name,
		AvatarURL:       util.SubUrlToFullUrl(admin.Avatar),
		IsEmailVerified: admin.IsEmailVerified,
	}

	return result
}
