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
	//text/xml 传输数据为Xml数据
	resp, err := http.Post(GwCaptureInformationUploadIpAddress, "application/xml", bytes.NewBuffer(*data))
	if err != nil {
		log.Println("post请求指标信息查询接口失败:", err)
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	Resp := new(dto.ResultRespXML)
	log.Println("前置机抓拍信息上传接口调用OK,前置机抓拍信息上传接口返回body OK")
	//
	unmerr := xml.Unmarshal(body, &Resp)
	if unmerr != nil {
		log.Println("前置机抓拍信息上传接口响应数据 xml.Unmarshal error：", unmerr)
		log.Println("body:", string(body))
		return nil, unmerr
	}

	ResultCount = ResultCount + 1
	log.Println("前置机抓拍信息上传接口 ok ResultCount:", ResultCount /*, time.Now().Format("2006-01-02 15:04:05")*/)
	log.Println("前置机抓拍信息上传接口 Post result:", Resp.Code, Resp.Msg)
	return Resp, nil
}

//2.获取token
func GetToken(deviceid string) (*dto.GetTokenRespXML, error) {
	log.Println("Gettoken + deviceid:", Gettoken+deviceid)
	//http://172.31.49.252/processor-control/collect/token/fe0442b5-2d40-486f-9682-d1043ceca4e5
	resp, err := http.Get(Gettoken + deviceid)
	if err != nil {
		log.Println("GetToken http error!", err)
		return nil, err
	}
	log.Println("GetToken http 请求 ok!")

	defer func() {
		_ = resp.Body.Close()
	}()

	body, _ := ioutil.ReadAll(resp.Body)

	Resp := new(dto.GetTokenRespXML)
	unmerr := xml.Unmarshal(body, Resp)
	if unmerr != nil {
		log.Println("GetToken http 请求 ok，但是xml.Unmarshal error!")

		log.Println("xml.Unmarshal error：", unmerr)
		log.Println("body:", string(body))
		return nil, unmerr
	}
	log.Println("Post request with  xml result:", Resp.Code, Resp.Msg)
	return Resp, nil
}

//3.根据token获取camera列表
func GetCameraList(token string) (*dto.GetCameraList, error) {
	log.Println("根据token获取camera列表的url:", GetCameraListip+token)
	resp, err := http.Get(GetCameraListip + token)
	if err != nil {
		log.Println("GetCameraList http error!:", err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := ioutil.ReadAll(resp.Body)
	Resp := new(dto.GetCameraList)
	unmerr := xml.Unmarshal(body, &Resp)
	if unmerr != nil {
		log.Println("根据token获取camera列表 xml.Unmarshal error:", unmerr)
		log.Println("根据token获取camera列表 body:", string(body))
		return nil, unmerr
	}
	return Resp, nil
}
