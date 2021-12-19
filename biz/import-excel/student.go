package import_excel

import (
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"fmt"
)

const (
	STUDENT_SHEET_NAME = "student"
)

type SheetStudentStruct struct {
	// Base import excel struct
	ExcelStruct

	ClassroomID uint

	// The final student array which will be inserted to database
	NewStudentArray []model.Student

	// Key: code | Value: Student model
	ExistedStudentCodeMap map[string]model.Student
}

func (SheetStudentStruct) Initialize(fullPath string, classroomID uint) *SheetStudentStruct {
	return &SheetStudentStruct{
		ExcelStruct: ExcelStruct{
			SheetName: STUDENT_SHEET_NAME,
			Cursor:    util.ReadXLXS(fullPath),
		},
		ClassroomID:           classroomID,
		ExistedStudentCodeMap: SheetStudentStruct{}.findExistedStudentCodeMap(classroomID),
	}
}

func (sheetStudent *SheetStudentStruct) Importing() (ok bool, importRowResponseArray []ImportRowResponse) {
	importRowResponseArray = make([]ImportRowResponse, 0)

	// Check the existence of sheet
	sheetIndex := sheetStudent.Cursor.GetSheetIndex(sheetStudent.SheetName)
	if sheetIndex < 0 {
		return false, importRowResponseArray
	}

	// Get all data
	importRowResponseArray = sheetStudent.fetchAllData()

	// Insert data to database
	importRowResponseArray = sheetStudent.insertData(importRowResponseArray)

	return true, importRowResponseArray
}

// ==============================================================
// ==============================================================
// ==============================================================
// PRIVATE FUNCTIONS
// ==============================================================
// ==============================================================
// ==============================================================
func (SheetStudentStruct) findExistedStudentCodeMap(classroomID uint) map[string]model.Student {
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

func (sheetStudent *SheetStudentStruct) fetchAllData() (importRowResponseArray []ImportRowResponse) {
	importRowResponseArray = make([]ImportRowResponse, 0)

	rowArray, _ := sheetStudent.Cursor.GetRows(sheetStudent.SheetName)
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
		var studentInfo model.Student
		studentInfo.ClassroomID = sheetStudent.ClassroomID

		// ---
		responseMessage := ""
		for columnIndex, columnData := range rowData {
			prefixErrorMsg := fmt.Sprintf("-- Failed: [Row %v][Col %v]", i+1, ColumnMapping[columnIndex])

			switch columnIndex {
			case 0: // Code
				columnData = util.StandardizedString(columnData)

				if len(columnData) > 0 {
					studentInfo.Code = columnData
				} else {
					ok = false
					responseMessage += fmt.Sprintf("%v Student code mustn't be empty. ", prefixErrorMsg)
				}
				break
			case 1: // Name
				if len(columnData) > 0 {
					studentInfo.Name = columnData
				} else {
					ok = false
					responseMessage += fmt.Sprintf("%v Student name mustn't be empty. ", prefixErrorMsg)
				}
				break
			}
		}

		if !ok {
			responseMessage += fmt.Sprintf("-- Failed: Skip inserting the new Student at row index %v. ", i+1)

			importRowResponse.Code = "fail"
			importRowResponse.Message = responseMessage
			importRowResponseArray = append(importRowResponseArray, importRowResponse)

			continue
		}

		sheetStudent.NewStudentArray = append(sheetStudent.NewStudentArray, studentInfo)
	}

	return importRowResponseArray
}

func (sheetStudent *SheetStudentStruct) insertData(importRowResponseArray []ImportRowResponse) []ImportRowResponse {
	for _, newStudent := range sheetStudent.NewStudentArray {
		responseMessage := ""
		var importRowResponse ImportRowResponse

		code := newStudent.Code
		dbCodeStudent, ok00 := sheetStudent.ExistedStudentCodeMap[code]

		// If not existed student code before
		if !ok00 {
			if err := model.DBInstance.Create(&newStudent).Error; err == nil {
				responseMessage = fmt.Sprintf("-- Success: New Student has been inserted (code = %v). ", code)
				sheetStudent.ExistedStudentCodeMap[newStudent.Code] = newStudent
			}
		} else {
			newStudent.ID = dbCodeStudent.ID
			if err := model.DBInstance.Save(&newStudent).Error; err == nil {
				responseMessage = fmt.Sprintf("-- Success: Info of existed Student has been updated (code = %v). ", code)
			}
		}

		importRowResponse.Code = "success"
		importRowResponse.Message = responseMessage
		importRowResponseArray = append(importRowResponseArray, importRowResponse)
	}

	return importRowResponseArray
}
