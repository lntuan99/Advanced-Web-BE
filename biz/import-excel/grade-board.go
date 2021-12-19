package import_excel

import (
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"fmt"
	"strconv"
)

const (
	GRADE_BOARD_SHEET_NAME = "grade_board"
)

type SheetGradeBoardStruct struct {
	// Base import excel struct
	ExcelStruct

	ClassroomID uint

	// The final student array which will be inserted to database
	NewStudentGradeMappingArray []model.StudentGradeMapping

	// Key: Column name | Value: Grade model
	ImportGradeMap map[string]model.Grade

	// Key: code | Value: Student model
	ExistedStudentCodeMap map[string]model.Student

	// Key: studentCode_gradeName | Value: StudentGradeMapping model
	ExistedStudentGradeMappingMap map[string]model.StudentGradeMapping
}

func (SheetGradeBoardStruct) Initialize(fullPath string, classroomID uint) *SheetGradeBoardStruct {
	sheetGradeBoardStruct := &SheetGradeBoardStruct{
		ExcelStruct: ExcelStruct{
			SheetName: GRADE_BOARD_SHEET_NAME,
			Cursor:    util.ReadXLXS(fullPath),
		},
		ClassroomID:                   classroomID,
		ExistedStudentCodeMap:         SheetGradeBoardStruct{}.findExistedStudentCodeMap(classroomID),
		ExistedStudentGradeMappingMap: SheetGradeBoardStruct{}.findExistedStudentGradeMappingMap(classroomID),
	}

	sheetGradeBoardStruct.findImportGradeMap(classroomID)

	return sheetGradeBoardStruct
}

func (sheetGradeBoardStruct *SheetGradeBoardStruct) Importing() (ok bool, importRowResponseArray []ImportRowResponse) {
	importRowResponseArray = make([]ImportRowResponse, 0)

	// Check the existence of sheet
	sheetIndex := sheetGradeBoardStruct.Cursor.GetSheetIndex(sheetGradeBoardStruct.SheetName)
	if sheetIndex < 0 {
		return false, importRowResponseArray
	}

	// Get all data
	importRowResponseArray = sheetGradeBoardStruct.fetchAllData()

	// Insert data to database
	importRowResponseArray = sheetGradeBoardStruct.insertData(importRowResponseArray)

	return true, importRowResponseArray
}

// ==============================================================
// ==============================================================
// ==============================================================
// PRIVATE FUNCTIONS
// ==============================================================
// ==============================================================
// ==============================================================
func (sheetGradeBoardStruct *SheetGradeBoardStruct) findImportGradeMap(classroomID uint) {
	sheetGradeBoardStruct.ImportGradeMap = make(map[string]model.Grade)

	colArray, _ := sheetGradeBoardStruct.Cursor.GetCols(GRADE_BOARD_SHEET_NAME)

	// Skip column index 0 (Code) and 1 (Name)
	for i, col := range colArray {
		if i == 0 || i == 1 || len(col) == 0 || util.EmptyOrBlankString(col[0]) {
			continue
		}

		// Find grade by name in this classroom with @classroomID
		gradeName := util.StandardizedString(col[0])
		var dbGrade model.Grade
		model.DBInstance.First(&dbGrade, "classroom_id = ? AND name = ?", classroomID, gradeName)

		if dbGrade.ID > 0 {
			colName := ConvertNumberToColumnName(i + 1)
			sheetGradeBoardStruct.ImportGradeMap[colName] = dbGrade
		}
	}
}

func (SheetGradeBoardStruct) findExistedStudentCodeMap(classroomID uint) map[string]model.Student {
	studentArray := make([]model.Student, 0)
	model.DBInstance.
		Preload("User").
		Where("classroom_id = ?", classroomID).
		Find(&studentArray)

	studentCodeMap := make(map[string]model.Student)
	for _, student := range studentArray {
		if len(student.Code) > 0 {
			studentCodeMap[student.Code] = student
		}
	}

	return studentCodeMap
}

func (SheetGradeBoardStruct) findExistedStudentGradeMappingMap(classroomID uint) map[string]model.StudentGradeMapping {
	studentGradeMappingArray := make([]model.StudentGradeMapping, 0)
	model.DBInstance.
		Preload("Student").
		Preload("Grade").
		Joins("INNER JOIN grades ON grades.id = student_grade_mappings.grade_id "+
			"AND grades.classroom_id = ?", classroomID).
		Find(&studentGradeMappingArray)

	studentGradeMappingMap := make(map[string]model.StudentGradeMapping)
	for _, mapping := range studentGradeMappingArray {
		key := SheetGradeBoardStruct{}.generateStudentGradeMappingKey(mapping.Student, mapping.Grade)
		studentGradeMappingMap[key] = mapping
	}

	return studentGradeMappingMap
}

func (SheetGradeBoardStruct) generateStudentGradeMappingKey(student model.Student, grade model.Grade) string {
	return fmt.Sprintf("%v_%v", student.Code, grade.Name)
}

func (sheetGradeBoardStruct *SheetGradeBoardStruct) fetchAllData() (importRowResponseArray []ImportRowResponse) {
	importRowResponseArray = make([]ImportRowResponse, 0)

	rowArray, _ := sheetGradeBoardStruct.Cursor.GetRows(sheetGradeBoardStruct.SheetName)
	startRowIndex := 1

	for i := startRowIndex; i < len(rowArray); i++ {
		rowData := rowArray[i]
		if rowData == nil || len(rowData) == 0 {
			continue
		}

		// Kiểm tra phần tử đầu tiên của mỗi dòng có là ô trống không.
		// Nếu là ô trống thì không tiếp tục xử lý.
		if len(rowData[0]) == 0 {
			continue
		}

		var importRowResponse ImportRowResponse

		ok := true
		var studentGradeMappingInfo model.StudentGradeMapping
		studentGradeMappingInfo.Student.ClassroomID = 14

		// ---
		responseMessage := ""
		for columnIndex, columnData := range rowData {
			prefixErrorMsg := fmt.Sprintf("-- Failed: [Row %v][Col %v]", i+1, ColumnMapping[columnIndex])

			excelColumnIndex := columnIndex + 1
			columnName := ConvertNumberToColumnName(excelColumnIndex)

			if columnName == "A" {
				columnData = util.StandardizedString(columnData)

				if len(columnData) > 0 {
					studentGradeMappingInfo.Student.Code = columnData
				} else {
					ok = false
					responseMessage += fmt.Sprintf("%v Student code mustn't be empty. ", prefixErrorMsg)
				}

				continue
			}

			if columnName == "B" {
				columnData = util.StandardizedString(columnData)

				if len(columnData) > 0 {
					studentGradeMappingInfo.Student.Name = columnData
				} else {
					ok = false
					responseMessage += fmt.Sprintf("%v Student name mustn't be empty. ", prefixErrorMsg)
				}

				continue
			}

			grade, ok1 := sheetGradeBoardStruct.ImportGradeMap[columnName]
			if !ok1 {
				ok = false
				continue
			}

			point, _ := strconv.ParseFloat(columnData, 32)

			studentGradeMappingInfo.Grade = grade
			studentGradeMappingInfo.GradeID = grade.ID
			studentGradeMappingInfo.Point = float32(point)

			sheetGradeBoardStruct.NewStudentGradeMappingArray = append(sheetGradeBoardStruct.NewStudentGradeMappingArray, studentGradeMappingInfo)
		}

		if !ok {
			responseMessage += fmt.Sprintf("-- Failed: Skip inserting the new Student at row index %v. ", i+1)

			importRowResponse.Code = "fail"
			importRowResponse.Message = responseMessage
			importRowResponseArray = append(importRowResponseArray, importRowResponse)

			continue
		}
	}

	return importRowResponseArray
}

func (sheetGradeBoardStruct *SheetGradeBoardStruct) insertData(importRowResponseArray []ImportRowResponse) []ImportRowResponse {
	for _, newStudentGradeMapping := range sheetGradeBoardStruct.NewStudentGradeMappingArray {
		responseMessage := ""
		var importRowResponse ImportRowResponse

		code := newStudentGradeMapping.Student.Code
		dbCodeStudent, ok00 := sheetGradeBoardStruct.ExistedStudentCodeMap[code]

		// If not existed student code before
		if !ok00 {
			if err := model.DBInstance.Create(&newStudentGradeMapping.Student).Error; err == nil {
				sheetGradeBoardStruct.ExistedStudentCodeMap[newStudentGradeMapping.Student.Code] = newStudentGradeMapping.Student

				newStudentGradeMapping.StudentID = newStudentGradeMapping.Student.ID
			}
		} else {
			newStudentGradeMapping.StudentID = dbCodeStudent.ID
			newStudentGradeMapping.Student.ID = dbCodeStudent.ID
			if err := model.DBInstance.Save(&newStudentGradeMapping.Student).Error; err == nil {
				//
			}
		}

		key := sheetGradeBoardStruct.generateStudentGradeMappingKey(newStudentGradeMapping.Student, newStudentGradeMapping.Grade)
		dbStudentGradeMapping, ok01 := sheetGradeBoardStruct.ExistedStudentGradeMappingMap[key]
		// If not existed student code before
		if !ok01 {
			if err := model.DBInstance.Create(&newStudentGradeMapping).Error; err == nil {
				responseMessage = fmt.Sprintf("-- Success: New grade mapping has been inserted (code = %v; grade name = %v). ", code, newStudentGradeMapping.Grade.Name)
				sheetGradeBoardStruct.ExistedStudentGradeMappingMap[key] = newStudentGradeMapping
			}
		} else {
			newStudentGradeMapping.ID = dbStudentGradeMapping.ID
			if err := model.DBInstance.Save(&newStudentGradeMapping).Error; err == nil {
				responseMessage = fmt.Sprintf("-- Success: Info of existed grade has been updated (code = %v; grade name = %v). ", code, newStudentGradeMapping.Grade.Name)
			}
		}

		importRowResponse.Code = "success"
		importRowResponse.Message = responseMessage
		importRowResponseArray = append(importRowResponseArray, importRowResponse)
	}

	return importRowResponseArray
}
