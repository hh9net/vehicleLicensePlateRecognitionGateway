package service

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"vehicleLicensePlateRecognitionGateway/dto"
)

var GwCaptureInformationUploadIpAddress string
var Gettoken string
var GetCameraListip string

//捕捉信息上传接口 gwCaptureInformationUploadService

//1.前置机抓拍信息上传接口 http://172.31.49.252/data-collect/report/frontend   POST
func GwCaptureInformationUploadPostWithXML(data *dto.DateXML) (*dto.ResultRespXML, error) {
	//post请求提交xml数据

	//data := dto.DateXML{}
	//MarshalIndent 有缩进 xml.Marshal ：无缩进
	ba, _ := xml.MarshalIndent(data, "  ", "  ")
	log.Println("+++++++++", string(ba))

	log.Println("Address:", GwCaptureInformationUploadIpAddress, data)
	//POST
	//text/xml 传输数据为Xml数据
	resp, err := http.Post(GwCaptureInformationUploadIpAddress, "text/xml", bytes.NewBuffer(ba))
	if err != nil {
		log.Println("post请求指标信息查询接口失败:", err)
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	Resp := new(dto.ResultRespXML)
	//
	unmerr := json.Unmarshal(body, &Resp)
	if unmerr != nil {
		log.Println("json.Unmarshal error", unmerr)
	}
	log.Println("Post request with json result:", Resp.Code, Resp.Msg)
	return Resp, nil
}

//2.获取token
func GetToken(deviceid string) (*dto.GetTokenRespXML, error) {

	resp, err := http.Get(Gettoken + deviceid)
	if err != nil {
		log.Println("GetToken http error!", err)
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, _ := ioutil.ReadAll(resp.Body)

	Resp := new(dto.GetTokenRespXML)
	unmerr := json.Unmarshal(body, &Resp)
	if unmerr != nil {
		log.Println("json.Unmarshal error", unmerr)
		return nil, unmerr
	}

	log.Println("Post request with json result:", Resp.Code, Resp.Msg)
	return Resp, nil
}

//3.根据token获取camera列表
func GetCameraList(token string) (*dto.GetCameraList, error) {

	resp, err := http.Get(GetCameraListip + token)
	if err != nil {
		log.Println("GetCameraList http error!", err)
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, _ := ioutil.ReadAll(resp.Body)

	Resp := new(dto.GetCameraList)
	unmerr := json.Unmarshal(body, &Resp)
	if unmerr != nil {
		log.Println("json.Unmarshal error", unmerr)
		return nil, unmerr
	}

	log.Println("Post request GetCameraList  with json result:", len(Resp.Data), Resp.Data[0].Name)
	return Resp, nil
}
