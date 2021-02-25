package service

import (
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"vehicleLicensePlateRecognitionGateway/dto"
	"vehicleLicensePlateRecognitionGateway/utils"
)

var (
	Parsexmlcount int
	files         chan string
	CapCnt        int //启动后抓拍总和
	CapZeroCnt    int //每日零点后抓拍的总和

	UploadRecordCnt     int //启动后上传oss以及xml的总和
	UploadRecordZeroCnt int //每日零点后上传的总和

	UploadImgCnt     int // 启动后上传的总和
	UploadImgZeroCnt int // 每日零点后上传的总和

	UploadFailCnt     int //启动后上传失败的总和 会再次上传
	UploadFailZeroCnt int //每日零点后上传的失败总和 会再次上传

	UploadFailImgCnt     int //启动后上传失败的总和
	UploadFailImgZeroCnt int //每日零点后上传的失败总和

	OSSCount      int //上传成功的oss数量
	NewOSSCount   int //上传成功的信路威新版图片数量
	Parsed        int //上传成功的删除OK的抓拍xml数量
	ResultCount   int
	ResultOKCount int //
	AgainCount    int //二次上传成功的数量

	Deviceid  string //网关设备id Token
	StationId map[string]string
	DeviceId  map[string]string

	LaneType                 map[string]string
	ImageType                map[string]string
	Name                     map[string]string
	CmeraId                  map[string]string
	Cmeraid                  map[string]string
	EngineId                 map[string]string
	CmeraCapCnt              map[string]int
	LatestSnapshotTime       map[string]string
	Token                    string
	NewDataNotificationCount int //新数据通知
	// HasUploadFile []string
	Pid map[string]string

	BacketName    string
	ObjectPrefix  string
	MainVersion   string
	MainStartTime string
	CameraCount   int
)

const (
	Signalway    string = "Signalway"    // 信路威
	HIKITS       string = "HIKITS"       // 海康ITS !!!!!!!!
	HUAWEI       string = "HUAWEI"       // 华为
	UNIVIEW      string = "UNIVIEW"      // 宇视    !!!!!!!!
	DaHua        string = "dahua"        // 大华
	HIK          string = "HIK"          // 海康
	JUDE         string = "JUDE"         // 聚德
	JINSANLI     string = "JINSANLI"     // 金三立
	DeYa         string = "DEYA"         // 德亚
	HWTC         string = "HWTC200"      // 汉王TC200
	SignalwayNew string = "SignalwayNew" // 信路威车型 有侧面图片
	GDPort       int    = 5000           //固定 进程向我拨号的的端口
	EngineType   string = "sjk-camera-lpa"
)

// 进程管理
func ProcessManagementService( /*ch chan int*/ ) {

	PorT := 6000 //固定端口
	var P int
	P = PorT
	//1、获取网关设备的token
CQ:
	resp, getTokenerr := GetGatawayToken()
	if getTokenerr != nil {
		log.Println("获取网关设备的token 失败,重新请求!") //getTokenerr 已打印
		time.Sleep(time.Second * 10)
		goto CQ
	}

	//全局 Token  BacketName 	ObjectPrefix 赋值
	if resp != nil {
		log.Println("获取网关设备的token OK", time.Now())
		Token = resp.Token
		BacketName = resp.Oss.BacketName
		ObjectPrefix = resp.Oss.ObjectPrefix
		utils.BacketName = BacketName
		log.Println("Token:", Token)
		log.Println("BacketName:", BacketName)

		log.Printf("ObjectPrefix:%s", ObjectPrefix)
	}

CmlistQ:
	//2、根据token获取camera列表
	CameraList, listerr := GetGatawayCameraList()
	if listerr != nil {
		log.Println("获取相机列表错误:", listerr)
		time.Sleep(time.Second * 10)
		goto CmlistQ
	}

	CameraCount = len(CameraList.Data)
	log.Println(" 相机列表数据的len:", CameraCount)
	log.Println(" 相机列表数据:", CameraList.Data)

	DeviceId = make(map[string]string, 0)
	StationId = make(map[string]string, 0)
	LaneType = make(map[string]string, 0)
	ImageType = make(map[string]string, 0)
	Name = make(map[string]string, 0)
	CmeraId = make(map[string]string, 0)
	EngineId = make(map[string]string, 0)
	Cmeraid = make(map[string]string, 0)
	CmeraCapCnt = make(map[string]int, 0)
	LatestSnapshotTime = make(map[string]string, 0)
	Pid = make(map[string]string, 0)

	uniview := make([]dto.CameraListData, 0) // 宇视的列表
	hikITS := make([]dto.CameraListData, 0)  //ITS列表

	for i, cmera := range CameraList.Data {
		//StationId
		StationId[cmera.Id] = cmera.StationId
		DeviceId[cmera.Id] = cmera.Gantryid //deviceid应该用gantryID
		LaneType[cmera.Id] = cmera.LaneType
		ImageType[cmera.Id] = cmera.Description
		EngineId[cmera.Id] = cmera.DevCompId //相机品牌
		Name[cmera.Id] = cmera.Name          //入口004
		Cmeraid[cmera.Id] = cmera.Id

		log.Println(i, "StationId:", StationId[cmera.Id], cmera.StationId)
		log.Println(i, "DeviceId:", DeviceId[cmera.Id], cmera.Gantryid)
		log.Println(i, "LaneType:", LaneType[cmera.Id], cmera.LaneType)
		log.Println(i, "ImageType:", ImageType[cmera.Id], cmera.Description)
		log.Println(i, "EngineId:", EngineId[cmera.Id], cmera.DevCompId)
		log.Println(i, "Name:", Name[cmera.Id], cmera.Name)
		log.Println(i, "Cmeraid:", Cmeraid[cmera.Id], cmera.Id)

		//	进程类型
		conflx := ""
		if cmera.DevCompId == UNIVIEW || cmera.DevCompId == HIKITS {
			conflx = "one2many"
		} else {
			conflx = "one2ont"
		}

		Configfname := ""
		//1、生成进程配置文件
		//ConfigPath:="abc"
		switch conflx {
		case "one2ont":
			log.Println("one2ont,相机品牌是：", cmera.DevCompId)

			PorT = PorT + 2
			P = P + 2
			confdata := new(OneToOneConfig)

			confdata.DevCompId = cmera.DevCompId //品牌名称
			strporrt := strconv.Itoa(PorT)
			confdata.Uuid = cmera.Id + "+" + strporrt //方便确定是哪一个进程发出的数据 我取相机id+进程端口号
			confdata.Udplistenport = PorT             //我向进程拨号的端口号
			confdata.Udptxport = PorT - 1             // 进程向我拨号的的端口
			confdata.Devlist.Dev.DevIp = cmera.DevIp  //相机IP
			confdata.Devlist.Dev.Port = cmera.Port    //相机端口号
			confdata.Devlist.Dev.UserName = cmera.UserName
			confdata.Devlist.Dev.Password = cmera.Password
			confdata.Devlist.Dev.Id = cmera.Id //相机id

			CmeraId[strconv.Itoa(PorT)] = confdata.Uuid        //方便确定是哪一个进程发出的数据 我取相机id+进程端口号
			confdata.Channellist.Channel.Id = cmera.Id         //相机id
			confdata.Channellist.Channel.Index = cmera.Channel //通道号

			//一对一生成配置文件
			fname := generateConfigToOne(confdata)
			if fname != "" {
				Configfname = fname
				//2、进程启动
				//传 一个配置文件的绝对路径 全局唯一
				go Runmain(Configfname)

			} else {
				log.Println("一对一生成配置文件 为空 Configfname", Configfname)
				return
			}

		case "one2many":
			log.Println("one2many,相机品牌是：", cmera.DevCompId)
			//HIKITS      UNIVIEW
			if cmera.DevCompId == UNIVIEW {
				uniview = append(uniview, cmera)
			}

			if cmera.DevCompId == HIKITS {
				hikITS = append(hikITS, cmera)
			}
		}
		continue
	}

	log.Println("DeviceId:", DeviceId)
	log.Println("StationId:", StationId)
	log.Println("LaneType:", LaneType)
	log.Println("ImageType:", ImageType)

	if len(hikITS) == 0 && len(uniview) == 0 {
		log.Println("该网关设备没有海康ITS相机和宇视相机")
		return
	}

	if len(uniview) != 0 {
		YSconfdata := new(MoreToMoreConfig)
		//多对多启动
		P = P + 2
		YSconfdata.Uuid = UNIVIEW + "+" + strconv.Itoa(P) //方便确定是哪一个进程发出的数据 我取品牌名称+进程端口号
		YSconfdata.Udplistenport = P                      //我向进程拨号的端口号
		YSconfdata.Udptxport = P - 1                      // 进程向我拨号的的端口

		for _, ys := range uniview {
			YSconfdata.DevCompId = ys.DevCompId //品牌名称
			devdata := new(MoreToMoreConfigDev)

			devdata.Id = ys.Id
			devdata.DevIp = ys.DevIp
			devdata.Port = ys.Port
			devdata.UserName = ys.UserName
			devdata.Password = ys.Password

			YSconfdata.Devlist.Dev = append(YSconfdata.Devlist.Dev, *devdata)
			Chan := new(MoreToMoreConfigChannel)
			Chan.Id = ys.Id
			Chan.Index = ys.Channel
			YSconfdata.Channellist.Channel = append(YSconfdata.Channellist.Channel, *Chan)
		}
		CmeraId[strconv.Itoa(P)] = YSconfdata.Uuid
		YSConfigfname := ""
		//宇视生成xml配置文件
		ysfname := generateYSConfig(YSconfdata)
		if ysfname != "" {
			YSConfigfname = ysfname
			//启动宇视的程序
			go Runmain(YSConfigfname)

		} else {
			log.Println("宇视生成xml配置文件为空,YSConfigfname:", YSConfigfname)
			return
		}
		//time.Sleep(time.Minute * 1)
	}

	if len(hikITS) == 0 {
		log.Println("该网关设备没有海康ITS相机")
		return
	}

	//根据海康ITS的ip生成配置文件   赋值
	itsmap := make(map[string][]OneToMoreConfigChannel, len(hikITS))

	for _, its := range hikITS {
		//its 根据ip 端口号生成 配置文件
		Chan := new(OneToMoreConfigChannel)
		Chan.Id = its.Id
		Chan.Index = its.Channel

		log.Println("its.DevCompId|its.DevIp|its.Port|its.UserName|its.Password|ProcessPort:", its.DevCompId+"|"+its.DevIp+"|"+its.Port+"|"+its.UserName+"|"+its.Password+"|"+its.ProcessPort)

		if val, ok := itsmap[its.DevCompId+"|"+its.DevIp+"|"+its.Port+"|"+its.UserName+"|"+its.Password+"|"+its.ProcessPort]; ok == true {
			log.Println("海康iTS的列表值已经存在", val, "｜海康iTS的map列表值已经存在", itsmap)

			val = append(val, *Chan)
			itsmap[its.DevCompId+"|"+its.DevIp+"|"+its.Port+"|"+its.UserName+"|"+its.Password+"|"+its.ProcessPort] = val

			log.Println("海康iTS的列表新存在值：", *Chan, "｜海康iTS的新map列表存在值 ：", itsmap)

		} else {
			log.Println("海康iTS的列表值空值:", val, "|海康iTS的map列表值空值:", itsmap[its.DevCompId+"|"+its.DevIp+"|"+its.Port+"|"+its.UserName+"|"+its.Password+"|"+its.ProcessPort])

			//新的ITS进程的配置文件
			itschan := make([]OneToMoreConfigChannel, 0)
			itschan = append(itschan, *Chan)

			itsmap[its.DevCompId+"|"+its.DevIp+"|"+its.Port+"|"+its.UserName+"|"+its.Password+"|"+its.ProcessPort] = itschan

			log.Println("海康iTS的列表空值:", val, "|海康iTS的map列表第一个值:", itsmap[its.DevCompId+"|"+its.DevIp+"|"+its.Port+"|"+its.UserName+"|"+its.Password+"|"+its.ProcessPort])

		}
	}
	log.Println("海康iTS的列表值:", itsmap)

	for key, itsone := range itsmap {
		log.Println(key, itsone)
		P = P + 2
		//生成配置文件
		//ITS 多对多启动   OneToMoreConfig
		ITSconfdata := new(OneToMoreConfig)
		k := strings.Split(key, "|") //its.DevCompId+"|"+its.DevIp+"|"+its.Port+"|"+its.UserName+"|"+its.Password+"|"+its.ProcessPort
		ITSconfdata.DevCompId = k[0]
		ITSconfdata.Uuid = HIKITS + k[1] + k[2] + "+" + strconv.Itoa(P) //方便确定是哪一个进程发出的数据 我取品牌名称+进程端口号
		ITSconfdata.Udplistenport = P
		ITSconfdata.Udptxport = P - 1
		ITSconfdata.Devlist.Dev.UserName = k[3]
		ITSconfdata.Devlist.Dev.Password = k[4]
		ITSconfdata.Devlist.Dev.DevIp = k[1]
		ITSconfdata.Devlist.Dev.Port = k[2]
		ITSconfdata.Devlist.Dev.ITSPort = k[5]
		//ITSconfdata.Devlist.Dev.Id   =
		ITSconfdata.Channellist.Channel = itsone
		CmeraId[strconv.Itoa(P)] = ITSconfdata.Uuid

		//生成启动进程的配置文件
		ITSConfigfname := ""
		itsfname := generateITSConfig(ITSconfdata)
		if itsfname != "" {
			ITSConfigfname = itsfname
			//启动海康的程序
			go Runmain(ITSConfigfname)

		} else {
			log.Println("generateITSConfig生成启动进程的配置文件 文件名称为空，", itsfname)
			return
		}

	}
	log.Println("进程管理ProcessManagementService 执行完毕！")
}

//1、启动进程
func Runmain(RunmainConfigfname string) {
	log.Println("启动进程Configfname:", RunmainConfigfname)
	port := strings.Split(RunmainConfigfname, "+")
	//心跳port
	xtport := strings.Split(port[2], ".")
	//cmd := exec.Command("udpmain.exe绝对路径")
	//cmd := exec.Command("./snap/udpmain")
	//udpmainDir是main所在的目录
	var udpmainDir string
	udpmainDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	var udpmainexePath = filepath.Join(udpmainDir, "snap", "udpmain.exe")
	log.Println("udpmain.exe绝对路径:", udpmainexePath)

	cmd := exec.Command(udpmainexePath)

	path := make([]string, 0)
	//configxmlpath要给udpmain的配置文件路径
	var configxmlpath = filepath.Join(udpmainDir, "cameraConfig", RunmainConfigfname)
	//configxmlpath启动进程的配置文件的绝对路径 cameraConfig+ Configfname
	path = append(path, configxmlpath)

	cmd.Args = path
	//	log.Println("cmd.Args:", cmd.Args)

	//if runtime.GOOS == "windows" {
	//	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	//}

	var err = cmd.Start()
	if err != nil {
		log.Println("++++++ Execute Command failed. ++++++++++++++", err)
		return
	} else {
		log.Println("启动进程 ok！！！ you see see you")
	}

	//心跳port
	go Heartbeat(xtport[0])

	//进程在进行心跳交互处是阻塞的
	/*
	   main函数已经阻塞不退出了
	   这里已经测试了
	   不需要执行无限循环
	   for {
	   		log.Println("运行", path, "正常")
	   		time.Sleep(time.Hour * 1)
	   	}*/
}

//1、获取网关设备的token
func GetGatawayToken() (*dto.GetTokenRespXML, error) {
	Resp, err := GetToken(Deviceid)
	if err != nil {
		return nil, err
	}
	log.Println(Resp.Token, err)
	return Resp, nil
}

//2、根据token获取camera列表
func GetGatawayCameraList() (*dto.GetCameraList, error) {

	Resp, err := GetCameraList(Token)
	if err != nil {
		log.Println("GetCameraList error:", err)
		return nil, err
	}
	log.Println("根据token获取camera列表成功！！！")
	return Resp, nil
}

//与抓拍进程交互心跳，得知抓拍进程程序死活
func Heartbeat(port string) {
	p, _ := strconv.Atoi(port)
	//监控抓拍进程的心跳
XT:
	address := "127.0.0.1" + ":" + strconv.Itoa(p-1) //SERVER_PORT
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Println("监控抓拍进程心跳 net.ResolveUDPAddr 时 err:", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println(address, "监控抓拍进程的心跳 net.ListenUDP err:", err)
		time.Sleep(time.Second * 3)
		goto XT
	}

	log.Println("管理抓拍进程心跳,UDP监听地址address:", address)

	defer func() {
		_ = conn.Close()
	}()

	readTimeout := 60 * time.Second
	//心跳开始时间
	xtbeginsj := time.Now()
	buffer := make([]byte, 2048)
	for {

		err = conn.SetReadDeadline(time.Now().Add(readTimeout)) // timeout
		if err != nil {
			log.Println("setReadDeadline failed:", err)
			time.Sleep(time.Second * 3)
			continue
		}
		//反序列化udp数据
		h := new(dto.Heartbeatbasic)
		_, err := conn.Read(buffer) //读
		if err != nil {
			log.Println("Read udp failed:", err)
			now := time.Now()
			//upd时间差
			updsjcstr := utils.TimeDifference(xtbeginsj, now)
			//超时推出
			if strings.Contains(updsjcstr, "m") {
				log.Println("时间差:", updsjcstr, "心跳时间差大于60秒，需要重启程序", port)
				ycdata := new(dto.ExcprptStuQeq)
				ycdata.GatewayId = Deviceid                                  //1	gatewayId		网关id
				ycdata.CameraId = CmeraId[port]                              //2	cameraId		摄像机id
				ycdata.ReportTime = time.Now().Format("2006-01-02 15:04:05") //3	reportTime	2020-12-21 12:05:12	上报时间

				ycdata.CamStatus = -5 //4	0	摄像机状态 0 : 正常 -1: 连接摄像机网络失败； -2：摄像机注册/登陆失败； -3：摄像机异常(接口返回)； -4：24小时无数据；-5心跳时间差大于60秒，需要重启程序

				ycdata.CamStatusDes = "心跳时间差大于60秒，需要重启程序" //5	camStatusDes	正常	摄像机状态描述
				log.Println(ycdata)
			ExcprptStuUp:
				ycsberr := ExcprptStuUploadPostWithJson(ycdata)
				if ycsberr != nil {
					time.Sleep(time.Second * 5)
					goto ExcprptStuUp
				}
				//重启程序
				rsudperr := RestartUpdmain(port)
				if rsudperr != nil {
					log.Println("重启程序时，error:", rsudperr)
					return
				}
			}
			xtbeginsj = now
			time.Sleep(time.Second * 3)
			continue
		}

		herr := xml.Unmarshal(buffer, h)
		if herr != nil {
			log.Println(address, "UDP接收时,xml.Unmarshal失败！", herr) //这样解析是肯定OK的
		} else {
			//接收到数据
			log.Println(address, "接收到数据Type:", h.Type, h.Uuid)
			//port 6002
			Pid[port] = h.Pid
		}
		now := time.Now()
		//upd时间差
		updsjcstr := utils.TimeDifference(xtbeginsj, now)
		//超时推出
		if strings.Contains(updsjcstr, "m") {
			log.Println("时间差:", updsjcstr, "心跳时间差大于60秒，需要重启程序")
		}
		xtbeginsj = now
		heartbeatresp := new(dto.Heartbeat)
		//   1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
		switch h.Type {
		//心跳
		case 1:
			h := new(dto.Heartbeat)
			herr := xml.Unmarshal(buffer, h)
			if herr != nil {
				log.Println(address, "心跳xml.Unmarshal error:", herr)
				continue
			} else {
				log.Println(address, "1心跳Type:", h.Type, h)
				heartbeatresp.Uuid = h.Uuid
				heartbeatresp.Type = h.Type       //<type>   1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
				heartbeatresp.Version = h.Version //<version>抓拍程序版本号
				heartbeatresp.Time = h.Time       //<time>   字符串2020-11-12 12:12:12
				heartbeatresp.Seq = h.Seq         //<seq>   消息序号累加
			}

		case 2:
			//2、新数据通知
			h := new(dto.Heartbeat)
			herr := xml.Unmarshal(buffer, h)
			if herr != nil {
				log.Println(herr)
			} else {
				NewDataNotificationCount = NewDataNotificationCount + 1
				log.Println(address, "2新数据通知Type:", h.Type, h.Uuid, "NewDataNotificationCount:", NewDataNotificationCount)
				heartbeatresp.Uuid = h.Uuid
				heartbeatresp.Type = h.Type       //<type>    1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
				heartbeatresp.Version = h.Version //<version>        抓拍程序版本号
				heartbeatresp.Time = h.Time       //<time>     字符串2020-11-12 12:12:12
				heartbeatresp.Seq = h.Seq         //<seq>   消息序号累加
			}

		case 5:
			//5、 摄像机状态
			h := new(dto.StatusResult)
			herr := xml.Unmarshal(buffer, h)
			if herr != nil {
				log.Println(herr)
				log.Println(string(buffer))
			} else {
				log.Println(address, "5抓拍进程的摄像机状态数据:", h.LastCaptime, h)
				//heartbeatresp.Uuid = h.Uuid
				heartbeatresp.Type = h.Type       //<type>    1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
				heartbeatresp.Version = h.VerNum  //<version>   抓拍程序版本号
				heartbeatresp.Time = h.ReportTime //<time>     字符串2020-11-12 12:12:12
			}

			Camrpt(h)
		default:
			continue
		}

		heartbeatresp.Content = time.Now().Format("2006-01-02 15:04:05")
		//回复udp的消息
		resp, hresperr := xml.Marshal(heartbeatresp)
		if hresperr != nil {
			log.Println(address, hresperr)
		}
		//回复udp数据
		hferr := Heartbeatclient(port, resp)
		if hferr != nil {
			//log已经打印过了
			continue
		} else {
			log.Println(address, "回复成功，但不退出，继续udp交互")
		}
	}
}

func Heartbeatclient(port string, toWrite []byte) error {

	serverAddr := "127.0.0.1" + ":" + port
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		log.Println(serverAddr, "管理平台主动给抓拍进程返回数据,net.Dial执行时err:", err)
		return err
	}
	log.Println("管理平台主动给抓拍进程返回数据 UDP net.Dial serverAddr:", serverAddr)

	defer func() {
		_ = conn.Close()
	}()
	_, err = conn.Write(toWrite)
	if err != nil {
		log.Println("管理平台主动给抓拍进程心跳UDP err:", err)
		return err
	}
	return nil
}

//每日零点后抓拍的总和等清零
func HandleDayZeroTasks() {
	for {
		now := time.Now()                                                                    //获取当前时间，放到now里面，要给next用
		next := now.Add(time.Hour * 24)                                                      //通过now偏移24小时
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location()) //获取下一个20点的日期
		t := time.NewTimer(next.Sub(now))                                                    //计算当前时间到凌晨的时间间隔，设置一个定时器
		<-t.C                                                                                //阻塞等待第二天到来才执行
		CapZeroCnt = 0                                                                       //每日零点后抓拍的总和
		UploadRecordZeroCnt = 0                                                              //每日零点后上传的总和
		UploadImgZeroCnt = 0                                                                 // 每日零点后上传的总和
		UploadFailZeroCnt = 0                                                                //每日零点后上传的失败总和 会再次上传
		UploadFailImgZeroCnt = 0                                                             //每日零点后上传的失败总和

	}
}

func HandleDayTasks() {
	for {
		now := time.Now()                                                                    //获取当前时间，放到now里面，要给next用
		next := now.Add(time.Hour * 24)                                                      //通过now偏移24小时
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location()) //获取下一个20点的日期
		t := time.NewTimer(next.Sub(now))                                                    //计算当前时间到凌晨的时间间隔，设置一个定时器
		<-t.C                                                                                //阻塞等待第二天到来才执行

		sj := time.Now().Format("2006-01-02T15:04:05")
		content := sj + "\n上传到oss成功数量，OSSCount=" + strconv.Itoa(OSSCount)
		content = content + "\n前置机抓拍信息第一次上传抓拍结果成功ok,并接收成功 ResultOKCount:" + strconv.Itoa(ResultOKCount)
		content = content + "\n第一次上传抓拍结果xml文件到云平台成功，进程抓拍结果xml移动到parsed成功,Parsed:" + strconv.Itoa(Parsed)
		content = content + "\n再次上传的抓拍结果成功,AgainCount:" + strconv.Itoa(AgainCount)
		content = content + "\n前置机抓拍信息上传接口 ok ResultCount:" + strconv.Itoa(ResultCount)
		content = content + "\n新版信路威上传到oss3图都成功，NewOSSCount:" + strconv.Itoa(NewOSSCount)

		StatisticalFile(content)
		NewOSSCount = 0
		OSSCount = 0
		ResultCount = 0
		AgainCount = 0
		ResultOKCount = 0
		//	log.Println("执行重置OSS上传数量与抓拍结果上传数量OSSCount, ResultCount,AgainCount ,ResultOKCount：", OSSCount, ResultCount, AgainCount, ResultOKCount, time.Now().Format("2006-01-02T15:04:05"))
		content = "\n重新开始计数时间:" + time.Now().Format("2006-01-02T15:04:05")
		StatisticalFile(content)
		//删除前几天日期文件夹中为空的文件夹
		log.Println("删除前几天日期文件夹中为空的文件夹")
		//2、处理文件
		dir, _ := os.Getwd()
		var DelsnapimagespathDir = filepath.Join(dir, "snap", "images")
		log.Println("/snap/images/绝对路径:", DelsnapimagespathDir+"/") //+"/"

		DirList, err := ioutil.ReadDir(DelsnapimagespathDir + "/") //也可以不加
		if err != nil {
			log.Println("扫描/snap/images/文件夹 读取文件信息 error:", err)
			time.Sleep(time.Second * 3)
			continue
		}

		if len(DirList) == 1 {
			log.Println("扫描该/snap/images/文件夹下可能没有文件夹") //有隐藏文件

		} else {
			if len(DirList) == 0 {
				log.Println("扫描该/snap/images/文件夹下没有文件夹")
				continue
			}
		}

		log.Println("扫描该/snap/images/文件夹下有文件的数量:", len(DirList))
		for i := range DirList {
			//判断文件的结尾名
			log.Println("DirList[i].Name():", DirList[i].Name())
			if DirList[i].IsDir() {

				log.Println("日期文件夹的绝对目录:", DelsnapimagespathDir+"/"+DirList[i].Name())
				fileList, err := ioutil.ReadDir(DelsnapimagespathDir + "/" + DirList[i].Name()) //可已不加"/"
				if err != nil {
					log.Println("扫描/snap/images/下文件夹读取文件信息 error:", err, DirList[i].Name())
					continue
				}

				if len(fileList) == 0 {
					log.Println("执行 扫描 该/snap/images/ 下文件夹的文件夹下没有需要解析的xml文件,是空文件夹，", DirList[i].Name())

					oldday := utils.OldData(7)
					for _, day := range oldday {
						if DirList[i].Name() == day {
							//删除空文件夹
							log.Println("删除空文件夹path:", DelsnapimagespathDir+"/"+DirList[i].Name())
							rmverr := os.RemoveAll(DelsnapimagespathDir + "/" + DirList[i].Name())
							if rmverr != nil {
								log.Println("删除空文件夹失败:", rmverr)

							} else {
								log.Println("删除空文件夹:", DirList[i].Name(), "ok!")
							}
						}
					}
				} else {
					log.Println("文件夹", DirList[i].Name(), "存在文件，文件名称：", fileList[0].Name())
					continue
				}
			}
		}
		log.Println("处理可能有要删除的空文件夹OK")
	}
}

func ChepZH(ys string) string {
	switch ys {

	case "黑":
		return "1"

	case "白":
		return "2"

	case "蓝":
		return "3"

	case "黄":
		return "4"

	case "绿":
		return "5"

	case "黄绿":
		return "6"

	default:
		log.Println("ys", ys)
		return "0"

	}

}
