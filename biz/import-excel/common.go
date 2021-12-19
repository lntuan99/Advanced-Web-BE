package import_excel

import (
	"github.com/xuri/excelize/v2"
	"math"
)

const (
	HOUR_IN_SECOND_12 = 12 * 60 * 60
	PAGE_SIZE         = 10
)

type ExcelStruct struct {
	SheetName string
	Cursor    *excelize.File
}

type ImportRowResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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

var ColumnMapping = map[int]string{
	0:  "A",
	1:  "B",
	2:  "C",
	3:  "D",
	4:  "E",
	5:  "F",
	6:  "G",
	7:  "H",
	8:  "I",
	9:  "J",
	10: "K",
	11: "L",
	12: "M",
	13: "N",
	14: "O",
	15: "P",
}
