package db

import (
	log "github.com/sirupsen/logrus"
	"vehicleLicensePlateRecognitionGateway/config"
	"vehicleLicensePlateRecognitionGateway/utils"

	"time"
)

//结算监控平台数据层：数据的增删改查
func Newdb() {
	conf := config.ConfigInit() //初始化配置文件
	log.Println("配置文件信息：", *conf)
	utils.InitLogrus(conf.LogPath, conf.LogFileName, time.Duration(24*conf.LogMaxAge)*time.Hour, conf.LogRotationSize, time.Duration(conf.LogRotationTime)*time.Hour, conf.RotationCount)

}

//1、查询表是否存在
func QueryTable(tablename string) {
	db := utils.GormClient.Client
	is := db.HasTable(tablename)
	if is == false {
		log.Println("不存在", tablename)
		return
	}
	log.Println("表存在：", tablename, is)
}
