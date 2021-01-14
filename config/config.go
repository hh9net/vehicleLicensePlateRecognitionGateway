package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
)

var conffilepath = "./conf/config.toml"      // go run gwWeb.go
var confGofilepath = "./conf/config_go.toml" // go run gwWeb.go
//var conffilepath = "../conf/config_c.toml"

type Config struct { //配置文件要通过tag来指定配置文件中的名称
	//日志
	LogPath         string `ini:"log_Path"`
	LogMaxAge       int64  `ini:"log_maxAge"`
	LogRotationTime int64  `ini:"log_rotationTime"` //日志切割时间间隔（小时）
	LogRotationSize int64  `ini:"log_rotationSize"` //日志切割大小（1024 KB）
	RotationCount   uint   `ini:"log_rotationCount"`

	LogFileName string `ini:"log_FileName"`

	//外网id
	IpAddress string `ini:"ip_address"`

	//前置机抓拍信息上传接口
	GwCaptureInformationUploadIpAddress string `ini:"gwCaptureInformationUpload_ip_address"`
	Gettoken                            string `ini:"get_token"`
	GetCameraList                       string `ini:"get_camera_list"`
	//网关设备id
	Deviceid string `ini:"deviceid"`
	//网关位置
	Gatewaylocation string `ini:"gatewaylocation"`

	//
	StatisticalReportIpAddress string `ini:"statistical_report_ipAddress"`
}

//读取配置文件并转成结构体
func ReadConfig(path string) (Config, error) {
	var config Config
	conf, err := ini.Load(path) //加载配置文件
	if err != nil {
		log.Println("load config file fail!")
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config) //解析成结构体
	if err != nil {
		log.Println("mapto config file fail!")
		return config, err
	}
	return config, nil
}

//获取mysql 配置文件信息
func ConfigInit() *Config {

	//生成新配置文件
	ConfigNewFile()
	//读配置文件
	config, err := ReadConfig(confGofilepath) //也可以通过os.arg或flag从命令行指定配置文件路径
	if err != nil {
		log.Println(err)
	}
	//log.Println(config)
	return &config
}

func ConfigNewFile() {
	//读配置文件
	content, err := ioutil.ReadFile(conffilepath)
	if err != nil {
		if fmt.Sprintf("%v", err) == "open "+conffilepath+": no such file or directory" {
			log.Println("生成新配置文件时，读配置的文件不存在:", err)
		} else {
			log.Println("生成新配置文件时，读配置文件错误信息:", err)
		}

		return
	}

	//创建一个目标文件，存储源文件数据，完成拷贝
	f_w, err := os.Create("./conf/config_go.toml")
	if err != nil {
		log.Println(err)
	}

	defer func() {
		_ = f_w.Close()
	}()

	//创建一个缓冲区buf 用于存储源文件的数据，数据为字节切片类型[]byte
	_, err = f_w.Write(content[16:])
	if err != nil {
		log.Println(err)
	}

	log.Println("生成新配置文件 ./conf/config_go.toml 成功！")
}
