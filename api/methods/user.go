package methods

import (
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
	updateUserProfileInfo.Phone = util.FormatPhoneNumber(updateUserProfileInfo.Phone)
	phone := util.FormatPhoneNumber(updateUserProfileInfo.Phone)
	if util.EmptyOrBlankString(phone) && !util.EmptyOrBlankString(updateUserProfileInfo.Phone) {
		return false, base.CodePhoneInvalid, nil
	}
	updateUserProfileInfo.Phone = phone
	_, isExpired, existedPhoneUser := model.User{}.FindUserByPhone(updateUserProfileInfo.Phone)
	if existedPhoneUser.ID > 0 && !isExpired && existedPhoneUser.ID != user.ID {
		return false, base.CodePhoneExisted, nil
	}

	// Check identity card of user valid
	if !util.EmptyOrBlankString(updateUserProfileInfo.IdentityCard) {
		_, isExpired, existedIdentityCardUser := model.User{}.FindUserByIdentityCard(updateUserProfileInfo.IdentityCard)
		if existedIdentityCardUser.ID > 0 && !isExpired && existedPhoneUser.ID != user.ID {
			return false, base.CodeIdentityCardExisted, nil
		}
	}

	var birthday time.Time
	if updateUserProfileInfo.Birthday > 0 {
		birthday = time.Unix(updateUserProfileInfo.Birthday, 0)
	}

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

	existed, isExpired, _ := model.User{}.FindUserByID(user.ID)
	if existed && !isExpired {
		return true, base.CodeSuccess, user.ToRes()
	}

	return false, base.CodeFindUserFail, nil
}
