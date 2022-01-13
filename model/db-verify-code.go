package model

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
