package req_res

type RespondUserLogin struct {
	Token     string `json:"token"`
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	RoleId    uint   `json:"roleId"`
	AvatarURL string `json:"avatarUrl"`
}
