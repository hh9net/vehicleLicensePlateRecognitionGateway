package main

import (
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	"vehicleLicensePlateRecognitionGateway/grpcSimulator/config"
	"vehicleLicensePlateRecognitionGateway/utils"
)

func ConfigInit() {
	conf := config.UdpConfigInit() //初始化配置文件
	log.Println("配置文件信息：", *conf)
	//初始化日志
	utils.InitLogrus(conf.LogPath, conf.LogFileName, time.Duration(24*conf.LogMaxAge)*time.Hour, conf.LogRotationSize, time.Duration(conf.LogRotationTime)*time.Hour, conf.RotationCount)
	//

}

//模拟程序，用于检测进程管理
func main() {
	ConfigInit()
	// 打印当前进程号

	log.Println("当前进程id：", syscall.Getpid())
	log.Println(len(os.Args))
	if len(os.Args) == 1 {
		log.Println("请输入配置文件的绝对路径++++++++++++os.Args:", os.Args)
	}
	if len(os.Args) == 2 {
		log.Println("len(os.Args) == 2 +++++++++++=os.Args:", os.Args)
		if os.Args[0] != "" {
			//cameraConfig/20201205T165922|sxjgl_shygs_321300_G2513_K101_415_3_2_1|6008.xml
			if strings.HasSuffix(os.Args[1], ".xml") {

				pxml := strings.Split(os.Args[1], "+")

				p := strings.Split(pxml[2], ".")
				log.Println("os.Args[1] port:", p)
				//go run udpmain.go cameraConfig/config3.xml
				s := "./cameraConfig/" + os.Args[1]
				d := "./cameraConfig/rm-rf/" + os.Args[1]
				log.Println("s:", s, "d:", d)
				frerr := FileRename(s, d)
				if frerr != nil {

					return
				}

				//发送给管理平台的心跳
				go Heartbeat(p[0])

				tiker := time.NewTicker(time.Second * 30) //每15秒执行一下
				for {
					logfile := ""
					logfile = logfile + os.Args[1]
					log.Println("执行定时任务,模拟程序，用于检测进程管理")
					//if configpath != "" {
					//	log.Println("执行定时任务,模拟程序，用于检测进程管理+++++++++++", "configpath:", configpath)
					//} else {
					//	logfile = logfile + "|configpath为空"
					//	log.Println("configpath为空:", configpath)
					//}

					generatefile(os.Args[1], "")

					log.Println(utils.DateTimeFormat(<-tiker.C), "执行定时任务已ok+++++++++++++++")
				}
			} else {
				log.Println("请输入正确的启动进程的配置文件")
			}

		} else {
			log.Println("os.Args 为空", os.Args)

		}

	}

}

func generatefile(configpath, logfile string) {
	log.Println("《改编的诗》昨夜西风凌，心头自一惊。身似黄叶老，梦如碧云青。")
	configdata := "《改编的诗》昨夜西风凌，心头自一惊。身似黄叶老，梦如碧云青。雪后空气冷，天高雁归声。酒中尘事溢，难平意绪更。configpath:" + configpath + logfile
	//使用MarshalIndent函数，生成的XML格式有缩进
	outputxml, err := xml.MarshalIndent(configdata, "  ", " ")
	if err != nil {
		log.Printf("执行线程1 打包原始记录消息包 xml.MarshalIndent error: %v\n", err)
		return
	}
	port := strings.Split(configpath, "+")
	xmlname := time.Now().Format("20060102150405") + "+" + port[2]

	createxml(xmlname, outputxml)
}

//创建xml文件
func createxml(xmlname string, outputxml []byte) string {

	fw, f_werr := os.Create("./grpcSimulator/xml/" + xmlname) //go run main.go
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

	return "grpcSimulator/xml/" + xmlname
}

//与抓拍进程管理中心交互心跳,  得知程序死活
func Heartbeat(port string) {
	zpxtp, _ := strconv.Atoi(port)
	//监听管理中心的心跳192.168.26.248
	go Heartbeatserver("192.168.150.164" + ":" + strconv.Itoa(zpxtp-1))

DGLPT:
	//给管理平台发送心跳
	serverAddr := "192.168.150.164" + ":" + port
	conn, err := net.Dial("udp", serverAddr)
	//checkError(err)
	if err != nil {
		time.Sleep(time.Second * 60)
		log.Println("UDP net.Dial err:", err, serverAddr)
		goto DGLPT
	}

	log.Println("UDP net.Dial serverAddr:", serverAddr)

	defer func() {
		_ = conn.Close()
	}()

	tiker := time.NewTicker(time.Second * 5) //每15秒执行一下
	for {
		var n int
		var toWrite string
		toWrite = "抓拍进程启动，端口号为：" + port

		n, err = conn.Write([]byte(toWrite))
		if err != nil {
			log.Println("抓拍进程 conn.Write error：", err)
			log.Println("给管理抓拍进程发送信息时错误，err:", err)
			//	checkError(err)
		}

		log.Println("抓拍进程发送的 Write:", toWrite, "｜n:", n)

		msg := make([]byte, 32)
		n, err = conn.Read(msg)
		//
		if err != nil {

			log.Println("抓拍进程 conn.Read error：", err)
		}

		log.Println("Response:", string(msg), "n:", n)
		log.Println(utils.DateTimeFormat(<-tiker.C), "+++++++++++++++")
	}
}

func Heartbeatserver(address string) {
ListenUDP:
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Println("抓拍进程监听管理平台 net.ResolveUDPAddr error:", err)
		time.Sleep(time.Second * 30)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println("抓拍进程 UDP监听 address:", address)
		log.Println("抓拍进程监听 端口 net.ListenUDP:", err)
		time.Sleep(time.Second * 30)
		goto ListenUDP
	}
	log.Println("抓拍进程 UDP监听 address:", address)

	defer func() {
		_ = conn.Close()
	}()

	for {
		//获取数据
		// Here must use make and give the lenth of buffer
		data := make([]byte, 32)

		//返回一个UDPAddr        ReadFromUDP从c读取一个UDP数据包，将有效负载拷贝到b，返回拷贝字节数和数据包来源地址。
		//ReadFromUDP方法会在超过一个固定的时间点之后超时，并返回一个错误。
		_, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println(err)
			continue
		}

		strData := address + "抓拍进程收到了你的信息！你的信息是：" + string(data)
		log.Println("Received:", strData)
		//转大写
		//	upper := strings.ToUpper(strData)

		_, err = conn.WriteToUDP([]byte(strData), rAddr)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("Send:", strData)

	}

}

func checkError(err error) {
	if err != nil {
		log.Println(err)
		//return
		os.Exit(1)
	}
}

//改文件名字
func FileRename(src string, des string) error {
	//err := os.Rename("./a", "/tmp/a")
	err := os.Rename(src, des)
	if err != nil {
		log.Println("移动文件错误：", err)
		return err
	}
	log.Printf("改文件名字%s to %s 成功", src, des)
	return nil
}
