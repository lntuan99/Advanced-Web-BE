package methods

import (
	"advanced-web.hcmus/config"
	"advanced-web.hcmus/services/smtp"
	"fmt"
	"path/filepath"
	"time"

	"advanced-web.hcmus/api/base"
	req_res "advanced-web.hcmus/api/req_res_struct"
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"github.com/gin-gonic/gin"
)

func MethodUpdateUserProfile(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	_ = c.Request.ParseMultipartForm(20971520)

	var updateUserProfileInfo req_res.PostCreateUpdateUserProfile
	if err := c.ShouldBind(&updateUserProfileInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	// Check username valid
	existedAccountUsername := model.Account{}.FindAccountByUsername(updateUserProfileInfo.Username)
	if existedAccountUsername.ID > 0 && existedAccountUsername.UserID > 0 && existedAccountUsername.UserID != user.ID {
		return false, base.CodeUsernameExisted, nil
	}

	// Check code of user valid
	_, isExpired, existedCodeUser := model.User{}.FindUserByCode(updateUserProfileInfo.Code)
	if existedCodeUser.ID > 0 && !isExpired && existedCodeUser.ID != user.ID {
		return false, base.CodeUserCodeExisted, nil
	}

	// Check email of user valid
	if util.EmptyOrBlankString(updateUserProfileInfo.Email) {
		return false, base.CodeEmptyEmail, nil
	}
	if !util.IsEmailValid(updateUserProfileInfo.Email) {
		return false, base.CodeInvalidEmailFormat, nil
	}
	_, isExpired, existedEmailUser := model.User{}.FindUserByEmail(updateUserProfileInfo.Email)
	if existedEmailUser.ID > 0 && !isExpired && existedEmailUser.ID != user.ID {
		return false, base.CodeEmailExisted, nil
	}

	// Check phone of user valid
	phone := util.FormatPhoneNumber(updateUserProfileInfo.Phone)
	if util.EmptyOrBlankString(phone) && !util.EmptyOrBlankString(updateUserProfileInfo.Phone) {
		return false, base.CodePhoneInvalid, nil
	}
	updateUserProfileInfo.Phone = phone
	if !util.EmptyOrBlankString(updateUserProfileInfo.Phone) {
		_, isExpired, existedPhoneUser := model.User{}.FindUserByPhone(updateUserProfileInfo.Phone)
		if existedPhoneUser.ID > 0 && !isExpired && existedPhoneUser.ID != user.ID {
			return false, base.CodePhoneExisted, nil
		}
	}

	// Check identity card of user valid
	if !util.EmptyOrBlankString(updateUserProfileInfo.IdentityCard) {
		_, isExpired, existedIdentityCardUser := model.User{}.FindUserByIdentityCard(updateUserProfileInfo.IdentityCard)
		if existedIdentityCardUser.ID > 0 && !isExpired && existedIdentityCardUser.ID != user.ID {
			return false, base.CodeIdentityCardExisted, nil
		}
	}

	var birthday = &time.Time{}
	if updateUserProfileInfo.Birthday != 0 {
		*birthday = time.Unix(updateUserProfileInfo.Birthday, 0)
	} else {
		birthday = nil
	}

	existedAccountUsername.Username = updateUserProfileInfo.Username
	user.Name = updateUserProfileInfo.Name
	user.Code = updateUserProfileInfo.Code
	user.Email = updateUserProfileInfo.Email
	user.Phone = updateUserProfileInfo.Phone
	user.Gender = updateUserProfileInfo.Gender
	user.Birthday = birthday
	user.IdentityCard = updateUserProfileInfo.IdentityCard

	// FormFile returns the first file for the given key `avatar`
	_, header, errFile := c.Request.FormFile("avatar")
	if errFile == nil {
		newFileName := fmt.Sprintf("%v%v", time.Now().Unix(), filepath.Ext(header.Filename))
		folderDst := fmt.Sprintf("/system/users/%v", user.ID)

		util.CreateFolder(folderDst)

		fileDst := fmt.Sprintf("%v/%v", folderDst, newFileName)
		if err := util.SaveUploadedFile(header, folderDst, newFileName); err == nil {
			user.Avatar = fileDst
		}
	}

	if err := model.DBInstance.Save(&user).Error; err != nil {
		return false, base.CodeUpdateUserProfileFail, nil
	}

	userLogin := generateUserToken(user)

	return true, base.CodeUpdateUserProfileSuccess, userLogin
}

func MethodGetUserProfile(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	return true, base.CodeSuccess, user.ToRes()
}

func MethodVerifyCode(c *gin.Context) (bool, string, interface{}) {
	verifyCode := c.Query("code")

	// Validate verify code
	if util.EmptyOrBlankString(verifyCode) {
		return false, base.CodeInvalidVerifyCode, nil
	}

	var dbVerifyCode model.VerifyCode
	model.DBInstance.
		Preload("User").
		First(&dbVerifyCode, "code = ?", verifyCode)

	if dbVerifyCode.ID == 0 || dbVerifyCode.UserID == 0 || dbVerifyCode.User.ID == 0 {
		return false, base.CodeInvalidVerifyCode, nil
	}

	var data interface{} = nil
	if dbVerifyCode.Type == model.VERIFY_CODE_TYPE_VERIFY_EMAIL {
		if err := model.DBInstance.
			Model(&model.User{}).
			Where("id = ?", dbVerifyCode.UserID).
			Updates(model.User{IsEmailVerified: true}).Error; err != nil {
			return false, base.CodeVerifyEmailFail, nil
		}
	}

	if dbVerifyCode.Type == model.VERIFY_CODE_TYPE_FORGOT_PASSWORD {
		userLogin := generateUserToken(dbVerifyCode.User)
		data = userLogin
	}

	model.DBInstance.Delete(&dbVerifyCode)

	return true, base.CodeSuccess, data
}

func MethodForgotPassword(c *gin.Context) (bool, string, interface{}) {
	email := c.Query("email")

	// validate existed user email
	if util.EmptyOrBlankString(email) {
		return false, base.CodeEmptyEmail, nil
	}
	if !util.IsEmailValid(email) {
		return false, base.CodeInvalidEmailFormat, nil
	}
	existed, isExpired, user := model.User{}.FindUserByEmail(email)
	if !existed || isExpired {
		return false, base.CodeFindUserFail, nil
	}

	if user.IsEmailVerified {
		verifyCode := model.VerifyCode{}.CreateVerifyCode(user, model.VERIFY_CODE_TYPE_FORGOT_PASSWORD)
		model.DBInstance.Create(&verifyCode)

		// Generate verify link
		verifyLink := fmt.Sprintf("%v/verify?code=%v", config.Config.FeDomain, verifyCode.Code)

		type TemplateData struct {
			URL string
		}
		forgotPasswordTemplate := TemplateData{
			URL: verifyLink,
		}

		// Send verify link
		r := smtp.NewRequest([]string{user.Email}, "RENEW YOUR PASSWORD", "RENEW YOUR PASSWORD")
		if err1 := r.ParseTemplate("./public/assets/email-template/forgot-password-template.html", forgotPasswordTemplate); err1 == nil {
			r.SendEmail()
		}
	}

	return true, base.CodeSuccess, nil
}

func MethodUpdatePassword(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	var forgotPassword req_res.PostForgotPassword
	if err := c.ShouldBind(&forgotPassword); err != nil {
		return false, base.CodeBadRequest, nil
	}

	// Check password valid
	if util.EmptyOrBlankString(util.StandardizedString(forgotPassword.Password)) {
		return false, base.CodeEmptyPassword, nil
	}

	if forgotPassword.Password != forgotPassword.RetypePassword {
		return false, base.CodePasswordAndRetypeDoesNotMatch, nil
	}

	// Find account mapping this user
	var dbAccount model.Account
	model.DBInstance.First(&dbAccount, "user_id = ?", user.ID)
	if dbAccount.ID == 0 {
		return false, base.CodeUpdatePasswordFail, nil
	}

	hashedPassword, ok := util.HashingPassword(forgotPassword.Password)
	if !ok {
		return false, base.CodeUpdatePasswordFail, nil
	}

	if err := model.DBInstance.
		Model(&dbAccount).
		Updates(model.Account{Password: hashedPassword}).
		Error; err != nil {
		return false, base.CodeUpdatePasswordFail, nil
	}

	return true, base.CodeSuccess, nil
}
