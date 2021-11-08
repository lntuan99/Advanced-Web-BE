package model

import "github.com/jinzhu/gorm"

type Account struct {
	gorm.Model
	Username string `gorm:"index:account_user_name_idx"`
	Password string
	UserID   uint
	User     User
}
