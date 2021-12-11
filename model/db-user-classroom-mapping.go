package model

import "github.com/jinzhu/gorm"

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
	tx.First(&mapping.UserRole, mapping.UserRoleID)

	if mapping.UserRole.JWTType == JWT_TYPE_STUDENT {
		tx.Preload("GradeArray").
			First(&mapping.Classroom, mapping.ClassroomID)

		for _, grade := range mapping.Classroom.GradeArray {
			var studentGradeMapping = UserGradeMapping{
				UserID:  mapping.UserID,
				GradeID: grade.ID,
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
			tx.Where("user_id = ? AND grade_id = ?", mapping.UserID, grade.ID).
				Delete(&UserGradeMapping{})
		}
	}

	return nil
}
