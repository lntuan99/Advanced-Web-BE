package req_res

type PostCreateUpdateUserProfile struct {
	Username     string `form:"username" json:"username"`
	Name         string `form:"name" json:"name"`
	Code         string `form:"code" json:"code"`
	Email        string `form:"email" json:"email"`
	Phone        string `form:"phone" json:"phone"`
	Birthday     int64  `form:"birthday" json:"birthday"`
	Gender       uint   `form:"gender" json:"gender"`
	IdentityCard string `form:"identityCard" json:"identityCard"`
}

type PostRegisterAccount struct {
	Password       string `form:"password" json:"password"`
	RetypePassword string `form:"retypePassword" json:"retypePassword"`
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
	ClassroomID   uint    `json:"classroomId"`
	Name          string  `json:"name"`
	MaxPoint      float32 `json:"maxPoint"`
	OrdinalNumber uint    `json:"ordinalNumber"`
}

type PostUpdateGrade struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	MaxPoint      float32 `json:"maxPoint"`
	OrdinalNumber uint    `json:"ordinalNumber"`
}
