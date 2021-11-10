package model

import "github.com/jinzhu/gorm"

type UserRole struct {
	gorm.Model
	Name       string `gorm:"index:user_role_name_idx"`
	Permission string `gorm:"type:jsonb"`
}
