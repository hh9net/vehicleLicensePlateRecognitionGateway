package dto

//json

//2.1.	网关状态上报
type GWStuStatisticalReportQeq struct {
	GatewayId            string `json:"gatewayId"`            //1	gatewayId		网关id
	VerDes               string `json:"verDes"`               //2	verDes	v1.0.21_20201221_gw	程序版本号
	VerNum               string `json:"verNum"`               //3	verNum	54	数字版本号[main程序版本]
	ReportTime           string `json:"reportTime"`           //4	reportTime	2020-12-21 12:05:12	上报时间
	ProgramStartTime     string `json:"programStartTime"`     //5	programStartTime	2020-12-21 01:01:01	程序启动时间
	CamCnt               int    `json:"camCnt"`               //6	camCnt	10	摄像机数量
	CamErrCnt            int    `json:"camErrCnt"`            //7	camErrCnt	1	有问题的摄像机数量
	IpAddr               string `json:"ipAddr"`               //8	ipAddr	10.132.12.42	网关IP地址
	OsVer                string `json:"osVer"`                //9	osVer	win10.111123.23	操作系统版本号
	CapCnt               int    `json:"capCnt"`               //10capCnt11231启动后抓拍总和
	CapZeroCnt           int    `json:"capZeroCnt"`           //11capZeroCnt112每日零点后抓拍的总和
	UploadRecordCnt      int    `json:"uploadRecordCnt"`      //12uploadRecordCnt1123启动后上传的总和
	UploadRecordZeroCnt  int    `json:"uploadRecordZeroCnt"`  //13uploadRecordZeroCnt112每日零点后上传的总和
	UploadImgCnt         int    `json:"uploadImgCnt"`         //14uploadImgCnt1123启动后上传的总和
	UploadImgZeroCnt     int    `json:"uploadImgZeroCnt"`     //15uploadImgZeroCnt112每日零点后上传的总和
	UploadFailCnt        int    `json:"uploadFailCnt"`        //16uploadFailCnt12启动后上传失败的总和
	UploadFailZeroCnt    int    `json:"uploadFailZeroCnt"`    //17uploadFailZeroCnt1每日零点后上传的失败总和
	UploadFailImgCnt     int    `json:"uploadFailImgCnt"`     //18uploadFailImgCnt12启动后上传失败的总和
	UploadFailImgZeroCnt int    `json:"uploadFailImgZeroCnt"` //19uploadFailImgZeroCnt1每日零点后上传的失败总和
	UnUploadCnt          int    `json:"unUploadCnt"`          //20unUploadCnt11当前未上传的数据总和
	DiskUsed             string `json:"diskUsed"`             //21diskUsed31使用当前硬盘盘符百分比
}

//2.2.	摄像机状态上报
type CameraStuQeq struct {
	GatewayId        string `json:"gatewayId"  example:"123"` //1	gatewayId		网关id
	CameraId         string `json:"cameraId"`                 //2	cameraId		摄像机id
	VerDes           string `json:"verDes"`                   //3	verDes	v1.0.21_20201221_cam	程序版本号
	VerNum           string `json:"verNum"`                   //4	verNum	54	数字版本号
	ReportTime       string `json:"reportTime"`               //5	reportTime	2020-12-21 12:05:12	上报时间
	ProgramStartTime string `json:"programStartTime"`         //6	programStartTime	2020-12-21 01:01:01	程序启动时间
	CamBrand         string `json:"camBrand"`                 //7	camBrand	华为	品牌信息
	CamStatus        int    `json:"camStatus"`                //8	camStatus	0	摄像机状态 0 : 正常； -1: 连接摄像机网络失败； -2：摄像机注册/登陆失败； -3：摄像机异常(接口返回)； -4：24小时无数据；
	CamStatusDes     string `json:"camStatusDes"`             //9	camStatusDes	正常	摄像机状态描述
	ReConnCnt        int    `json:"reConnCnt"`                //10	reConnCnt	15	进程启动到目前为止，网络重连次数
	CapCnt           int    `json:"capCnt"`                   //11	capCnt	1221	启动后摄像机抓拍总和
	CapZeroCnt       int    `json:"capZeroCnt"`               //12	capZeroCnt	112	启动后每日零时统计的总和
	LastCaptime      string `json:"lastCaptime"`              //13	lastCaptime	2020-05-10 15:01:02	最近一次抓拍的时间
}

//2.3.	异常上报接口
type ExcprptStuQeq struct {
	GatewayId    string `json:"gatewayId"  example:"123"` //1	gatewayId		网关id
	CameraId     string `json:"cameraId"`                 //2	cameraId		摄像机id
	ReportTime   string `json:"reportTime"`               //3	reportTime	2020-12-21 12:05:12	上报时间
	CamStatus    int    `json:"camStatus"`                //4	camStatus	0	摄像机状态 0 : 正常 -1: 连接摄像机网络失败； -2：摄像机注册/登陆失败； -3：摄像机异常(接口返回)； -4：24小时无数据；
	CamStatusDes string `json:"camStatusDes"`             //5	camStatusDes	正常	摄像机状态描述
}

//2.5.	版本查询接口
type VersionQeq struct {
	GatewayId string `json:"gatewayId"  example:"123"` //1	gatewayId		网关id
	Curver    string `json:"curver"`                   //2	curver	v1.0.21_20201221gateway	版本号
	CurverNum int    `json:"curverNum"`                //3	curverNum	21	数字版本号，打包一次版本号+1
	ReqTime   string `json:"reqTime"`                  //4	reqTime	2020-02-04 15:01:04	版本请求时间
}
