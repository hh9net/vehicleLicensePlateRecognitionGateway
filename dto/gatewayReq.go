package dto

import "encoding/xml"

//车牌识别云抓拍数据 请求对象
//抓拍信息
type DateXML struct {
	XMLName   xml.Name  `xml:"data"`
	Token     string    `xml:"token"` // token
	LprInfo   LprInfo   `xml:"lprInfo"`
	LpaResult LpaResult `xml:"lpaResult"`
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
	ImageType      int      `xml:"imageType"`      //	imageType> 图片类型
	UploadStamp    int64    `xml:"uploadStamp"`    //	uploadStamp> 上传时间
	Stationid      string   `xml:"stationid"`      //	stationid>站点编号
	LaneType       int      `xml:"laneType"`       //	laneType> 出入口类型 0:入口；1：出口

}

type LpaResult struct {
	XMLName         xml.Name `xml:"lpaResult"`
	PassId          string   `xml:"passId"`          //passId>     过车编号
	EngineType      string   `xml:"engineType"`      //engineType>   引擎类型
	EngineId        string   `xml:"engineId"`        //engineId>     引擎编号
	PlateNo         string   `xml:"plateNo"`         //plateNo>     车牌编号
	PlateColor      string   `xml:"plateColor"`      //plateColor>     车牌颜色
	ComputeInterval int64    `xml:"computeInterval"` //computeInterval>  计算时间
	VehicleColor    string   `xml:"vehicleColor"`    //vehicleColor>       车辆颜色
	VehicleType     string   `xml:"vehicleType"`     //vehicleType>       车辆类型
	VehicleBrand    string   `xml:"vehicleBrand"`    //vehicleBrand>       车辆品牌
	VehicleYear     int      `xml:"vehicleYear"`     //vehicleYear>     车辆年份
	LprFrameEntity  LprFrameEntity
}

type LprFrameEntity struct {
	XMLName     xml.Name `xml:"lprFrameEntity"`
	PlateLeft   int      `xml:"plateLeft"`   //plateLeft>        车牌左坐标
	PlateTop    int      `xml:"plateTop"`    //plateTop>        车牌上坐标
	PlateRight  int      `xml:"plateRight"`  //plateRight>        车牌右坐标
	PlateBottom int      `xml:"plateBottom"` //plateBottom>     车牌下坐标
}

/*
<data>
    <token>fddac3a1fe31ef3c1eeb</token>  token
    <lprInfo>
        <passId>sxjgl_shygs_320800011141003034288_20180705_001369
        </passId>     过车编号
        <camId>sxjgl_shygs_320800011141003034288</camId>    摄像机编号
        <deviceId>86a85d73-c381-4f19-a291-f5a2eab3b18e</deviceId>前置机编号
        <passTime>2018-07-05 09:50:49</passTime>     过车编号
        <vehicleImgPath>/shygs/sxjgl_shygs_320800011141003034288/20180705/sxjgl_shygs_320800011141003034288_20180705_001389.jpg
        </vehicleImgPath>     过车图片地址
        <plateImgPath/>     车牌图片地址
        <bucketId>042e8dab-8019-47b6-8606-adf6cf8eaa2b</bucketId>bucket编号
        <imageType>1</imageType> 图片类型
        <uploadStamp>1530755269798</uploadStamp> 上传时间
        <stationid>86a85d73-c381-4f19-a291-f5a2eab3b18e</stationid>站点编号
        <laneType>1</laneType> 出入口类型 0:入口；1：出口
    </lprInfo>
    <lpaResult>
        <passId>sxjgl_shygs_320800011141003034288_20180705_001369</passId>     过车编号
        <engineType>sjk-camera-lpa</engineType>   引擎类型
        <engineId>Signalway</engineId>     引擎编号
        <plateNo>苏A09D08</plateNo>     车牌编号
        <plateColor>3</plateColor>     车牌颜色
        <computeInterval>0</computeInterval>  计算时间
        <vehicleColor>white</vehicleColor>       车辆颜色
        <vehicleType>car</vehicleType>       车辆类型
        <vehicleBrand>HAVAl</vehicleBrand>       车辆品牌
        <vehicleYear>1</vehicleYear>     车辆年份
        <lprFrameEntity>
            <plateLeft>0</plateLeft>        车牌左坐标
            <plateTop>0</plateTop>        车牌上坐标
            <plateRight>0</plateRight>        车牌右坐标
            <plateBottom>0</plateBottom>     车牌下坐标
        </lprFrameEntity>
    </lpaResult>
    <vehicleInfo>
        <sideImgPath>/shygs/sxjgl_shygs_320800011141003034288/20180705/sxjgl_shygs_320800011141003034288_20180705_001389.jpg
        </sideImgPath> 侧面图片地址
        <tailImgPath>/shygs/sxjgl_shygs_320800011141003034288/20180705/sxjgl_shygs_320800011141003034288_20180705_001389.jpg
        </tailImgPath> 车尾图片地址
        <CarType>1</CarType>  车辆型号
        <AxleNum>11</AxleNum>  轴数
        <AxleType>1</AxleType>  轴型
        <WheelNumber>2</WheelNumber> 轮胎数量
        <AxleDist>2</AxleDist>  轴距
        <CarLengthMeter>2</CarLengthMeter> 车长
        <VideoScaleSpeed>200</VideoScaleSpeed> 车速
        <WXPCharIndex>爆</WXPCharIndex>  危险品标识
        <ZXType>专1</ZXType> 专项作业车标识
    </vehicleInfo>
</data>*/

//特别的一个品牌相机抓拍信息 需要上传侧面信息   信路威车型
type TBXJDateXML struct {
	XMLName     xml.Name        `xml:"data"`
	Token       string          `xml:"token"` // token
	LprInfo     TBXJLprInfo     `xml:"lprInfo"`
	LpaResult   TBXJLpaResult   `xml:"lpaResult"`
	VehicleInfo TBXJVehicleInfo `xml:"vehicleInfo"`
}

type TBXJLprInfo struct {
	XMLName        xml.Name `xml:"lprInfo"`
	PassId         string   `xml:"passId"`         // 过车编号
	CamId          string   `xml:"camId"`          //camId>    摄像机编号
	DeviceId       string   `xml:"deviceId"`       //deviceId>前置机编号
	PassTime       string   `xml:"passTime"`       //passTime>过车编号
	VehicleImgPath string   `xml:"vehicleImgPath"` //vehicleImgPath> 过车图片地址
	PlateImgPath   string   `xml:"plateImgPath"`   //<plateImgPath/>车牌图片地址
	BucketId       string   `xml:"bucketId"`       //bucketId>   bucket编号
	ImageType      int      `xml:"imageType"`      //	imageType> 图片类型
	UploadStamp    int64    `xml:"uploadStamp"`    //	uploadStamp> 上传时间
	Stationid      string   `xml:"stationid"`      //	stationid>站点编号
	LaneType       int      `xml:"laneType"`       //	laneType> 出入口类型 0:入口；1：出口

}

type TBXJLpaResult struct {
	XMLName         xml.Name `xml:"lpaResult"`
	PassId          string   `xml:"passId"`          //passId>     过车编号
	EngineType      string   `xml:"engineType"`      //engineType>   引擎类型
	EngineId        string   `xml:"engineId"`        //engineId>     引擎编号
	PlateNo         string   `xml:"plateNo"`         //plateNo>     车牌编号
	PlateColor      string   `xml:"plateColor"`      //plateColor>     车牌颜色
	ComputeInterval int64    `xml:"computeInterval"` //computeInterval>  计算时间
	VehicleColor    string   `xml:"vehicleColor"`    //vehicleColor>       车辆颜色
	VehicleType     string   `xml:"vehicleType"`     //vehicleType>       车辆类型
	VehicleBrand    string   `xml:"vehicleBrand"`    //vehicleBrand>       车辆品牌
	VehicleYear     int      `xml:"vehicleYear"`     //vehicleYear>     车辆年份
	LprFrameEntity  LprFrameEntity
}

type TBXJLprFrameEntity struct {
	XMLName     xml.Name `xml:"lprFrameEntity"`
	PlateLeft   int      `xml:"plateLeft"`   //plateLeft>        车牌左坐标
	PlateTop    int      `xml:"plateTop"`    //plateTop>        车牌上坐标
	PlateRight  int      `xml:"plateRight"`  //plateRight>        车牌右坐标
	PlateBottom int      `xml:"plateBottom"` //plateBottom>     车牌下坐标
}

type TBXJVehicleInfo struct {
	XMLName         xml.Name `xml:"vehicleInfo"`
	SideImgPath     string   `xml:"sideImgPath"` //sideImgPath> 侧面图片地址
	TailImgPath     string   `xml:"tailImgPath"` //tailImgPath> 车尾图片地址
	CarType         string   //CarType>  车辆型号
	AxleNum         int      //AxleNum>  轴数
	AxleType        string   //AxleType>  轴型
	WheelNumber     string   //WheelNumber> 轮胎数量
	AxleDist        string   //AxleDist>  轴距
	CarLengthMeter  string   //CarLengthMeter> 车长
	VideoScaleSpeed string   //VideoScaleSpeed> 车速
	WXPCharIndex    string   //WXPCharIndex>  危险品标识
	ZXType          string   //ZXType> 专项作业车标识
}

//车牌抓拍的结果xml
type CaptureDateXML struct {
	XMLName         xml.Name `xml:"PlateInfo"`
	PassId          string   `xml:"passId"`          //<passId>sxjgl_yzjtd_320200_G2_K1071_2_0_004_20201124_000031</passId>
	CamId           string   `xml:"camId"`           //<camId>sxjgl_yzjtd_320200_G2_K1071_2_0_004</camId>摄像机编号
	PassTime        string   `xml:"passTime"`        //<passTime>2020-11-24 14:35:17</passTime>
	VehicleImgPath  string   `xml:"vehicleImgPath"`  //<vehicleImgPath>C:\Users\Administrator\Desktop\HSJDEBUG\images\20201124\sxjgl_yzjtd_320200_G2_K1071_2_0_004_20201124143417_000031.jpg</vehicleImgPath>
	VehicleImgPath1 string   `xml:"vehicleImgPath1"` //<vehicleImgPath1></vehicleImgPath1>
	VehicleImgPath2 string   `xml:"vehicleImgPath2"` //<vehicleImgPath2></vehicleImgPath2>
	Channel         string   `xml:"channel"`         //<channel>0</channel>
	PlateColor      string   `xml:"plateColor"`      //<plateColor>蓝</plateColor>
	PlateNo         string   `xml:"plateNo"`         //<plateNo>苏CA0E37</plateNo>
	AppedInfo       AppedInfo
}

type AppedInfo struct {
	AxleDist        string
	AxleNum         string
	AxleType        string
	CarLengthMeter  string
	CarType         string
	VideoScaleSpeed string
	WheelNumber     string
	WXPCharIndex    string
	ZXType          string
}

/*<PlateInfo>
    <passId>sxjgl_yzjtd_320200_G2_K1071_2_0_004_20201124_000031</passId>
    <camId>sxjgl_yzjtd_320200_G2_K1071_2_0_004</camId>
    <passTime>2020-11-24 14:35:17</passTime>
    <vehicleImgPath>C:\Users\Administrator\Desktop\HSJDEBUG\images\20201124\sxjgl_yzjtd_320200_G2_K1071_2_0_004_20201124143417_000031.jpg</vehicleImgPath>
    <vehicleImgPath1></vehicleImgPath1>
    <vehicleImgPath2></vehicleImgPath2>
    <channel>0</channel>
    <plateColor>蓝</plateColor>
    <plateNo>苏CA0E37</plateNo>
    <AppedInfo>
        <AxleDist></AxleDist>
        <AxleNum></AxleNum>
        <AxleType></AxleType>
        <CarLengthMeter></CarLengthMeter>
        <CarType></CarType>
        <VideoScaleSpeed></VideoScaleSpeed>
        <WheelNumber></WheelNumber>
        <WXPCharIndex></WXPCharIndex>
        <ZXType></ZXType>
    </AppedInfo>
</PlateInfo>*/

//发送与接受的心跳
type Heartbeatbasic struct {
	XMLName         xml.Name `xml:"message"`
	Uuid            string   `xml:"uuid"`            //<uuid>
	Type            int      `xml:"type"`            //<type>    1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
	Version         string   `xml:"version"`         //<version> 抓拍程序版本号
	Time            string   `xml:"time"`            //<time>     字符串2020-11-12 12:12:12
	Seq             int      `xml:"seq"`             //<seq>   消息序号累加
	Pid             string   `xml:"pid"`             //进程ID
	RunParamXMLFile string   `xml:"RunParamXMLFile"` //运行XML文件名
}

//1、心跳
type Heartbeat struct {
	Heartbeatbasic
	Content string `xml:"content"` //<content>    内容
}

//3、日志
type HeartbeatLog struct {
	Heartbeatbasic
	Content HeartbeatLogContent `xml:"content"` //<content>    内容
}

type HeartbeatLogContent struct {
	Level string `xml:"level"` // <level> //级别
	Data  string `xml:"data"`  //	<data> //字符串
}

// 心跳
/*
<message>
<uuid></uuid>
<type>1</type>//1
<version></version>
<time></time>
<seq></seq>
<pid>123999</pid>
<RunParamXMLFile></RunParamXMLFile>
<content></content>
</message>

eHeartbeatMsg  = 1,       // 心跳
eDataNotifyMsg = 2,       // 新数据通知
eLogMsg        = 3,       // 日志
eCloseApp      = 4,       // 采集进程被动关闭命令


// 心跳
<message>
<uuid></uuid>
<type>1</type>
<version></version>
<time></time>
<seq></seq>
<content></content>
</message>


//新数据通知
<message>
<uuid></uuid>
<type>2</type>
<version></version>
<time></time>
<seq></seq>
<content>d:\xxx\xxx.xml</content>//xml的完整路径
</message>


// 日志
<message>
<uuid></uuid>
<type>3</type>
<version></version>
<time></time>
<seq></seq>
<content>
	<level></level>//级别
	<data></data>//字符串
</content>
</message>


//返回给抓拍进程
//采集进程被动关闭命令
<message>
<uuid></uuid>
<type>4</type>
<version></version>
<time></time>
<seq></seq>
<content></content>
</message>









*/
