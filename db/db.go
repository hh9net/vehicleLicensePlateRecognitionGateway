package db

import (
	"vehicleLicensePlateRecognitionGateway/utils"
)

var GormClient *utils.GormDB

//数据库的初始化
func DBInit(mstr string) {
	GormClient = utils.InitGormDB(&utils.DBConfig{
		DBAddr:       mstr,
		MaxIdleConns: 30,
		LogMode:      utils.Uint8ToBool(1),
	})
}
