package req_res

type PostCreateUpdateUserProfile struct {
	Username     string `form:"username" json:"username" binding:"required"`
	Name         string `form:"name" json:"name" binding:"required"`
	Code         string `form:"code" json:"code" binding:"required"`
	Email        string `form:"email" json:"email" binding:"required"`
	Phone        string `form:"phone" json:"phone"`
	Birthday     int64  `form:"birthday" json:"birthday"`
	Gender       uint   `form:"gender" json:"gender"`
	IdentityCard string `form:"identityCard" json:"identityCard"`
}

type PostRegisterAccount struct {
	Password       string `form:"password" json:"password" binding:"required"`
	RetypePassword string `form:"retypePassword" json:"retypePassword" binding:"required"`
	GoogleID       string `form:"googleId" json:"googleId"`
	PostCreateUpdateUserProfile
}

type PostLoginAccount struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type PostGoogleLogin struct {
	GoogleID string `form:"googleId" json:"googleId"`
}

type PostCreateClassroom struct {
	Name        string `form:"name" json:"name"`
	Code        string `form:"code" json:"code" `
	Description string `form:"description" json:"description" `
}

type PostInviteToClassroom struct {
	ClassroomID       uint     `json:"classroomId"`
	TeacherEmailArray []string `json:"teacherEmailArray"`
	StudentEmailArray []string `json:"studentEmailArray"`
}

type PostCreateGrade struct {
	ClassroomID uint    `json:"classroomId"`
	Name        string  `json:"name"`
	MaxPoint    float32 `json:"maxPoint"`
}

type PostUpdateGrade struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	MaxPoint      float32 `json:"maxPoint"`
	OrdinalNumber uint    `json:"ordinalNumber"`
	IsFinalized   bool    `json:"isFinalized"`
}

type PostInputGradeForAStudent struct {
	StudentID string  `json:"studentId"`
	GradeID   uint    `json:"gradeId"`
	Point     float32 `json:"point"`
}

type PostExportGradeBoard struct {
	GradeIDArray []uint `json:"gradeIdArray"`
}
