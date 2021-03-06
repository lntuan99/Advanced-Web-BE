package export_excel

import (
	"advanced-web.hcmus/config"
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"fmt"
	"os"
	"time"
)

const (
	STUDENT_SHEET_NAME                  = "student"
	STUDENT_FOLDER                      = "student"
	STUDENT_NORMAL_EXPORT_TEMPLATE_NAME = "export-student-template.xlsx"
	STUDENT_PERMISSION_FILE             = 0755
)

type StudentExcel ExcelStruct

func (excel *StudentExcel) GetWorkingDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return workingDir
}

func (excel *StudentExcel) CreateFolder(folderName string) {
	workingDir := excel.GetWorkingDir()
	folderQor := fmt.Sprintf("%s%s/system/%s", workingDir, config.Config.MediaDir, folderName)
	if _, err := os.Stat(folderQor); err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(folderQor, STUDENT_PERMISSION_FILE)
		if err != nil {
			panic(err)
		}
	}
}

func (excel *StudentExcel) Save() error {
	if err := excel.cursor.Save(); err != nil {
		return err
	}

	return nil
}

func (excel *StudentExcel) WriteLine(data interface{}, rowIndex int) *StudentExcel {
	// -----------------------------
	// -----------------------------
	// APPEND DATA
	// -----------------------------
	// -----------------------------
	studentData := data.(model.Student)

	name := studentData.Name
	if util.EmptyOrBlankString(name) {
		name = studentData.User.Name
	}

	code := studentData.Code
	if util.EmptyOrBlankString(code) {
		code = studentData.User.Code
	}

	values := make([]interface{}, 0)
	values = append(values, code)
	values = append(values, name)

	// -----------------------------
	// -----------------------------
	// FORMAT COLUMN WIDTH
	// -----------------------------
	// -----------------------------
	if width, _ := excel.cursor.GetColWidth(STUDENT_SHEET_NAME, "A"); int(width) < len(code) {
		_ = excel.cursor.SetColWidth(STUDENT_SHEET_NAME, "A", "A", float64(len(code)))
	}

	if width, _ := excel.cursor.GetColWidth(STUDENT_SHEET_NAME, "B"); int(width) < len(name) {
		_ = excel.cursor.SetColWidth(STUDENT_SHEET_NAME, "B", "B", float64(len(name)))
	}

	//------------------------------
	// -----------------------------
	// WRITE DATA
	// -----------------------------
	// -----------------------------
	_ = excel.cursor.SetSheetRow(STUDENT_SHEET_NAME, fmt.Sprintf("A%v", rowIndex), &values)

	return excel
}

func NewStudentExcelFile() *StudentExcel {
	studentExcel := &StudentExcel{}
	studentExcel.CreateFolder(STUDENT_FOLDER)

	workingDir, _ := os.Getwd()
	studentExcelFileTemplate := fmt.Sprintf("%s%s/assets/export-template/%v", workingDir, config.Config.MediaDir, STUDENT_NORMAL_EXPORT_TEMPLATE_NAME)

	nameSheet := fmt.Sprintf("DANH_SACH_HOC_SINH_%v.xlsx", time.Now().Unix())
	mediaPath := workingDir + config.Config.MediaDir + "/system/" + STUDENT_FOLDER
	filePath := mediaPath + "/" + nameSheet

	_ = util.CopyFile(studentExcelFileTemplate, filePath)

	excel := util.ReadXLXS(filePath)

	studentExcel.cursor = excel
	studentExcel.FileName = nameSheet

	return studentExcel
}

func ProcessExportStudent(studentArray []model.Student) string {
	// Initialize the necessary excel
	excelFile := NewStudentExcelFile()

	// Write each data line to the table.
	// Start writing data at row 2th.
	for index, data := range studentArray {
		rowIndex := index + 2
		excelFile.WriteLine(data, rowIndex)
	}

	// Save file
	err := excelFile.Save()
	util.CheckErr(err)

	return fmt.Sprintf("%v/export-data/%v/%v", config.Config.ApiDomain, STUDENT_FOLDER, excelFile.FileName)
}

//===================================================
//===================================================
//PRIVATE FUNCTION
//===================================================
//===================================================
