package service

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
	"vehicleLicensePlateRecognitionGateway/dto"
)

var StatisticalReportIpAddress string

//这里是车牌抓拍结果的统计上报接口
//1. 网关状态上报接口  POST
func GwStatusUploadPostWithJson(GWStudata *dto.GWStuStatisticalReportQeq) error {
	//post请求提交json数据
	data, _ := json.Marshal(*GWStudata)

	resp, err := http.Post(StatisticalReportIpAddress+"/report/gatewayrpt", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("post请求网关状态上报接口失败:", err)
		return err
	} else {
		log.Println("网关状态上传接口调用OK")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	Resp := new(dto.Resp)
	log.Println("网关状态上传接口 返回 body:", string(body))
	//
	unmerr := json.Unmarshal(body, &Resp)
	if unmerr != nil {
		log.Println("网关状态上传接口响应数据 json.Unmarshal error：", unmerr)
		return unmerr
	}
	log.Println("网关状态上传接口 ok")
	log.Println("网关状态上传接口 Post request with  json result:", Resp.Code, Resp.Msg)
	return nil
}

//2.2.	摄像机状态上报接口
func CameraStuUploadPostWithJson(CameraStudata *dto.CameraStuQeq) error {
	//post请求提交json数据
	data, _ := json.Marshal(*CameraStudata)

	resp, err := http.Post(StatisticalReportIpAddress+"/report/camrpt", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("post请求摄像机状态上报接口失败:", err)
		return err
	} else {
		log.Println("摄像机状态上报接口调用OK")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	Resp := new(dto.Resp)
	log.Println("摄像机状态上报接口返回 body:", string(body))
	//
	unmerr := json.Unmarshal(body, &Resp)
	if unmerr != nil {
		log.Println("摄像机状态上报接口响应数据 json.Unmarshal error：", unmerr)
		return unmerr
	}
	log.Println("摄像机状态上报接口 ok")
	log.Println("摄像机状态上报接口 Post request with json result:", Resp.Code, Resp.Msg)
	return nil
}

//2.3.	异常上报接口
func ExcprptStuUploadPostWithJson(ExcprptStudata *dto.ExcprptStuQeq) error {
	//post请求提交json数据
	data, _ := json.Marshal(*ExcprptStudata)

	resp, err := http.Post(StatisticalReportIpAddress+"/report/excprpt", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("post请求异常上报接口失败:", err)
		return err
	} else {
		log.Println("异常上报接口调用OK")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	Resp := new(dto.Resp)
	log.Println("异常上报接口返回 body:", string(body))
	//
	unmerr := json.Unmarshal(body, &Resp)
	if unmerr != nil {
		log.Println("异常上报接口响应数据 json.Unmarshal error：", unmerr)
		return unmerr
	}
	log.Println("异常上报接口 ok")
	log.Println("异常上报接口 Post request with  json result:", Resp.Code, Resp.Msg)
	return nil
}

//2.5.	版本查询接口  网关每隔10分钟轮询请求服务器的版本
func VersionQeqUploadPostWithJson(VersionQeqdata *dto.VersionQeq) error {
	//post请求提交json数据
	data, _ := json.Marshal(*VersionQeqdata)

	resp, err := http.Post(StatisticalReportIpAddress+"/report/versionquery", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("post请求版本查询接口失败:", err)
		return err
	} else {
		log.Println("版本查询接口调用OK")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	Resp := new(dto.Resp)
	log.Println("版本查询接口 body:", string(body))
	//
	unmerr := json.Unmarshal(body, &Resp)
	if unmerr != nil {
		log.Println("版本查询接口响应数据 json.Unmarshal error：", unmerr)
		return unmerr
	}
	log.Println("版本查询接口 ok")
	log.Println("版本查询接口 Post request with  json result:", Resp.Code, Resp.Msg)
	return nil
}

//处理统计数据上传到云平台
//定时20秒网关上报自身状态、摄像机状态状态至平台
func StatisticalReport() {
	tiker := time.NewTicker(time.Second * 20)
	for {
		<-tiker.C
		Gatewayrpt()
		//Camrpt()
		log.Println("定时20秒上报自身状态、摄像机状态状态至平台 ok")
	}

}

//网关每隔10分钟轮询请求服务器的版本
func VersionQeq() {
	tiker := time.NewTicker(time.Minute * 10) //暂时3分钟
	for {
		<-tiker.C
		vs := new(dto.VersionQeq)
		vs.GatewayId = Deviceid //1	gatewayId		网关id
		vs.Curver = ""          //2	curver	v1.0.21_20201221gateway	版本号
		vs.CurverNum = 0        //3	curverNum	21	数字版本号，打包一次版本号+1
		vs.ReqTime = ""         //4	reqTime	2020-02-04 15:01:04	版本请求时间
		qverr := VersionQeqUploadPostWithJson(vs)
		if qverr != nil {
			log.Println(qverr)
			return
		}

		//若服务器返回版本号高于当前版本号，则网关从oss地址上下载文件，自行更新。

		//把结果返回给魏俊一
		log.Println("网关每隔10分钟轮询请求服务器的版本 ok")
	}

}

//定时20秒上报网关自身状态
func Gatewayrpt() {
	//获取数据

	gwstudata := new(dto.GWStuStatisticalReportQeq)
	//数据赋值
	gwstudata.GatewayId = Deviceid //1	gatewayId		网关id
	gwstudata.VerDes = MainVersion //2	verDes	v1.0.21_20201221_gw	程序版本号
	gwstudata.VerNum = "1.5.678.9" //3	verNum	54	数字版本号

	gwstudata.ReportTime = time.Now().Format("2006-01-02 15:04:05") //4	reportTime	2020-12-21 12:05:12	上报时间

	gwstudata.ProgramStartTime = MainStartTime //5	programStartTime	2020-12-21 01:01:01	程序启动时间

	gwstudata.CamCnt = CameraCount    //6	camCnt	10	摄像机数量
	gwstudata.CamErrCnt = 0           //7	camErrCnt	1	有问题的摄像机数量
	gwstudata.IpAddr = ""             //8	ipAddr	10.132.12.42	网关IP地址
	gwstudata.OsVer = "windows"       //9	osVer	win10.111123.23	操作系统版本号
	gwstudata.CapCnt = CapCnt         //10capCnt11231启动后抓拍总和
	gwstudata.CapZeroCnt = CapZeroCnt //11capZeroCnt112每日零点后抓拍的总和

	gwstudata.UploadRecordCnt = UploadRecordCnt         //12uploadRecordCnt1123启动后上传的总和
	gwstudata.UploadRecordZeroCnt = UploadRecordZeroCnt //13uploadRecordZeroCnt112每日零点后上传的总和

	gwstudata.UploadImgCnt = UploadImgCnt         //14uploadImgCnt1123启动后上传的总和
	gwstudata.UploadImgZeroCnt = UploadImgZeroCnt //15uploadImgZeroCnt112每日零点后上传的总和

	gwstudata.UploadFailCnt = UploadFailCnt         //16uploadFailCnt12启动后上传失败的总和
	gwstudata.UploadFailZeroCnt = UploadFailZeroCnt //17uploadFailZeroCnt1每日零点后上传的失败总和

	gwstudata.UploadFailImgCnt = UploadFailImgCnt         //18uploadFailImgCnt12启动后上传失败的总和
	gwstudata.UploadFailImgZeroCnt = UploadFailImgZeroCnt //19uploadFailImgZeroCnt1每日零点后上传的失败总和
	gwstudata.UnUploadCnt = CapCnt - UploadRecordCnt      //当前未上传的数据总和

	gwstudata.DiskUsed = "70.54" //utils.GetDisk() //21diskUsed31使用当前硬盘盘符百分比
	log.Println("网关上报数据:", gwstudata)
	//上报数据
	sbaoerr := GwStatusUploadPostWithJson(gwstudata)
	if sbaoerr != nil {
		log.Println(sbaoerr)
		return
	}
}

//摄像机状态状态至平台
func Camrpt(data *dto.StatusResult) {
	//获取数据

	Camerastudata := new(dto.CameraStuQeq)
	//数据赋值
	Camerastudata.GatewayId = Deviceid //1	gatewayId		网关id
	CameraId := data.Uuid
	Camerastudata.CameraId = CameraId  //2	cameraId		摄像机id
	Camerastudata.VerDes = MainVersion //3	verDes	v1.0.21_20201221_cam	程序版本号
	Camerastudata.VerNum = data.VerNum //4	verNum	54	数字版本号

	Camerastudata.ReportTime = data.ReportTime //time.Now().Format("2006-01-02 15:04:05") //5	reportTime	2020-12-21 12:05:12	上报时间

	Camerastudata.ProgramStartTime = data.ProgramStartTime //6	programStartTime	2020-12-21 01:01:01	程序启动时间
	Camerastudata.CamBrand = data.CamBrand                 //7	camBrand	华为	品牌信息
	Camerastudata.CamStatus = data.CamStatus               //8	camStatus	0	摄像机状态 0 : 正常； -1: 连接摄像机网络失败； -2：摄像机注册/登陆失败； -3：摄像机异常(接口返回)； -4：24小时无数据；
	Camerastudata.CamStatusDes = data.CamStatusDes         //9	camStatusDes	正常	摄像机状态描述
	Camerastudata.ReConnCnt = data.ReConnCnt               //10	reConnCnt	15	进程启动到目前为止，网络重连次数

	Camerastudata.CapCnt = data.CapCnt //11	capCnt	1221	启动后摄像机抓拍总和
	CmeraCapCnt[CameraId] = data.CapCnt
	Camerastudata.CapZeroCnt = data.CapZeroCnt //12	capZeroCnt	112	启动后每日零时统计的总和

	LatestSnapshotTime[CameraId] = data.LastCaptime
	Camerastudata.LastCaptime = data.LastCaptime //13	lastCaptime	2020-05-10 15:01:02	最近一次抓拍的时间

	//上报数据
	log.Println("摄像机上报数据:", Camerastudata)
	sbaoerr := CameraStuUploadPostWithJson(Camerastudata)
	if sbaoerr != nil {
		//已经打印了不要重复记录 log.Println(sbaoerr)
		return
	}
}
