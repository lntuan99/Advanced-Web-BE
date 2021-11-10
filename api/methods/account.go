package methods

import (
	"advanced-web.hcmus/api/base"
	req_res "advanced-web.hcmus/api/req_res_struct"
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"github.com/gin-gonic/gin"
)

func MethodRegisterAccount(c *gin.Context) (bool, string, interface{}) {
	var registerAccountInfo req_res.PostRegisterAccount
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

	if err := model.DBInstance.Create(&newAccount).Error; err != nil {
		return false, base.CodeRegisterAccountFail, nil
	}

	return true, base.CodeSuccess, nil
}