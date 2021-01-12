package gatewayWeb

var (
	Gatewaylocation string
)

type GwResponse struct {
	Code    int         `json:"code"  example:"200"`
	Data    interface{} `json:"data"`
	Message string      `json:"message" example:"响应成功信息"`
}

//失败响应
type GWResponseFailure struct {
	Code    int         `json:"code"  example:"404"`
	Data    interface{} `json:"data"`
	Message string      `json:"message" example:"响应失败信息"`
}

//1、网关基本数据
type GatewayBasicDataResp struct {
	Gatewayid   string `json:"gateway_id"`   //1、网关id
	Version     string `json:"version"`      //2、版本号
	StartTime   string `json:"start_time"`   //3、启动时间
	CameraCnt   int    `json:"cameraCnt"`    //4、摄像头数量
	GatewayType string `json:"gateway_type"` //5、网关类型  1门架、2、服务区[默认] 3、收费站
	//StationCameraCnt int    `json:"station_cameraCnt"` //6、收费站摄像头数量
	//ServiceAreaCameraCnt int `json:"service_area_cameraCnt"` //7、服务区摄像头数量
}

//2、网关动态数据
type GatewayDynamicDataResp struct {
	IMgCnt            int `json:"imgCnt"`             //1、网关启动后共抓拍照片数量
	CameraNormalCnt   int `json:"camera_normalCnt"`   //2、正常摄像头数量
	CameraAbnormalCnt int `json:"camera_abnormalCnt"` //3、异常摄像头数量
	OfflineCameraCnt  int `json:"offline_cameraCnt"`  //4、离线摄像头数量
}

//3、摄像头基本信息数据
type CameraInfo struct {
	CameraId           string `json:"cameraId"`             //1、摄像头id
	BrandName          string `json:"brand_name"`           //2、品牌名称
	CameraIMGCnt       int    `json:"camera_imgCnt"`        //3、抓拍统计图片的数量
	LatestSnapshotTime string `json:"latest_snapshot_time"` //4、最近抓拍时间
	MainrestartCnt     int    `json:"main_restartCnt"`      //5、进程重启次数
	Location           string `json:"location"`             //6、所在位置 1门架、2、服务区 3、收费站
}

//4'摄像头基本信息列表
type CameraInfosResp struct {
	CameraInfos []CameraInfo `json:"camerainfos"`
}
