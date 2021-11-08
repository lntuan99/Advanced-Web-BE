package model

import (
	"advanced-web.hcmus/util"
	"github.com/jinzhu/gorm"
)

type Classroom struct {
	gorm.Model
	Name          string `gorm:"index:classroom_name_idx"`
	CoverImageURL string
	Code          string `gorm:"index:classroom_code_idx"`
	Description   string
	InviteLink    string `gorm:"index:classroom_invite_link_idx"`
	Users         []User `gorm:"many2many:user_classroom_mappings"`
}

type ClassroomRes struct {
	ID            uint
	Name          string `json:"name"`
	CoverImageURL string `json:"coverImageUrl"`
	Code          string `json:"code"`
	InviteLink    string `json:"inviteLink"`
	Description   string `json:"description"`
}

func (classroom Classroom) ToRes() ClassroomRes {
	return ClassroomRes{
		ID:            classroom.ID,
		Name:          classroom.Name,
		CoverImageURL: util.SubUrlToFullUrl(classroom.CoverImageURL),
		Code:          classroom.Code,
		InviteLink:    classroom.InviteLink,
		Description:   classroom.Description,
	}
}

//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
func (classroom Classroom) InitializeTableConfig() {
	// "gin" means: The column must be of tsvector type
	DBInstance.Exec(`CREATE INDEX IF NOT EXISTS search_field
    ON classrooms USING
    gin(search_field)`)

	DBInstance.Exec(`CREATE INDEX IF NOT EXISTS classroom_name_idx 
    ON classrooms
    USING gin (f_unaccent(name) gin_trgm_ops)`)
}

func (classroom Classroom) FindClassroomByCode(code string) Classroom {
	var res Classroom
	DBInstance.First(&res, "code = ?", code)

	return res
}
