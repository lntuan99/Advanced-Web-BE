package req_res

type RespondUserLogin struct {
	Token           string `json:"token"`
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	AvatarURL       string `json:"avatarUrl"`
	IsEmailVerified bool   `json:"isEmailVerified"`
}
