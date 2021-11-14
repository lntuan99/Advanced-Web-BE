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

func (Classroom) FindUsersByIDClassroom(id string) ([]UserInfor, []UserInfor) {
	var mappingArray = make([]UserClassroomMapping, 0)
	DBInstance.First(&mappingArray, "ClassroomID = ?", id)
	var students = make([]UserInfor, 0)
	var teachers = make([]UserInfor, 0)
	for _, classMapping := range mappingArray {
		if classMapping.UserRole.Permission == "teacher" {
			students = append(students, classMapping.User.ToGetInfor("teacher"))
		}
		teachers = append(teachers, classMapping.User.ToGetInfor("student"))
	}
	return students, teachers
}
