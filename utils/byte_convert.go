package utils

import (
	"github.com/shopspring/decimal"
	"log"
	"strconv"
)

//bytes To MB
func ByteToMB(bytes float64) float64 {

	//bytes  kb  mb
	b := (bytes / 1024) / 1024
	log.Println("bytes:", bytes, "bytes to MB:", b)
	//保留两位小数
	d := decimal.New(1, 0) // exp->    0:原样 2：除100
	result := decimal.NewFromFloat(b).DivRound(d, 2).StringFixed(2)
	log.Println("输入值为：", b, "精度为二的结果为：", result, "MB")
	value, _ := strconv.ParseFloat(result, 64)
	return value
}

//bytes To GB
func ByteToGB(bytes float64) float64 {
	//bytes       kb      mb     GB
	b := ((bytes / 1000) / 1000) / 1000
	//保留两位小数
	d := decimal.New(1, 0) // exp->    0:原样 2：除100
	result := decimal.NewFromFloat(b).DivRound(d, 2).StringFixed(2)
	log.Println("输入值为：", b, "精度为二的结果为：", result, "GB")

	//string转成float64
	value, _ := strconv.ParseFloat(result, 64)
	return value
}

//bytes To 字符串MB
func ByteToStrMB(bytes float64) string {

	//bytes  kb  mb
	b := (bytes / 1000) / 1000
	log.Println("bytes:", bytes, "bytes to MB:", b)
	//保留两位小数
	d := decimal.New(1, 0) // exp->    0:原样 2：除100
	result := decimal.NewFromFloat(b).DivRound(d, 2).StringFixed(2)
	log.Println("输入值为：", b, "精度为二的结果为：", result, "MB")
	//value, _ := strconv.ParseFloat(result, 64)
	return result
}

//bytes转字符串Gb
func ByteToStrGB(bytes float64) string {
	//bytes       kb      mb     GB
	b := ((bytes / 1000) / 1000) / 1000
	//保留两位小数
	d := decimal.New(1, 0) // exp->    0:原样 2：除100
	result := decimal.NewFromFloat(b).DivRound(d, 2).StringFixed(2)
	log.Println("输入值为：", b, "精度为二的结果为：", result, "GB")

	//float保留两位小数
	//value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", vfloat), 64)
	//fmt.Println(value, reflect.TypeOf(value))
	value, _ := strconv.ParseFloat(result, 64)
	log.Println(value, result)
	return result
}
