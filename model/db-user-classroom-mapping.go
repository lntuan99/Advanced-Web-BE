package model

import "github.com/jinzhu/gorm"

type UserClassroomMapping struct {
	gorm.Model
	ClassroomID uint `gorm:"unique_index:user_classroom_unique_idx"`
	Classroom   Classroom
	UserID      uint `gorm:"unique_index:user_classroom_unique_idx"`
	User        User
}
