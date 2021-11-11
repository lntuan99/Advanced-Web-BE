package req_res

type PostRegisterAccount struct {
	Username       string `form:"username" json:"username"`
	Password       string `form:"password" json:"password"`
	RetypePassword string `form:"retypePassword" json:"retypePassword"`
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
