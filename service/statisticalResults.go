package service

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type StatisticalResult struct { //配置文件要通过tag来指定配置文件中的名称
	ImgAllCount    int `ini:"ImgAllCount"`    //所有图片数量
	ResultAllCount int `ini:"ResultAllCount"` //所有抓拍记录数量

	ImgDayAllCount    int `ini:"ImgDayAllCount"`    //一天的所有图片数量
	ResultDayAllCount int `ini:"ResultDayAllCount"` //一天的所有抓拍结果数量

	ImgDayUpOkCount    int `ini:"ImgDayUpOkCount"`    //一天上传OK的图片数量
	ResultDayUpOkCount int `ini:"ResultDayUpOkCount"` //一天上传OK的抓拍结果数量

	ImgDayUpErrorCount    int `ini:"ImgDayUpErrorCount"`    //一天上传失败的图片数量
	ResultDayUpErrorCount int `ini:"ResultDayUpErrorCount"` //一天上传失败的抓拍结果数量

	ImgAllCountStr    string `ini:"ImgAllCountStr"`    //所有图片数量
	ResultAllCountStr string `ini:"ResultAllCountStr"` //所有抓拍记录数量

	ImgDayAllCountStr    string `ini:"ImgDayAllCountStr"`    //一天的所有图片数量
	ResultDayAllCountStr string `ini:"ResultDayAllCountStr"` //一天的所有抓拍结果数量

	ImgDayUpOkCountStr    string `ini:"ImgDayUpOkCountStr"`    //一天上传OK的图片数量
	ResultDayUpOkCountStr string `ini:"ResultDayUpOkCountStr"` //一天上传OK的抓拍结果数量

	ImgDayUpErrorCountStr    string `ini:"ImgDayUpErrorCountStr"`    //一天上传失败的图片数量
	ResultDayUpErrorCountStr string `ini:"ResultDayUpErrorCountStr"` //一天上传失败的抓拍结果数量
}

var (
	m *sync.RWMutex

	newvalue  string
	newvalue1 string
	newvalue2 string
	newvalue3 string
	newvalue4 string
	newvalue5 string
	newvalue6 string
	newvalue7 string
)

//注意读取参数时，要做加锁处理
func StatisticalResults(ImgAllCount, ResultAllCount, ImgDayAllCount, ResultDayAllCount, ImgDayUpOkCount, ResultDayUpOkCount, ImgDayUpErrorCount, ResultDayUpErrorCount int) {
	dir, _ := os.Getwd()
	log.Println("当前路径：", dir)
	var statisticalResultsfnamePath = filepath.Join(dir, "statisticalResults")
	// check
	if _, err := os.Stat(statisticalResultsfnamePath); err == nil {
		log.Println("path exists 1", statisticalResultsfnamePath)
	} else {
		log.Println("path not exists ", statisticalResultsfnamePath)
		err := os.MkdirAll(statisticalResultsfnamePath, 0711)
		if err != nil {
			log.Println("Error creating directory")
			log.Println(err)
			return
		}
	}
	m = new(sync.RWMutex)

	cfg, err := ini.Load("../statisticalResults/statisticalResults.txt") //读配置文件  goland不能使用 ./ 方式  go run main.go 可以
	if err != nil {
		log.Print("Fail to read file:", err)
		return
	}
	m.Lock()
	SET := 0
	//所有图片数量
	ImgAllCountvalue, ImgAllCountvalueErr := GetValue(cfg, "ImgAllCount") //取值
	if ImgAllCountvalueErr != nil {
		log.Println("GetValue error:", ImgAllCountvalueErr)
		return
	}
	log.Println("read File ImgAllCountvalue:", ImgAllCountvalue) //取值
	//使用值，比较数据
	if ImgAllCountvalue == ImgAllCount {
		SET = SET + 1
	} else {

	}
	ImgAllCountvalue = ImgAllCount
	newvalue = strconv.Itoa(ImgAllCountvalue)

	//所有抓拍记录数量
	ResultAllCountvalue, ResultAllCountErr := GetValue(cfg, "ResultAllCount") //取值
	if ResultAllCountErr != nil {
		log.Println("GetValue error:", ResultAllCountErr)
		return
	}
	log.Println("read File ResultAllCountvalue:", ResultAllCountvalue) //取值
	//使用值，比较数据
	if ResultAllCountvalue == ResultAllCount {
		SET = SET + 1
	} else {

	}
	ResultAllCountvalue = ResultAllCount
	newvalue1 = strconv.Itoa(ResultAllCountvalue)

	//一天的所有图片数量
	ImgDayAllCountvalue, ImgDayAllCountvalueErr := GetValue(cfg, "ImgDayAllCount") //取值
	if ImgDayAllCountvalueErr != nil {
		log.Println("GetValue error:", ImgDayAllCountvalueErr)
		return
	}
	log.Println("read File ImgDayAllCountvalue:", ImgDayAllCountvalue) //取值
	//使用值，比较数据
	if ImgDayAllCountvalue == ImgDayAllCount {
		SET = SET + 1
	} else {

	}
	ImgDayAllCountvalue = ImgDayAllCount
	newvalue2 = strconv.Itoa(ImgDayAllCountvalue)

	//一天的所有抓拍结果数量
	ResultDayAllCountvalue, ResultDayAllCountvalueErr := GetValue(cfg, "ResultDayAllCount") //取值
	if ResultDayAllCountvalueErr != nil {
		log.Println("GetValue error:", ResultDayAllCountvalueErr)
		return
	}
	log.Println("read File ResultDayAllCountvalue:", ResultDayAllCountvalue) //取值
	//使用值，比较数据
	if ResultDayAllCountvalue == ResultDayAllCount {
		SET = SET + 1
	} else {

	}
	ResultDayAllCountvalue = ResultDayAllCount
	newvalue3 = strconv.Itoa(ResultDayAllCountvalue)

	//一天上传OK的图片数量
	ImgDayUpOkCountvalue, ImgDayUpOkCountvalueErr := GetValue(cfg, "ImgDayUpOkCount") //取值
	if ImgDayUpOkCountvalueErr != nil {
		log.Println("GetValue error:", ImgDayUpOkCountvalueErr)
		return
	}
	log.Println("read File ImgDayUpOkCountvalue:", ImgDayUpOkCountvalue) //取值
	//使用值，比较数据
	if ImgDayUpOkCountvalue == ImgDayUpOkCount {
		SET = SET + 1
	} else {

	}
	ImgDayUpOkCountvalue = ImgDayUpOkCount
	newvalue4 = strconv.Itoa(ImgDayUpOkCountvalue)

	//一天上传OK的抓拍结果数量
	ResultDayUpOkCountvalue, ResultDayUpOkCountvalueErr := GetValue(cfg, "ResultDayUpOkCount") //取值
	if ResultDayUpOkCountvalueErr != nil {
		log.Println("GetValue error:", ResultDayUpOkCountvalueErr)
		return
	}
	log.Println("read File ResultDayUpOkCountvalue:", ResultDayUpOkCountvalue) //取值
	//使用值，比较数据
	if ResultDayUpOkCountvalue == ResultDayUpOkCount {
		SET = SET + 1
	} else {

	}
	ResultDayUpOkCountvalue = ResultDayUpOkCount
	newvalue5 = strconv.Itoa(ResultDayUpOkCountvalue)

	//一天上传失败的图片数量
	ImgDayUpErrorCountvalue, ImgDayUpErrorCountvalueErr := GetValue(cfg, "ImgDayUpErrorCount") //取值
	if ImgDayUpErrorCountvalueErr != nil {
		log.Println("GetValue error:", ImgDayUpErrorCountvalueErr)
		return
	}
	log.Println("read File ImgDayUpErrorCountvalue:", ImgDayUpErrorCountvalue) //取值
	//使用值，比较数据
	if ImgDayUpErrorCountvalue == ImgDayUpErrorCount {
		SET = SET + 1
	} else {

	}
	ImgDayUpErrorCountvalue = ImgDayUpErrorCount
	newvalue6 = strconv.Itoa(ImgDayUpErrorCountvalue)

	//一天上传失败的抓拍结果数量
	ResultDayUpErrorCountvalue, ResultDayUpErrorCountvalueErr := GetValue(cfg, "ResultDayUpErrorCount") //取值
	if ResultDayUpErrorCountvalueErr != nil {
		log.Println("GetValue error:", ResultDayUpErrorCountvalueErr)
		return
	}
	log.Println("read File ImgAllCount:", ResultDayUpErrorCountvalue) //取值
	//使用值，比较数据
	if ResultDayUpErrorCountvalue == ResultDayUpErrorCount {
		SET = SET + 1
	} else {

	}
	ResultDayUpErrorCountvalue = ResultDayUpErrorCount
	newvalue7 = strconv.Itoa(ResultDayUpErrorCountvalue)

	if SET == 8 {
		m.Unlock()
		return
	}

	//新值存文件

	cfg.Section("").Key("ImgAllCount").SetValue(newvalue)                 //  修改后值然后进行保存
	cfg.Section("").Key("ResultAllCount").SetValue(newvalue1)             //  修改后值然后进行保存
	cfg.Section("").Key("ImgDayAllCount").SetValue(newvalue2)             //  修改后值然后进行保存
	cfg.Section("").Key("ResultDayAllCount").SetValue(newvalue3)          //  修改后值然后进行保存
	cfg.Section("").Key("ImgDayUpOkCount").SetValue(newvalue4)            //  修改后值然后进行保存
	cfg.Section("").Key("ResultDayUpOkCount").SetValue(newvalue5)         //  修改后值然后进行保存
	cfg.Section("").Key("ImgDayUpErrorCount").SetValue(newvalue6)         //  修改后值然后进行保存
	cfg.Section("").Key("ResultDayUpErrorCount").SetValue(newvalue7)      //  修改后值然后进行保存
	cfg.Section("").Key("ImgAllCountStr").SetValue(newvalue)              //  修改后值然后进行保存
	cfg.Section("").Key("ResultAllCountStr").SetValue(newvalue)           //  修改后值然后进行保存
	cfg.Section("").Key("ImgDayAllCountStr").SetValue(newvalue)           //  修改后值然后进行保存
	cfg.Section("").Key("ResultDayAllCountStr").SetValue(newvalue)        //  修改后值然后进行保存
	cfg.Section("").Key("ImgDayUpOkCountStr").SetValue(newvalue)          //  修改后值然后进行保存
	cfg.Section("").Key("ResultDayUpOkCountStr").SetValue(newvalue)       //  修改后值然后进行保存
	cfg.Section("").Key("ImgDayUpErrorCountStr").SetValue(newvalue)       //  修改后值然后进行保存
	cfg.Section("").Key("ResultDayUpErrorCountStr").SetValue(newvalue)    //  修改后值然后进行保存
	Saveerr := cfg.SaveTo("../statisticalResults/statisticalResults.txt") //读配置文件  goland不能使用 ./ 方式  go run main.go 可以
	m.Unlock()
	if Saveerr != nil {
		log.Print("Fail to SaveTo file:", Saveerr)
		return
	}
	log.Println("修改后保存的值: ", newvalue) //  修改后保存的值
}

func GetValue(cfg *ini.File, key string) (int, error) {
	return cfg.Section("").Key(key).Int()
}

func GetStrValue(cfg *ini.File, key string) string {
	return cfg.Section("").Key(key).String()
}

func StatisticalFile(val string) {
	//用OpenFile创建一个可读可写的文件
	f, err := os.OpenFile("../statisticalResults/statisticalResultsFile.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		_ = f.Close()
	}()

	//写入字符串
	_, err = f.WriteString(val)
	//n,err:=f.WriteString("hello word,你好世界，I love you")
	if err != nil {
		log.Println(err)
	}
}
