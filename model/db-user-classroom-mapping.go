package model

type UserClassroomMapping struct {
	ClassroomID uint `gorm:"unique_index:user_classroom_role_unique_idx"`
	Classroom   Classroom
	UserID      uint `gorm:"unique_index:user_classroom_role_unique_idx"`
	User        User
	UserRoleID  uint `gorm:"unique_index:user_classroom_role_unique_idx"`
	UserRole    UserRole
}
