package main

import (
	log "github.com/sirupsen/logrus"
	"time"
	"vehicleLicensePlateRecognitionGateway/config"
	"vehicleLicensePlateRecognitionGateway/router"
	"vehicleLicensePlateRecognitionGateway/service"
	"vehicleLicensePlateRecognitionGateway/utils"
)

func main() {

	conf := config.ConfigInit() //初始化配置文件
	log.Println("配置文件信息：", *conf)
	utils.InitLogrus(conf.LogPath, conf.LogFileName, time.Duration(24*conf.LogMaxAge)*time.Hour, conf.LogRotationSize, time.Duration(conf.LogRotationTime)*time.Hour, conf.RotationCount)
	IpAddress := conf.IpAddress
	service.GwCaptureInformationUploadIpAddress = conf.GwCaptureInformationUploadIpAddress

	service.UserName = conf.UserName
	service.Password = conf.Password
	service.Deviceid = conf.Deviceid

	router.RouteInit(IpAddress)

	tiker := time.NewTicker(time.Microsecond * 1) //每15秒执行一下
	for {
		<-tiker.C
		log.Println("配置文件信息")
	}

}