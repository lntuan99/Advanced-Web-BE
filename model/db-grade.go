package model

import "github.com/jinzhu/gorm"

type Grade struct {
	gorm.Model
	ClassroomID   uint
	Classroom     Classroom
	Name          string
	MaxPoint      float32
	Percent       uint
	OrdinalNumber uint
}
