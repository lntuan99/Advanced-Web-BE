package model

import (
	"advanced-web.hcmus/util"
	"github.com/jinzhu/gorm"
)

type UserClassroomMapping struct {
	gorm.Model
	ClassroomID uint `gorm:"unique_index:user_classroom_role_unique_idx"`
	Classroom   Classroom
	UserID      uint `gorm:"unique_index:user_classroom_role_unique_idx"`
	User        User
	UserRoleID  uint `gorm:"unique_index:user_classroom_role_unique_idx"`
	UserRole    UserRole
}

func (mapping *UserClassroomMapping) AfterCreate(tx *gorm.DB) error {
	// If create mapping for student
	// Create mapping all grade in class
	// Create new record in students
	tx.First(&mapping.UserRole, mapping.UserRoleID)
	tx.First(&mapping.User, mapping.UserID)

	if mapping.UserRole.JWTType == JWT_TYPE_STUDENT {
		tx.Preload("GradeArray").
			First(&mapping.Classroom, mapping.ClassroomID)

		var existedStudent Student
		tx.First(&existedStudent, "classroom_id = ? AND code = ?", mapping.ClassroomID, mapping.User.Code)

		existedStudent.ClassroomID = mapping.ClassroomID
		existedStudent.Code = mapping.User.Code

		if util.EmptyOrBlankString(existedStudent.Name) {
			existedStudent.Name = mapping.User.Name
		}

		// Use save for create new if not existed or update if existed
		tx.Save(&existedStudent)

		for _, grade := range mapping.Classroom.GradeArray {
			var studentGradeMapping = StudentGradeMapping{
				StudentID: existedStudent.ID,
				GradeID:   grade.ID,
			}

			tx.Create(&studentGradeMapping)
		}
	}

	return nil
}

func (mapping *UserClassroomMapping) AfterDelete(tx *gorm.DB) error {
	// If create mapping for student
	// Create mapping all grade in class
	tx.First(&mapping.UserRole, mapping.UserRoleID)

	if mapping.UserRole.JWTType == JWT_TYPE_STUDENT {
		tx.Preload("GradeArray").
			First(&mapping.Classroom, mapping.ClassroomID)

		for _, grade := range mapping.Classroom.GradeArray {
			tx.Where("student_id = ? AND grade_id = ?", mapping.UserID, grade.ID).
				Delete(&StudentGradeMapping{})
		}
	}

	return nil
}
