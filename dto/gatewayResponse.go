package dto

import "encoding/xml"

//1、抓拍信息响应信息    响应对象
type ResultRespXML struct {
	XMLName xml.Name `xml:"result"`
	Code    int      `xml:"code"`
	Msg     string   `xml:"msg"`
}

/*
<result>
 <code>0</code>
 <msg>接收成功</msg>
</result>
*/

//1、获取token
type GetTokenRespXML struct {
	XMLName xml.Name `xml:"data"`
	Token   string
	Expire  string
	Code    int
	Msg     string
	Oss     Oss
}
type Oss struct {
	XMLName      xml.Name `xml:"oss"`
	BacketName   string   `xml:"backetName"`
	ObjectPrefix string   `xml:"objectPrefix"`
}

/*
第一版
<?xml version="1.0" encoding="utf-8"?>

<data>
  <Token>fddac3a1fe31ef3c1eeb</Token>
  <Expire>2h</Expire>
  <Code>0</Code>
  <Msg>成功</Msg>
  <oss>aqwreqwf</oss>
</data>
*/

type GetCameraList struct {
	XMLName xml.Name `xml:"ListData"`
	Data    []CameraListData
}
type CameraListData struct {
	XMLName     xml.Name `xml:"Data"`
	Id          string   `xml:"id"`
	Name        string   `xml:"name"`
	StationId   string   `xml:"stationId"`
	LaneType    string   `xml:"laneType"`
	DevCompId   string   `xml:"devCompId"`
	Description string   `xml:"description"`
	DevIp       string   `xml:"devIp"`
	Port        string   `xml:"port"`
	UserName    string   `xml:"userName"`
	Password    string   `xml:"password"`
	Channel     string   `xml:"channel"` //新增
	LaneNo      string   `xml:"laneNo"`  //新增
}
