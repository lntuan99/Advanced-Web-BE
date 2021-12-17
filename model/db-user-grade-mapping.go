package model

type StudentGradeMapping struct {
	ID        uint `gorm:"primary_key"`
	StudentID uint `gorm:"index:student_grade_in_mapping_idx"`
	Student   Student
	GradeID   uint `gorm:"index:student_grade_in_mapping_idx"`
	Grade     Grade
	Point     float32
}

type StudentGradeMappingRes struct {
	GradeID       uint    `json:"id"`
	Name          string  `json:"name"`
	Point         float32 `json:"point"`
	MaxPoint      float32 `json:"maxPoint"`
	OrdinalNumber uint    `json:"ordinalNumber"`
}

type ResponseStudentGradeInClassroom struct {
	StudentRes
	StudentName   string                   `json:"studentName"`
	StudentCode   string                   `json:"studentCode"`
	TotalGrade    float32                  `json:"totalGrade"`
	MaxTotalGrade float32                  `json:"maxTotalGrade"`
	GradeArray    []StudentGradeMappingRes `json:"gradeArray"`
}

func (mapping StudentGradeMapping) ToRes() StudentGradeMappingRes {
	DBInstance.First(&mapping.Grade, mapping.GradeID)

	return StudentGradeMappingRes{
		GradeID:       mapping.GradeID,
		Name:          mapping.Grade.Name,
		Point:         mapping.Point,
		MaxPoint:      mapping.Grade.MaxPoint,
		OrdinalNumber: mapping.Grade.OrdinalNumber,
	}
}
