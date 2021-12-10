package model

type UserGradeMapping struct {
	ID      uint `gorm:"primary_key"`
	UserID  uint `gorm:"index:user_grade_in_mapping_idx"`
	User    User
	GradeID uint `gorm:"index:user_grade_in_mapping_idx"`
	Grade   Grade
	Point   float32
}
