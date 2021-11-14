package model

import (
	"advanced-web.hcmus/util"
	"github.com/jinzhu/gorm"
)

type Classroom struct {
	gorm.Model
	Name              string `gorm:"index:classroom_name_idx"`
	CoverImageURL     string
	Code              string `gorm:"index:classroom_code_idx"`
	Description       string
	InviteTeacherLink string `gorm:"index:classroom_invite_teacher_link_idx"`
	InviteStudentLink string `gorm:"index:classroom_invite_student_link_idx"`
	StudentArray      []User `gorm:"many2many:user_classroom_mappings"`
	TeacherArray      []User `gorm:"many2many:user_classroom_mappings"`
}

type ClassroomRes struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	CoverImageURL     string    `json:"coverImageUrl"`
	Code              string    `json:"code"`
	InviteTeacherLink string    `json:"inviteTeacherLink"`
	InviteStudentLink string    `json:"inviteStudentLink"`
	Description       string    `json:"description"`
	StudentResArray   []UserRes `json:"studentArray"`
	TeacherResArray   []UserRes `json:"teacherArray"`
}

func (classroom Classroom) ToRes() ClassroomRes {
	var studentResArray = make([]UserRes, 0)
	for _, student := range classroom.StudentArray {
		studentResArray = append(studentResArray, student.ToRes())
	}

	var teacherResArray = make([]UserRes, 0)
	for _, teacher := range classroom.TeacherArray {
		teacherResArray = append(teacherResArray, teacher.ToRes())
	}

	return ClassroomRes{
		ID:                classroom.ID,
		Name:              classroom.Name,
		CoverImageURL:     util.SubUrlToFullUrl(classroom.CoverImageURL),
		Code:              classroom.Code,
		InviteTeacherLink: classroom.InviteTeacherLink,
		InviteStudentLink: classroom.InviteStudentLink,
		Description:       classroom.Description,
		StudentResArray:   studentResArray,
		TeacherResArray:   teacherResArray,
	}
}

//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
func (Classroom) InitializeTableConfig() {
	// "gin" means: The column must be of tsvector type
	DBInstance.Exec(`CREATE INDEX IF NOT EXISTS search_field
    ON classrooms USING
    gin(search_field)`)

	DBInstance.Exec(`CREATE INDEX IF NOT EXISTS classroom_name_idx 
    ON classrooms
    USING gin (f_unaccent(name) gin_trgm_ops)`)
}

func (Classroom) FindClassroomByCode(code string) Classroom {
	var res Classroom
	DBInstance.First(&res, "code = ?", code)

	return res
}

func (Classroom) FindClassroomByID(id uint) Classroom {
	var res Classroom
	DBInstance.First(&res, id)

	return res
}

func (classroom Classroom) GetListUserByJWTType(userRole UserRole) []User {
	var userArray = make([]User, 0)

	var mappingArray = make([]UserClassroomMapping, 0)
	DBInstance.
		Preload("User").
		Preload("UserRole").
		Where("classroom_id = ?", classroom.ID).
		Where("user_role_id = ?", userRole.ID).
		Find(&mappingArray)

	for _, mapping := range mappingArray {
		userArray = append(userArray, mapping.User)
	}

	return userArray
}
