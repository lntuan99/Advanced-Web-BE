package model

import "github.com/jinzhu/gorm"

type Student struct {
	gorm.Model
	ClassroomID uint
	Classroom   Classroom
	Code        string `gorm:"index"` // Mean code in User
	User        User   `gorm:"foreignkey:Code;association_foreignkey:Code"`
	Name        string
}

type StudentRes struct {
	// Add new field if needed
	UserRes
}

func (student Student) ToRes() StudentRes {
	DBInstance.First(&student.User, student.Code)

	return StudentRes{
		student.User.ToRes(),
	}
}

func (student Student) MappedStudentInformationToResponseStudentGradeInClassroom(classroomID uint) ResponseStudentGradeInClassroom {
	totalGrade, maxTotalGrade, gradeArray := student.GetAllGradeInClassroom(classroomID)

	return ResponseStudentGradeInClassroom{
		StudentRes:    student.ToRes(),
		StudentName:   student.Name,
		StudentCode:   student.Code,
		TotalGrade:    totalGrade,
		MaxTotalGrade: maxTotalGrade,
		GradeArray:    gradeArray,
	}
}

func (student Student) GetAllGradeInClassroom(classroomID uint) (totalGrade float32, maxTotalGrade float32, result []StudentGradeMappingRes) {
	result = make([]StudentGradeMappingRes, 0)

	// Find all grade in class
	var gradeArray = make([]Grade, 0)
	DBInstance.Find(&gradeArray, "classroom_id = ?", classroomID)

	// Check student is mapped with all grade
	var studentGradeMappingArray = make([]StudentGradeMapping, 0)
	for _, grade := range gradeArray {
		var dbStudentGradeMapping StudentGradeMapping
		DBInstance.First(&dbStudentGradeMapping, "student_id = ? AND grade_id = ?", student.ID, grade.ID)

		// if not existed => create new
		if dbStudentGradeMapping.ID == 0 {
			dbStudentGradeMapping.StudentID = Student{}.ID
			dbStudentGradeMapping.GradeID = grade.ID
			dbStudentGradeMapping.Point = 0
			DBInstance.Create(&dbStudentGradeMapping)
		}

		totalGrade += dbStudentGradeMapping.Point
		maxTotalGrade += grade.MaxPoint
		studentGradeMappingArray = append(studentGradeMappingArray, dbStudentGradeMapping)
	}

	for _, mapping := range studentGradeMappingArray {
		result = append(result, mapping.ToRes())
	}

	return totalGrade, maxTotalGrade, result
}
