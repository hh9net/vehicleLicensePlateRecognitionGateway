package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/go-kratos/kratos/pkg/sync/errgroup"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"vehicleLicensePlateRecognitionGateway/dto"
	"vehicleLicensePlateRecognitionGateway/utils"
)

var (
	Parsexmlcount            int
	files                    chan string
	OSSCount                 int
	NewOSSCount              int
	Parsed                   int
	ResultCount              int
	ResultOKCount            int
	AgainCount               int
	Deviceid                 string //网关设备id Token
	StationId                map[string]string
	DeviceId                 map[string]string
	CmeraId                  map[string]string
	LaneType                 map[string]string
	ImageType                map[string]string
	EngineId                 map[string]string
	Token                    string
	NewDataNotificationCount int //新数据通知
	// HasUploadFile []string
	Pid map[string]string

	BacketName   string
	ObjectPrefix string
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
		log.Println("获取网关设备的token 失败,重新请求!", time.Now()) //getTokenerr 已打印
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

	DeviceId = make(map[string]string, len(CameraList.Data))
	StationId = make(map[string]string, len(CameraList.Data))
	LaneType = make(map[string]string, len(CameraList.Data))
	ImageType = make(map[string]string, len(CameraList.Data))
	EngineId = make(map[string]string, len(CameraList.Data))
	CmeraId = make(map[string]string, len(CameraList.Data))
	Pid = make(map[string]string, len(CameraList.Data))
	log.Println(" 相机列表数据的len:", len(CameraList.Data))
	log.Println(" 相机列表数据:", CameraList.Data)
	uniview := make([]dto.CameraListData, 0) // 宇视的列表

	hikITS := make([]dto.CameraListData, 0) //ITS列表

	for i, cmera := range CameraList.Data {
		//StationId
		//deviceid应该用gantryID
		StationId[cmera.Id] = cmera.StationId
		DeviceId[cmera.Id] = cmera.Gantryid //deviceid应该用gantryID
		LaneType[cmera.Id] = cmera.LaneType
		ImageType[cmera.Id] = cmera.Description
		EngineId[cmera.Id] = cmera.DevCompId //相机品牌

		log.Println(i, "StationId:", StationId[cmera.Id], cmera.StationId)
		log.Println(i, "DeviceId:", DeviceId[cmera.Id], cmera.Gantryid)
		log.Println(i, "LaneType:", LaneType[cmera.Id], cmera.LaneType)
		log.Println(i, "ImageType:", ImageType[cmera.Id], cmera.Description)
		log.Println(i, "EngineId:", EngineId[cmera.Id], cmera.DevCompId)

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
		log.Println("+++该网关设备没有海康ITS相机和宇视相机")
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

}

//1、启动进程
func Runmain(Configfname string) {

	log.Println("启动进程Configfname:", Configfname)

	port := strings.Split(Configfname, "+")

	//心跳port
	xtpt := strings.Split(port[2], ".")

	//cmd := exec.Command("udpmain.exe绝对路径")
	//cmd := exec.Command("./snap/udpmain")
	var additionalBilldataDir string
	additionalBilldataDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	var billoutputDir = filepath.Join(additionalBilldataDir, "snap", "udpmain.exe")
	log.Println("udpmain.exe绝对路径:", billoutputDir)

	cmd := exec.Command(billoutputDir)

	path := make([]string, 0)

	var configxmlpath = filepath.Join(additionalBilldataDir, "cameraConfig", Configfname)
	//configxmlpath启动进程的配置文件的绝对路径 cameraConfig+ Configfname
	path = append(path, configxmlpath)

	cmd.Args = path
	log.Println("cmd.Args:", cmd.Args)
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
	go Heartbeat(xtpt[0])
	//进程在进行心跳交互处是阻塞的
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
		log.Println("GetCameraList error:", err)
		return nil, err
	}
	log.Println("根据token获取camera列表成功！！！")
	return Resp, nil
}

//上传文件  开线程读取xml文件 上传图片到oss  上传抓拍结果到车牌识别云端服务器
func UploadFile() {
	Parsexmlcount = 0
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()
	files = make(chan string, 100)
	eg := errgroup.WithContext(ctx)
	eg.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go parse(ctx)
	}
	for {
		//上传图片以及抓拍结果到车牌识别云端服务器
		hferr := HandleFile(ctx)
		if hferr != nil {
			log.Println("执行HandleFile上传图片以及抓拍结果 到 车牌识别云端服务器hferr error:", hferr)
			time.Sleep(time.Second * 1)
		}

	}
}

func HandleFile(ctx context.Context) error {
	//起go程提取xml文件
	err, hasNew := extract(ctx)
	if err != nil {
		log.Printf("extract: %v\n", err)
	}
	if !hasNew {
		time.Sleep(time.Second * 1)
	}
	return nil
}

//提取文件
func extract(Ctx context.Context) (err error, hasNewFile bool) {
	dir, _ := os.Getwd()
	log.Println("当前路径：", dir)
	var snapxmlPathDir = filepath.Join(dir, "snap", "xml")
	// check
	if _, err := os.Stat(snapxmlPathDir); err == nil {
		fmt.Println("path exists 1", snapxmlPathDir)
	} else {
		fmt.Println("path not exists ", snapxmlPathDir)
		err := os.MkdirAll(snapxmlPathDir, 0711)
		if err != nil {
			log.Println("Error creating directory")
			log.Println(err)
		}
	}
	//	for {
	//提取xml文件夹下文件
	fileList, err := ioutil.ReadDir(snapxmlPathDir)
	if err != nil {
		log.Println("扫描 snapxml 文件夹 读取文件信息 error:", err)
		time.Sleep(time.Second * 1)
		return err, false
	}
	log.Println("执行 扫描 该snap/xml/文件夹下有文件的数量 ：", len(fileList))
	if len(fileList) == 1 {
		fmt.Println("执行 扫描 该snap/xml文件夹下可能没有需要解析的xml文件") //有隐藏文件
		time.Sleep(time.Second * 1)
	} else {
		if len(fileList) == 0 {
			fmt.Println("执行 扫描 该snap/xml/文件夹下没有需要解析的xml文件")
			time.Sleep(time.Second * 1)
			return nil, false
		}
	}

	for i := range fileList {
		//判断文件的结尾名
		if strings.HasSuffix(fileList[i].Name(), ".xml") {

			//这里对文件做一个处理先解析文件，以防xml中数据没有写完
			content, err := ioutil.ReadFile(snapxmlPathDir + "/" + fileList[i].Name())
			if err != nil {
				log.Println("防止抓拍程序写xml中数据没有写完时。读取这文件失败:", err)
				continue
			}
			//将xml文件转换为对象
			var result dto.CaptureDateXML
			uerr := xml.Unmarshal(content, &result)
			if uerr != nil {
				log.Println("该snap/xml/文件夹下需要解析的xml文件内容时，这个xml文件还没有写完，error:", uerr)
				continue
			}

			log.Println("获取抓拍结果中，图片路径result.VehicleImgPath:", result.VehicleImgPath)

			log.Println("该snap/xml/文件夹下需要解析的xml文件名字为:", fileList[i].Name())

			Renameerr := RenameFile(snapxmlPathDir+"/"+fileList[i].Name(), snapxmlPathDir+"/"+fileList[i].Name()+"_suffix")
			if Renameerr != nil {
				log.Println("该snap/xml/文件夹下需要解析的xml，改文件名时错误，文件名字为:", fileList[i].Name())
				time.Sleep(time.Second * 1)
				continue
			}
			select {
			case files <- snapxmlPathDir + "/" + fileList[i].Name() + "_suffix":

			case <-Ctx.Done():
				//return
			}

		} // if .xml
	}

	log.Println("extract()提取文件完成。休息3秒中")
	time.Sleep(time.Second * 3)
	return nil, true
}

//解析xml文件
func parse(Ctx context.Context) {

	dir, _ := os.Getwd()
	log.Println("当前路径：", dir)
	var snapxmlPathdir = filepath.Join(dir, "snap", "xml")

	for {
		select {
		case filePath := <-files:
			if err := UploadFileToOSS(snapxmlPathdir, filePath); err == nil {
				log.Println("执行parse处理xml数据包解析以及oss上传以及抓拍结果上传完成")
			}
		case <-Ctx.Done():
			return
		}
	}

}

func UploadFileToOSS(snapxmlPathdir, xmlnamepath string) (err error) {
	log.Println("获取抓拍的结果的xml的文件夹路径，snapxmlPathdir:", snapxmlPathdir)
	log.Println("获取抓拍的结果xml文件绝对路径，xmlnamepath", xmlnamepath)
	content, err := ioutil.ReadFile(xmlnamepath)
	if err != nil {
		log.Println("执行读取抓拍结果xml文件的位置错误信息:", err)
		return err
	}
	//将xml文件转换为对象
	var result dto.CaptureDateXML
	uerr := xml.Unmarshal(content, &result)
	if uerr != nil {
		log.Println("执行UploadFileToOSS扫描 该snap/xml/文件夹下需要解析的xml文件内容时，错误信息为：", uerr)
		return uerr
	}

	log.Println("获取抓拍结果中，图片路径result.VehicleImgPath:", result.VehicleImgPath)

	//处理新版本的信路威有三种图片上传情况
	if result.VehicleImgPath1 != "" {
		swlNewerr := SignalwayNewUpload(result, xmlnamepath, snapxmlPathdir)
		if swlNewerr != nil {
			return swlNewerr
		}

		return nil
	}

	//把图片上传到oss上
	strfname := strings.Split(result.VehicleImgPath, "\\") //windows
	//上传到oss                    日期文件夹     图片名称               前缀"/jiangsu/suhuaiyangs"
	log.Println("上传到oss图片的地址", result.VehicleImgPath)
	log.Println("上传到oss图片的名称", strfname[len(strfname)-1])
	log.Println("上传到oss的前缀", ObjectPrefix)

	ImgPath := strings.Split(result.VehicleImgPath, strfname[len(strfname)-1])
	// check，防止被删除文件夹
	//新建图片文件夹
	if _, err := os.Stat(ImgPath[0]); err == nil {
		fmt.Println("path exists 1", ImgPath[0])
	} else {
		fmt.Println("path not exists ", ImgPath[0])
		err := os.MkdirAll(ImgPath[0], 0711)

		if err != nil {
			log.Println("Error creating directory")
			log.Println(err)
		}
	}
	//xml的文件路径
	strxmlnamepath := strings.Split(xmlnamepath, "/")
	//获取文件名称
	Xmlname := strxmlnamepath[len(strxmlnamepath)-1]
	log.Println("xml的文件路径中Xmlname:", Xmlname)
	//上传oss图片
	Stationid := ""
	if val, ok := StationId[result.CamId]; ok == true {
		Stationid = val //   string   `xml:"stationid"`//	stationid站点编号
	}
	log.Println("站点IdStationid:", Stationid)
	//前缀/站点Id/摄像机ID/日期/passid
	Pname := ObjectPrefix + "/" + Stationid + "/" + result.CamId + "/" + time.Now().Format("2006-01-02") + "/" + strfname[len(strfname)-1]
	log.Printf("前缀/站点Id/摄像机ID/日期/passid==:%s", Pname)
	code, scsj, ossDZ := utils.QingStorUpload(result.VehicleImgPath, strfname[len(strfname)-1], Pname)

	if code == utils.UPloadOK {
		OSSCount = OSSCount + 1
		log.Println("上传到oss   成功，开始返回抓拍结果给云平台")
		log.Println("上传到oss   成功，OSSCount:", OSSCount, time.Now().Format("2006-01-02 15:04:05"))
		//删除本地图片 result.VehicleImgPath
		utils.DelFile(result.VehicleImgPath)
		//生产xml返回给云平台 [暂时上传到模拟云平台]
		// check
		if _, err := os.Stat(snapxmlPathdir + "/error/upload/"); err == nil {
			fmt.Println("path exists 1", snapxmlPathdir+"/error/upload/")
		} else {
			fmt.Println("path not exists ", snapxmlPathdir+"/error/upload/")
			err := os.MkdirAll(snapxmlPathdir+"/error/upload/", 0711)

			if err != nil {
				log.Println("Error creating directory")
				log.Println(err)
			}
		}
		ossDZ3 := ""
		ossDZ2 := ossDZ3
		//第一次上传失败的抓拍结果存储于【errorpathname】：snapxmlPathDir+"/error/upload/"+fileList[i].Name()
		uploaderr := GwCaptureInforUpload(&result, scsj, ossDZ, ossDZ2, ossDZ3, snapxmlPathdir+"/error/upload/"+Xmlname)
		if uploaderr != nil {
			//删除抓拍xml文件
			//xml/error
			source := snapxmlPathdir + "/" + Xmlname
			d := snapxmlPathdir + "/error/" + Xmlname
			mverr := utils.MoveFile(source, d)
			if mverr != nil {
				log.Println("第一次上传抓拍结果xml文件到云平台失败，进程抓拍结果的xml文件移动到error文件夹失败！")
				log.Println(mverr)
				return mverr
			}
			log.Println("第一次上传抓拍结果xml文件到云平台失败,xml文件移动到error文件夹成功")
			return nil
		} else {
			//上传抓拍结果到云平台成功
			DelFile(xmlnamepath)
			Parsed = Parsed + 1
			Parsexmlcount = Parsexmlcount + 1
			log.Println("Parsexmlcount:", Parsexmlcount)
			log.Println("第一次上传抓拍结果xml文件到云平台成功，进程抓拍结果xml移动到parsed 成功,Parsed:", Parsed, time.Now())
		}
	} else {
		log.Println("上传oss失败", code)
		//上传oss失败
		//删除抓拍xml文件
		//xml/error
		// check
		// ossError 图片不存在或者是上传oos的其他问题
		if _, err := os.Stat(snapxmlPathdir + "/error/ossError/"); err == nil {
			fmt.Println("path exists 1", snapxmlPathdir+"/error/ossError/")
		} else {
			fmt.Println("path not exists ", snapxmlPathdir+"/error/ossError/")
			err := os.MkdirAll(snapxmlPathdir+"/error/ossError/", 0711)

			if err != nil {
				log.Println("Error creating directory")
				log.Println(err)
			}
		}
		source := snapxmlPathdir + "/" + Xmlname
		d := snapxmlPathdir + "/error/ossError/" + Xmlname
		mverr := utils.MoveFile(source, d)
		if mverr != nil {
			log.Println("上传oss失败，进程抓拍结果的xml文件移动到error文件夹失败")
			log.Println(mverr)
			return mverr
		}
		log.Println("上传oss失败，进程抓拍结果的xml文件移动到error文件夹成功")
		return nil
	}
	return nil
}

//新版本的信路威有三种图片
func SignalwayNewUpload(result dto.CaptureDateXML, xmlnamepath, snapxmlPathdir string) error {

	//把图片上传到oss上
	strfname := strings.Split(result.VehicleImgPath, "\\") //windows
	//上传到oss                    日期文件夹     图片名称               前缀"/jiangsu/suhuaiyangs"
	log.Println("上传到oss图片1的地址，result.VehicleImgPath1:", result.VehicleImgPath)
	log.Println("上传到oss图片的名称", strfname[len(strfname)-1])

	//
	str2fname := strings.Split(result.VehicleImgPath1, "\\") //windows
	log.Println("上传到oss图片2的地址，result.VehicleImgPath2:", result.VehicleImgPath1)
	log.Println("上传到oss图片的名称", str2fname[len(str2fname)-1])

	//
	str3fname := strings.Split(result.VehicleImgPath2, "\\") //windows
	log.Println("上传到oss图片3的地址，result.VehicleImgPath3:", result.VehicleImgPath2)
	log.Println("上传到oss图片的名称", str3fname[len(str3fname)-1])

	log.Println("上传到oss的前缀", ObjectPrefix)

	ImgPath := strings.Split(result.VehicleImgPath, strfname[len(strfname)-1])
	// check，防止被删除文件夹
	//新建图片文件夹
	if _, err := os.Stat(ImgPath[0]); err == nil {
		fmt.Println("path exists 1", ImgPath[0])
	} else {
		fmt.Println("path not exists ", ImgPath[0])
		err := os.MkdirAll(ImgPath[0], 0711)

		if err != nil {
			log.Println("Error creating directory")
			log.Println(err)
		}
	}

	//xml的文件路径
	strxmlnamepath := strings.Split(xmlnamepath, "/")
	//获取文件名称
	Xmlname := strxmlnamepath[len(strxmlnamepath)-1]
	log.Println("xml的文件路径中Xmlname:", Xmlname)
	//上传oss图片
	Stationid := ""
	if val, ok := StationId[result.CamId]; ok == true {
		Stationid = val //   string   `xml:"stationid"`//	stationid站点编号
	}
	//log.Println("站点IdStationid:", Stationid)
	//前缀/站点Id/摄像机ID/日期/passid
	Pname := ObjectPrefix + "/" + Stationid + "/" + result.CamId + "/" + time.Now().Format("2006-01-02") + "/" + strfname[len(strfname)-1]
	log.Printf("前缀/站点Id/摄像机ID/日期/passid==:%s", Pname)
	code, scsj, ossDZ := utils.QingStorUpload(result.VehicleImgPath, strfname[len(strfname)-1], Pname)

	Pname2 := ObjectPrefix + "/" + Stationid + "/" + result.CamId + "/" + time.Now().Format("2006-01-02") + "/" + str2fname[len(str2fname)-1]
	log.Printf("前缀/站点Id/摄像机ID/日期/passid==:%s", Pname2)
	code2, scsj2, ossDZ2 := utils.QingStorUpload(result.VehicleImgPath1, str2fname[len(str2fname)-1], Pname2)
	log.Printf("第二张图片上传时间:%v", scsj2)

	Pname3 := ObjectPrefix + "/" + Stationid + "/" + result.CamId + "/" + time.Now().Format("2006-01-02") + "/" + str3fname[len(str3fname)-1]
	log.Printf("前缀/站点Id/摄像机ID/日期/passid==:%s", Pname3)
	code3, scsj3, ossDZ3 := utils.QingStorUpload(result.VehicleImgPath2, str3fname[len(str3fname)-1], Pname3)
	log.Printf("第二张图片上传时间:%v", scsj3)

	if code == utils.UPloadOK && code2 == utils.UPloadOK && code3 == utils.UPloadOK {
		NewOSSCount = NewOSSCount + 3
		log.Println("新版信路威上传到oss 3图都成功，开始返回抓拍结果给云平台")
		log.Println("新版信路威上传到oss 3图都成功，NewOSSCount:", NewOSSCount)
		//删除本地图片 result.VehicleImgPath
		utils.DelFile(result.VehicleImgPath)
		utils.DelFile(result.VehicleImgPath1)
		utils.DelFile(result.VehicleImgPath2)
		//生产xml返回给云平台 [暂时上传到模拟云平台]
		// check
		if _, err := os.Stat(snapxmlPathdir + "/error/upload/"); err == nil {
			fmt.Println("path exists 1", snapxmlPathdir+"/error/upload/")
		} else {
			fmt.Println("path not exists ", snapxmlPathdir+"/error/upload/")
			err := os.MkdirAll(snapxmlPathdir+"/error/upload/", 0711)

			if err != nil {
				log.Println("Error creating directory")
				log.Println(err)
			}
		}

		//第一次上传失败的抓拍结果存储于【errorpathname】：snapxmlPathDir+"/error/upload/"+fileList[i].Name()
		uploaderr := GwCaptureInforUpload(&result, scsj, ossDZ, ossDZ2, ossDZ3, snapxmlPathdir+"/error/upload/"+Xmlname)
		if uploaderr != nil {
			//删除抓拍xml文件
			//xml/error
			source := snapxmlPathdir + "/" + Xmlname
			d := snapxmlPathdir + "/error/" + Xmlname
			mverr := utils.MoveFile(source, d)
			if mverr != nil {
				log.Println("第一次上传3图抓拍结果xml文件到云平台失败，进程抓拍结果的xml文件移动到error文件夹失败！")
				log.Println(mverr)
				return mverr
			}
			log.Println("第一次上传3图抓拍结果xml文件到云平台失败，xml文件移动到error文件夹成功")
			return nil
		} else {
			//删除抓拍xml文件
			//xml/parsed
			// check
			//if _, err := os.Stat(snapxmlPathDir + "/parsed/"); err == nil {
			//	log.Println("path exists 1", snapxmlPathDir+"/parsed/")
			//} else {
			//	log.Println("path not exists ", snapxmlPathDir+"/parsed/")
			//	err := os.MkdirAll(snapxmlPathDir+"/parsed/", 0711)
			//
			//	if err != nil {
			//		log.Println("Error creating directory")
			//		log.Println(err)
			//	}
			//}

			//// check again
			//if _, err := os.Stat(snapxmlPathDir + "/parsed/"); err == nil {
			//	log.Println("path exists 2", snapxmlPathDir+"/parsed/")
			//}
			//	source := snapxmlPathdir + "/" + Xmlname
			//d := snapxmlPathDir + "/parsed/" + fileList[i].Name()
			DelFile(xmlnamepath)
			Parsed = Parsed + 1
			Parsexmlcount = Parsexmlcount + 1
			log.Println("Parsexmlcount:", Parsexmlcount)
			log.Println("第一次上传3图抓拍结果xml文件到云平台成功，进程抓拍结果xml移动到parsed 成功,Parsed:", Parsed, time.Now())
		}
	} else {
		log.Println("上传oss失败", code)
		//上传oss失败
		//删除抓拍xml文件
		//xml/error
		// check
		// ossError 图片不存在或者是上传oos的其他问题
		if _, err := os.Stat(snapxmlPathdir + "/error/3tuossError/"); err == nil {
			fmt.Println("path exists 1", snapxmlPathdir+"/error/3tuossError/")
		} else {
			fmt.Println("path not exists ", snapxmlPathdir+"/error/3tuossError/")
			err := os.MkdirAll(snapxmlPathdir+"/error/3tuossError/", 0711)

			if err != nil {
				fmt.Println("Error creating directory")
				log.Println(err)
			}
		}
		source := snapxmlPathdir + "/" + Xmlname
		d := snapxmlPathdir + "/error/3tuossError/" + Xmlname
		mverr := utils.MoveFile(source, d)
		if mverr != nil {
			log.Println("上传3图的oss失败，进程抓拍结果的xml文件移动到error文件夹失败")
			log.Println(mverr)
			return mverr
		}
		log.Println("上传3图的oss失败，进程抓拍结果的xml文件移动到error文件夹成功")
		return nil
	}
	return nil
}

func DelFile(src string) {
	//"./1.txt"
	del := os.Remove(src)
	if del != nil {
		log.Println("删除失败", del)
		return
	}
	//time.Sleep(time.Millisecond * 100)
	log.Println("删除xmlok", src)
}

func RenameFile(src string, des string) error {
	//err := os.Rename("./a", "/tmp/a")
	err := os.Rename(src, des)
	if err != nil {
		log.Println("Rename错误:", err)
		return err
	}
	log.Printf("Rename文件：%s to： %s 成功", src, des)
	return nil
}

func HandleFileAgainUpload() {
	//定期检查抓拍文件夹文件夹 captureXml
	log.Println(" HandleFileAgainUpload 执行处理xml数据包解析以及抓拍结果再次上传")
	//2、处理文件
	//扫描 captureXml 文件夹 读取文件信息
	dir, _ := os.Getwd()
	log.Println("++当前路径：", dir)

	var AgainUpsnapxmlpathDir = filepath.Join(dir, "snap", "xml", "error", "upload")
	log.Println("/snap/xml/error/upload/绝对路径:", AgainUpsnapxmlpathDir) //可以不需要加"/"

	for {
		// check
		if _, err := os.Stat(AgainUpsnapxmlpathDir); err == nil {
			fmt.Println("path exists 1", AgainUpsnapxmlpathDir)
		} else {
			fmt.Println("path not exists ", AgainUpsnapxmlpathDir)
			err := os.MkdirAll(AgainUpsnapxmlpathDir, 0711)
			if err != nil {
				log.Println("Error creating directory")
				log.Println(err)
			}
		}

		/*// check again
		if _, err := os.Stat(AgainUpsnapxmlpathDir); err == nil {
			fmt .Println("path exists 2", AgainUpsnapxmlpathDir)
		}*/

		fileList, err := ioutil.ReadDir(AgainUpsnapxmlpathDir) //不需要加"/"
		if err != nil {
			log.Println("扫描/snap/xml/error/upload/文件夹,读取文件信息error:", err)
			return
		}
		log.Println("扫描该/snap/xml/error/upload/文件夹下有文件的数量:", len(fileList))
		if len(fileList) == 1 {
			fmt.Println("扫描该/snap/xml/error/upload/文件夹下可能没有需要解析的xml文件") //有隐藏文件
		} else {
			if len(fileList) == 0 {
				fmt.Println("扫描该/snap/xml/error/upload/文件夹下没有需要解析的xml文件")
				time.Sleep(time.Second * 5)
				continue
			}
		}

		for i := range fileList {
			//判断文件的结尾名+ "_suffix"
			if strings.HasSuffix(fileList[i].Name(), ".xml_suffix") {
				log.Println("扫描该/snap/xml/error/upload/文件夹下需要解析的xml文件名字为:", fileList[i].Name())
				//error/upload/fname
				content, err := ioutil.ReadFile(AgainUpsnapxmlpathDir + "/" + fileList[i].Name())
				if err != nil {
					log.Println("读/upload/文件夹中文件错误信息:", err)
					continue
				}

				result, UploadPostWithXMLerr := GwCaptureInformationUploadPostWithXML(&content)
				if UploadPostWithXMLerr != nil {
					log.Println("需要再再次上传的抓拍结果xml文件pathname:", AgainUpsnapxmlpathDir+"/"+fileList[i].Name())
					log.Println("需要再次上传的抓拍结果xml文件失败：", UploadPostWithXMLerr)
					continue
				} else {
					//删除抓拍xml文件
					//xml/error/upload/
					source := AgainUpsnapxmlpathDir + "/" + fileList[i].Name()
					utils.DelFile(source)
					//log.Println("再次上传的抓拍结果成功,已经删除再次上传的抓拍结果的xml",source)
					//再次上传的数量或者说第一次上传失败的
					AgainCount = AgainCount + 1
					log.Println("再次上传的数量或者说第一次上传失败的数量")
					log.Println("再次上传的抓拍结果成功,AgainCount:", AgainCount)
				}
				if (*result).Code == 0 {
					fmt.Println("再次上传的抓拍结果成功")
					continue
				} /*else {
					fmt .Println("再次上传的抓拍结果失败")
					continue
				}*/
			}
		}
		time.Sleep(time.Minute * 3)
	}
}

//errorpathname：snapxmlPathDir+"/error/upload/"+fileList[i].Name()
func GwCaptureInforUpload(Result *dto.CaptureDateXML, scsj int64, ossDZ, ossDZ2, ossDZ3, errorpathname string) error {
	//判断哪一种品牌相机
	//Result.
	var ba []byte

	if Result.AppedInfo.AxleDist != "" {
		data := new(dto.TBXJDateXML)
		//抓拍结果的赋值
		data.Token = Token //抓拍结果上传

		if val, ok := StationId[Result.CamId]; ok == true {
			data.LprInfo.Stationid = val //   string   `xml:"stationid"`//	stationid站点编号
		} else {
			data.LprInfo.Stationid = ""
		}
		//	log.Println("data.LprInfo.Stationid:", data.LprInfo.Stationid)
		if val, ok := DeviceId[Result.CamId]; ok == true {
			data.LprInfo.DeviceId = val //    string   `xml:"deviceId"`//deviceId>前置机编号  deviceid应该用gantryID
		} else {
			data.LprInfo.DeviceId = ""
		}
		//	log.Println("data.LprInfo.DeviceId:", data.LprInfo.DeviceId)
		data.LprInfo.PassId = Result.PassId //    string   `xml:"passId"`         // 过车编号
		data.LprInfo.CamId = Result.CamId   //    string   `xml:"camId"`          //camId>    摄像机编号

		data.LprInfo.PassTime = Result.PassTime //    string   `xml:"passTime"`       //passTime>     过车时间
		data.LprInfo.VehicleImgPath = ossDZ     //    string   `xml:"vehicleImgPath"` //vehicleImgPath>  "oss地址"   过车图片地址
		data.LprInfo.PlateImgPath = ""          //无 string   `xml:"plateImgPath"`   //<plateImgPath/>     车牌图片地址【无】
		data.LprInfo.BucketId = BacketName      //   string   `xml:"bucketId"`       //bucketId>   bucket编号

		if val, ok := ImageType[Result.CamId]; ok == true {
			v, _ := strconv.Atoi(val)
			data.LprInfo.ImageType = v //   int      `xml:"imageType"`      //	imageType> 图片类型
		} else {
			data.LprInfo.ImageType = 0
		}

		//	log.Println("data.LprInfo.ImageType:", data.LprInfo.ImageType)
		data.LprInfo.UploadStamp = scsj //   int64    `xml:"uploadStamp"`    //	uploadStamp> 上传时间

		data.LprInfo.LaneType = 0 //   int      `xml:"laneType"`       //	laneType> 出入口类型 0:入口；1：出口

		if val, ok := LaneType[Result.CamId]; ok == true {
			v, _ := strconv.Atoi(val)
			data.LprInfo.LaneType = v //   int      `xml:"laneType"`       //	laneType> 出入口类型 0:入口；1：出口
		} else {
			data.LprInfo.LaneType = 0
		}
		//log.Println("data.LprInfo.LaneType:", data.LprInfo.LaneType)

		data.LpaResult.PassId = Result.PassId  //passId>过车编号
		data.LpaResult.EngineType = EngineType //`xml:"engineType"`      //engineType>   引擎类型  写死sjk-camera-lpa
		//log.Println("data.LpaResult.EngineType:", data.LpaResult.EngineType)

		//EngineId
		if val, ok := EngineId[Result.CamId]; ok == true {

			data.LpaResult.EngineId = val //`xml:"engineId"`        //engineId> 引擎编号 相机品牌名称
		} else {
			data.LpaResult.EngineId = "0"

		}
		//	log.Println("data.LpaResult.EngineId:", data.LpaResult.EngineId)
		if "" == Result.PlateNo {
			log.Println("车牌编号不能为空呀，需要抓拍程序定位")
			data.LpaResult.PlateNo = "unrecognized"
		} else {
			data.LpaResult.PlateNo = Result.PlateNo //`xml:"plateNo"`         //plateNo>     车牌编号
		}

		data.LpaResult.PlateColor = ChepZH(Result.PlateColor) // `xml:"plateColor"`      //plateColor>     车牌颜色
		data.LpaResult.ComputeInterval = 0                    //int64 `xml:"computeInterval"` //computeInterval>  计算时间

		data.LpaResult.VehicleColor = "" //`xml:"vehicleColor"`    //vehicleColor>       车辆颜色
		data.LpaResult.VehicleType = ""  //`xml:"vehicleType"`     //vehicleType>       车辆类型
		data.LpaResult.VehicleBrand = "" //`xml:"vehicleBrand"`    //vehicleBrand>       车辆品牌
		data.LpaResult.VehicleYear = 0   //int`xml:"vehicleYear"`     //vehicleYear>     车辆年份

		data.LpaResult.LprFrameEntity.PlateLeft = 0   // int      `xml:"plateLeft"`   //plateLeft>        车牌左坐标
		data.LpaResult.LprFrameEntity.PlateTop = 0    //  int      `xml:"plateTop"`    //plateTop>        车牌上坐标
		data.LpaResult.LprFrameEntity.PlateRight = 0  // int      `xml:"plateRight"`  //plateRight>        车牌右坐标
		data.LpaResult.LprFrameEntity.PlateBottom = 0 // int      `xml:"plateBottom"` //plateBottom>     车牌下坐标

		data.VehicleInfo.SideImgPath = ossDZ2               // string   `xml:"sideImgPath"` //sideImgPath> 侧面图片地址
		data.VehicleInfo.TailImgPath = ossDZ3               //  string   `xml:"tailImgPath"` //tailImgPath> 车尾图片地址
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
		log.Println("前置机抓拍信息上传数据 +++++++++", string(ba))

	} else {
		data := new(dto.DateXML)
		//抓拍结果的赋值
		data.Token = Token //抓拍结果上传

		if val, ok := StationId[Result.CamId]; ok == true {
			data.LprInfo.Stationid = val //   string   `xml:"stationid"`//	stationid站点编号
		} else {
			data.LprInfo.Stationid = ""
		}
		//	log.Println("data.LprInfo.Stationid:", data.LprInfo.Stationid)
		if val, ok := DeviceId[Result.CamId]; ok == true {
			data.LprInfo.DeviceId = val //    string   `xml:"deviceId"`//deviceId>前置机编号  deviceid应该用gantryID
		} else {
			data.LprInfo.DeviceId = ""
		}
		//		log.Println("data.LprInfo.DeviceId:", data.LprInfo.DeviceId)
		data.LprInfo.PassId = Result.PassId //    string   `xml:"passId"`         // 过车编号
		data.LprInfo.CamId = Result.CamId   //    string   `xml:"camId"`          //camId>    摄像机编号

		data.LprInfo.PassTime = Result.PassTime //    string   `xml:"passTime"`       //passTime>     过车时间
		data.LprInfo.VehicleImgPath = ossDZ     //   string   `xml:"vehicleImgPath"` //vehicleImgPath>  "oss地址"   过车图片地址
		data.LprInfo.PlateImgPath = ""          //     string   `xml:"plateImgPath"`   //<plateImgPath/>     车牌图片地址【无】
		data.LprInfo.BucketId = BacketName      //     string   `xml:"bucketId"`       //bucketId>   bucket编号
		if val, ok := ImageType[Result.CamId]; ok == true {
			v, _ := strconv.Atoi(val)
			data.LprInfo.ImageType = v //   int      `xml:"imageType"`      //	imageType> 图片类型
		} else {
			data.LprInfo.ImageType = 0
		}
		//	log.Println("data.LprInfo.ImageType:", data.LprInfo.ImageType)
		data.LprInfo.UploadStamp = scsj //     int64    `xml:"uploadStamp"`    //	uploadStamp> 上传时间

		if val, ok := LaneType[Result.CamId]; ok == true {
			v, _ := strconv.Atoi(val)
			data.LprInfo.LaneType = v //   int      `xml:"laneType"`       //	laneType> 出入口类型 0:入口；1：出口
		} else {
			data.LprInfo.LaneType = 0
		}
		//	log.Println("data.LprInfo.LaneType:", data.LprInfo.LaneType)
		data.LpaResult.PassId = Result.PassId  //passId>     过车编号
		data.LpaResult.EngineType = EngineType //`xml:"engineType"`      //engineType>   引擎类型 写死 sjk-camera-lpa

		//	log.Println("data.LpaResult.EngineType:", data.LpaResult.EngineType)

		//EngineId
		if val, ok := EngineId[Result.CamId]; ok == true {

			data.LpaResult.EngineId = val //`xml:"engineId"`  //engineId> 引擎编号 相机品牌名称
		} else {
			data.LpaResult.EngineId = "0"

		}
		//log.Println("data.LpaResult.EngineId:", data.LpaResult.EngineId)
		if "" == Result.PlateNo {
			log.Println("车牌编号不能为空呀，需要抓拍程序定位")
			data.LpaResult.PlateNo = "unrecognized"
		} else {
			data.LpaResult.PlateNo = Result.PlateNo //`xml:"plateNo"`         //plateNo>     车牌编号
		}

		data.LpaResult.PlateColor = ChepZH(Result.PlateColor) // `xml:"plateColor"`      //plateColor>     车牌颜色
		data.LpaResult.ComputeInterval = 0                    //int64 `xml:"computeInterval"` //computeInterval>  计算时间
		data.LpaResult.VehicleColor = ""                      //`xml:"vehicleColor"`    //vehicleColor>       车辆颜色
		data.LpaResult.VehicleType = ""                       //`xml:"vehicleType"`     //vehicleType>       车辆类型
		data.LpaResult.VehicleBrand = ""                      //`xml:"vehicleBrand"`    //vehicleBrand>       车辆品牌
		data.LpaResult.VehicleYear = 0                        //int`xml:"vehicleYear"`     //vehicleYear>     车辆年份

		data.LpaResult.LprFrameEntity.PlateLeft = 0   // int      `xml:"plateLeft"`   //plateLeft>        车牌左坐标
		data.LpaResult.LprFrameEntity.PlateTop = 0    //  int      `xml:"plateTop"`    //plateTop>        车牌上坐标
		data.LpaResult.LprFrameEntity.PlateRight = 0  // int      `xml:"plateRight"`  //plateRight>        车牌右坐标
		data.LpaResult.LprFrameEntity.PlateBottom = 0 // int      `xml:"plateBottom"` //plateBottom>     车牌下坐标
		//MarshalIndent 有缩进 xml.Marshal ：无缩进

		ba, _ = xml.MarshalIndent(data, "  ", "  ")
		log.Println("前置机抓拍信息上传数据 +++++++++", string(ba))
	}

	log.Println("前置机抓拍信息上传接口 Address:", GwCaptureInformationUploadIpAddress)
	//调用云平台接口
	result, err := GwCaptureInformationUploadPostWithXML(&ba)
	if err != nil {
		//需要再次上传的抓拍结果，所以需要把抓拍结果保存下来
		uploadagainxml := createXml(errorpathname, ba)
		log.Println("需要再次上传的抓拍结果xml文件uploadagainxml:", uploadagainxml)
		log.Println("第一次上传抓拍结果失败,需要再次上传的抓拍结果xml文件生成成功")
		return err
	}

	if (*result).Code == 0 {
		if result.Msg == "接收成功" {
			ResultOKCount = ResultOKCount + 1
			log.Println("前置机抓拍信息第一次上传抓拍结果成功ok,并接收成功 ResultOKCount:", ResultOKCount, time.Now().Format("2006-01-02 15:04:05"))
		}
		log.Println("第一次上传抓拍结果成功 ")
	}
	return nil
}

//创建xml文件
func createXml(xmlname string, outputxml []byte) string {

	fw, f_werr := os.Create(xmlname) //go run main.go
	if f_werr != nil {
		log.Println("Read:", f_werr)
		return ""
	}
	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, outputxml...)

	_, ferr := fw.Write((xmlOutPutData))
	if ferr != nil {
		log.Printf("  Write xml file error: %v\n", ferr)
		return ""
	}

	defer func() {
		_ = fw.Close()
	}()

	return xmlname
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
		_, err := conn.Read(buffer)
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
				ycsberr := ExcprptStuUploadPostWithJson(ycdata)
				if ycsberr != nil {

				}
				//重启程序
				rsudperr := RestartUpdmain(port)
				if rsudperr != nil {
					log.Println("重启程序时，error:", rsudperr)
				}
			}
			xtbeginsj = now
			time.Sleep(time.Second * 3)
			continue
		}

		/*	//获取数据
			// Here must use make and give the lenth of buffer
			//返回一个UDPAddr ReadFromUDP从c读取一个UDP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址。
			//ReadFromUDP方法会在超过一个固定的时间点之后超时，并返回一个错误。
			//log.Println("conn.ReadFromUDP address:", address)
			//_, rAddr, err := conn.ReadFromUDP(buffer)
			//if err != nil {
			//	log.Println(address, "conn.ReadFromUDP error:", err)
			//	continue
			//}
			//log.Println(address, "conn.ReadFromUDP ok！rAddr:", rAddr)*/

		herr := xml.Unmarshal(buffer, h)
		if herr != nil {
			log.Println(address, "UDP接收时,xml.Unmarshal失败！", herr) //这样解析是肯定OK的
			//log.Println(address, "UDP接收数据data:", string(data[:256]))
		} else {
			//接收到数据
			log.Println(address, "接收到数据1、心跳,2、新数据通知;h.Type:", h.Type, h.Uuid)
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
				log.Println(address, "1、心跳h.Type:", h.Type, h)
				heartbeatresp.Uuid = h.Uuid
				heartbeatresp.Type = h.Type       //<type>   1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
				heartbeatresp.Version = h.Version //<version>抓拍程序版本号
				heartbeatresp.Time = h.Time       //<time>   字符串2020-11-12 12:12:12
				heartbeatresp.Seq = h.Seq         //<seq>   消息序号累加

				/*				old := utils.WindowsStrTimeTotime(heartbeatresp.Time)

								now := time.Now()
								sjcstr := utils.TimeDifference(old, now)

								SJC := strings.Split(sjcstr, "s")
								sjc, _ := strconv.Atoi(SJC[0])
								//超时推出
								if sjc > 10 {
									log.Println(address, "心跳时间差大于10秒，需要重启程序")
									// 4、采集进程被动关闭命令
									h := new(dto.Heartbeat)
									heartbeatresp.Uuid = h.Uuid
									heartbeatresp.Type = 4            //<type> 1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
									heartbeatresp.Version = h.Version //<version>  抓拍程序版本号
									heartbeatresp.Time = h.Time       //<time>     字符串2020-11-12 12:12:12
									heartbeatresp.Seq = h.Seq         //<seq>   消息序号累加
									log.Println(address, "4、采集进程被动关闭命令 h.Type:", h.Type, h)
								}*/
			}

		case 2:
			//2、新数据通知
			h := new(dto.Heartbeat)
			herr := xml.Unmarshal(buffer, h)
			if herr != nil {
				log.Println(herr)
			} else {
				NewDataNotificationCount = NewDataNotificationCount + 1
				log.Println(address, "新数据通知:", h.Type, h.Uuid, "NewDataNotificationCount:", NewDataNotificationCount)
				heartbeatresp.Uuid = h.Uuid
				heartbeatresp.Type = h.Type       //<type>    1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
				heartbeatresp.Version = h.Version //<version>        抓拍程序版本号
				heartbeatresp.Time = h.Time       //<time>     字符串2020-11-12 12:12:12
				heartbeatresp.Seq = h.Seq         //<seq>   消息序号累加
			}

		/*case 3:
		//3、 日志
		h := new(dto.HeartbeatLog)
		herr := xml.Unmarshal(data, h)
		if herr != nil {
			log.Println(herr)
		} else {
			log.Println(address, "抓拍进程的日志：", h)
			heartbeatresp.Uuid = h.Uuid
			heartbeatresp.Type = h.Type       //<type>    1、心跳   2、新数据通知  3、 日志  4、采集进程被动关闭命令
			heartbeatresp.Version = h.Version //<version>        抓拍程序版本号
			heartbeatresp.Time = h.Time       //<time>     字符串2020-11-12 12:12:12
			heartbeatresp.Seq = h.Seq         //<seq>   消息序号累加
		}*/
		default:
			continue
		}

		heartbeatresp.Content = time.Now().Format("2006-01-02 15:04:05")
		//回复udp的消息
		resp, hresperr := xml.Marshal(heartbeatresp)
		if hresperr != nil {
			log.Println(address, hresperr)
		} else {
			//	log.Println(address, "xml.Marshal ok! 管理平台收到抓拍进程的 h.Type:", h.Type, "1、心跳|2、新数据通知|3、 日志|4、采集进程被动关闭命令")
		}

		////回复udp数据
		//_, err = conn.WriteToUDP(resp, rAddr)
		//if err != nil {
		//	log.Println(err)
		//	continue
		//}
		//log.Println("管理平台回复抓拍进程的udp响应 Send:", heartbeatresp)

		//回复udp数据
		hferr := Heartbeatclient(port, resp)
		if hferr != nil {
			//	log.Println(address, "此log已经打印过了")
			//log已经打印过了
			continue
		} else {
			log.Println(address, "回复成功，但是不退出，继续udp交互，")
		}

		//不管结果如何，我重启它
		//log.Println("++++++++++++++++不管结果如何，我重启它")
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

	//var n int
	_, err = conn.Write(toWrite)
	if err != nil {
		log.Println("管理平台主动给抓拍进程心跳UDP err:", err)
		return err
	}
	//log.Println(" 管理平台 主动给抓拍进程心跳 UDP 写的字节数n:", n)

	//msg := make([]byte, 32)
	//n, err = conn.Read(msg)
	//if err != nil {
	//	log.Println("管理平台 主动给抓拍进程心跳后，要收（Read）抓拍进程响应的信息 error:", err)
	//	return err
	//}
	//log.Println("管理平台 主动给抓拍进程心跳后，收到抓拍进程响应的信息，Response:", string(msg), "n:", n)
	//
	//log.Println("管理平台 主动给抓拍进程心跳后，收到抓拍进程响应的信息，Response 字节数量n:", n)
	return nil
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
		//return
		os.Exit(1)
	}
}

func HandleDayTasks() {
	for {
		now := time.Now()                                                                    //获取当前时间，放到now里面，要给next用
		next := now.Add(time.Hour * 24)                                                      //通过now偏移24小时
		next = time.Date(next.Year(), next.Month(), next.Day(), 1, 0, 0, 0, next.Location()) //获取下一个20点的日期
		t := time.NewTimer(next.Sub(now))                                                    //计算当前时间到凌晨的时间间隔，设置一个定时器
		<-t.C                                                                                //阻塞等待第二天到来才执行

		sj := time.Now().Format("2006-01-02T15:04:05")
		content := sj + "上传到oss成功数量，OSSCount=" + strconv.Itoa(OSSCount)
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
		content = "重新计数时间:" + time.Now().Format("2006-01-02T15:04:05")
		StatisticalFile(content)
		//删除前几天日期文件夹中为空的文件夹
		log.Println("执行删除前几天日期文件夹中为空的文件夹", time.Now())
		//2、处理文件
		//扫描 captureXml 文件夹 读取文件信息
		dir, _ := os.Getwd()
		log.Println("+++++++++++++++++++++++++当前路径：", dir)

		//var snapimagespath string
		//snapimagespath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

		var DelsnapimagespathDir = filepath.Join(dir, "snap", "images")
		log.Println("/snap/images/绝对路径:", DelsnapimagespathDir+"/") //+"/"
		//pwd := snapimagespathDir
		DirList, err := ioutil.ReadDir(DelsnapimagespathDir + "/") //也可以不加
		if err != nil {
			log.Println("扫描/snap/images/文件夹 读取文件信息 error:", err)
			time.Sleep(time.Second * 3)
			continue //DirListP
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
		log.Println("执行线程，处理一天一次的定时任务【完成】")
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
