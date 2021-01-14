package service

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/pkg/sync/errgroup"
	log "github.com/sirupsen/logrus"
	//"google.golang.org/genproto/googleapis/ads/googleads/v1/errors"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"vehicleLicensePlateRecognitionGateway/dto"
	"vehicleLicensePlateRecognitionGateway/utils"
)

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
			log.Println("执行HandleFile上传图片以及抓拍结果 到 车牌识别云端服务器 error:", hferr)
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
		log.Println("扫描snap/xml/文件夹,读取文件夹信息 error:", err)
		time.Sleep(time.Second * 1)
		return err, false
	}
	log.Println("扫描该snap/xml/文件夹下有文件的数量 ：", len(fileList))
	if len(fileList) == 1 {
		fmt.Println("扫描该snap/xml文件夹下可能没有需要解析的xml文件") //有隐藏文件
		time.Sleep(time.Second * 1)
	} else {
		if len(fileList) == 0 {
			fmt.Println("扫描该snap/xml/文件夹下没有需要解析的xml文件")
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
			CapCnt = CapCnt + 1
			log.Println("启动后抓拍总和:", CapCnt)
			CapZeroCnt = CapZeroCnt + 1
			log.Println("每日凌晨1点后抓拍的总和:", CapZeroCnt)

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

		}
	}

	log.Println("extract()提取文件完成。休息3秒中")
	time.Sleep(time.Second * 3)
	return nil, true
}

//解析xml文件
func parse(Ctx context.Context) {
	dir, _ := os.Getwd()
	var parsesnapxmlPathdir = filepath.Join(dir, "snap", "xml")
	for {
		select {
		case filePath := <-files:
			if err := UploadFileToOSS(parsesnapxmlPathdir, filePath); err == nil {
				log.Println("执行parse处理xml数据包解析以及oss上传以及抓拍结果上传完成")
			}
		case <-Ctx.Done():
			return
		}
	}

}

//处理xml上传
func UploadFileToOSS(snapxmlPathdir, xmlnamepath string) error {
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
		log.Println("执行UploadFileToOSS扫描snap/xml/文件夹下需要解析的xml文件内容时，错误信息为:", uerr)
		return uerr
	}

	log.Println("获取抓拍结果中，图片路径result.VehicleImgPath:", result.VehicleImgPath)

	//处理新版本的信路威有三种图片上传情况
	if result.BrandName == "SignalwayNew" {
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
		UploadRecordCnt = UploadRecordCnt + 1
		UploadImgCnt = UploadImgCnt + 1
		UploadRecordZeroCnt = UploadRecordZeroCnt + 1
		UploadImgZeroCnt = UploadImgZeroCnt + 1
		OSSCount = OSSCount + 1
		log.Println("上传到oss   成功，开始返回抓拍结果给云平台")
		log.Println("上传到oss   成功，OSSCount:", OSSCount)

		//生产xml返回给云平台 [暂时上传到模拟云平台]
		// check
		if _, err := os.Stat(snapxmlPathdir + "/error/upload/"); err == nil {
			fmt.Println("path exists 1", snapxmlPathdir+"/error/upload/")
		} else {
			log.Println("path not exists ", snapxmlPathdir+"/error/upload/")
			err := os.MkdirAll(snapxmlPathdir+"/error/upload/", 0711)
			if err != nil {
				log.Println("Error creating directory", err)
				log.Println(err)
			}
		}
		ossDZ3 := ""
		ossDZ2 := ""
		//第一次上传失败的抓拍结果存储于【errorpathname】：snapxmlPathDir+"/error/upload/"+fileList[i].Name()
		uploaderr := GwCaptureInforUpload(&result, scsj, ossDZ, ossDZ2, ossDZ3, snapxmlPathdir+"/error/upload/"+Xmlname)
		if uploaderr != nil {
			UploadFailCnt = UploadFailCnt + 1
			UploadFailZeroCnt = UploadFailZeroCnt + 1
			//删除抓拍xml文件
			//上传抓拍结果到云平台，早晚都会成功
			DelFile(xmlnamepath)
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
		//删除本地图片 result.VehicleImgPath
		utils.DelFile(result.VehicleImgPath)
	} else {
		log.Println("上传oss失败", code)
		UploadFailImgCnt = UploadFailImgCnt + 1
		UploadFailImgZeroCnt = UploadFailImgZeroCnt + 1
		//上传oss失败
		// ossError 图片不存在或者是上传oos的其他问题
		if _, err := os.Stat(snapxmlPathdir + "/error/ossError/"); err == nil {
			fmt.Println("path exists 1", snapxmlPathdir+"/error/ossError/")
		} else {
			log.Println("path not exists ", snapxmlPathdir+"/error/ossError/")
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
		//需要对上传oss失败的进行二次上传
		return nil
	}
	return nil
}

//第一次上传 新版本的信路威有三种图片的
func SignalwayNewUpload(result dto.CaptureDateXML, xmlnamepath, snapxmlPathdir string) error {
	strfname := make([]string, 0)
	str2fname := make([]string, 0)
	str3fname := make([]string, 0)
	code := 0
	code2 := code
	code3 := code
	scsj := int64(0)
	scsj2 := int64(0)
	scsj3 := int64(0)
	ossDZ := ""
	ossDZ2 := ""
	ossDZ3 := ""
	Pathis := make([]string, 0)
	//把图片上传到oss上
	if result.VehicleImgPath != "" {
		strfname = strings.Split(result.VehicleImgPath, "\\") //windows
		//上传到oss                    日期文件夹     图片名称               前缀"/jiangsu/suhuaiyangs"
		log.Println("上传到oss图片1的地址，result.VehicleImgPath1:", result.VehicleImgPath)
		log.Println("上传到oss图片的名称", strfname[len(strfname)-1])
		Pathis = append(Pathis, "VehicleImgPathok")
	} else {
		code = utils.UPloadOK
	}
	if result.VehicleImgPath1 != "" {
		str2fname = strings.Split(result.VehicleImgPath1, "\\") //windows
		log.Println("上传到oss图片2的地址，result.VehicleImgPath2:", result.VehicleImgPath1)
		log.Println("上传到oss图片的名称", str2fname[len(str2fname)-1])

		Pathis = append(Pathis, "VehicleImgPath1ok")
	} else {
		code2 = utils.UPloadOK
	}
	if result.VehicleImgPath1 != "" {
		str3fname = strings.Split(result.VehicleImgPath2, "\\") //windows
		log.Println("上传到oss图片3的地址，result.VehicleImgPath3:", result.VehicleImgPath2)
		log.Println("上传到oss图片的名称", str3fname[len(str3fname)-1])

		Pathis = append(Pathis, "VehicleImgPath2ok")
	} else {
		code3 = utils.UPloadOK
	}

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

	if code == utils.UPloadOK && code2 == utils.UPloadOK && code3 == utils.UPloadOK {
		log.Println("SignalwayNew上传时，都没有图片++++++肯定有错！")
		return errors.New("SignalwayNew上传时，3图都没有图片！++++++肯定有错！")
	}
	for _, p := range Pathis {
		switch p {
		case "VehicleImgPathok":
			//log.Println("站点IdStationid:", Stationid)
			//前缀/站点Id/摄像机ID/日期/passid
			Pname := ObjectPrefix + "/" + Stationid + "/" + result.CamId + "/" + time.Now().Format("2006-01-02") + "/" + strfname[len(strfname)-1]
			log.Printf("前缀/站点Id/摄像机ID/日期/passid==:%s", Pname)
			code, scsj, ossDZ = utils.QingStorUpload(result.VehicleImgPath, strfname[len(strfname)-1], Pname)
		case "VehicleImgPath1ok":
			Pname2 := ObjectPrefix + "/" + Stationid + "/" + result.CamId + "/" + time.Now().Format("2006-01-02") + "/" + str2fname[len(str2fname)-1]
			log.Printf("前缀/站点Id/摄像机ID/日期/passid==:%s", Pname2)
			code2, scsj2, ossDZ2 = utils.QingStorUpload(result.VehicleImgPath1, str2fname[len(str2fname)-1], Pname2)
			log.Printf("第二张图片上传时间:%v", scsj2)
		case "VehicleImgPath2ok":
			Pname3 := ObjectPrefix + "/" + Stationid + "/" + result.CamId + "/" + time.Now().Format("2006-01-02") + "/" + str3fname[len(str3fname)-1]
			log.Printf("前缀/站点Id/摄像机ID/日期/passid==:%s", Pname3)
			code3, scsj3, ossDZ3 = utils.QingStorUpload(result.VehicleImgPath2, str3fname[len(str3fname)-1], Pname3)
			log.Printf("第二张图片上传时间:%v", scsj3)

		default:

		}
	}

	if scsj == 0 {
		scsj = scsj2
		if scsj == 0 {
			scsj = scsj3
		}
	}

	if code == utils.UPloadOK && code2 == utils.UPloadOK && code3 == utils.UPloadOK {
		UploadRecordCnt = UploadRecordCnt + 1
		UploadImgCnt = UploadImgCnt + 1
		UploadImgZeroCnt = UploadImgZeroCnt + 1
		UploadRecordZeroCnt = UploadRecordZeroCnt + 1
		NewOSSCount = NewOSSCount + 3
		log.Println("新版信路威上传到oss 3图都成功，开始返回抓拍结果给云平台")
		log.Println("新版信路威上传到oss 3图都成功，NewOSSCount:", NewOSSCount)

		//生产xml返回给云平台 [暂时上传到模拟云平台]
		// check
		if _, err := os.Stat(snapxmlPathdir + "/error/upload/"); err == nil {
			fmt.Println("path exists 1", snapxmlPathdir+"/error/upload/")
		} else {
			log.Println("path not exists ", snapxmlPathdir+"/error/upload/")
			err := os.MkdirAll(snapxmlPathdir+"/error/upload/", 0711)
			if err != nil {
				log.Println("Error creating directory")
				log.Println(err)
			}
		}
		//第一次上传失败的抓拍结果存储于【errorpathname】：snapxmlPathDir+"/error/upload/"+fileList[i].Name()
		uploaderr := GwCaptureInforUpload(&result, scsj, ossDZ, ossDZ2, ossDZ3, snapxmlPathdir+"/error/upload/"+Xmlname)
		if uploaderr != nil {
			UploadFailCnt = UploadFailCnt + 1
			UploadFailZeroCnt = UploadFailZeroCnt + 1
			log.Println("第一次上传3图抓拍结果xml文件到云平台失败，xml文件移动到error文件夹成功")
			return nil
		} else {
			//删除抓拍xml文件
			Parsed = Parsed + 1
			Parsexmlcount = Parsexmlcount + 1
			log.Println("Parsexmlcount:", Parsexmlcount)
			log.Println("第一次上传3图抓拍结果xml文件到云平台成功，进程抓拍结果xml移动到parsed 成功,Parsed:", Parsed, time.Now())
		}

		//最终都是要删除的
		DelFile(xmlnamepath)
		for _, p := range Pathis {
			switch p {
			case "VehicleImgPathok":
				//删除本地图片 result.VehicleImgPath
				utils.DelFile(result.VehicleImgPath)
			case "VehicleImgPath1ok":
				utils.DelFile(result.VehicleImgPath1)
			case "VehicleImgPath2ok":
				utils.DelFile(result.VehicleImgPath2)
			default:

			}
		}
	} else {
		log.Println("上传oss失败", code)
		UploadFailImgCnt = UploadFailImgCnt + 1
		UploadFailImgZeroCnt = UploadFailImgZeroCnt + 1
		//上传oss失败
		//删除抓拍xml文件
		// ossError 图片不存在或者是上传oos的其他问题
		if _, err := os.Stat(snapxmlPathdir + "/error/3TuOssError/"); err == nil {
			fmt.Println("path exists 1", snapxmlPathdir+"/error/3TuOssError/")
		} else {
			fmt.Println("path not exists ", snapxmlPathdir+"/error/3TuOssError/")
			err := os.MkdirAll(snapxmlPathdir+"/error/3TuOssError/", 0711)

			if err != nil {
				fmt.Println("Error creating directory")
				log.Println(err)
			}
		}
		source := snapxmlPathdir + "/" + Xmlname
		d := snapxmlPathdir + "/error/3TuOssError/" + Xmlname
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

//goroutine3 处理文件再次上传
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

//goroutine4 OssError中抓拍结果再次上传
func HandleOssAgainUpload() {

	//定期检查抓拍文件夹文件夹  ossError
	log.Println("HandleOssAgainUpload 处理xml_suffix数据包解析以及抓拍结果再次上传")
	//2、处理文件
	//扫描 captureXml 文件夹 读取文件信息
	dir, _ := os.Getwd()
	var OssAgainUpsnapxmlpathDir = filepath.Join(dir, "snap", "xml", "error", "ossError")
	log.Println("/snap/xml/error/ossError/绝对路径:", OssAgainUpsnapxmlpathDir) //可以不需要加"/"
	for {
		// check
		if _, err := os.Stat(OssAgainUpsnapxmlpathDir); err == nil {
			fmt.Println("path exists 1", OssAgainUpsnapxmlpathDir)
		} else {
			log.Println("path not exists ", OssAgainUpsnapxmlpathDir)
			err := os.MkdirAll(OssAgainUpsnapxmlpathDir, 0711)
			if err != nil {
				log.Println("Error creating directory")
				log.Println(err)
			}
		}

		fileList, err := ioutil.ReadDir(OssAgainUpsnapxmlpathDir) //不需要加"/"
		if err != nil {
			log.Println("扫描/snap/xml/error/ossError/文件夹,读取文件信息error:", err)
			return
		}
		log.Println("扫描该/snap/xml/error/ossError/文件夹下有文件的数量:", len(fileList))
		if len(fileList) == 1 {
			fmt.Println("扫描该/snap/xml/error/ossError/文件夹下可能没有需要解析的xml文件") //有隐藏文件
		} else {
			if len(fileList) == 0 {
				fmt.Println("扫描该/snap/xml/error/ossError/文件夹下没有需要解析的xml文件")
				time.Sleep(time.Minute * 10)
				continue
			}
		}

		for i := range fileList {
			//判断文件的结尾名+ "_suffix"
			if strings.HasSuffix(fileList[i].Name(), ".xml_suffix") {
				log.Println("扫描该/snap/xml/error/ossError/文件夹下需要解析的xml文件名字为:", fileList[i].Name())
				//error/upload/fname
				content, err := ioutil.ReadFile(OssAgainUpsnapxmlpathDir + "/" + fileList[i].Name())
				if err != nil {
					log.Println("读/ossError/文件夹中文件错误信息:", err)
					continue
				}
				var xmlresult dto.CaptureDateXML
				uerr := xml.Unmarshal(content, &xmlresult)
				if uerr != nil {
					log.Println("执行HandleOssAgainUpload扫描 该/snap/xml/error/ossError/文件夹下需要解析的xml文件内容时，错误信息为：", uerr)
					log.Println(string(content))
					continue
				}

				log.Println("获取抓拍结果中，图片路径result.VehicleImgPath:", xmlresult.VehicleImgPath)

				//1、读取图片内容
				_, err = ioutil.ReadFile(xmlresult.VehicleImgPath)
				//2、判断图片是否存在
				if err != nil {
					//4、如果不存在，需要去getoss
					if fmt.Sprintf("%v", err) == "open "+xmlresult.VehicleImgPath+": no such file or directory" {
						log.Println("判断图片是否存在时，读的图片不存在:", err)
					} else {
						log.Println("判断图片是否存在时，读图片的错误信息:", err)
					}

				} else {
					//3、如果存在直接上传
					log.Println("再次上传oss时，os.Open imgfname ok:", xmlresult.VehicleImgPath)
				}

				//把图片上传到oss上
				strfname := strings.Split(xmlresult.VehicleImgPath, "\\") //windows
				//上传到oss                    日期文件夹     图片名称               前缀"/jiangsu/suhuaiyangs"
				log.Println("上传到oss图片的地址", xmlresult.VehicleImgPath)
				log.Println("上传到oss图片的名称", strfname[len(strfname)-1])
				log.Println("上传到oss的前缀", ObjectPrefix)
				//获取文件名称
				Xmlname := fileList[i].Name()
				log.Println("xml的文件路径中Xmlname:", Xmlname)
				//上传oss图片
				Stationid := ""
				if val, ok := StationId[xmlresult.CamId]; ok == true {
					Stationid = val //   string   `xml:"stationid"`//	stationid站点编号
				}
				log.Println("站点IdStationid:", Stationid)
				//前缀/站点Id/摄像机ID/日期/passid
				Pname := ObjectPrefix + "/" + Stationid + "/" + xmlresult.CamId + "/" + time.Now().Format("2006-01-02") + "/" + strfname[len(strfname)-1]
				log.Printf("前缀/站点Id/摄像机ID/日期/passid==:%s", Pname)

				code, scsj, ossDZ := utils.QingStorUpload(xmlresult.VehicleImgPath, strfname[len(strfname)-1], Pname)

				if code == utils.UPloadOK {
					UploadRecordCnt = UploadRecordCnt + 1
					UploadImgCnt = UploadImgCnt + 1
					UploadRecordZeroCnt = UploadRecordZeroCnt + 1
					UploadImgZeroCnt = UploadImgZeroCnt + 1
					OSSCount = OSSCount + 1
					log.Println("上传到oss   成功，开始返回抓拍结果给云平台")
					log.Println("上传到oss   成功，OSSCount:", OSSCount, time.Now().Format("2006-01-02 15:04:05"))

					ossDZ3 := ""
					ossDZ2 := ossDZ3
					//第一次上传失败的抓拍结果存储于【errorpathname】：snapxmlPathDir+"/error/upload/"+fileList[i].Name()
					uploaderr := GwCaptureInforUpload(&xmlresult, scsj, ossDZ, ossDZ2, ossDZ3, "./snap/xml/error/upload/"+Xmlname)
					if uploaderr != nil {
						UploadFailCnt = UploadFailCnt + 1
						UploadFailZeroCnt = UploadFailZeroCnt + 1
						log.Println("上传抓拍结果xml文件到云平台失败,xml文件移动到error文件夹成功")
						//删除抓拍xml文件
						//上传抓拍结果到云平台，早晚都会成功
						DelFile(OssAgainUpsnapxmlpathDir + "/" + fileList[i].Name())
						continue
					} else {
						//上传抓拍结果到云平台成功

						Parsed = Parsed + 1
						Parsexmlcount = Parsexmlcount + 1
						log.Println("Parsexmlcount:", Parsexmlcount)
						log.Println("ossError中的xml，第一次上传抓拍结果xml文件到云平台成功，进程抓拍结果xml移动到parsed 成功,Parsed:", Parsed, time.Now())
						DelFile(OssAgainUpsnapxmlpathDir + "/" + fileList[i].Name())
					}
					//删除本地图片 result.VehicleImgPath
					utils.DelFile(xmlresult.VehicleImgPath)
				} else {
					log.Println("ossError中的进程抓拍结果的xml，上传oss失败", code)
					//UploadFailImgCnt = UploadFailImgCnt + 1
					//UploadFailImgZeroCnt = UploadFailImgZeroCnt + 1
					//上传oss失败
					// ossError 图片不存在或者是上传oos的其他问题
					log.Println("上传oss失败，ossError中的进程抓拍结果的xml文件移动到ossError文件夹成功")
					//需要对上传oss失败的进行继续上传
					continue
				}
			}
		}
		time.Sleep(time.Minute * 2)
	}
}

//再次上传3图车型
func HandleSignalwayNewOssAgainUpload() {
	//定期检查抓拍文件夹文件夹 3TuOssError
	log.Println("HandleSignalwayNewOssAgainUpload 处理xml_suffix数据包解析以及抓拍结果再次上传")
	//2、处理文件
	//扫描 captureXml 文件夹 读取文件信息
	dir, _ := os.Getwd()
	var SignalwayNewOssAgainUpPathDir = filepath.Join(dir, "snap", "xml", "error", "3TuOssError")
	log.Println("/snap/xml/error/3TuOssError/绝对路径:", SignalwayNewOssAgainUpPathDir) //可以不需要加"/"
	for {
		// check
		if _, err := os.Stat(SignalwayNewOssAgainUpPathDir); err == nil {
			fmt.Println("path exists 1", SignalwayNewOssAgainUpPathDir)
		} else {
			log.Println("path not exists ", SignalwayNewOssAgainUpPathDir)
			err := os.MkdirAll(SignalwayNewOssAgainUpPathDir, 0711)
			if err != nil {
				log.Println("Error creating directory")
				log.Println(err)
			}
		}

		fileList, err := ioutil.ReadDir(SignalwayNewOssAgainUpPathDir) //不需要加"/"
		if err != nil {
			log.Println("扫描/snap/xml/error/3TuOssError/文件夹,读取文件信息error:", err)
			return
		}
		log.Println("扫描该/snap/xml/error/3TuOssError/文件夹下有文件的数量:", len(fileList))
		if len(fileList) == 1 {
			fmt.Println("扫描该/snap/xml/error/3TuOssError/文件夹下可能没有需要解析的xml文件") //有隐藏文件
		} else {
			if len(fileList) == 0 {
				fmt.Println("扫描该/snap/xml/error/3TuOssError/文件夹下没有需要解析的xml文件")
				time.Sleep(time.Minute * 10)
				continue
			}
		}

		for i := range fileList {
			//判断文件的结尾名+ "_suffix"
			if strings.HasSuffix(fileList[i].Name(), ".xml_suffix") {
				log.Println("扫描该/snap/xml/error/3TuOssError/文件夹下需要解析的xml文件名字为:", fileList[i].Name())
				//error/upload/fname
				content, err := ioutil.ReadFile(SignalwayNewOssAgainUpPathDir + "/" + fileList[i].Name())
				if err != nil {
					log.Println("读/3TuOssError/文件夹中文件错误信息:", err)
					continue
				}
				var xmlresult dto.CaptureDateXML
				uerr := xml.Unmarshal(content, &xmlresult)
				if uerr != nil {
					log.Println("执行HandleOssAgainUpload扫描 该/snap/xml/error/3TuOssError/文件夹下需要解析的xml文件内容时，错误信息为：", uerr)
					log.Println(string(content))
					continue
				}

				log.Println("获取抓拍结果中，图片路径result.VehicleImgPath:", xmlresult.VehicleImgPath)

				//处理新版本的信路威有三种图片   再次上传的情况
				if xmlresult.BrandName == "SignalwayNew" {
					swlNewerr := SignalwayNewOssErrorUpload(xmlresult, SignalwayNewOssAgainUpPathDir+"/"+fileList[i].Name())
					if swlNewerr != nil {
						return
					}
					return
				}
				/*
					//1、读取图片内容
					_, err = ioutil.ReadFile(xmlresult.VehicleImgPath)
					//2、判断图片是否存在
					if err != nil {
						//4、如果不存在，需要去getoss
						if fmt.Sprintf("%v", err) == "open "+xmlresult.VehicleImgPath+": no such file or directory" {
							log.Println("判断图片是否存在时，读的图片不存在:", err)
						} else {
							log.Println("判断图片是否存在时，读图片的错误信息:", err)
						}

					} else {
						//3、如果存在直接上传
						log.Println("再次上传oss时，os.Open imgfname ok:", xmlresult.VehicleImgPath)
					}

					//把图片上传到oss上
					strfname := strings.Split(xmlresult.VehicleImgPath, "\\") //windows
					//上传到oss                    日期文件夹     图片名称               前缀"/jiangsu/suhuaiyangs"
					log.Println("上传到oss图片的地址", xmlresult.VehicleImgPath)
					log.Println("上传到oss图片的名称", strfname[len(strfname)-1])
					log.Println("上传到oss的前缀", ObjectPrefix)
					//获取文件名称
					Xmlname := fileList[i].Name()
					log.Println("xml的文件路径中Xmlname:", Xmlname)
					//上传oss图片
					Stationid := ""
					if val, ok := StationId[xmlresult.CamId]; ok == true {
						Stationid = val //   string   `xml:"stationid"`//	stationid站点编号
					}
					log.Println("站点IdStationid:", Stationid)
					//前缀/站点Id/摄像机ID/日期/passid
					Pname := ObjectPrefix + "/" + Stationid + "/" + xmlresult.CamId + "/" + time.Now().Format("2006-01-02") + "/" + strfname[len(strfname)-1]
					log.Printf("前缀/站点Id/摄像机ID/日期/passid==:%s", Pname)

					code, scsj, ossDZ := utils.QingStorUpload(xmlresult.VehicleImgPath, strfname[len(strfname)-1], Pname)

					if code == utils.UPloadOK {
						UploadRecordCnt = UploadRecordCnt + 1
						UploadImgCnt = UploadImgCnt + 1
						UploadRecordZeroCnt = UploadRecordZeroCnt + 1
						UploadImgZeroCnt = UploadImgZeroCnt + 1
						OSSCount = OSSCount + 1
						log.Println("上传到oss   成功，开始返回抓拍结果给云平台")
						log.Println("上传到oss   成功，OSSCount:", OSSCount, time.Now().Format("2006-01-02 15:04:05"))

						ossDZ3 := ""
						ossDZ2 := ossDZ3
						//第一次上传失败的抓拍结果存储于【errorpathname】：snapxmlPathDir+"/error/upload/"+fileList[i].Name()
						uploaderr := GwCaptureInforUpload(&xmlresult, scsj, ossDZ, ossDZ2, ossDZ3, "./snap/xml/error/upload/"+Xmlname)
						if uploaderr != nil {
							UploadFailCnt = UploadFailCnt + 1
							UploadFailZeroCnt = UploadFailZeroCnt + 1
							log.Println("上传抓拍结果xml文件到云平台失败,xml文件移动到error文件夹成功")
							//删除抓拍xml文件
							//上传抓拍结果到云平台，早晚都会成功
							DelFile(OssAgainUpsnapxmlpathDir + "/" + fileList[i].Name())
							continue
						} else {
							//上传抓拍结果到云平台成功

							Parsed = Parsed + 1
							Parsexmlcount = Parsexmlcount + 1
							log.Println("Parsexmlcount:", Parsexmlcount)
							log.Println("ossError中的xml，第一次上传抓拍结果xml文件到云平台成功，进程抓拍结果xml移动到parsed 成功,Parsed:", Parsed, time.Now())
							DelFile(OssAgainUpsnapxmlpathDir + "/" + fileList[i].Name())
						}
						//删除本地图片 result.VehicleImgPath
						utils.DelFile(xmlresult.VehicleImgPath)
					} else {
						log.Println("ossError中的进程抓拍结果的xml，上传oss失败", code)
						//UploadFailImgCnt = UploadFailImgCnt + 1
						//UploadFailImgZeroCnt = UploadFailImgZeroCnt + 1
						//上传oss失败
						// ossError 图片不存在或者是上传oos的其他问题
						log.Println("上传oss失败，ossError中的进程抓拍结果的xml文件移动到ossError文件夹成功")
						//需要对上传oss失败的进行继续上传
						continue
					}*/
			}
		}
		time.Sleep(time.Minute * 2)
	}
}

//新版本的信路威有三种图片的再次上传  xmlnamepath=3TuOssError+"/"+fileList[i].Name()
func SignalwayNewOssErrorUpload(result dto.CaptureDateXML, xmlnamepath string) error {
	strfname := make([]string, 0)
	str2fname := make([]string, 0)
	str3fname := make([]string, 0)
	code := 0
	code2 := code
	code3 := code
	scsj := int64(0)
	scsj2 := int64(0)
	scsj3 := int64(0)
	ossDZ := ""
	ossDZ2 := ""
	ossDZ3 := ""
	Pathis := make([]string, 0)
	//把图片上传到oss上
	if result.VehicleImgPath != "" {
		strfname = strings.Split(result.VehicleImgPath, "\\") //windows
		//上传到oss                    日期文件夹     图片名称               前缀"/jiangsu/suhuaiyangs"
		log.Println("上传到oss图片1的地址，result.VehicleImgPath1:", result.VehicleImgPath)
		log.Println("上传到oss图片的名称", strfname[len(strfname)-1])
		Pathis = append(Pathis, "VehicleImgPathok")
	} else {
		code = utils.UPloadOK
	}
	if result.VehicleImgPath1 != "" {
		str2fname = strings.Split(result.VehicleImgPath1, "\\") //windows
		log.Println("上传到oss图片2的地址，result.VehicleImgPath2:", result.VehicleImgPath1)
		log.Println("上传到oss图片的名称", str2fname[len(str2fname)-1])

		Pathis = append(Pathis, "VehicleImgPath1ok")
	} else {
		code2 = utils.UPloadOK
	}
	if result.VehicleImgPath1 != "" {
		str3fname = strings.Split(result.VehicleImgPath2, "\\") //windows
		log.Println("上传到oss图片3的地址，result.VehicleImgPath3:", result.VehicleImgPath2)
		log.Println("上传到oss图片的名称", str3fname[len(str3fname)-1])

		Pathis = append(Pathis, "VehicleImgPath2ok")
	} else {
		code3 = utils.UPloadOK
	}

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

	if code == utils.UPloadOK && code2 == utils.UPloadOK && code3 == utils.UPloadOK {
		log.Println("SignalwayNew上传时，都没有图片++++++肯定有错！")
		return errors.New("SignalwayNew上传时，3图都没有图片！++++++肯定有错！")
	}
	for _, p := range Pathis {
		switch p {
		case "VehicleImgPathok":
			//log.Println("站点IdStationid:", Stationid)
			//前缀/站点Id/摄像机ID/日期/passid
			Pname := ObjectPrefix + "/" + Stationid + "/" + result.CamId + "/" + time.Now().Format("2006-01-02") + "/" + strfname[len(strfname)-1]
			log.Printf("前缀/站点Id/摄像机ID/日期/passid==:%s", Pname)
			code, scsj, ossDZ = utils.QingStorUpload(result.VehicleImgPath, strfname[len(strfname)-1], Pname)
		case "VehicleImgPath1ok":
			Pname2 := ObjectPrefix + "/" + Stationid + "/" + result.CamId + "/" + time.Now().Format("2006-01-02") + "/" + str2fname[len(str2fname)-1]
			log.Printf("前缀/站点Id/摄像机ID/日期/passid==:%s", Pname2)
			code2, scsj2, ossDZ2 = utils.QingStorUpload(result.VehicleImgPath1, str2fname[len(str2fname)-1], Pname2)
			log.Printf("第二张图片上传时间:%v", scsj2)
		case "VehicleImgPath2ok":
			Pname3 := ObjectPrefix + "/" + Stationid + "/" + result.CamId + "/" + time.Now().Format("2006-01-02") + "/" + str3fname[len(str3fname)-1]
			log.Printf("前缀/站点Id/摄像机ID/日期/passid==:%s", Pname3)
			code3, scsj3, ossDZ3 = utils.QingStorUpload(result.VehicleImgPath2, str3fname[len(str3fname)-1], Pname3)
			log.Printf("第二张图片上传时间:%v", scsj3)

		default:

		}
	}

	if scsj == 0 {
		scsj = scsj2
		if scsj == 0 {
			scsj = scsj3
		}
	}

	if code == utils.UPloadOK && code2 == utils.UPloadOK && code3 == utils.UPloadOK {
		UploadRecordCnt = UploadRecordCnt + 1
		UploadImgCnt = UploadImgCnt + 1
		UploadImgZeroCnt = UploadImgZeroCnt + 1
		UploadRecordZeroCnt = UploadRecordZeroCnt + 1
		NewOSSCount = NewOSSCount + 3
		log.Println("新版信路威上传到oss 3图都成功，开始返回抓拍结果给云平台")
		log.Println("新版信路威上传到oss 3图都成功，NewOSSCount:", NewOSSCount)
		if _, err := os.Stat("./snap/xml/error/upload/"); err == nil {
			fmt.Println("path exists 1", "./snap/xml/error/upload/")
		} else {
			log.Println("path not exists ", "./snap/xml/error/upload/")
			err := os.MkdirAll("./snap/xml/error/upload/", 0711)
			if err != nil {
				log.Println("Error creating directory")
				log.Println(err)
			}
		}
		//第一次上传失败的抓拍结果存储于【errorpathname】：snapxmlPathDir+"/error/upload/"+fileList[i].Name()
		uploaderr := GwCaptureInforUpload(&result, scsj, ossDZ, ossDZ2, ossDZ3, "./snap/xml/error/upload/"+Xmlname)
		if uploaderr != nil {
			UploadFailCnt = UploadFailCnt + 1
			UploadFailZeroCnt = UploadFailZeroCnt + 1
			log.Println("第一次上传3图抓拍结果xml文件到云平台失败，xml文件移动到error文件夹成功")
			return nil
		} else {
			//删除抓拍xml文件
			Parsed = Parsed + 1
			Parsexmlcount = Parsexmlcount + 1
			log.Println("Parsexmlcount:", Parsexmlcount)
			log.Println("第一次上传3图抓拍结果xml文件到云平台成功，进程抓拍结果xml移动到parsed 成功,Parsed:", Parsed, time.Now())
		}

		//最终都是要删除的
		DelFile(xmlnamepath)
		for _, p := range Pathis {
			switch p {
			case "VehicleImgPathok":
				//删除本地图片 result.VehicleImgPath
				utils.DelFile(result.VehicleImgPath)
			case "VehicleImgPath1ok":
				utils.DelFile(result.VehicleImgPath1)
			case "VehicleImgPath2ok":
				utils.DelFile(result.VehicleImgPath2)
			default:

			}
		}
	} else {
		log.Println("上传oss失败", code)
		UploadFailImgCnt = UploadFailImgCnt + 1
		UploadFailImgZeroCnt = UploadFailImgZeroCnt + 1
		//上传oss失败
		log.Println("上传3图的oss失败，进程抓拍结果的xml文件移动到error文件夹成功")
		return nil
	}
	return nil
}

//errorpathname：snapxmlPathDir+"/error/upload/"+fileList[i].Name()
func GwCaptureInforUpload(Result *dto.CaptureDateXML, scsj int64, ossDZ, ossDZ2, ossDZ3, errorpathname string) error {
	//判断哪一种品牌相机
	//Result.
	var ba []byte

	if Result.BrandName == "SignalwayNew" {
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
		log.Println("SignalwayNew 前置机抓拍信息上传数据", string(ba))

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
	}
	return nil
}

//创建xml文件
func createXml(xmlname string, outputxml []byte) string {

	fw, f_werr := os.Create(xmlname) //go run gwWeb.go
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
