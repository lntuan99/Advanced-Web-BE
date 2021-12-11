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

func (grade *Grade) AfterCreate(tx *gorm.DB) error {
	// Create mapping for all student in classroom
	var dbClassroom = Classroom{}.FindClassroomByID(grade.ClassroomID)
	dbClassroom.StudentArray = Classroom{}.GetListUserByJWTType(JWT_TYPE_STUDENT)

	for _, student := range dbClassroom.StudentArray {
		var dbStudentGradeMapping = UserGradeMapping{
			UserID:  student.ID,
			GradeID: grade.ID,
		}
		tx.First(&dbStudentGradeMapping, "user_id = ? AND grade_id = ?", student.ID, grade.ID)

		if dbStudentGradeMapping.ID == 0 {
			tx.Create(&dbStudentGradeMapping)
		}
	}

	return nil
}

func (grade *Grade) AfterDelete(tx *gorm.DB) error {
	// clear all mapping for all student in classroom
	tx.Where("grade_id = ?", grade.ID).
		Delete(&UserGradeMapping{})

	return nil
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
