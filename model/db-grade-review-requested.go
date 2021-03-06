package model

import "github.com/jinzhu/gorm"

type GradeReviewRequested struct {
	gorm.Model
	StudentGradeMappingID uint
	StudentGradeMapping   StudentGradeMapping
	StudentExpectation    float32
	StudentExplanation    string
	Comments              []GradeReviewRequestedComment
	FinalPoint            *float32 `gorm:"default:null"`
	IsProcessed           bool
}

type GradeReviewRequestedRes struct {
	ID                 uint                             `json:"id"`
	StudentRes         StudentRes                       `json:"student"`
	GradeRes           GradeRes                         `json:"grade"`
	CurrentPoint       float32                          `json:"currentPoint"`
	StudentExpectation float32                          `json:"studentExpectation"`
	StudentExplanation string                           `json:"studentExplanation"`
	FinalPoint         float32                          `json:"finalPoint"`
	Comments           []GradeReviewRequestedCommentRes `json:"comments"`
	IsProcessed        bool                             `json:"isProcessed"`
}

func (review GradeReviewRequested) ToRes() GradeReviewRequestedRes {
	// Find student grade mapping
	if review.StudentGradeMapping.ID == 0 {
		DBInstance.
			Preload("Student").
			Preload("Grade").
			First(&review.StudentGradeMapping, review.StudentGradeMappingID)
	}

	var comments = make([]GradeReviewRequestedCommentRes, 0)
	for _, comment := range review.Comments {
		comments = append(comments, comment.ToRes())
	}

	finalPoint := float32(0)
	if review.FinalPoint != nil {
		finalPoint = *review.FinalPoint
	}

	return GradeReviewRequestedRes{
		ID:                 review.ID,
		StudentRes:         review.StudentGradeMapping.Student.ToRes(),
		GradeRes:           review.StudentGradeMapping.Grade.ToRes(),
		CurrentPoint:       review.StudentGradeMapping.Point,
		StudentExpectation: review.StudentExpectation,
		StudentExplanation: review.StudentExplanation,
		FinalPoint:         finalPoint,
		Comments:           comments,
		IsProcessed:        review.IsProcessed,
	}

}
