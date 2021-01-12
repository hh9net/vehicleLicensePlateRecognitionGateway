package gatewayWeb

import (
	log "github.com/sirupsen/logrus"
	"time"
	"vehicleLicensePlateRecognitionGateway/service"
	"vehicleLicensePlateRecognitionGateway/utils"
)

func GatawayWeb() {
	GatawayWebData()
	conf := ConfigInit() //初始化配置文件
	log.Println("GatawayWeb 配置文件信息：", *conf)
	utils.InitLogrus(conf.WebLogPath, conf.WebLogFileName, time.Duration(24*conf.WebLogMaxAge)*time.Hour, conf.WebLogRotationSize, time.Duration(conf.WebLogRotationTime)*time.Hour, conf.WebRotationCount)

	IpAddress := conf.IpAddress

	UserName = conf.UserName
	Password = conf.Password

	RouteInit(IpAddress)

	tiker := time.NewTicker(time.Second * 15) //每15秒执行一下
	for {
		<-tiker.C
		log.Println("gwWeb程序")
	}

}

//测试使用
func GatawayWebData() {
	service.Deviceid = "aaaaaaa"      //1、网关id
	service.MainVersion = "main——001" //2、版本号

	service.MainStartTime = time.Now().Format("2006-01-02 15:04:05") //3、启动时间

	service.CameraCount = 4 //4、摄像头数量
	Gatewaylocation = "2"   //5、网关类型  1门架、2、服务区[默认] 3、收费站 +站点

	service.CapCnt = 100    //1、网关启动后共抓拍照片数量
	service.CameraCount = 4 //2、正常摄像头数量

}
