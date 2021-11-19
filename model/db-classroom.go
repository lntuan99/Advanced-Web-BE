package model

import (
	"advanced-web.hcmus/config/constants"
	"advanced-web.hcmus/util"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	STUDENT_INVITE_LINK_PREFIX_DOMAIN = "/student/join/"
	TEACHER_INVITE_LINK_PREFIX_DOMAIN = "/teacher/join/"
)

type Classroom struct {
	gorm.Model
	OwnerID           uint
	Owner             User
	Name              string `gorm:"index:classroom_name_idx"`
	CoverImageURL     string
	Code              string `gorm:"index:classroom_code_idx"`
	Description       string
	InviteTeacherCode string `gorm:"index:classroom_invite_teacher_link_idx"`
	InviteStudentCode string `gorm:"index:classroom_invite_student_link_idx"`
	StudentArray      []User `gorm:"many2many:user_classroom_mappings"`
	TeacherArray      []User `gorm:"many2many:user_classroom_mappings"`
}

type ClassroomRes struct {
	ID                uint      `json:"id"`
	OwnerName         string    `json:"ownerName"`
	OwnerAvatar       string    `json:"ownerAvatar"`
	Name              string    `json:"name"`
	CoverImageURL     string    `json:"coverImageUrl"`
	Code              string    `json:"code"`
	InviteTeacherCode string    `json:"inviteTeacherCode"`
	InviteStudentCode string    `json:"inviteStudentCode"`
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
		OwnerName:         classroom.Owner.Name,
		OwnerAvatar:       util.SubUrlToFullUrl(classroom.Owner.Avatar),
		Name:              classroom.Name,
		CoverImageURL:     util.SubUrlToFullUrl(classroom.CoverImageURL),
		Code:              classroom.Code,
		InviteTeacherCode: classroom.InviteTeacherCode,
		InviteStudentCode: classroom.InviteStudentCode,
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

func (Classroom) FindClassroomByCodeBelongToOwner(code string, ownerID uint) Classroom {
	var res Classroom
	DBInstance.First(&res, "code = ? and owner_id = ?", code, ownerID)

	return res
}

func (Classroom) FindClassroomByID(id uint) Classroom {
	var res Classroom
	DBInstance.
		Preload("Owner").
		First(&res, id)

	return res
}

func (classroom Classroom) GetListUserByJWTType(JWTType uint) []User {
	var userRole = UserRole{}.GetRoleByJWTType(JWTType)

	var userArray = make([]User, 0)

	if userRole.ID > 0 {
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
	}

	return userArray
}

func (classroom *Classroom) ClassroomGenerateInviteCode() {
	classroom.InviteTeacherCode = GenerateInviteCode(*classroom, JWT_TYPE_TEACHER)
	classroom.InviteStudentCode = GenerateInviteCode(*classroom, JWT_TYPE_STUDENT)
}

func GenerateInviteCode(classroom Classroom, jwtType uint) string {
	inviteCode := fmt.Sprintf("%v_%v_%v_%v", classroom.Code, classroom.OwnerID, jwtType, time.Now().Unix())
	inviteCode = util.HexSha256String([]byte(inviteCode))
	inviteCode += fmt.Sprintf("%v", time.Now().Unix()%constants.PRIME_NUMBER_FOR_MOD)

	return inviteCode
}

func (Classroom) GetClassroomByInviteCode(inviteCode string) (existed bool, classroom Classroom, jwtType uint) {
	classroom = Classroom{}

	DBInstance.First(&classroom, "invite_teacher_code = ?", inviteCode)
	if classroom.ID > 0 {
		return true, classroom, JWT_TYPE_TEACHER
	}

	DBInstance.First(&classroom, "invite_student_code = ?", inviteCode)
	if classroom.ID > 0 {
		return true, classroom, JWT_TYPE_STUDENT
	}

	return false, classroom, 0
}
