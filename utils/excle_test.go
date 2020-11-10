package utils

import "testing"

// ExportExcel 导出Excel文件
// sheetName 工作表名称
// columns 列名切片
// rows 数据切片，是一个二维数组
//func ExportExcel(sheetName string, columns []string, rows [][]interface{})
func TestExportExcel(t *testing.T) {
	columns := []string{"name", "age"}
	s1 := []interface{}{"小红", "18"}
	s2 := []interface{}{"小信", "27"}
	rows := [][]interface{}{s1, s2}
	ExportExcel("estExportExcel", columns, rows)
}
