package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	//ConfigInit()
	//tiker := time.NewTicker(time.Second * 10) //每15秒执行一下
	//for {
	//	log.Println("执行定时任务,模拟程序，用于检测进程管理")
	var addr string
	flag.StringVar(&addr, "addr", ":8088", "example':8087'")
	flag.Parse() //在执行这个服务时，就可以通过命令行，来指定addr的值，如果没指定，则默认是8087   go run main.go -addr ':8087'端口指定为8099

	if addr == ":8099" {
		fmt.Println("addr==8099")
	} else {
		fmt.Println("默认addr=", addr)
	}
	fmt.Println("执行定时任务,模拟程序，用于检测进程管理")
	src := "./service/data.xml"
	des := "./service/data.xml"
	utils.MoveFile(src, des)
	fmt.Println("执行定时任务已ok")
	//log.Println(utils.DateTimeFormat(<-tiker.C), "+++++++++++++++")
	//}

}
