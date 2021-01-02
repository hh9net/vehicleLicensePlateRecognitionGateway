package service

import (
	"github.com/sirupsen/logrus"
	"log"
	"testing"
	"time"
	"vehicleLicensePlateRecognitionGateway/config"
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
	//Runmain()
	Heartbeat("5000")

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
