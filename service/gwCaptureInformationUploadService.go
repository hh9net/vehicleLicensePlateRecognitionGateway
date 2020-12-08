package service

import (
	"bytes"
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
func GwCaptureInformationUploadPostWithXML(data *[]byte) (*dto.ResultRespXML, error) {
	//post请求提交xml数据

	//POST
	//text/xml 传输数据为Xml数据
	resp, err := http.Post(GwCaptureInformationUploadIpAddress, "application/xml", bytes.NewBuffer(*data))
	if err != nil {
		log.Println("post请求指标信息查询接口失败:", err)
		return nil, err
	} else {
		log.Println("前置机抓拍信息上传接口调用OK")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	Resp := new(dto.ResultRespXML)
	log.Println("前置机抓拍信息上传接口 返回 body:", string(body))
	//
	unmerr := xml.Unmarshal(body, &Resp)
	if unmerr != nil {
		log.Println("前置机抓拍信息上传接口响应数据 xml.Unmarshal error：", unmerr)
		return nil, unmerr
	}
	log.Println("前置机抓拍信息上传接口 Post request with  xml result:", Resp.Code, Resp.Msg)
	return Resp, nil
}

//2.获取token
func GetToken(deviceid string) (*dto.GetTokenRespXML, error) {

	//http://172.31.49.252/processor-control/collect/token/fe0442b5-2d40-486f-9682-d1043ceca4e5
	resp, err := http.Get(Gettoken + deviceid)
	if err != nil {
		log.Println("GetToken http error!", err)
		return nil, err
	}
	log.Println("GetToken http ok!")

	defer func() {
		_ = resp.Body.Close()
	}()

	body, _ := ioutil.ReadAll(resp.Body)

	Resp := new(dto.GetTokenRespXML)
	unmerr := xml.Unmarshal(body, Resp)
	if unmerr != nil {
		log.Println("xml.Unmarshal error", unmerr)
		return nil, unmerr
	}

	log.Println("Post request with  xml result:", Resp.Code, Resp.Msg)
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
	unmerr := xml.Unmarshal(body, &Resp)
	if unmerr != nil {
		log.Println(" xml.Unmarshal error", unmerr)
		return nil, unmerr
	}

	log.Println("Post request GetCameraList  with  xml result:", len(Resp.Data), Resp.Data[0].Name)
	return Resp, nil
}
