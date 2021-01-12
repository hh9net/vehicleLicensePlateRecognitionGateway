package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
	"vehicleLicensePlateRecognitionGateway/config"
	web "vehicleLicensePlateRecognitionGateway/gatewayWeb"
	"vehicleLicensePlateRecognitionGateway/service"
	"vehicleLicensePlateRecognitionGateway/utils"
)

func Init() {
	conf := config.ConfigInit() //初始化配置文件
	log.Println("配置文件信息：", *conf)
	//初始化日志
	utils.InitLogrus(conf.LogPath, conf.LogFileName, time.Duration(24*conf.LogMaxAge)*time.Hour, conf.LogRotationSize, time.Duration(conf.LogRotationTime)*time.Hour, conf.RotationCount)
	service.DelcameraConfigDir() //删除旧cameraConfigDir文件夹

	//
	go func() {
		Listenerr := http.ListenAndServe("0.0.0.0:3020", nil)
		if Listenerr != nil {
			//进程互斥
			log.Println("监控gc内存 Listen:", Listenerr)
			log.Println("监控gc内存 该端口已经启动，无法运行新进程！")
			os.Exit(0)
		}
	}()

	web.Gatewaylocation = conf.Gatewaylocation //1门架、2、服务区 3、收费站
	service.GwCaptureInformationUploadIpAddress = conf.GwCaptureInformationUploadIpAddress
	service.Gettoken = conf.Gettoken             //http://172.31.49.252/processor-control/collect/token/
	service.GetCameraListip = conf.GetCameraList //http://172.31.49.252/processor-control/collect/cameras/

	service.Deviceid = conf.Deviceid //fe0442b5-2d40-486f-9682-d1043ceca4e5
	service.StatisticalReportIpAddress = conf.StatisticalReportIpAddress
	service.MainStartTime = time.Now().Format("2006-01-02 15:04:05")
	//作为一个每次发布的一个版本记录
	service.MainVersion = "2021-01-12T17h30m00s_build"
	vs := "\n" + service.MainVersion + ""
	service.VersionFile(vs)

	service.OSSCount = 0
	service.NewOSSCount = 0
	service.ResultCount = 0
	service.AgainCount = 0
	service.ResultOKCount = 0
}

func main() {

	//初始化配置文件
	Init()
	//进程管理
	//service.ProcessManagementService()
	//goroutine1
	//开线程读取xml文件 上传图片到oss  上传抓拍结果到车牌识别云端服务器
	//go service.UploadFile()
	//goroutine2
	//	go service.HandleDayTasks()
	//goroutine3 抓拍结果再次上传
	//	go service.HandleFileAgainUpload()
	//goroutine4 定时20秒网关上报自身状态、摄像机状态状态至平台
	//go service.StatisticalReport()
	//goroutine5 网关每隔10分钟轮询请求服务器的版本
	//	go service.VersionQeq()
	//goroutine6 凌晨零点清零
	//	go service.HandleDayZeroTasks()

	go web.GatawayWeb()

	tiker := time.NewTicker(time.Minute * 5) //每15秒执行一下
	for {
		log.Println("主go程执行 抓拍进程管理程序 OK呢！")
		<-tiker.C

	}
}
