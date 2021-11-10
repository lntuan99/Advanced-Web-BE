package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

const (
	GENDER_MALE    = 1
	GENDER_FEMALE  = 2
	GENDER_UNKNOWN = 0
)

type User struct {
	gorm.Model
	UserRoleID   uint
	UserRole     UserRole
	Name         string `gorm:"index:user_name_idx"`
	Code         string `gorm:"index:user_code_unique_idx"`
	Email        string `gorm:"index:user_email_unique_idx"`
	Phone        string `gorm:"index:user_phone_unique_idx"`
	Birthday     time.Time
	Gender       uint
	Avatar       string
	IdentityCard string
	Classrooms []Classroom `gorm:"many2many:user_classroom_mappings"`
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