package service

import (
	"github.com/sirupsen/logrus"
	"log"
	"testing"
	"time"
	"vehicleLicensePlateRecognitionGateway/config"
	"vehicleLicensePlateRecognitionGateway/dto"
	"vehicleLicensePlateRecognitionGateway/utils"
)

//
func TestGwCaptureInformationUploadPostWithXML(t *testing.T) {

	conf := config.ConfigInit() //初始化配置文件
	log.Println("配置文件信息：", *conf)
	utils.InitLogrus(conf.LogPath, conf.LogFileName, time.Duration(24*conf.LogMaxAge)*time.Hour, conf.LogRotationSize, time.Duration(conf.LogRotationTime)*time.Hour, conf.RotationCount)
	GwCaptureInformationUploadIpAddress = conf.GwCaptureInformationUploadIpAddress

	//data := new(dto.TBXJDateXML)
	//data.Token = "asgdgsajhdfgsajkg"
	//data.LprInfo.PassId = "sdknasasf"
	//data.LpaResult.PassId = "dfasfasasd"
	//data.VehicleInfo.AxleDist = "dfasfadafsd"
	//
	//
	//err, resp := GwCaptureInformationUploadPostWithXML(data)
	//if err != nil {
	//	logrus.Print("查询，失败:", err)
	//}
	//logrus.Println("查询结果：", resp)

}

//service.GetCameraListip=conf.GetCameraList
func TestGetCameraList(t *testing.T) {

	conf := config.ConfigInit() //初始化配置文件
	log.Println("配置文件信息：", *conf)
	utils.InitLogrus(conf.LogPath, conf.LogFileName, time.Duration(24*conf.LogMaxAge)*time.Hour, conf.LogRotationSize, time.Duration(conf.LogRotationTime)*time.Hour, conf.RotationCount)

	//
	GetCameraListip = conf.GetCameraList

	data := ""

	err, resp := GetCameraList(data)
	if err != nil {
		logrus.Print("查询，失败:", err)
	}
	logrus.Println("查询结果：", resp)
}

func TestHeartbeat(t *testing.T) {
	StatisticalReportIpAddress = "https://newydcpsbxt.jchc.cn/gateway-report"
	//Runmain()
	Heartbeat("6002")

}

func TestHeartbeatClientbb(t *testing.T) {
	hferr := Heartbeatclient("6001", []byte(`<message>
	<gatewayId>TTTT</gatewayId>
	<type>5</type>
	<verNum>1.5.678.9</verNum>
	<reportTime>1970-01-01 00:00:01</reportTime>
	<programStartTime>1970-01-01 00:00:01</programStartTime>
	<camBrand>UNIVIEW</camBrand>  
	<camStatus>0</camStatus>
	<camStatusDes>normal</camStatusDes> 				 
	<reConnCnt>80</reConnCnt>
	<capCnt>1221</capCnt>
	<capZeroCnt>112</capZeroCnt>
	<lastCaptime>1970-01-01 00:00:01</lastCaptime>
  </message>`))
	if hferr != nil {
		//	log.Println(address, "此log已经打印过了")
		//log已经打印过了

	} else {
		log.Println("回复成功，但是不退出，继续udp交互，")
	}

}

func TestHeartbeatclient(t *testing.T) {
	//Runmain()

	Heartbeatclient("4999", []byte("haha"))
}

func TestStatisticalResults(t *testing.T) {

	StatisticalResults(0, 0, 0, 0, 0, 0, 0, 0)
}

//
func TestStatisticalFile(t *testing.T) {

	StatisticalFile("abv\n")
}

//
func TestVersionFile(t *testing.T) {
	vs := "20210102T21h00m00s_build"
	vs = "\n" + vs
	VersionFile(vs)

}

func TestExcprptStuUploadPostWithJson(t *testing.T) {
	StatisticalReportIpAddress = "https://newydcpsbxt.jchc.cn/gateway-report"
	ycdata := new(dto.ExcprptStuQeq)
	ycdata.GatewayId = "ceshibbbbb"                              //1	gatewayId		网关id
	ycdata.CameraId = "ceshiccccc"                               //2	cameraId		摄像机id
	ycdata.ReportTime = time.Now().Format("2006-01-02 15:04:05") //3	reportTime	2020-12-21 12:05:12	上报时间
	ycdata.CamStatus = -5                                        //4	camStatus	0	摄像机状态 0 : 正常 -1: 连接摄像机网络失败； -2：摄像机注册/登陆失败； -3：摄像机异常(接口返回)； -4：24小时无数据；-5 心跳超时
	ycdata.CamStatusDes = "心跳时间差大于60秒，需要重启程序"                    //5	camStatusDes	正常	摄像机状态描述
	ycsberr := ExcprptStuUploadPostWithJson(ycdata)
	if ycsberr != nil {

	}

}

//()
func TestStatisticalReport(t *testing.T) {
	StatisticalReportIpAddress = "https://newydcpsbxt.jchc.cn/gateway-report"
	Deviceid = "aasffa"
	MainVersion = "main.exe"
	MainStartTime = "2021-01-07 12:12:12"
	CameraCount = 100
	CapCnt = 100
	CapZeroCnt = 100
	UploadRecordCnt = 100
	UploadRecordZeroCnt = 100
	UploadImgCnt = 100
	UploadImgZeroCnt = 100
	UploadFailCnt = 100

	UploadFailZeroCnt = 100
	UploadFailImgCnt = 100
	UploadFailImgZeroCnt = 100
	StatisticalReport()

}
