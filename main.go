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
	service.Gettoken = conf.Gettoken             //http://172.31.49.252/processor-control/collect/token/
	service.GetCameraListip = conf.GetCameraList //http://172.31.49.252/processor-control/collect/cameras/

	service.Deviceid = conf.Deviceid //fe0442b5-2d40-486f-9682-d1043ceca4e5

}

func main() {

	//初始化配置文件
	ConfigInit()

	//进程管理
	service.ProcessManagementService()

	//开线程读取xml文件 上传图片到oss  上传抓拍结果到车牌识别云端服务器
	service.UploadFile()

	tiker := time.NewTicker(time.Second * 10) //每15秒执行一下
	for {
		<-tiker.C
		log.Println(" 测试进程管理 ")
	}

}
