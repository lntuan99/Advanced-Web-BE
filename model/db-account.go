package model

import (
	req_res "advanced-web.hcmus/api/req_res_struct"
	"advanced-web.hcmus/util"
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	Username string `gorm:"index:account_user_name_idx"`
	Password string
	UserID   uint
	User     User
}

//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
func (Account) FindAccountByUsername(username string) Account {
	var res Account
	DBInstance.First(&res, "username = ?", username)

	return res
}

func (Account) ConvertPostRegisterAccountToModelAccount(postAccount req_res.PostRegisterAccount) (Account, bool) {
	hashPassword, success := util.HashingPassword(postAccount.Password)

	newAccount := Account{
		Username: postAccount.Username,
		Password: hashPassword,
	}

	return newAccount, success
}