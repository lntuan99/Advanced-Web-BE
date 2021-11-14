package model

import (
	"time"

	"advanced-web.hcmus/util"
	"github.com/jinzhu/gorm"
)

const (
	GENDER_MALE    = 1
	GENDER_FEMALE  = 2
	GENDER_UNKNOWN = 0
)

type User struct {
	gorm.Model
	Name         string `gorm:"index:user_name_idx"`
	Code         string `gorm:"index:user_code_unique_idx"`
	Email        string `gorm:"index:user_email_unique_idx"`
	Phone        string `gorm:"index:user_phone_unique_idx"`
	Birthday     time.Time
	Gender       uint
	Avatar       string
	IdentityCard string
	Enabled      bool
	ExpiredAt    *time.Time
	Classrooms   []Classroom `gorm:"many2many:user_classroom_mappings"`
}
type UserInfor struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
}

func (user User) ToGetInfor(role string) UserInfor {
	return UserInfor{
		ID:     user.ID,
		Name:   user.Name,
		Code:   user.Code,
		Avatar: util.SubUrlToFullUrl(user.Avatar),
		Role:   role,
	}
}

//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
func (user User) InitializeTableConfig() {
	// "gin" means: The column must be of tsvector type
	DBInstance.Exec(`CREATE INDEX IF NOT EXISTS search_field
    ON users USING
    gin(search_field)`)

	DBInstance.Exec(`CREATE INDEX IF NOT EXISTS user_name_idx 
    ON users
    USING gin (f_unaccent(name) gin_trgm_ops)`)
}

func (User) FindUserByID(ID uint) (existed bool, isExpired bool, user User) {
	// Response: EXISTED_USER, IS_EXPIRED / DISABLED, USER INFO
	user = User{}
	DBInstance.
		Where("id = ? ", ID).
		First(&user)

	if user.ID == 0 {
		return false, false, user
	} else {
		if user.Enabled == false || (user.ExpiredAt != nil && user.ExpiredAt.Unix() <= time.Now().Unix()) {
			return true, true, user
		}
		return true, false, user
	}
}
