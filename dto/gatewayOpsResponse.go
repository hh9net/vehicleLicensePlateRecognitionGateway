package dto

//网关运维监控响应对象

type QueryResponse struct {
	Code    int `json:"code"  example:"200"` //3000
	CodeMsg string
	Data    interface{} `json:"data"`
	Message string      `json:"message" example:"响应成功信息"`
}

//网关运维成功响应
type GatewayOPSResponse struct {
	Code    int         `json:"code"  example:"200"`
	Data    interface{} `json:"data"`
	Message string      `json:"message" example:"响应成功信息"`
}

//网关运维失败响应
type GatewayOPSResponseFailure struct {
	Code    int         `json:"code"  example:"404"`
	Data    interface{} `json:"data"`
	Message string      `json:"message" example:"响应失败信息"`
}

//网关设备详情
type GatewayDeviceDetails struct {
	GatewayNumber int         `json:"gw_number"  example:"200"`
	Data          interface{} `json:"data"`
	Message       string      `json:"message" example:"响应成功信息"`
}

//网关基本信息接口响应
type GatewayDeviceMsgResp struct {
	Code int                `json:"code"`
	Date []GatewayDeviceMsg `json:"data"`
	Msg  string             `json:"msg"`
}

type GatewayDeviceMsg struct {
	MsgHead                   GatewayDeviceMsgHead `json:"msghead"`
	UpdateTime                string               `json:"updatetime"`                 //采集时间
	NetWorkDelay              string               `json:"networkdelay"`               //网络延迟 ms
	ProgrameRuntime           string               `json:"programe_runtime"`           //   网关运行时间，秒
	Deviceid                  string               `json:"deviceid"`                   //   设备ID
	Gatewayip                 string               `json:"gatewayip"`                  //   网关IP地址，多个地址则用”, ”分隔
	GetwayVersion             string               `json:"getway_version"`             //   场内网关版本号
	LastversionUpdatedatetime string               `json:"lastversion_updatedatetime"` //   场内网关版本最后更新成功时间
	AntennaInfos              []AntennaInfo        `json:"antenna_infos"`
}
type GatewayDeviceMsgHead struct {
	Parkid     string `json:"parkid"`      // 停车场ID
	CompanyId  string `json:"companyid"`   // 公司ID
	TerminalId string `json:"terminal_id"` // 设备ID，如CE4C37043A520C93
	Msgtype    string `json:"msgtype"`     //网关基本信息，值固定”3”
}

type AntennaInfo struct {
	Laneid                  string `json:"laneid"`   //  车道ID
	Rsuip                   string `json:"rsuip"`    // 天线IP地址
	Rsuport                 string `json:"rsuport"`  //  天线连接端口
	Power                   string `json:"power"`    // 天线功率
	Waittime                string `json:"waittime"` // 天线等待时间
	AllowRepeattime         string `json:"allow_repeattime"`
	Isregister              string `json:"isregister"`                //车道是否注册：1 表示已注册, 其它未注册
	AntennaStatus           string `json:"antenna_status"`            //   车道天线状态：1正常，其它值异常
	AntennaStatusUpdatetime string `json:"antenna_status_updatetime"` //   车道天线状态更新时间，如果与当前时间相差很大，说明天线状态也是异常的。
}

/*{
    "code": 0,
    "data": [{
        "msghead": {
            "parkid": "2002009998",停车场ID
            "companyid": "3202999999",公司ID
            "terminal_id": "CE4C37043A520C93",设备ID，如CE4C37043A520C93
            "msgtype": "3"网关基本信息，值固定”3”
        },
        "programe_runtime": "1912",网关运行时间，秒
        "deviceid": "CE4C37043A520C93"设备ID,
        "gatewayip": "192.168.200.215",网关IP地址，多个地址则用”,”分
        "getway_version": "build2020-06-29 08:56:41|ver1",场内网关版本号
        "lastversion_updatedatetime": "",场内网关最后更新成功时间
        "antenna_infos": [{
            "laneid": "1101",车道ID
            "rsuip": "192.168.200.248",天线IP地址
            "rsuport": "21003",天线连接端口
            "power": "12",天线功率
            "waittime": "12",天线等待时间
            "allow_repeattime": "0",【重启时间】
            "isregister": "1",车道是否注册：1 表示已注册,其它未注册
            "antenna_status": "0",车道天线状态：1正常，其它值异常
            "antenna_status_updatetime": "2020-09-17 16:02:35"车道天线状态更新时间
        }, {
            "laneid": "1102",
            "rsuip": "192.168.200.248",
            "rsuport": "21003",
            "power": "20",
            "waittime": "20",
            "allow_repeattime": "0",
            "isregister": "",
            "antenna_status": "",
            "antenna_status_updatetime": ""
        }]
    }],
    "msg": "SUCCESS"
}

*/

//告警信息接口响应
type ErrorMsgResp struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Date []ErrorMsg `json:"data"`
}

type ErrorMsg struct {
	Id            int    `json:"id"`              //告警的id、n9e_mon数据库中event表的id
	NodePath      string `json:"node_path"`       //对应报警策略中的节点
	Value         string `json:"value"`           //产生报警的事件，上报到n9e中的键值
	Etime         string `json:"etime"`           //"etime": "1600683220",报警的事件产生时的时间戳
	Name          string `json:"name"`            //"	"name": "网关日志总大小超过10GB",对应报警策略中的配置的策略名称
	Priority      string `json:"priority"`        //"	"priority": "1", 优先级、报警优先级
	Endpoint      string `json:"endpoint"`        //"	"endpoint": "4E944368CEE82941",节点
	EndpointAlias string `json:"endpoint_alias"`  //"	"endpoint_alias": "南京禄口机场T1-1",节点别名
	Event_type    string `json:"event_type"`      //"	"event_type": "recovery",事件类型：对应值为alert或recovery
	Status        string `json:"status"`          //"	"status": "1",状态: 整数值
	StatusName    string `json:"status_name"`     //"	"status_name": "已发送",【状态名称，可能值为处理中、已发送等】
	EventTypeName string `json:"event_type_name"` //"	"event_type_name": "恢复"【事件类型名称：报警、恢复】
}

/*{
    "code": 0,
    "msg": "SUCCESS",
    "data": [
        {
            "id": 892,
            "node_path": "cop.parkgateway",
            "value": "gateway.park.log.totalsize: 7669312",
            "etime": "1600683220",
            "name": "网关日志总大小超过10GB",
            "priority": "1",
            "endpoint": "4E944368CEE82941",
            "endpoint_alias": "南京禄口机场T1-1",
            "event_type": "recovery",
            "status": "1",
            "status_name": "已发送",
            "event_type_name": "恢复"
        },
        {
            "id": 891,
            "node_path": "cop.parkgateway",
            "value": "gateway.park.log.totalsize: 7713408",
            "etime": "1600682920",
            "name": "网关日志总大小超过10GB",
            "priority": "1",
            "endpoint": "4E944368CEE82941",
            "endpoint_alias": "南京禄口机场T1-1",
            "event_type": "alert",
            "status": "1",
            "status_name": "已发送",
            "event_type_name": "报警"
        },
    ]
}
*/

//查询最近metric指标值接口响应
//指标信息接口响应
type MetricMsgResp struct {
	Code          int       `json:"code"`
	Msg           string    `json:"msg"`
	MetricMsgDate MetricMsg `json:"data"` //指标数据
}

type MetricMsg struct {
	Metric string   `json:"metric"`
	Date   []Metric `json:"data"` //多个设备的指标
}

type Metric struct {
	Time     string  `json:"time"`     //"time": "2020-09-22 14:31:30",//采集时间
	Endpoint string  `json:"endpoint"` //"endpoint": "0A4924F0F82DDF19",设备id
	Value    float64 `json:"value"`    //"value": 3.8 指标值
}

//var MetricMsgMapList map[string]MetricMsgmap
//
//type MetricMsgmap struct {
//	Time     string `json:"time"`
//	Endpoint string `json:"endpoint"`
//	Value    string `json:"value"`
//}

/*{
    "code": 0,
    "msg": "SUCCESS",
    "data": {
        "metric": "gateway.park.gateway.mempercent",
        "data": {
            "03854A46E9875454": {
                "time": "2020-09-21 16:33:30",
                "endpoint": "03854A46E9875454",
                "value": 0.5
            },
            "E6EA54314E3AFEAF": {
                "time": "2020-09-21 16:32:30",
                "endpoint": "E6EA54314E3AFEAF",
                "value": 0.5
            },
            "FFBDBFD2D769D023": {
                "time": "2020-09-21 16:32:30",
                "endpoint": "FFBDBFD2D769D023",
                "value": 0.6
            }
        }
    }
}
*/

//网关基本信息查询
type QueryGatewayListResp struct {
	TerminalId                string  `json:"terminal_id"`   // 设备ID，如CE4C37043A520C93
	Parkid                    string  `json:"parkid"`        // 停车场ID
	ParkName                  string  `json:"park_name"`     // 停车场名称
	CompanyId                 string  `json:"companyid"`     // 公司ID
	CompanyName               string  `json:"company_Name"`  // 公司ID
	OnlineStatus              int     `json:"online_status"` //"	"status": "1"： 在线状态 0 :离线
	Gatewayip                 string  `json:"gatewayip"`     //   网关IP地址，多个地址则用”, ”分隔
	CPU                       float64 `json:"cpu_percent"`
	MEMpercent                float64 `json:"mem_percent"`
	MEM                       float64 `json:"mem"`
	DISKpercent               float64 `json:"disk_percent"`
	DISK                      float64 `json:"disk"`
	UnprocessedErrors         int     `json:"unprocessed_errors"`
	Errors                    int     `json:"errors"`
	Restarts                  int     `json:"restarts"`
	GetwayVersion             string  `json:"getway_version"`             //   场内网关版本号
	LastversionUpdatedatetime string  `json:"lastversion_updatedatetime"` //   场内网关最后更新成功时间
	RsuNum                    int     `json:"offline_rsu_num"`
	RsuALLNum                 int     `json:"rsu_all_num"`
	Network                   int64   `json:"net_work"`
	Flag                      bool    `json:"flag"`     //前端需要的一个标记
	SortFlag                  bool    `json:"sortflag"` //前端需要的一个标记
}

//告警信息查询
type QueryErrorListResp struct {
	TerminalId    string `json:"terminal_id"`    // 设备ID，如CE4C37043A520C93
	ErrorTime     string `json:"error_time"`     //
	ErrorDescribe string `json:"error_describe"` //
	ManId         string `json:"man_id"`
	ManName       string `json:"man_name"` //
	Time          string `json:"time"`
}

//重启信息查询
type QueryRestartListResp struct {
	TerminalId  string `json:"terminal_id"` // 设备ID，如CE4C37043A520C93
	RestartTime string `json:"error_time"`  //
	WorkTime    string `json:"work_time"`   //
	//ManName     string `json:"man_name"`    //
	//Type        string `json:"type"`        //"	 重启类型
}

//重启信息接口响应
type RestartMsgResp struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Date RestartMsgData `json:"data"`
}
type RestartMsgData struct {
	Metric  string    `json:"metric"`
	Datamsg []DataMsg `json:"data"` //"etime": "1600683220",报警的事件产生时的时间戳
}
type DataMsg struct {
	Time     string `json:"time"`     //"etime": "1600683220",报警的事件产生时的时间戳
	Endpoint string `json:"endpoint"` //"	"endpoint": "4E944368CEE82941",节点
	Value    int    `json:"value"`    //产生报警的事件，上报到n9e中的键值
}

//天线列表请求信息
type QueryRSUMsgListResp struct {
	TerminalId              string `json:"terminal_id"`               // 设备ID，如CE4C37043A520C93	//网关id
	RSUIP                   string `json:"rsu_ip"`                    // 天线ip
	WorkTime                string `json:"work_time"`                 // 连续工作时长lane
	Lane                    string `json:"lane"`                      // 车道
	Isregister              string `json:"isregister"`                //车道是否注册：1 表示已注册, 其它未注册
	AntennaStatus           string `json:"antenna_status"`            //   车道天线状态：1正常，其它值异常
	AntennaStatusUpdatetime string `json:"antenna_status_updatetime"` //   车道天线状态更新时间，如果与当前时间相差很大，说明天线状态也是异常的。

}

//网关基础信息
type QueryGatewayOneResp struct {
	TerminalId    string  `json:"terminal_id"`    // 设备ID，如CE4C37043A520C93
	ParkName      string  `json:"park_name"`      // 停车场名称
	Gatewayip     string  `json:"gateway_ip"`     //   网关IP地址，多个地址则用”, ”分隔
	GetwayVersion string  `json:"getway_version"` //   场内网关版本号
	WorkTime      string  `json:"work_time"`
	RestartTime   string  `json:"restart_time"`
	Restarts      int     `json:"restarts"`
	CPU           float64 `json:"cpu_percent"`
	MEMpercent    float64 `json:"mem_percent"`
	MEM           float64 `json:"mem"`
	DISKpercent   float64 `json:"disk_percent"`
	DISK          float64 `json:"disk"`
	Network       int64   `json:"net_work"`

	//UnprocessedErrors         int     `json:"unprocessed_errors"`
	//	Errors                    int     `json:"errors"`
	//LastversionUpdatedatetime string  `json:"lastversion_updatedatetime"` //   场内网关最后更新成功时间
	//RsuNum                    int     `json:"rsu_num"`
	//CompanyId                 string  `json:"companyid"`     // 公司ID
	//CompanyName               string  `json:"company_Name"`  // 公司ID
	//OnlineStatus              int     `json:"online_status"` //"	"status": "1"： 在线状态 0 :离线
	//	Parkid                    string  `json:"parkid"`        // 停车场ID

}

//查询软件版本列表信息
type QueryVersionListResp struct {
	Version     string `json:"version"`
	Time        string `json:"time"`
	VersionNote string `json:"version_note"`
	Num         int    `json:"num"` //运行版本数据
	Name        string `json:"name"`
}

type PerformVersionUpdateResp struct {
	TerminalId       string `json:"terminal_id"`      //设备ID，如CE4C37043A520C93
	Upgrade          string `json:"upgrade"`          //是否需要升级:0不需要 1需要升级
	GatewayVersion   string `json:"gateway_version"`  //网关新版本号
	Download_url     string `json:"download_url"`     //网关下载的URL
	CurrversionMd5   string `json:"currversion_md5"`  //场内网关gateway文件MD5值
	Decrypt_password string `json:"decrypt_password"` //版本文件zip的解压密码(暂不加密压缩)
}

type QueryVersionsResp struct {
	Versions []VersionMsg `json:"versions"` // 软件版本号
}

type VersionMsg struct {
	Version     string `json:"version"`      // 软件版本号
	VersionNote string `json:"version_note"` // 软件版本号
}

type QueryGatewaysResp struct {
	TerminalId []string `json:"terminal_id"` // 软件版本号
}

type QueryParkNamesResp struct {
	Parkmsg []ParkMSG `json:"parkmsg"`
}

type ParkMSG struct {
	ParkName string `json:"park_name"` //停车场名称
	ParkNum  string `json:"park_id"`   // 停车场编号
}
