package service

import (
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func generateConfigToOne(configdata *OneToOneConfig) string {

	//使用MarshalIndent函数，生成的XML格式有缩进
	outputxml, err := xml.MarshalIndent(*configdata, "  ", "     ")
	if err != nil {
		log.Printf("执行启动一对一进程配置文件 xml.MarshalIndent error: %v\n", err)
		return ""
	}
	xmlname := time.Now().Format("20060102T150405") + "+" + configdata.Uuid

	//20060102T150405+configdata.Uuid[HIKITS+6002].xml
	fname := createxml(xmlname, outputxml)
	if fname != "" {
		log.Println("启动一对一进程配置文件生成OK，可以启动进程，fname=", fname)
	}
	return fname
}

//宇视生成配置文件
func generateYSConfig(configdata *MoreToMoreConfig) string {
	//使用MarshalIndent函数，生成的XML格式有缩进
	outputxml, err := xml.MarshalIndent(*configdata, "     ", "     ")
	if err != nil {
		log.Printf("执行宇视生成配置文件xml.MarshalIndent error: %v\n", err)
		return ""
	}

	xmlname := time.Now().Format("20060102T150405") + "+" + configdata.Uuid

	//20060102T150405+configdata.Uuid[HIKITS+6002].xml
	fname := createxml(xmlname, outputxml)
	if fname != "" {
		log.Println("启动宇视进程配置文件生成OK，可以启动进程", fname)
	}

	return fname
}

//海康ITS生成配置文件
func generateITSConfig(configdata *OneToMoreConfig) string {
	//使用MarshalIndent函数，生成的XML格式有缩进
	outputxml, err := xml.MarshalIndent(*configdata, "     ", "     ")
	if err != nil {
		log.Printf("执行海康ITS生成配置文件 xml.MarshalIndent error: %v\n", err)
		return ""
	}

	xmlname := time.Now().Format("20060102T150405") + "+" + configdata.Uuid

	//20060102T150405+configdata.Uuid[HIKITS+6002].xml
	fname := createxml(xmlname, outputxml)
	if fname != "" {
		log.Println("启动海康ITS进程配置文件生成OK，可以启动进程", fname)
	}

	return fname
}

//创建xml文件
func createxml(xmlname string, outputxml []byte) string {

	fw, f_werr := os.Create("./cameraConfig/" + xmlname + ".xml") //go run main.go
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
	//20060102T150405+configdata.Uuid[HIKITS+6002].xml
	return xmlname + ".xml"
}
