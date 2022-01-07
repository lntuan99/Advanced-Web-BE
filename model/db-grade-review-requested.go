package model

import "github.com/jinzhu/gorm"

type GradeReviewRequested struct {
	gorm.Model
	StudentGradeMappingID uint
	StudentGradeMapping   StudentGradeMapping
	StudentExpectation    float32
	StudentExplanation    string
	TeacherComment        string
}

type GradeReviewRequestedRes struct {
	StudentRes         StudentRes `json:"student"`
	GradeRes           GradeRes   `json:"grade"`
	CurrentPoint       float32    `json:"currentPoint"`
	StudentExpectation float32    `json:"studentExpectation"`
	StudentExplanation string     `json:"studentExplanation"`
	TeacherComment     string     `json:"teacherComment"`
}

func (review GradeReviewRequested) ToRes() GradeReviewRequestedRes {
	// Find student grade mapping
	return GradeReviewRequestedRes{
		StudentRes:         review.StudentGradeMapping.Student.ToRes(),
		GradeRes:           review.StudentGradeMapping.Grade.ToRes(),
		CurrentPoint:       review.StudentGradeMapping.Point,
		StudentExpectation: review.StudentExpectation,
		StudentExplanation: review.StudentExplanation,
		TeacherComment:     review.TeacherComment,
	}

}
