package dto

import "encoding/xml"

//抓拍信息响应信息    响应对象
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

type GetTokenRespXML struct {
	XMLName xml.Name `xml:"data"`
	Token   string
	Expire  string
	Code    int
	Msg     string
	Oss     string `xml:"oss"`
}

/*
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
	XMLName     xml.Name `xml:"data"`
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
}

//<?xml version="1.0" encoding="utf-8"?>
//
//<Data>
//<id>sxjgl_dongbugs_320600_G15_K1103_600_1_1_2_0</id> 相机id
//<name>行车2车道</name>
//<stationId>7dd2277d943a4f33a883f831f4a308ae</stationId>
//<laneType>0</laneType>
//<devCompId>aqwreqwf</devCompId>
//<description>3</description>
//<devIp>172.31.0.251</devIp>
//<port>80</port>
//<userName>admin</userName>
//<password>Admin12345</password>
//</Data>
