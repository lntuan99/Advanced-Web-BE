package model

import "github.com/jinzhu/gorm"

const (
	JWT_TYPE_TEACHER = 1
	JWT_TYPE_STUDENT = 2
)

type UserRole struct {
	gorm.Model
	JWTType    uint
	Name       string `gorm:"index:user_role_name_idx"`
	Permission string `gorm:"type:jsonb"`
}

func (UserRole) GetRoleByJWTType(JWTType uint) UserRole {
	var res UserRole
	DBInstance.First(&res, "jwt_type = ?", JWTType)

	return res
}