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

	data := new(dto.DateXML)
	data.Token = "asgdgsajhdfgsajkg"
	data.LprInfo.PassId = "sdknasasf"
	data.LpaResult.PassId = "dfasfasasd"
	data.VehicleInfo.AxleDist = "dfasfadafsd"

	err, resp := GwCaptureInformationUploadPostWithXML(data)
	if err != nil {
		logrus.Print("查询，失败:", err)
	}
	logrus.Println("查询结果：", resp)

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

func TestRunmain(t *testing.T) {
	//Runmain()
}
