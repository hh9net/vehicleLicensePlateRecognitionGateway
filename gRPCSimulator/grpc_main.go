package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
	"vehicleLicensePlateRecognitionGateway/config"
	"vehicleLicensePlateRecognitionGateway/utils"
)

func ConfigInit() {
	conf := config.ConfigInit() //初始化配置文件
	log.Println("配置文件信息：", *conf)
	//初始化日志
	utils.InitLogrus(conf.LogPath, conf.LogFileName, time.Duration(24*conf.LogMaxAge)*time.Hour, conf.LogRotationSize, time.Duration(conf.LogRotationTime)*time.Hour, conf.RotationCount)
	//

}

//模拟程序，用于检测进程管理
func main() {
	var configpath string
	var addr string
	flag.StringVar(&addr, "addr", ":8088", "example':8087'") //就可以通过命令行，来指定addr的值，如果没指定，则默认是8088   go run main.go -addr ':8087'端口指定为8087
	//flag.StringVar(&configpath, "configpath", "", "example'service/data1.xml'") //go run  grpc_main.go -configpath 'path'
	flag.StringVar(&configpath, " ", "", "example'service/data1.xml'") //go run  grpc_main.go  'path'

	flag.Parse() //在执行这个服务时，

	if addr == ":8099" {
		fmt.Println("addr==8099")
	} else {
		fmt.Println("默认addr=", addr)
	}

	tiker := time.NewTicker(time.Second * 30) //每15秒执行一下
	for {
		log.Println("执行定时任务,模拟程序，用于检测进程管理")
		if configpath != "" {
			log.Println("执行定时任务,模拟程序，用于检测进程管理+++++++++++", "configpath:", configpath)
		} else {
			log.Println("configpath为空:", configpath)
		}

		generatefile(configpath)

		time.Sleep(time.Second * 5)
		fmt.Println("执行定时任务已ok")
		log.Println(utils.DateTimeFormat(<-tiker.C), "+++++++++++++++")
	}

}

func generatefile(configpath string) {
	configdata := "《改编父亲的诗》昨夜西风凌，心头自一惊。身似黄叶老，梦如碧云青。雪后空气冷，天高雁归声。酒中尘事溢，难平意绪更。" + configpath
	//使用MarshalIndent函数，生成的XML格式有缩进
	outputxml, err := xml.MarshalIndent(configdata, "  ", " ")
	if err != nil {

		log.Printf("执行线程1 打包原始记录消息包 xml.MarshalIndent error: %v\n", err)
		return
	}
	xmlname := time.Now().Format("20060102150405")
	createxml(xmlname, outputxml)
}

//创建xml文件
func createxml(xmlname string, outputxml []byte) string {

	fw, f_werr := os.Create("./grpcSimulator/xml/" + xmlname + ".xml") //go run main.go
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

	//_ = fw.Close()

	defer func() {
		_ = fw.Close()
	}()

	return "/grpcSimulator/xml/" + xmlname + ".xml"
}
