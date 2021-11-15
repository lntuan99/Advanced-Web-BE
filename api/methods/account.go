package methods

import (
	api_jwt "advanced-web.hcmus/api/api-jwt"
	"advanced-web.hcmus/api/base"
	req_res "advanced-web.hcmus/api/req_res_struct"
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"time"
)

func MethodRegisterAccount(c *gin.Context) (bool, string, interface{}) {
	var registerAccountInfo req_res.PostRegisterAccount

	_ = c.Request.ParseMultipartForm(20971520)

	if err := c.ShouldBind(&registerAccountInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	existedAccountUsername := model.Account{}.FindAccountByUsername(registerAccountInfo.Username)
	if existedAccountUsername.ID > 0 && existedAccountUsername.UserID > 0 {
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

	newAccount.ID = existedAccountUsername.ID
	if err := model.DBInstance.Save(&newAccount).Error; err != nil {
		return false, base.CodeRegisterAccountFail, nil
	}

	if util.EmptyOrBlankString(registerAccountInfo.Name) {
		return false, base.CodeNameOfUserEmpty, nil
	}

	existedCodeUser, isExpired, _ := model.User{}.FindUserByCode(registerAccountInfo.Code)
	if existedCodeUser && !isExpired{
		return false, base.CodeUserCodeExisted, nil
	}

	if util.EmptyOrBlankString(registerAccountInfo.Email) {
		return false, base.CodeEmptyEmail, nil
	}
	if !util.IsEmailValid(registerAccountInfo.Email) {
		return false, base.CodeInvalidEmailFormat, nil
	}
	existedEmailUser, isExpired, _ := model.User{}.FindUserByEmail(registerAccountInfo.Email)
	if existedEmailUser && !isExpired {
		return false, base.CodeEmailExisted, nil
	}

	registerAccountInfo.Phone = util.FormatPhoneNumber(registerAccountInfo.Phone)
	existedPhoneUser, isExpired, _ := model.User{}.FindUserByPhone(registerAccountInfo.Phone)
	if existedPhoneUser && !isExpired{
		return false, base.CodePhoneExisted, nil
	}

	if !util.EmptyOrBlankString(registerAccountInfo.IdentityCard) {
		existedIdentityCardUser, isExpired, _ := model.User{}.FindUserByIdentityCard(registerAccountInfo.IdentityCard)
		if existedIdentityCardUser && !isExpired{
			return false, base.CodeIdentityCardExisted, nil
		}
	}

	var newUser = model.User{
		Name:         registerAccountInfo.Name,
		Code:         registerAccountInfo.Code,
		Email:        registerAccountInfo.Email,
		Phone:        registerAccountInfo.Phone,
		Gender:       registerAccountInfo.Gender,
		Birthday:     time.Unix(registerAccountInfo.Birthday, 0),
		IdentityCard: registerAccountInfo.IdentityCard,
		Enabled:      true,
		ExpiredAt:    nil,
	}

	if err := model.DBInstance.Create(&newUser).Error; err != nil {
		return false, base.CodeRegisterAccountFail, nil
	}

	model.DBInstance.
		Model(&newAccount).
		Update("user_id", newUser.ID)

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

	userLogin := generateUserToken(newUser)
	return true, base.CodeSuccess, userLogin
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

	userLogin := req_res.RespondUserLogin{
		Token:     "",
		ID:        0,
		Name:      "",
		AvatarURL: "",
	}

	_, isExpired, user := model.User{}.FindUserByID(existedUsername.UserID)
	if isExpired {
		return false, base.CodeExpiredUserAccount, userLogin
	}

	userLogin = generateUserToken(user)

	return true, base.CodeSuccess, userLogin
}

func generateUserToken(user model.User) req_res.RespondUserLogin {
	result := req_res.RespondUserLogin{
		Token:     "",
		ID:        0,
		Name:      "",
		AvatarURL: "",
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
		return result
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

	return result
}
