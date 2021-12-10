package methods

import (
	"advanced-web.hcmus/model"
)

func MiddlewareImplementUserInClassroom(userID uint, classroomID uint) (ok bool, mapping model.UserClassroomMapping) {
	mapping = model.UserClassroomMapping{}
	model.DBInstance.
		Preload("UserRole").
		First(&mapping, "user_id = ? AND classroom_id = ?", userID, classroomID)

	if mapping.ID == 0 {
		return false, mapping
	}

	return true, mapping
}

func MiddlewareImplementUserIsATeacherInClassroom(userID uint, classroomID uint) (ok bool, mapping model.UserClassroomMapping) {
	mapping = model.UserClassroomMapping{}

	model.DBInstance.
		Preload("UserRole").
		First(&mapping, "user_id = ? AND classroom_id = ?", userID, classroomID)

	if mapping.ID == 0 {
		return false, mapping
	}

	if mapping.UserRole.JWTType != model.JWT_TYPE_TEACHER {
		return false, mapping
	}

	return true, mapping
}

func MiddlewareImplementUserIsAnOwnerOfClassroom(userID uint, classroom model.Classroom) (ok bool) {
	return classroom.OwnerID == userID
}
