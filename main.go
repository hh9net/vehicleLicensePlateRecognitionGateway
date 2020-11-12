package main

import (
	log "github.com/sirupsen/logrus"
	"time"
	"vehicleLicensePlateRecognitionGateway/config"

	"vehicleLicensePlateRecognitionGateway/service"
	"vehicleLicensePlateRecognitionGateway/utils"
)

func ConfigInit() {
	conf := config.ConfigInit() //初始化配置文件
	log.Println("配置文件信息：", *conf)
	//初始化日志
	utils.InitLogrus(conf.LogPath, conf.LogFileName, time.Duration(24*conf.LogMaxAge)*time.Hour, conf.LogRotationSize, time.Duration(conf.LogRotationTime)*time.Hour, conf.RotationCount)
	//
	service.GwCaptureInformationUploadIpAddress = conf.GwCaptureInformationUploadIpAddress
	service.Gettoken = conf.Gettoken
	service.GetCameraListip = conf.GetCameraList

	service.Deviceid = conf.Deviceid

}

func main() {
	ConfigInit()
	if err := service.Runmain(); err != nil {
		err = service.Runmain()
		log.Println("重启")
	}
	tiker := time.NewTicker(time.Second * 10) //每15秒执行一下
	for {
		<-tiker.C
		log.Println(" 测试进程管理")
	}

}
