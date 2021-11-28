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
type GradeRes struct {
	ID            uint    `json:"id"`
	ClassroomID   uint    `json:"classroomId"`
	Name          string  `json:"name"`
	MaxPoint      float32 `json:"maxPoint"`
	Percent       uint    `json:"percent"`
	OrdinalNumber uint    `json:"ordinalNumber"`
}

func (grade Grade) ToRes() GradeRes {
	return GradeRes{
		ID:            grade.ID,
		ClassroomID:   grade.ClassroomID,
		Name:          grade.Name,
		MaxPoint:      grade.MaxPoint,
		Percent:       grade.Percent,
		OrdinalNumber: grade.OrdinalNumber,
	}

}
