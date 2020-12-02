package service

import (
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"vehicleLicensePlateRecognitionGateway/dto"
	"vehicleLicensePlateRecognitionGateway/utils"
)

var Deviceid string //网关设备id Token
var IpAddress string

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
	SignalwayNew string = "SignalwayNew" // 信路威车型
	GDPort       int    = 5000           //固定 进程向我拨号的的端口
)

// 进程管理
func ProcessManagementService() {
	PorT := 6000 //固定端口
	//1、获取网关设备的token
CQ:
	resp, getTokenerr := GetGatawayToken()
	if getTokenerr != nil {
		log.Println("获取网关设备的token 失败") //getTokenerr 已打印

		time.Sleep(time.Minute * 1)
		goto CQ
		//cqresp, cqgetTokenerr := 	GetGatawayToken()
		//if cqgetTokenerr!=nil{
		//	log.Println("获取网关设备的token 重启失败")
		//	return
		//}
		//
		//if cqresp != nil {
		//	Token = resp.Token
		//	BacketName = resp.Oss.BacketName
		//	ObjectPrefix = resp.Oss.ObjectPrefix
		//}

	}
	//全局 Token  BacketName 	ObjectPrefix 赋值
	if resp != nil {
		Token = resp.Token
		BacketName = resp.Oss.BacketName
		ObjectPrefix = resp.Oss.ObjectPrefix
	}

	//2、根据token获取camera列表
	CameraList, listerr := GetGatawayCameraList()
	if listerr != nil {
		log.Println("获取相机列表错误", listerr)
		return
	}
	log.Println(" 相机列表数据 ：", CameraList)

	for i, cmera := range CameraList.Data {

		conflx := ""
		if cmera.DevCompId == UNIVIEW || cmera.DevCompId == HIKITS {
			conflx = "one2many"

		} else {
			conflx = "one2ont"
		}
		ConfigPath := ""
		//1、生成进程配置文件
		//ConfigPath:="abc"
		switch conflx {
		case "one2ont":
			PorT = PorT + 2
			confdata := new(OneToOneConfig)

			confdata.DevCompId = cmera.DevCompId //品牌名称
			strporrt := strconv.Itoa(PorT)
			confdata.Uuid = cmera.Id + "+" + strporrt //方便确定是哪一个进程发出的数据 我取相机id+进程端口号
			confdata.Udplistenport = PorT             //我向进程拨号的端口号
			confdata.Udptxport = PorT - 1             //固定 进程向我拨号的的端口
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
				ConfigPath = fname
			}

		case "one2many":
			log.Println("one2many,相机品牌是：", cmera.DevCompId)
			//HIKITS
			//UNIVIEW

			//	generateConfig()

			//case "many2many":
			//
			//	generateConfig()

		}

		//2、进程启动
		//A:
		//传 一个配置文件的绝对路径 全局唯一
		if err := Runmain(ConfigPath); err != nil {
			log.Println("重启")

			//var a int
			////2、进程重启
			//Rerr := Runmain(ConfigPath)
			//a = a + 1
			//log.Println("重启数a：", a)
			//
			//if Rerr != nil {
			//	log.Println("重启 error!", Rerr)
			//	goto A
			//}
		}

		log.Println("进程已启动", i+1)
		continue
	}

}

//1、启动进程
func Runmain(ConfigPath string) error {

	//cmd := exec.Command("../grpcSimulator/grpc_main", "test_file")
	//命令行参数是配置文件的绝对路径 +文件名【全局唯一】
	//cmd := exec.Command("./grpcSimulator/grpc_main",  "-configpath",  ConfigPath)//模拟器方式一
	//cmd := exec.Command("./grpcSimulator/grpc_main","" ,ConfigPath) //进程程序方式 模拟器方式二
	//buf, err := cmd.Output()
	//fmt.Printf("output: %s\n", string(buf))
	//if err != nil {
	//	fmt.Printf("err: %v\n", err)
	//}

	//与抓拍进程交互心跳 [cameraConfig/2020-12-02T11:35:03+sxjgl_shygs_321300_G2513_K101_415_3_2_1+6004.xml]
	port := strings.Split(ConfigPath, "+")
	xtpt := strings.Split(port[2], ".")

	//cmd := exec.Command("capture.exe绝对路径")
	cmd := exec.Command("./grpcSimulator/grpc_main")

	path := make([]string, 0)
	path = append(path, ConfigPath)

	cmd.Args = path
	log.Println("cmd.Args:", cmd.Args)

	//if runtime.GOOS == "windows" {
	//	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	//}

	var err = cmd.Start()
	if err != nil {
		log.Println("++++++ Execute Command failed. ++", err)
		return err
	}

	//心跳port
	go Heartbeat(xtpt[0])

	//执行Cmd中包含的命令，阻塞直到命令执行完成
	//Runerr := cmd.Run()
	//if Runerr != nil {
	//	log.Println("++++++ Execute Command failed. ++", "+++++++ Runerr:", Runerr.Error())
	//	return Runerr
	//}
	//log.Println("Execute Command finished.")

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
	//删除前几天日期文件夹中为空的文件夹

	//  上传图片以及抓拍结果到车牌识别云端服务器
	HandleFile()
}

func HandleFile() {
	//定期检查抓拍文件夹文件夹 captureXml
	//tiker := time.NewTicker(time.Minute * 1)
	tiker := time.NewTicker(time.Second * 30)
	for {
		log.Println("执行处理数据包")
		//2、处理文件
		//扫描 captureXml 文件夹 读取文件信息
		pwd := "./captureXml/"
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
				content, err := ioutil.ReadFile("./captureXml/" + fileList[i].Name())
				if err != nil {
					log.Println("执行  读文件位置错误信息：", err)
					return
				}

				//将xml文件转换为对象
				var result dto.CaptureDateXML
				uerr := xml.Unmarshal(content, &result)
				if uerr != nil {
					log.Println("执行 扫描 该captureXml文件夹下需要解析的xml文件内容时，错误信息为：", uerr)
				}

				log.Println("获取抓拍结果，result:", result.VehicleImgPath)

				//把图片上传到oss上
				c := strings.Split(result.VehicleImgPath, ":")
				str2 := strings.Replace(c[1], "\\", "/", -1) //linux
				log.Println(str2)
				strfname := strings.Split(str2, "/")
				//上传到oss                    日期文件夹     图片名称               前缀
				code := utils.QingStorUpload(strfname[6], strfname[7], "/jiangsu/suhuaiyangs/")
				if code == utils.UPloadOK {
					//删除本地图片
					utils.DelFile("./images/" + strfname[6] + "/" + strfname[7])

					//生产xml返回给云平台 [暂时上传到模拟云平台]
					GwCaptureInforUpload(&result)
					//删除抓拍xml文件

				} else {
					continue
				}

			}

		}
		log.Println("执行 处理数据包，休息1分钟+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

		log.Println("执行处理数据包", (<-tiker.C).Format("2006-01-02 15:04:05"))
	}
}

func GwCaptureInforUpload(Result *dto.CaptureDateXML) {
	//判断哪一种品牌相机
	//Result.
	data := new(dto.DateXML)
	//抓拍结果的赋值

	//MarshalIndent 有缩进 xml.Marshal ：无缩进
	ba, _ := xml.MarshalIndent(data, "  ", "  ")
	log.Println("+++++++++", string(ba))

	log.Println("Address:", GwCaptureInformationUploadIpAddress)
	result, err := GwCaptureInformationUploadPostWithXML(&ba)
	if err != nil {
		return
	}

	if (*result).Code == 0 {
		log.Println("上传抓拍结果成功")
		return
	}

}

//与抓拍进程交互心跳，得知抓拍进程程序死活
func Heartbeat(port string) {

	//SERVER_IP = "127.0.0.1" IpAddress
	time.Sleep(time.Second * 10)
	//dport, _ := strconv.Atoi(port)

	//主动给抓拍进程心跳  10秒
	//	go HeartbeatClient(strconv.Itoa(dport - 1))

	//监控抓拍进程的心跳
XT:
	ip := strings.Split(IpAddress, ":")
	address := ip[0] + ":" + port //SERVER_PORT
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Println("监控抓拍进程心跳 net.ResolveUDPAddr 时 err:", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println("监控抓拍进程的心跳 net.ListenUDP err:", err)
		time.Sleep(time.Second * 30)
		goto XT
	}
	log.Println("管理平台 UDP监听 address:", address)

	defer func() {
		_ = conn.Close()
	}()

	for {
		//获取数据
		// Here must use make and give the lenth of buffer
		data := make([]byte, 512)

		//返回一个UDPAddr        ReadFromUDP从c读取一个UDP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址。
		//ReadFromUDP方法会在超过一个固定的时间点之后超时，并返回一个错误。
		_, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println(err)
			continue
		}

		strData := string(data)
		log.Println("Received:", strData)
		//转大写
		//	upper := strings.ToUpper(strData)
		strData = "管理平台收到抓拍进程的信息 ｜" + strData
		_, err = conn.WriteToUDP([]byte(strData), rAddr)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("管理平台 Send:", strData)
	}
}

func HeartbeatClient(port string) {
	tiker := time.NewTicker(time.Second * 10) //每15秒执行一下
	for {
		log.Println(utils.DateTimeFormat(<-tiker.C), "管理平台要发送心跳给抓拍进程++++++++++++")

		Heartbeatclient(port)

	}
}

func Heartbeatclient(port string) {

	serverAddr := "192.168.26.248" + ":" + port
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

	msg := make([]byte, 512)
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
