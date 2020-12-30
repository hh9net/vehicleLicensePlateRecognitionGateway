package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
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
	service.OSSCount = 0
	service.ResultCount = 0
	service.AgainCount = 0
	service.ResultOKCount = 0
	//进程管理
	service.ProcessManagementService()
	//	goroutine1
	//	开线程读取xml文件 上传图片到oss  上传抓拍结果到车牌识别云端服务器
	go service.UploadFile()
	//goroutine2
	go service.HandleDayTasks()
	//goroutine3 抓拍结果再次上传
	go service.HandleFileAgainUpload()

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()
	//
	tiker := time.NewTicker(time.Minute * 3) //每15秒执行一下
	for {
		<-tiker.C
		log.Println("主go程执行 抓拍进程管理程序 OK呢！")
	}
}
