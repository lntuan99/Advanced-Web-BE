package export_excel

import (
	"github.com/xuri/excelize/v2"
	"math"
)

const (
	HOUR_IN_SECOND_12 = 12 * 60 * 60
	PAGE_SIZE         = 10
)

type ExcelStruct struct {
	cursor     *excelize.File
	NumberLine int
	FileName   string
	Header     []interface{}
}

type IExcel interface {
	GetWorkingDir() string
	CreateFolder(folderName string) string
	WriteHeader() *IExcel
	Save(fileName string) (string, error)
	WriteLine(data interface{}, rowIndex int) *IExcel
}

func ConvertNumberToColumnName(number int) string {
	// initialize output string as empty
	result := ""

	for number > 0 {
		// find the index of the next letter and concatenate the letter
		// to the solution

		// here index 0 corresponds to `A`, and 25 corresponds to `Z`
		index := (number - 1) % 26
		result = string(rune(index+'A')) + result
		number = (number - 1) / 26
	}

	return result
}

func GetTotalPage(totalSize int, pageSize int) int {
	return int(math.Ceil(float64(totalSize) / float64(pageSize)))
}
