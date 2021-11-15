package methods

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	api_jwt "advanced-web.hcmus/api/api-jwt"
	"advanced-web.hcmus/api/base"
	req_res "advanced-web.hcmus/api/req_res_struct"
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func MethodRegisterAccount(c *gin.Context) (bool, string, interface{}) {
	var registerAccountInfo req_res.PostRegisterAccount

	_ = c.Request.ParseMultipartForm(20971520)

	if err := c.ShouldBindJSON(&registerAccountInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	existedUsername := model.Account{}.FindAccountByUsername(registerAccountInfo.Username)
	if existedUsername.ID > 0 {
		return false, base.CodeUsernameExisted, nil
	}

	registerAccountInfo.Password = util.StandardizedString(registerAccountInfo.Password)
	if util.EmptyOrBlankString(registerAccountInfo.Password) {
		return false, base.CodeEmptyPassword, nil
	}

	if registerAccountInfo.Password != registerAccountInfo.RetypePassword {
		return false, base.CodePasswordAndRetypeDoesNotMatch, nil
	}

	newAccount, success := model.Account{}.ConvertPostRegisterAccountToModelAccount(registerAccountInfo)
	if !success {
		return false, base.CodeRegisterAccountFail, nil
	}

	// if err := model.DBInstance.Create(&newAccount).Error; err != nil {
	// 	return false, base.CodeRegisterAccountFail, nil
	// }
	gender := 0
	if registerAccountInfo.Gender == "2" {
		gender = model.GENDER_FEMALE
	} else if registerAccountInfo.Gender == "1" {
		gender = model.GENDER_MALE
	} else {
		gender = model.GENDER_UNKNOWN
	}

	var birthday time.Time
	if registerAccountInfo.Birthday != "" {
		tm, err := strconv.ParseInt(registerAccountInfo.Birthday, 10, 64)
		if err != nil {
			return false, base.CodeInvalidDate, nil
		}
		birthday = time.Unix(tm, 0)
	}

	if !util.IsEmailValid(registerAccountInfo.Email) {
		return false, base.CodeInvalidEmailFormat, nil
	}

	var newUser = model.User{
		Name:         registerAccountInfo.Name,
		Code:         registerAccountInfo.Code,
		Email:        registerAccountInfo.Email,
		Phone:        registerAccountInfo.Phone,
		Gender:       uint(gender),
		Birthday:     birthday,
		IdentityCard: registerAccountInfo.IdentityCard,
		Enabled:      true,
	}
	if err := model.DBInstance.Create(&newUser).Error; err != nil {
		return false, base.CodeRegisterAccountFail, nil
	}
	model.DBInstance.Model(&newAccount).Update("user_id", newUser.ID)

	// FormFile returns the first file for the given key `avatar`
	_, header, errFile := c.Request.FormFile("avatar")
	if errFile == nil {
		newFileName := fmt.Sprintf("%v%v", time.Now().Unix(), filepath.Ext(header.Filename))
		folderDst := fmt.Sprintf("/system/users/%v", newUser.ID)

		util.CreateFolder(folderDst)

		fileDst := fmt.Sprintf("%v/%v", folderDst, newFileName)
		if err := util.SaveUploadedFile(header, folderDst, newFileName); err == nil {
			model.DBInstance.
				Model(&newUser).
				Update("avatar", fileDst)
		}
	}
	return true, base.CodeSuccess, nil
}

func MethodLoginAccount(c *gin.Context) (bool, string, interface{}) {
	var loginAccountInfo req_res.PostLoginAccount
	if err := c.ShouldBindJSON(&loginAccountInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	existedUsername := model.Account{}.FindAccountByUsername(loginAccountInfo.Username)

	if existedUsername.ID == 0 {
		return false, base.CodeUsernameNotExisted, nil
	}

	loginAccountInfo.Password = util.StandardizedString(loginAccountInfo.Password)
	if util.EmptyOrBlankString(loginAccountInfo.Password) {
		return false, base.CodeEmptyPassword, nil
	}

	if !util.CompareHashingPasswordAndPassWord(existedUsername.Password, loginAccountInfo.Password) {
		return false, base.CodeWrongPassword, nil
	}

	result := req_res.RespondUserLogin{
		Token:     "",
		ID:        0,
		Name:      "",
		AvatarURL: "",
	}

	_, isExpired, user := model.User{}.FindUserByID(existedUsername.UserID)
	if isExpired {
		return false, base.CodeExpiredUserAccount, result
	}

	mw := api_jwt.GwtAuthMiddleware
	_ = mw.MiddlewareInit()

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)
	if claims["id"] == nil {
		claims["id"] = user.ID
	}
	expire := mw.TimeFunc().Add(mw.Timeout)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.SignedString(token)

	if err != nil {
		return false, base.CodeLoginAccountFail, result
	}

	avatarURL := ""
	if user.Avatar != "" {
		avatarURL = util.SubUrlToFullUrl(avatarURL)
	}

	result = req_res.RespondUserLogin{
		Token:     tokenString,
		ID:        user.ID,
		Name:      user.Name,
		AvatarURL: avatarURL,
	}

	return true, base.CodeSuccess, result
}
