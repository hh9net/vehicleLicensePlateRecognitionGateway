package utils

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

//处理进程互斥
func ProcessMutexBegin() {
	//读取文件
	f, err := ioutil.ReadFile("./MainMutexFile.txt")
	if err != nil {
		log.Println("扫描 snapxml 文件夹 读取文件信息 error:", err)
		//如果不存在则创建
		MainMutexFileCreate()
	} else {
		log.Println(string(f))
	}

	//如果存在则退出程序

}

//main函数中defer执行删除文件
func ProcessMutexEnd() {
	//读取文件

	//如果不存在则创建

}

func MainMutexFileCreate() {
	//用OpenFile创建一个可读可写的文件
	f, err := os.OpenFile("./MainMutexFile.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		_ = f.Close()
	}()
	n, _ := f.Seek(0, 2)
	// 从末尾的偏移量开始写入内容
	_, err = f.WriteAt([]byte("\nMainMutexFile.txt|"+time.Now().Format("2006-01-02T15:04:05")), n)
}
