package model

import (
	"advanced-web.hcmus/util"
	"time"

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

type UserRes struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Birthday     int64  `json:"birthday"` //Unix
	Gender       uint   `json:"gender"`
	Avatar       string `json:"avatar"`
	IdentityCard string `json:"identityCard"`
	Enabled      bool   `json:"enabled"`
	ExpiredAt    int64  `json:"expiredAt"`
}

func (user User) ToRes() UserRes {
	expiredAt := int64(0)
	if user.ExpiredAt != nil {
		expiredAt = user.ExpiredAt.Unix()
	}

	var userAccount Account
	DBInstance.First(&userAccount, "user_id = ?", user.ID)

	return UserRes {
		Username:     userAccount.Username,
		Name:         user.Name,
		Code:         user.Code,
		Email:        user.Email,
		Phone:        user.Phone,
		Birthday:     user.Birthday.Unix(),
		Gender:       user.Gender,
		Avatar:       util.SubUrlToFullUrl(user.Avatar),
		IdentityCard: user.IdentityCard,
		Enabled:      user.Enabled,
		ExpiredAt:    expiredAt,
	}
}

//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
func (user User) GetExpiredAt() int64 {
	if user.ExpiredAt == nil {
		return 0
	} else {
		return user.ExpiredAt.Unix()
	}
}

func (user *User) SetExpiredAt(expiredAt int64) {
	user.ExpiredAt = new(time.Time)
	if expiredAt <= 0 {
		user.ExpiredAt = nil
	} else {
		*user.ExpiredAt = time.Unix(expiredAt, 0)
	}
}

func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	user.Phone = util.FormatPhoneNumber(user.Phone)

	// ---
	if user.ID > 0 {
		tx.Model(&User{}).
			Where("library_user_id = ?", user.ID)
	}
	// ---
	return nil
}

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

func (User) FindUserByCode(code string) (existed bool, isExpired bool, user User) {
	// Response: EXISTED_USER, IS_EXPIRED / DISABLED, USER INFO
	user = User{}
	DBInstance.
		Where("code = ? ", code).
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

func (User) FindUserByEmail(email string) (existed bool, isExpired bool, user User) {
	// Response: EXISTED_USER, IS_EXPIRED / DISABLED, USER INFO
	user = User{}
	DBInstance.
		Where("email = ? ", email).
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

func (User) FindUserByPhone(phone string) (existed bool, isExpired bool, user User) {
	// Response: EXISTED_USER, IS_EXPIRED / DISABLED, USER INFO
	user = User{}
	DBInstance.
		Where("phone = ? ", phone).
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

func (User) FindUserByIdentityCard(identityCard string) (existed bool, isExpired bool, user User) {
	// Response: EXISTED_USER, IS_EXPIRED / DISABLED, USER INFO
	user = User{}
	DBInstance.
		Where("identity_card = ? ", identityCard).
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
