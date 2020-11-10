package dto

import "encoding/xml"

//车牌识别云抓拍数据 请求对象

//抓拍信息
type DateXML struct {
	XMLName     xml.Name    `xml:"data"`
	Token       string      `xml:"token"` // token
	LprInfo     LprInfo     `xml:"lprInfo"`
	LpaResult   LpaResult   `xml:"lpaResult"`
	VehicleInfo VehicleInfo `xml:"vehicleInfo"`
}

type LprInfo struct {
	XMLName        xml.Name `xml:"lprInfo"`
	PassId         string   `xml:"passId"`         // 过车编号
	CamId          string   `xml:"camId"`          //camId>    摄像机编号
	DeviceId       string   `xml:"deviceId"`       //deviceId>前置机编号
	PassTime       string   `xml:"passTime"`       //passTime>     过车编号
	VehicleImgPath string   `xml:"vehicleImgPath"` //vehicleImgPath>     过车图片地址
	PlateImgPath   string   `xml:"plateImgPath"`   //<plateImgPath/>     车牌图片地址
	BucketId       string   `xml:"bucketId"`       //bucketId>   bucket编号
	ImageType      string   `xml:"imageType"`      //	imageType> 图片类型
	UploadStamp    string   `xml:"uploadStamp"`    //	uploadStamp> 上传时间
	Stationid      string   `xml:"stationid"`      //	stationid>站点编号
	LaneType       string   `xml:"laneType"`       //	laneType> 出入口类型 0:入口；1：出口

}

type LpaResult struct {
	XMLName         xml.Name `xml:"lpaResult"`
	PassId          string   `xml:"passId"`          //passId>     过车编号
	EngineType      string   `xml:"engineType"`      //engineType>   引擎类型
	EngineId        string   `xml:"engineId"`        //engineId>     引擎编号
	PlateNo         string   `xml:"plateNo"`         //plateNo>     车牌编号
	PlateColor      string   `xml:"plateColor"`      //plateColor>     车牌颜色
	ComputeInterval string   `xml:"computeInterval"` //computeInterval>  计算时间
	VehicleColor    string   `xml:"vehicleColor"`    //vehicleColor>       车辆颜色
	VehicleType     string   `xml:"vehicleType"`     //vehicleType>       车辆类型
	VehicleBrand    string   `xml:"vehicleBrand"`    //vehicleBrand>       车辆品牌
	VehicleYear     string   `xml:"vehicleYear"`     //vehicleYear>     车辆年份
	LprFrameEntity  LprFrameEntity
}

type LprFrameEntity struct {
	XMLName     xml.Name `xml:"lprFrameEntity"`
	PlateLeft   string   `xml:"plateLeft"`   //plateLeft>        车牌左坐标
	PlateTop    string   `xml:"plateTop"`    //plateTop>        车牌上坐标
	PlateRight  string   `xml:"plateRight"`  //plateRight>        车牌右坐标
	PlateBottom string   `xml:"plateBottom"` //plateBottom>     车牌下坐标
}

type VehicleInfo struct {
	XMLName         xml.Name `xml:"vehicleInfo"`
	SideImgPath     string   `xml:"sideImgPath"` //sideImgPath> 侧面图片地址
	TailImgPath     string   `xml:"tailImgPath"` //tailImgPath> 车尾图片地址
	CarType         string   //CarType>  车辆型号
	AxleNum         string   //AxleNum>  轴数
	AxleType        string   //AxleType>  轴型
	WheelNumber     string   //WheelNumber> 轮胎数量
	AxleDist        string   //AxleDist>  轴距
	CarLengthMeter  string   //CarLengthMeter> 车长
	VideoScaleSpeed string   //VideoScaleSpeed> 车速
	WXPCharIndex    string   //WXPCharIndex>  危险品标识
	ZXType          string   //ZXType> 专项作业车标识
}

//
