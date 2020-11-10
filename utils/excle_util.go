package utils

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"time"
)

// maxCharCount 最多26个字符A-Z
const maxCharCount = 26

// ExportExcel 导出Excel文件
// sheetName 工作表名称
// columns 列名切片
// rows 数据切片，是一个二维数组
func ExportExcel(sheetName string, columns []string, rows [][]interface{}) string {
	f := excelize.NewFile()
	sheetIndex := f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")
	maxColumnRowNameLen := 1 + len(strconv.Itoa(len(rows)))
	columnCount := len(columns)
	if columnCount > maxCharCount {
		maxColumnRowNameLen++
	} else if columnCount > maxCharCount*maxCharCount {
		maxColumnRowNameLen += 2
	}
	columnNames := make([][]byte, 0, columnCount)
	for i, column := range columns {
		columnName := getColumnName(i, maxColumnRowNameLen)
		columnNames = append(columnNames, columnName)
		// 首行
		f.SetCellValue(sheetName, getColumnRowName(columnName, 1), column)
	}
	for rowIndex, row := range rows {
		for columnIndex, columnName := range columnNames {
			// 从第二行开始
			f.SetCellValue(sheetName, getColumnRowName(columnName, rowIndex+2), row[columnIndex])
		}
	}
	f.SetActiveSheet(sheetIndex)
	// Save xlsx file by the given path.
	err := f.SaveAs("./" + time.Now().Format("2006-01-02 15:04:05") + ".xlsx")
	if err != nil {
		fmt.Println(err)
	}

	return time.Now().Format("2006-01-02 15:04:05") + ".xlsx"
}

// getColumnName 生成列名
// Excel的列名规则是从A-Z往后排;超过Z以后用两个字母表示，比如AA,AB,AC;两个字母不够以后用三个字母表示，比如AAA,AAB,AAC
// 这里做数字到列名的映射：0 -> A, 1 -> B, 2 -> C
// maxColumnRowNameLen 表示名称框的最大长度，假设数据是10行，1000列，则最后一个名称框是J1000(如果有表头，则是J1001),是4位
// 这里根据 maxColumnRowNameLen 生成切片，后面生成名称框的时候可以复用这个切面，而无需扩容
func getColumnName(column, maxColumnRowNameLen int) []byte {
	const A = 'A'
	if column < maxCharCount {
		// 第一次就分配好切片的容量
		slice := make([]byte, 0, maxColumnRowNameLen)
		return append(slice, byte(A+column))
	} else {
		// 递归生成类似AA,AB,AAA,AAB这种形式的列名
		return append(getColumnName(column/maxCharCount-1, maxColumnRowNameLen), byte(A+column%maxCharCount))
	}
}

// getColumnRowName 生成名称框
// Excel的名称框是用A1,A2,B1,B2来表示的，这里需要传入前一步生成的列名切片，然后直接加上行索引来生成名称框，就无需每次分配内存
func getColumnRowName(columnName []byte, rowIndex int) (columnRowName string) {
	l := len(columnName)
	columnName = strconv.AppendInt(columnName, int64(rowIndex), 10)
	columnRowName = string(columnName)
	columnName = columnName[:l]
	return
}
