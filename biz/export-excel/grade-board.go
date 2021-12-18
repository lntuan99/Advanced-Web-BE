package export_excel

import (
	"advanced-web.hcmus/config"
	"advanced-web.hcmus/model"
	"advanced-web.hcmus/util"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"time"
)

const (
	GRADE_BOARD_SHEET_NAME                  = "grade_board"
	GRADE_BOARD_FOLDER                      = "grade_board"
	GRADE_BOARD_NORMAL_EXPORT_TEMPLATE_NAME = "export-grade-board-template.xlsx"
	GRADE_BOARD_PERMISSION_FILE             = 0755
)

type GradeBoardExcel ExcelStruct

func (gradeBoardExcel *GradeBoardExcel) GetWorkingDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return workingDir
}

func (gradeBoardExcel *GradeBoardExcel) CreateFolder(folderName string) {
	workingDir := gradeBoardExcel.GetWorkingDir()
	folderQor := fmt.Sprintf("%s%s/system/%s", workingDir, config.Config.MediaDir, folderName)
	if _, err := os.Stat(folderQor); err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(folderQor, GRADE_BOARD_PERMISSION_FILE)
		if err != nil {
			panic(err)
		}
	}
}

func (gradeBoardExcel *GradeBoardExcel) Save() error {
	if err := gradeBoardExcel.cursor.Save(); err != nil {
		return err
	}

	return nil
}

func (gradeBoardExcel *GradeBoardExcel) WriteLine(data interface{}, rowIndex int) *GradeBoardExcel {
	// -----------------------------
	// -----------------------------
	// APPEND DATA
	// -----------------------------
	// -----------------------------
	studentGradeData := data.(model.ResponseStudentGradeInClassroom)

	name := studentGradeData.StudentName
	if util.EmptyOrBlankString(name) {
		name = studentGradeData.Name
	}

	code := studentGradeData.StudentCode
	if util.EmptyOrBlankString(code) {
		code = studentGradeData.Code
	}

	values := make([]interface{}, 0)
	values = append(values, code)
	values = append(values, name)

	for _, grade := range studentGradeData.GradeArray {
		values = append(values, grade.Point)
	}

	// -----------------------------
	// -----------------------------
	// FORMAT COLUMN WIDTH
	// -----------------------------
	// -----------------------------
	if width, _ := gradeBoardExcel.cursor.GetColWidth(GRADE_BOARD_SHEET_NAME, "A"); int(width) < len(code) {
		_ = gradeBoardExcel.cursor.SetColWidth(GRADE_BOARD_SHEET_NAME, "A", "A", float64(len(code)))
	}

	if width, _ := gradeBoardExcel.cursor.GetColWidth(GRADE_BOARD_SHEET_NAME, "B"); int(width) < len(name) {
		_ = gradeBoardExcel.cursor.SetColWidth(GRADE_BOARD_SHEET_NAME, "B", "B", float64(len(name)))
	}

	//------------------------------
	// -----------------------------
	// WRITE DATA
	// -----------------------------
	// -----------------------------
	_ = gradeBoardExcel.cursor.SetSheetRow(GRADE_BOARD_SHEET_NAME, fmt.Sprintf("A%v", rowIndex), &values)

	return gradeBoardExcel
}

func NewGradeBoardExcelFile() *GradeBoardExcel {
	gradeBoardExcel := &GradeBoardExcel{}
	gradeBoardExcel.CreateFolder(GRADE_BOARD_FOLDER)

	workingDir, _ := os.Getwd()
	gradeBoardExcelFileTemplate := fmt.Sprintf("%s%s/assets/export-template/%v", workingDir, config.Config.MediaDir, GRADE_BOARD_NORMAL_EXPORT_TEMPLATE_NAME)

	nameSheet := fmt.Sprintf("DANH_SACH_BANG_DIEM_%v.xlsx", time.Now().Unix())
	mediaPath := workingDir + config.Config.MediaDir + "/system/" + GRADE_BOARD_FOLDER
	filePath := mediaPath + "/" + nameSheet

	_ = util.CopyFile(gradeBoardExcelFileTemplate, filePath)

	excel := util.ReadXLXS(filePath)

	gradeBoardExcel.cursor = excel
	gradeBoardExcel.FileName = nameSheet

	return gradeBoardExcel
}

func (gradeBoardExcel *GradeBoardExcel) WriteHeader(okeGradeArray []model.Grade) {
	// Declare Header style
	headerStyle, _ := gradeBoardExcel.cursor.NewStyle(
		&excelize.Style{
			Border: []excelize.Border{
				{
					Type:  "left",
					Color: "#000000",
					Style: 1,
				},
				{
					Type:  "top",
					Color: "#000000",
					Style: 1,
				},
				{
					Type:  "right",
					Color: "#000000",
					Style: 1,
				},
				{
					Type:  "bottom",
					Color: "#000000",
					Style: 1,
				},
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Pattern: 1,
				Color:   []string{"#DBDBDB"},
			},
			Font: &excelize.Font{
				Bold:   true,
				Italic: false,
				Family: "Calibri",
				Size:   11,
				Color:  "#000000",
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
				WrapText:   true,
			},
		},
	)

	startGradeColumnIndex := 3 // Mean column C in excel
	for i, grade := range okeGradeArray {
		columnName := ConvertNumberToColumnName(startGradeColumnIndex + i)
		gradeAxis := fmt.Sprintf("%v1", columnName)
		_ = gradeBoardExcel.cursor.SetCellStyle(GRADE_BOARD_SHEET_NAME, gradeAxis, gradeAxis, headerStyle)
		_ = gradeBoardExcel.cursor.SetCellStr(GRADE_BOARD_SHEET_NAME, gradeAxis, grade.Name)
	}
}

func ProcessExportGradeBoard(responseStudentGradeInClassroomArray []model.ResponseStudentGradeInClassroom, okeGradeArray []model.Grade) string {
	// Initialize the necessary excel
	excelFile := NewGradeBoardExcelFile()

	// Write Header
	excelFile.WriteHeader(okeGradeArray)

	// Write each data line to the table.
	// Start writing data at row 2th.
	for index, data := range responseStudentGradeInClassroomArray {
		rowIndex := index + 2
		excelFile.WriteLine(data, rowIndex)
	}

	// Save file
	err := excelFile.Save()
	util.CheckErr(err)

	return fmt.Sprintf("%v/export-data/%v/%v", config.Config.ApiDomain, GRADE_BOARD_FOLDER, excelFile.FileName)
}

//===================================================
//===================================================
//PRIVATE FUNCTION
//===================================================
//===================================================
