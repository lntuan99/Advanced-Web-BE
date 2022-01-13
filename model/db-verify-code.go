package model

import (
	"advanced-web.hcmus/config/constants"
	"advanced-web.hcmus/util"
	"fmt"
	"time"
)

const (
	VERIFY_CODE_TYPE_VERIFY_EMAIL    = 1
	VERIFY_CODE_TYPE_FORGOT_PASSWORD = 2
)

type VerifyCode struct {
	ID     uint   `gorm:"primary_key"`
	Code   string `gorm:"unique_index:user_verify_code_uid"`
	UserID uint   `gorm:"unique_index:user_verify_code_uid"`
	User   User
	Type   uint
}

func (VerifyCode) CreateVerifyCode(user User, verifyCodeType uint) VerifyCode {
	// Generate verify code
	verifyCode := fmt.Sprintf("%v_%v_%v_%v", user.Name, user.Code, user.Email, time.Now().Unix())
	verifyCode = util.HexSha256String([]byte(verifyCode))
	verifyCode += fmt.Sprintf("%v", time.Now().Unix()%constants.PRIME_NUMBER_FOR_MOD)

	return VerifyCode{
		Code:   verifyCode,
		UserID: user.ID,
		Type:   verifyCodeType,
	}
}
