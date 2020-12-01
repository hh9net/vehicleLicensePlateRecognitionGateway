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
	Expire  string //过期时间 假的
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
  <Token>a2caedfcb22b21bedafe</Token>
  <Expire>2h</Expire>
  <Code>0</Code>
  <Msg>请求成功</Msg>
  <oss>
    <backetName>ydcpsbxt</backetName>
    <objectPrefix>cloud_lpr/jiangsu/suhuaiyangs</objectPrefix>
  </oss>
</data>
*/

//Signalway          信路威
//HIKITS             海康ITS
//HUAWEI    华为
//UNIVIEW    宇视
//dahua             大华
//HIK             海康
//JUDE             聚德
//JINSANLI 金三立
//DEYA             德亚
//HWTC200 汉王TC200
//SignalwayNew 信路威车型

//获取的相机列表
type GetCameraList struct {
	XMLName xml.Name `xml:"ListData"`
	Data    []CameraListData
}
type CameraListData struct {
	XMLName     xml.Name `xml:"Data"`
	Id          string   `xml:"id"`          //相机id    sxjgl_yzjtd_320200_G2_K1071_2_0_004
	Name        string   `xml:"name"`        //入口004
	StationId   string   `xml:"stationId"`   //站ID 3a6e449a18ed435e80bff3782709e6dd
	LaneType    string   `xml:"laneType"`    //车道 0
	DevCompId   string   `xml:"devCompId"`   //HIK 相机品牌
	Description string   `xml:"description"` //描述 1
	DevIp       string   `xml:"devIp"`       //devIp 相机IP 10.25.50.94
	Port        string   `xml:"port"`        //8000 相机P
	UserName    string   `xml:"userName"`    //admin
	Password    string   `xml:"password"`    //12345
	Channel     string   `xml:"channel"`     //新增0   通道 HIK 一对多用
	LaneNo      string   `xml:"laneNo"`      //新增1   车道编号
	Gantryid    string   `xml:"gantryid"`    //新增    门架id
}

/*
<ListData>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_1_0</id>
   <name>南区入口(卡口)</name>
   <stationId>9c667aef8bb64a1e99ba328e76cb1a65</stationId>
   <laneType>0</laneType>
   <devCompId>UNIVIEW</devCompId>
   <description>2</description>
   <devIp>10.113.1.37</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>0</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_1_1</id>
   <name>南区出口(卡口)</name>
   <stationId>9c667aef8bb64a1e99ba328e76cb1a65</stationId>
   <laneType>1</laneType>
   <devCompId>UNIVIEW</devCompId>
   <description>2</description>
   <devIp>10.113.1.36</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>0</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_2_0</id>
   <name>北区入口(卡口)</name>
   <stationId>20ebf41475174ff7a8ed46fc902aa3a4</stationId>
   <laneType>0</laneType>
   <devCompId>UNIVIEW</devCompId>
   <description>2</description>
   <devIp>10.113.1.16</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>0</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_2_1</id>
   <name>北区出口(卡口)</name>
   <stationId>20ebf41475174ff7a8ed46fc902aa3a4</stationId>
   <laneType>1</laneType>
   <devCompId>UNIVIEW</devCompId>
   <description>2</description>
   <devIp>10.113.1.17</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>0</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
</ListData>
*/
