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
	"github.com/jinzhu/gorm"
	"path/filepath"
	"time"
)

func MethodRegisterAccount(c *gin.Context) (bool, string, interface{}) {
	_ = c.Request.ParseMultipartForm(20971520)

	var registerAccountInfo req_res.PostRegisterAccount
	if err := c.ShouldBind(&registerAccountInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	// Check username valid
	existedAccountUsername := model.Account{}.FindAccountByUsername(registerAccountInfo.Username)
	if existedAccountUsername.ID > 0 && existedAccountUsername.UserID > 0 {
		return false, base.CodeUsernameExisted, nil
	}

	// Check password valid
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

	if util.EmptyOrBlankString(registerAccountInfo.Name) {
		return false, base.CodeNameOfUserEmpty, nil
	}

	// Check code of user valid
	existedCodeUser, isExpired, _ := model.User{}.FindUserByCode(registerAccountInfo.Code)
	if existedCodeUser && !isExpired {
		return false, base.CodeUserCodeExisted, nil
	}

	// Check email of user valid
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

	// Check phone of user valid
	phone := util.FormatPhoneNumber(registerAccountInfo.Phone)
	if util.EmptyOrBlankString(phone) && !util.EmptyOrBlankString(registerAccountInfo.Phone){
		return false, base.CodePhoneInvalid, nil
	}
	registerAccountInfo.Phone = phone
	existedPhoneUser, isExpired, _ := model.User{}.FindUserByPhone(registerAccountInfo.Phone)
	if existedPhoneUser && !isExpired {
		return false, base.CodePhoneExisted, nil
	}

	// Check identity card of user valid
	if !util.EmptyOrBlankString(registerAccountInfo.IdentityCard) {
		existedIdentityCardUser, isExpired, _ := model.User{}.FindUserByIdentityCard(registerAccountInfo.IdentityCard)
		if existedIdentityCardUser && !isExpired {
			return false, base.CodeIdentityCardExisted, nil
		}
	}

	var newUser = model.User {
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

	code := base.CodeSuccess
	err := model.DBInstance.Transaction(func(tx *gorm.DB) error {
		if err1 := tx.Create(&newUser).Error; err1 != nil {
			code = base.CodeRegisterAccountFail
			return err1
		}

		// For case username existed but not user link with this account
		newAccount.ID = existedAccountUsername.ID
		newAccount.UserID = newUser.ID
		if err1 := tx.Save(&newAccount).Error; err1 != nil {
			code = base.CodeRegisterAccountFail
			return err1
		}

		// FormFile returns the first file for the given key `avatar`
		_, header, errFile := c.Request.FormFile("avatar")
		if errFile == nil {
			newFileName := fmt.Sprintf("%v%v", time.Now().Unix(), filepath.Ext(header.Filename))
			folderDst := fmt.Sprintf("/system/users/%v", newUser.ID)

			util.CreateFolder(folderDst)

			fileDst := fmt.Sprintf("%v/%v", folderDst, newFileName)
			if err := util.SaveUploadedFile(header, folderDst, newFileName); err == nil {
				tx.Model(&newUser).
					Update("avatar", fileDst)
			}
		}

		return nil
	})

	if err != nil {
		return false, code, nil
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

func MethodGoogleLogin(c *gin.Context) (bool, string, interface{}) {
	var googleLoginInfo req_res.PostGoogleLogin
	if err := c.ShouldBindJSON(&googleLoginInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	existedGoogleID := model.Account{}.FindAccountByGoogleID(googleLoginInfo.GoogleID)
	if existedGoogleID.ID > 0 {
		userLogin := req_res.RespondUserLogin{
			Token:     "",
			ID:        0,
			Name:      "",
			AvatarURL: "",
		}

		_, isExpired, user := model.User{}.FindUserByID(existedGoogleID.UserID)
		if isExpired {
			return false, base.CodeExpiredUserAccount, userLogin
		}

		userLogin = generateUserToken(user)

		return true, base.CodeSuccess, userLogin
	}

	return true, base.CodeGoogleIDNotExisted, nil
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
