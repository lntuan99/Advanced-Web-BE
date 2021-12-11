package model

type UserGradeMapping struct {
	ID      uint `gorm:"primary_key"`
	UserID  uint `gorm:"index:user_grade_in_mapping_idx"`
	User    User
	GradeID uint `gorm:"index:user_grade_in_mapping_idx"`
	Grade   Grade
	Point   float32
}

type UserGradeMappingRes struct {
	GradeID       uint    `json:"id"`
	Name          string  `json:"name"`
	Point         float32 `json:"point"`
	MaxPoint      float32 `json:"maxPoint"`
	OrdinalNumber uint    `json:"ordinalNumber"`
}

type ResponseStudentGradeInClassroom struct {
	UserRes
	GradeArray []UserGradeMappingRes
}

func (mapping UserGradeMapping) ToRes() UserGradeMappingRes {
	DBInstance.First(&mapping.Grade, mapping.GradeID)

	return UserGradeMappingRes{
		GradeID:       mapping.GradeID,
		Name:          mapping.Grade.Name,
		Point:         mapping.Point,
		MaxPoint:      mapping.Grade.MaxPoint,
		OrdinalNumber: mapping.Grade.OrdinalNumber,
	}
}
