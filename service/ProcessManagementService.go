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

var Deviceid string //网关设备id Token
//var IpAddress string

var Token string

var BacketName string
var ObjectPrefix string

const SERVER_PORT = "5000"
const SERVER_RECV_LEN = 10

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
)

// 进程管理
func ProcessManagementService(ch chan int) {

	PorT := 6000 //固定端口
	var P int
	P = PorT
	//1、获取网关设备的token
CQ:
	resp, getTokenerr := GetGatawayToken()
	if getTokenerr != nil {
		log.Println("获取网关设备的token 失败") //getTokenerr 已打印
		time.Sleep(time.Minute * 1)
		goto CQ
	}

	//全局 Token  BacketName 	ObjectPrefix 赋值
	if resp != nil {
		Token = resp.Token
		BacketName = resp.Oss.BacketName
		ObjectPrefix = resp.Oss.ObjectPrefix
		utils.BacketName = BacketName
		log.Println("Token:", Token)
		log.Println("BacketName:", BacketName)

		log.Println("ObjectPrefix:", ObjectPrefix)

	}

	//2、根据token获取camera列表
	CameraList, listerr := GetGatawayCameraList()
	if listerr != nil {
		log.Println("获取相机列表错误", listerr)
		return
	}

	log.Println(" 相机列表数据的len（） ：", len(CameraList.Data))

	uniview := make([]dto.CameraListData, 0) // 宇视的列表

	hikITS := make([]dto.CameraListData, 0) //ITS列表

	for i, cmera := range CameraList.Data {
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

			confdata.Channellist.Channel.Id = cmera.Id         //相机id
			confdata.Channellist.Channel.Index = cmera.Channel //通道号

			//一对一生成配置文件
			fname := generateConfigToOne(confdata)
			if fname != "" {
				Configfname = fname
			}

			//2、进程启动
			//传 一个配置文件的绝对路径 全局唯一
			if err := Runmain(Configfname); err != nil {
				log.Println("需要重启")
			} else {
				log.Println("一对一的进程已启动ok", i+1)
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

	if len(hikITS) == 0 && len(uniview) == 0 {
		log.Println("++++++++++++++++++++++++++++++该网关设备没有海康ITS相机和宇视相机")
		return
	}

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
	YSConfigfname := ""
	//宇视生成xml配置文件
	ysfname := generateYSConfig(YSconfdata)
	if ysfname != "" {
		YSConfigfname = ysfname
	}
	time.Sleep(time.Minute * 1)

	//启动宇视的程序
	if err := Runmain(YSConfigfname); err != nil {
		log.Println("宇视需要重启")
	} else {
		log.Println("启动宇视的程序ok")
	}

	if len(hikITS) == 0 {
		log.Println("++++++++++++++++++++++++++++++该网关设备没有海康ITS相机")
		return
	}

	//根据海康ITS的ip生成配置文件   赋值
	itsmap := make(map[string][]OneToMoreConfigChannel, len(hikITS))

	for _, its := range hikITS {
		//its 根据ip 端口号生成 配置文件
		Chan := new(OneToMoreConfigChannel)
		Chan.Id = its.Id
		Chan.Index = its.Channel

		log.Println("its.DevCompId++its.DevIp++its.Port++its.UserName++its.Password:", its.DevCompId+"|"+its.DevIp+"|"+its.Port+"|"+its.UserName+"|"+its.Password)

		if val, ok := itsmap[its.DevCompId+"+"+its.DevIp+"+"+its.Port+"+"+its.UserName+"+"+its.Password]; ok == true {
			log.Println("海康iTS的列表值已经存在", val, "｜海康iTS的map列表值已经存在", itsmap)

			val = append(val, *Chan)
			itsmap[its.DevCompId+"+"+its.DevIp+"+"+its.Port+"+"+its.UserName+"+"+its.Password] = val

			log.Println("海康iTS的列表新存在值：", *Chan, "｜海康iTS的新map列表存在值 ：", itsmap)

		} else {
			log.Println("海康iTS的列表值空值:", val, "+海康iTS的map列表值空值:", itsmap[its.DevCompId+"+"+its.DevIp+"+"+its.Port+"+"+its.UserName+"+"+its.Password])

			//新的ITS进程的配置文件
			itschan := make([]OneToMoreConfigChannel, 0)
			itschan = append(itschan, *Chan)

			itsmap[its.DevCompId+"+"+its.DevIp+"+"+its.Port+"+"+its.UserName+"+"+its.Password] = itschan

			log.Println("海康iTS的列表空值:", val, "+海康iTS的map列表第一个值:", itsmap[its.DevCompId+"+"+its.DevIp+"+"+its.Port+"+"+its.UserName+"+"+its.Password])

		}
	}
	log.Println("海康iTS的列表值:", itsmap)

	for key, itsone := range itsmap {
		log.Println(key, itsone)
		P = P + 2
		//生成配置文件
		//ITS 多对多启动   OneToMoreConfig
		ITSconfdata := new(OneToMoreConfig)
		k := strings.Split(key, "+") //its.DevCompId+"+"+its.DevIp+"+"+its.Port+"+"+its.UserName+"+"+its.Password
		ITSconfdata.DevCompId = k[0]
		ITSconfdata.Uuid = HIKITS + k[1] + k[2] + "+" + strconv.Itoa(P) //方便确定是哪一个进程发出的数据 我取品牌名称+进程端口号
		ITSconfdata.Udplistenport = P
		ITSconfdata.Udptxport = P - 1
		ITSconfdata.Devlist.Dev.UserName = k[3]
		ITSconfdata.Devlist.Dev.Password = k[4]
		ITSconfdata.Devlist.Dev.DevIp = k[1]
		ITSconfdata.Devlist.Dev.Port = k[2]
		//ITSconfdata.Devlist.Dev.Id   =
		ITSconfdata.Channellist.Channel = itsone

		//生成启动进程的配置文件
		ITSConfigfname := ""
		itsfname := generateITSConfig(ITSconfdata)
		if itsfname != "" {
			ITSConfigfname = itsfname
		}
		//启动海康的程序
		if err := Runmain(ITSConfigfname); err != nil {
			log.Println("海康ITS需要重启")

		} else {
			log.Println("启动海康的程序ok")
		}

	}
	ch <- 1
}

//1、启动进程
func Runmain(Configfname string) error {

	log.Println("Configfname:", Configfname)
	//与抓拍进程交互心跳 [ ]
	port := strings.Split(Configfname, "+")

	//心跳port
	xtpt := strings.Split(port[2], ".")

	//cmd := exec.Command("capture.exe绝对路径")
	//cmd := exec.Command("./snap/udpmain")
	var additionalBilldataDir string
	additionalBilldataDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	var billoutputDir = filepath.Join(additionalBilldataDir, "snap", "udpmain.exe")
	log.Println("capture.exe绝对路径:", billoutputDir)

	cmd := exec.Command(billoutputDir)

	path := make([]string, 0)
	//	path = append(path, "ConfigPath:")
	var configxmlpath = filepath.Join(additionalBilldataDir, "cameraConfig", Configfname)
	//绝对路径
	path = append(path, configxmlpath)

	cmd.Args = path
	log.Println("cmd.Args:", cmd.Args)

	//if runtime.GOOS == "windows" {
	//	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	//}

	var err = cmd.Start()
	if err != nil {
		log.Println("++++++ Execute Command failed. ++++++++++++++", err)
		return err
	} else {
		log.Println("启动进程 ok！！！ you see see you")
	}

	//心跳port
	go Heartbeat(xtpt[0])

	return nil
}

//1、获取网关设备的token
func GetGatawayToken() (*dto.GetTokenRespXML, error) {

	Resp, err := GetToken(Deviceid)
	if err != nil {
		log.Println("GetToken error:", err)
		return nil, err
	}
	log.Println(Resp.Token, err)
	return Resp, nil
}

//2、根据token获取camera列表
func GetGatawayCameraList() (*dto.GetCameraList, error) {

	Resp, err := GetCameraList(Token)
	if err != nil {
		log.Println("GetToken error:", err)
		return nil, err
	}
	log.Println("根据token获取camera列表成功！！！")
	return Resp, nil
}

//上传文件  开线程读取xml文件 上传图片到oss  上传抓拍结果到车牌识别云端服务器
func UploadFile() {

	tiker := time.NewTicker(time.Second * 5) //每5秒执行一下

DirListP:

	for {
		//删除前几天日期文件夹中为空的文件夹
		log.Println("执行删除前几天日期文件夹中为空的文件夹")
		//2、处理文件
		//扫描 captureXml 文件夹 读取文件信息

		pwd := "./snap/images/"
		DirList, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Println("扫描 /snap/images/文件夹 读取文件信息 error:", err)
			time.Sleep(time.Second * 5)
			break DirListP
		}

		log.Println("执行 扫描 该/snap/images/文件夹下有文件的数量 ：", len(DirList))

		if len(DirList) == 1 {
			log.Println("执行 扫描 该/snap/images/ 文件夹下可能没有文件夹") //有隐藏文件

		} else {
			if len(DirList) == 0 {
				log.Println("执行 扫描 该/snap/images/ 文件夹下没有文件夹")
				break
			}
		}

		for i := range DirList {
			//判断文件的结尾名
			if IsDir("./snap/images/" + DirList[i].Name()) {

				fileList, err := ioutil.ReadDir("./snap/images/" + DirList[i].Name())
				if err != nil {
					log.Println("扫描 /snap/images/下文件夹 读取文件信息 error:", err, DirList[i].Name())
					break
				}

				if len(fileList) == 0 {
					log.Println("执行 扫描 该/snap/images/ 下文件夹的文件夹下没有需要解析的xml文件", DirList[i].Name())

					oldday := utils.OldData(7)
					for _, day := range oldday {
						if DirList[i].Name() == day {
							//删除空文件夹
							rmverr := os.RemoveAll("./snap/images/" + DirList[i].Name() + "/")
							if rmverr != nil {
								log.Println("删除空文件夹失败:", rmverr)
							} else {
								log.Println("删除空文件夹", DirList[i].Name(), "ok")
							}

						}

					}

				}

			}

		}
		log.Println("处理可能有要删除的空文件夹OK")

		//上传图片以及抓拍结果到车牌识别云端服务器
		HandleFile()
		log.Println(<-tiker.C, "执行 上传图片以及抓拍结果到车牌识别云端服务器  ")
	}

}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func HandleFile() {
	//定期检查抓拍文件夹文件夹 captureXml
	//tiker := time.NewTicker(time.Minute * 1)
	//tiker := time.NewTicker(time.Second * 30)
	//for {
	log.Println("执行处理xml数据包解析以及oss上传以及抓拍结果上传")
	//2、处理文件
	//扫描 captureXml 文件夹 读取文件信息
	pwd := "./snap/xml/"
	fileList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Println("扫描 captureXml 文件夹 读取文件信息 error:", err)
		return
	}
	log.Println("执行 扫描 该captureXml文件夹下有文件的数量 ：", len(fileList))
	if len(fileList) == 1 {
		log.Println("执行 扫描 该captureXml 文件夹下可能没有需要解析的xml文件") //有隐藏文件

	} else {
		if len(fileList) == 0 {
			log.Println("执行 扫描 该captureXml 文件夹下没有需要解析的xml文件")
			return
		}
	}

	for i := range fileList {
		//判断文件的结尾名
		if strings.HasSuffix(fileList[i].Name(), ".xml") {
			log.Println("执行 扫描 该captureXml文件夹下需要解析的xml文件名字为:", fileList[i].Name())
			content, err := ioutil.ReadFile("./snap/xml/" + fileList[i].Name())
			if err != nil {
				log.Println("执行  读文件位置错误信息：", err)
				continue
			}

			//将xml文件转换为对象
			var result dto.CaptureDateXML
			uerr := xml.Unmarshal(content, &result)
			if uerr != nil {
				log.Println("执行 扫描 该captureXml文件夹下需要解析的xml文件内容时，错误信息为：", uerr)
				continue
			}

			log.Println("获取抓拍结果，result:", result.VehicleImgPath)

			//把图片上传到oss上
			//result.VehicleImgPath C:\Users\Administrator\Desktop\HSJDEBUG\images\20201124\sxjgl_yzjtd_320200_G2_K1071_2_0_004_20201124143417_000031.jpg 图片路径
			//c := strings.Split(result.VehicleImgPath, ":")
			//str2 := strings.Replace(c[1], "\\", "/", -1) //linux

			//log.Println(str2)
			//strfname := strings.Split(str2, "/")
			strfname := strings.Split(result.VehicleImgPath, "\\") //windows
			//上传到oss                    日期文件夹     图片名称               前缀"/jiangsu/suhuaiyangs"
			log.Println("上传到oss     图片地址     图片名称   前缀", result.VehicleImgPath, strfname[7], ObjectPrefix)
			code, scsj, ossDZ := utils.QingStorUpload(result.VehicleImgPath, strfname[7], ObjectPrefix)

			if code == utils.UPloadOK {
				log.Println("上传到oss   成功，开始返回抓拍结果给云平台")
				//删除本地图片result.VehicleImgPath
				//utils.DelFile("./images/" + strfname[6] + "/" + strfname[7])
				utils.DelFile(result.VehicleImgPath)

				//生产xml返回给云平台 [暂时上传到模拟云平台]
				uploaderr := GwCaptureInforUpload(&result, scsj, ossDZ)
				if uploaderr != nil {
					continue
				} else {
					//删除抓拍xml文件
					source := "./snap/xml/" + fileList[i].Name()
					d := "./snap/xml/parsed/" + fileList[i].Name()
					mverr := utils.MoveFile(source, d)
					if mverr != nil {
						log.Println(mverr)
						continue
					}
				}

			} else {
				log.Println("上传oss失败")
				//上传oss失败
				continue
			}

		}

	}
	//log.Println("执行处理xml数据包", (<-tiker.C).Format("2006-01-02 15:04:05"))
	//}
}

func GwCaptureInforUpload(Result *dto.CaptureDateXML, scsj int64, ossDZ string) error {
	//判断哪一种品牌相机
	//Result.
	var ba []byte

	if Result.AppedInfo.AxleDist != "" {
		data := new(dto.TBXJDateXML)
		//抓拍结果的赋值
		data.Token = Token //抓拍结果上传

		data.LprInfo.PassId = Result.PassId     //    string   `xml:"passId"`         // 过车编号
		data.LprInfo.CamId = Result.CamId       //    string   `xml:"camId"`          //camId>    摄像机编号
		data.LprInfo.DeviceId = Deviceid        //    string   `xml:"deviceId"`       //deviceId>前置机编号
		data.LprInfo.PassTime = Result.PassTime //    string   `xml:"passTime"`       //passTime>     过车时间
		data.LprInfo.VehicleImgPath = ossDZ     //    string   `xml:"vehicleImgPath"` //vehicleImgPath>  "oss地址"   过车图片地址
		data.LprInfo.PlateImgPath = ""          //无 string   `xml:"plateImgPath"`   //<plateImgPath/>     车牌图片地址【无】
		data.LprInfo.BucketId = BacketName      //   string   `xml:"bucketId"`       //bucketId>   bucket编号
		data.LprInfo.ImageType = 0              //   int      `xml:"imageType"`      //	imageType> 图片类型
		data.LprInfo.UploadStamp = scsj         //   int64    `xml:"uploadStamp"`    //	uploadStamp> 上传时间
		data.LprInfo.Stationid = ""             //   string   `xml:"stationid"`      //	stationid>站点编号
		data.LprInfo.LaneType = 0               //   int      `xml:"laneType"`       //	laneType> 出入口类型 0:入口；1：出口

		data.LpaResult.PassId = Result.PassId         //passId>过车编号
		data.LpaResult.EngineType = ""                //`xml:"engineType"`      //engineType>   引擎类型
		data.LpaResult.EngineId = ""                  //`xml:"engineId"`        //engineId>     引擎编号
		data.LpaResult.PlateNo = Result.PlateNo       //`xml:"plateNo"`         //plateNo>     车牌编号
		data.LpaResult.PlateColor = Result.PlateColor // `xml:"plateColor"`      //plateColor>     车牌颜色
		data.LpaResult.ComputeInterval = 0            //int64 `xml:"computeInterval"` //computeInterval>  计算时间
		data.LpaResult.VehicleColor = ""              //`xml:"vehicleColor"`    //vehicleColor>       车辆颜色
		data.LpaResult.VehicleType = ""               //`xml:"vehicleType"`     //vehicleType>       车辆类型
		data.LpaResult.VehicleBrand = ""              //`xml:"vehicleBrand"`    //vehicleBrand>       车辆品牌
		data.LpaResult.VehicleYear = 0                //int`xml:"vehicleYear"`     //vehicleYear>     车辆年份

		data.LpaResult.LprFrameEntity.PlateLeft = 0   // int      `xml:"plateLeft"`   //plateLeft>        车牌左坐标
		data.LpaResult.LprFrameEntity.PlateTop = 0    //  int      `xml:"plateTop"`    //plateTop>        车牌上坐标
		data.LpaResult.LprFrameEntity.PlateRight = 0  // int      `xml:"plateRight"`  //plateRight>        车牌右坐标
		data.LpaResult.LprFrameEntity.PlateBottom = 0 // int      `xml:"plateBottom"` //plateBottom>     车牌下坐标

		//data.VehicleInfo.SideImgPath  =   // string   `xml:"sideImgPath"` //sideImgPath> 侧面图片地址
		//	data.VehicleInfo.TailImgPath    = //  string   `xml:"tailImgPath"` //tailImgPath> 车尾图片地址
		data.VehicleInfo.CarType = Result.AppedInfo.CarType //  string   //CarType>  车辆型号
		AxleNum, _ := strconv.Atoi(Result.AppedInfo.AxleNum)
		data.VehicleInfo.AxleNum = AxleNum                                  //  int      //AxleNum>  轴数
		data.VehicleInfo.AxleType = Result.AppedInfo.AxleType               // string   //AxleType>  轴型
		data.VehicleInfo.WheelNumber = Result.AppedInfo.WheelNumber         // string   //WheelNumber> 轮胎数量
		data.VehicleInfo.AxleDist = Result.AppedInfo.AxleDist               // string   //AxleDist>  轴距
		data.VehicleInfo.CarLengthMeter = Result.AppedInfo.CarLengthMeter   //string   //CarLengthMeter> 车长
		data.VehicleInfo.VideoScaleSpeed = Result.AppedInfo.VideoScaleSpeed //string   //VideoScaleSpeed> 车速
		data.VehicleInfo.WXPCharIndex = Result.AppedInfo.WXPCharIndex       // string   //WXPCharIndex>  危险品标识
		data.VehicleInfo.ZXType = Result.AppedInfo.ZXType                   // string   //ZXType> 专项作业车标识

		//MarshalIndent 有缩进 xml.Marshal ：无缩进

		ba, _ = xml.MarshalIndent(data, "  ", "  ")
		log.Println("+++++++++", string(ba))

	} else {
		data := new(dto.DateXML)
		//抓拍结果的赋值
		data.Token = Token //抓拍结果上传

		data.LprInfo.PassId = Result.PassId     //    string   `xml:"passId"`         // 过车编号
		data.LprInfo.CamId = Result.CamId       //      string   `xml:"camId"`          //camId>    摄像机编号
		data.LprInfo.DeviceId = Deviceid        //      string   `xml:"deviceId"`       //deviceId>前置机编号
		data.LprInfo.PassTime = Result.PassTime //    string   `xml:"passTime"`       //passTime>     过车时间
		data.LprInfo.VehicleImgPath = ossDZ     //   string   `xml:"vehicleImgPath"` //vehicleImgPath>  "oss地址"   过车图片地址
		data.LprInfo.PlateImgPath = ""          //     string   `xml:"plateImgPath"`   //<plateImgPath/>     车牌图片地址【无】
		data.LprInfo.BucketId = BacketName      //     string   `xml:"bucketId"`       //bucketId>   bucket编号
		data.LprInfo.ImageType = 0              //     int      `xml:"imageType"`      //	imageType> 图片类型
		data.LprInfo.UploadStamp = scsj         //     int64    `xml:"uploadStamp"`    //	uploadStamp> 上传时间
		data.LprInfo.Stationid = ""             //     string   `xml:"stationid"`      //	stationid>站点编号
		data.LprInfo.LaneType = 0               //     int      `xml:"laneType"`       //	laneType> 出入口类型 0:入口；1：出口

		data.LpaResult.PassId = Result.PassId         //passId>     过车编号
		data.LpaResult.EngineType = ""                //`xml:"engineType"`      //engineType>   引擎类型
		data.LpaResult.EngineId = ""                  //`xml:"engineId"`        //engineId>     引擎编号
		data.LpaResult.PlateNo = Result.PlateNo       //`xml:"plateNo"`         //plateNo>     车牌编号
		data.LpaResult.PlateColor = Result.PlateColor // `xml:"plateColor"`      //plateColor>     车牌颜色
		data.LpaResult.ComputeInterval = 0            //int64 `xml:"computeInterval"` //computeInterval>  计算时间
		data.LpaResult.VehicleColor = ""              //`xml:"vehicleColor"`    //vehicleColor>       车辆颜色
		data.LpaResult.VehicleType = ""               //`xml:"vehicleType"`     //vehicleType>       车辆类型
		data.LpaResult.VehicleBrand = ""              //`xml:"vehicleBrand"`    //vehicleBrand>       车辆品牌
		data.LpaResult.VehicleYear = 0                //int`xml:"vehicleYear"`     //vehicleYear>     车辆年份

		data.LpaResult.LprFrameEntity.PlateLeft = 0   // int      `xml:"plateLeft"`   //plateLeft>        车牌左坐标
		data.LpaResult.LprFrameEntity.PlateTop = 0    //  int      `xml:"plateTop"`    //plateTop>        车牌上坐标
		data.LpaResult.LprFrameEntity.PlateRight = 0  // int      `xml:"plateRight"`  //plateRight>        车牌右坐标
		data.LpaResult.LprFrameEntity.PlateBottom = 0 // int      `xml:"plateBottom"` //plateBottom>     车牌下坐标
		//MarshalIndent 有缩进 xml.Marshal ：无缩进

		ba, _ = xml.MarshalIndent(data, "  ", "  ")
		log.Println("+++++++++", string(ba))
	}

	log.Println("前置机抓拍信息上传接口 Address:", GwCaptureInformationUploadIpAddress)
	result, err := GwCaptureInformationUploadPostWithXML(&ba)
	if err != nil {
		return err
	}

	if (*result).Code == 0 {
		log.Println("上传抓拍结果成功")
		return nil
	} else {
		log.Println("上传抓拍结果失败")
		return err
	}
}

//与抓拍进程交互心跳，得知抓拍进程程序死活
func Heartbeat(port string) {

	//监控抓拍进程的心跳
XT:
	ip := strings.Split("127.0.0.1", ":")
	address := ip[0] + ":" + port //SERVER_PORT
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Println("监控抓拍进程心跳 net.ResolveUDPAddr 时 err:", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println("监控抓拍进程的心跳 net.ListenUDP err:", err)
		time.Sleep(time.Second * 3)
		goto XT
	}

	log.Println("管理平台 UDP监听 address:", address)

	defer func() {
		_ = conn.Close()
	}()
	data := make([]byte, 4096)
	for {
		//获取数据
		// Here must use make and give the lenth of buffer

		//返回一个UDPAddr        ReadFromUDP从c读取一个UDP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址。
		//ReadFromUDP方法会在超过一个固定的时间点之后超时，并返回一个错误。
		_, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println(err)
			continue
		}

		//反序列化udp数据
		h := new(dto.Heartbeatbasic)
		herr := xml.Unmarshal(data, h)
		if herr != nil {
			log.Println(herr)
		} else {
			log.Println("h.Type:", h.Type)
		}

		heartbeatresp := new(dto.Heartbeat)
		//   1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
		switch h.Type {

		case 1:
			//   1、心跳
			h := new(dto.Heartbeat)
			herr := xml.Unmarshal(data, h)
			if herr != nil {
				log.Println(herr)
			} else {
				log.Println("h.Type:", h.Type, h)
				heartbeatresp.Uuid = h.Uuid
				heartbeatresp.Type = h.Type       //<type>    1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
				heartbeatresp.Version = h.Version //<version>        抓拍程序版本号
				heartbeatresp.Time = h.Time       //<time>     字符串2020-11-12 12:12:12
				heartbeatresp.Seq = h.Seq         //<seq>   消息序号累加
			}

		case 2:
			//2、新数据通知
			h := new(dto.Heartbeat)
			herr := xml.Unmarshal(data, h)
			if herr != nil {
				log.Println(herr)
			} else {
				log.Println("h.Type:", h.Type, h)
				heartbeatresp.Uuid = h.Uuid
				heartbeatresp.Type = h.Type       //<type>    1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
				heartbeatresp.Version = h.Version //<version>        抓拍程序版本号
				heartbeatresp.Time = h.Time       //<time>     字符串2020-11-12 12:12:12
				heartbeatresp.Seq = h.Seq         //<seq>   消息序号累加
			}

		case 3:
			//3、 日志
			h := new(dto.HeartbeatLog)
			herr := xml.Unmarshal(data, h)
			if herr != nil {
				log.Println(herr)
			} else {
				log.Println("抓拍进程的日志：", h)
				heartbeatresp.Uuid = h.Uuid
				heartbeatresp.Type = h.Type       //<type>    1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
				heartbeatresp.Version = h.Version //<version>        抓拍程序版本号
				heartbeatresp.Time = h.Time       //<time>     字符串2020-11-12 12:12:12
				heartbeatresp.Seq = h.Seq         //<seq>   消息序号累加
			}

		case 4:
			// 4、采集进程被动关闭命令
			h := new(dto.Heartbeat)
			herr := xml.Unmarshal(data, h)
			if herr != nil {
				log.Println(herr)
			} else {
				log.Println("h.Type:", h.Type, h)

				heartbeatresp.Uuid = h.Uuid
				heartbeatresp.Type = h.Type       //<type>    1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
				heartbeatresp.Version = h.Version //<version>        抓拍程序版本号
				heartbeatresp.Time = h.Time       //<time>     字符串2020-11-12 12:12:12
				heartbeatresp.Seq = h.Seq         //<seq>   消息序号累加
			}
		}

		heartbeatresp.Content = time.Now().Format("2006-01-02 15:04:05")
		resp, hresperr := xml.Marshal(heartbeatresp)
		if hresperr != nil {
			log.Println(hresperr)
		} else {
			log.Println("xml.Marshal ok! 管理平台收到抓拍进程的信息 h.Type:", h.Type)
		}
		//回复udp数据
		_, err = conn.WriteToUDP(resp, rAddr)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("管理平台 Send:", heartbeatresp)
	}
}

//【不做了】
func HeartbeatClient(port string) {
	tiker := time.NewTicker(time.Second * 10) //每15秒执行一下
	for {
		log.Println(utils.DateTimeFormat(<-tiker.C), "管理平台要发送心跳给抓拍进程++++++++++++")

		Heartbeatclient(port)

	}
}

//【不做了】
func Heartbeatclient(port string) {

	serverAddr := "127.0.0.1" + ":" + port
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		log.Println(serverAddr, "管理平台 主动给抓拍进程心跳,net.Dial执行时", "err:", err)
		time.Sleep(time.Second * 10)
		return
	}
	log.Println("管理平台 主动给抓拍进程心跳 UDP net.Dial serverAddr:", serverAddr)

	defer func() {
		_ = conn.Close()
	}()

	var n int
	var toWrite string
	toWrite = serverAddr + "管理平台细心问候：你启动了么，是否活着呀!"

	n, err = conn.Write([]byte(toWrite))
	if err != nil {
		log.Println("err", err)
		return
	}

	log.Println("Write:", toWrite, "n:", n)

	msg := make([]byte, 32)
	n, err = conn.Read(msg)
	if err != nil {
		log.Println("err:", err)
		return
	}

	log.Println("抓拍进程给的响应，Response:", string(msg), "n:", n)

}

func checkError(err error) {
	if err != nil {
		log.Println(err)
		//return
		os.Exit(1)
	}
}
