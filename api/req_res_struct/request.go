package req_res

type PostRegisterAccount struct {
	Username       string `form:"username" json:"username"`
	Password       string `form:"password" json:"password"`
	RetypePassword string `form:"retypePassword" json:"retypePassword"`
	Name           string `form:"name" json:"name"`
	Code           string `form:"code" json:"code"`
	Email          string `form:"email" json:"email"`
	Phone          string `form:"phone" json:"phone"`
	Birthday       string `form:"birthday" json:"birthday"`
	Gender         string `form:"gender" json:"gender"`
	IdentityCard   string `form:"identityCard" json:"identityCard"`
}

type PostLoginAccount struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type PostCreateClassroom struct {
	Name        string `form:"name" json:"name"`
	Code        string `form:"code" json:"code" `
	Description string `form:"description" json:"description" `
}
