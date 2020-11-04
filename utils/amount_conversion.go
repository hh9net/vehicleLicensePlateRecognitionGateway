package utils

import (
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

//分转元
func Fen2Yuan(price int64) string {

	d := decimal.New(1, 2) //分除以100得到元
	result := decimal.NewFromFloat(float64(price)).DivRound(d, 2).StringFixed(2)
	//log.Printf("输入值为：%d,分转元后，精度为二的结果为：%s\n", price, result)
	return result
}

//元转分,乘以100后，保留整数部分
func Yuan2Fen(price float64) int64 {
	d := decimal.New(1, 2)  //分转元乘以100
	d1 := decimal.New(1, 0) //乘完之后，保留2为小数，需要这么一个中间参数
	df := decimal.NewFromFloat(price).Mul(d).IntPart()

	//如下是满足，当乘以100后，仍然有小数位，取四舍五入法后，再取整数部分
	dff := decimal.NewFromFloat(price).Mul(d).DivRound(d1, 0).IntPart()
	log.Printf("输入值为：%f, 简单的元转分后，取整数部分：%d\n", price, df)
	log.Printf("输入值为：%f, 元转分后，若还有小数，需做四舍五入后，再取整数：%d\n", price, dff)
	//return df
	return dff
}
