package model

import (
	req_res "advanced-web.hcmus/api/req_res_struct"
	"advanced-web.hcmus/util"
	"github.com/jinzhu/gorm"
)

type AdminUser struct {
	gorm.Model
	Username        string `gorm:"index:admin_user_name_unique_idx"`
	Password        string
	Name            string `gorm:"index:admin_name_idx"`
	Email           string `gorm:"index:admin_email_unique_idx"`
	Phone           string `gorm:"index:admin_phone_unique_idx"`
	Avatar          string
	IsEmailVerified bool `gorm:"default:false"`
}

type AdminUserRes struct {
	ID              uint   `json:"id"`
	CreatedAt       int64  `json:"createdAt"`
	Username        string `json:"username"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Avatar          string `json:"avatar"`
	IsEmailVerified bool   `json:"isEmailVerified"`
}

func (admin AdminUser) ToRes() AdminUserRes {
	return AdminUserRes{
		ID:              admin.ID,
		CreatedAt:       admin.CreatedAt.Unix(),
		Username:        admin.Username,
		Name:            admin.Name,
		Email:           admin.Email,
		Phone:           admin.Phone,
		Avatar:          util.SubUrlToFullUrl(admin.Avatar),
		IsEmailVerified: admin.IsEmailVerified,
	}
}

//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
func (AdminUser) FindAccountByUsername(username string) AdminUser {
	var res AdminUser
	DBInstance.First(&res, "username = ?", username)

	return res
}

func (AdminUser) FindAdminUserByID(ID uint) AdminUser {
	var res AdminUser
	DBInstance.First(&res, ID)

	return res
}

func (AdminUser) ConvertPostRegisterAccountToModelAccount(postCreateAdminUser req_res.PostCreateAdminUser) (AdminUser, bool) {
	hashPassword, success := util.HashingPassword(postCreateAdminUser.Password)

	newAccount := AdminUser{
		Username:        postCreateAdminUser.Username,
		Password:        hashPassword,
		Name:            postCreateAdminUser.Name,
		Email:           postCreateAdminUser.Email,
		Phone:           postCreateAdminUser.Phone,
		IsEmailVerified: false,
	}

	return newAccount, success
}

func (AdminUser) FindAdminUserByEmail(email string) AdminUser {
	res := AdminUser{}
	DBInstance.
		Where("email = ? ", email).
		First(&res)

	return res
}

func (AdminUser) FindAdminUserByPhone(phone string) AdminUser {
	res := AdminUser{}
	DBInstance.
		Where("phone = ? ", phone).
		First(&res)

	return res
}
