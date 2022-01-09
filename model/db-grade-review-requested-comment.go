package model

import "github.com/jinzhu/gorm"

type GradeReviewRequestedComment struct {
	gorm.Model
	GradeReviewRequestedID uint
	GradeReviewRequested   GradeReviewRequested
	UserID                 uint
	User                   User
	Comment                string
}

type GradeReviewRequestedCommentRes struct {
	UserRes UserRes `json:"user"`
	Comment string  `json:"comment"`
}

func (comment GradeReviewRequestedComment) ToRes() GradeReviewRequestedCommentRes {
	if comment.User.ID == 0 {
		DBInstance.First(&comment.User, comment.UserID)
	}

	return GradeReviewRequestedCommentRes{
		UserRes: comment.User.ToRes(),
		Comment: comment.Comment,
	}
}
