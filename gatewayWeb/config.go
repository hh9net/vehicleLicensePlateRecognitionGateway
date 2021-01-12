package gatewayWeb

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var conffilepath = "./conf/config.toml" // go run gwWeb.go
//var conffilepath = "../conf/config.toml"

type Config struct { //配置文件要通过tag来指定配置文件中的名称
	//日志
	WebLogPath         string `ini:"weblog_Path"`
	WebLogMaxAge       int64  `ini:"weblog_maxAge"`
	WebLogRotationTime int64  `ini:"weblog_rotationTime"` //日志切割时间间隔（小时）
	WebLogRotationSize int64  `ini:"weblog_rotationSize"` //日志切割大小（1024 KB）
	WebRotationCount   uint   `ini:"weblog_rotationCount"`
	WebLogFileName     string `ini:"weblog_FileName"`

	//外网id
	IpAddress string `ini:"ip_address"`

	UserName string `ini:"user_name"`
	Password string `ini:"password"`
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
	//读配置文件
	config, err := ReadConfig(conffilepath) //也可以通过os.arg或flag从命令行指定配置文件路径
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(config)
	return &config
}
