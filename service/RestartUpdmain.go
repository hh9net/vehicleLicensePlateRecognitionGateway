package service

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func RestartUpdmain(port string) error {
	c := exec.Command("taskkill", "/f", "/pid", Pid[port])
	err := c.Start()
	if err != nil {
		log.Println("Kill pid error:", Pid[port])
		return err
	}
	log.Println("Kill pid ok:", Pid[port])

	dir, _ := os.Getwd()
	log.Println("RestartUpdmain时，当前路径：", dir)
	var cameraConfigxmlPathDir = filepath.Join(dir, "cameraConfig")
	//提取文件夹下文件
	fileList, err := ioutil.ReadDir(cameraConfigxmlPathDir)
	if err != nil {
		log.Println("扫描cameraConfig文件夹 读取文件信息 error:", err)
		time.Sleep(time.Second * 1)
		return err
	}
	log.Println("扫描该cameraConfig/文件夹下有文件的数量 ：", len(fileList))

	if len(fileList) == 1 {
		fmt.Println("该cameraConfig/文件夹下可能没有需要解析的xml文件") //有隐藏文件
		time.Sleep(time.Second * 1)
	} else {
		if len(fileList) == 0 {
			fmt.Println("扫描cameraConfig/文件夹下没有需要解析的xml文件")
			time.Sleep(time.Second * 1)
			return errors.New("扫描cameraConfig/文件夹下没有需要重启的的xml文件")
		}
	}

	for i := range fileList {
		//判断文件的结尾名
		if strings.Contains(fileList[i].Name(), "+"+port+".xml") {
			log.Println("扫描该cameraConfig/文件夹下需要解析的xml文件名字为:", fileList[i].Name())
			go Runmain(fileList[i].Name())
		}
	}
	return nil
}

func DelcameraConfigDir() {
	dir, _ := os.Getwd()
	fmt.Println("当前路径：", dir)
	var cameraConfigxmlPathDir = filepath.Join(dir, "cameraConfig")
	// check
	if _, err := os.Stat(cameraConfigxmlPathDir); err == nil {
		fmt.Println("path exists 1", cameraConfigxmlPathDir)

		RemoveAllerr := os.RemoveAll(cameraConfigxmlPathDir)
		if RemoveAllerr != nil {
			log.Println("cameraConfigxmlPathDir RemoveAllerr", RemoveAllerr)
		} else {
			fmt.Println("path exists 1", cameraConfigxmlPathDir, "delete ok!")

		}
	} else {
		fmt.Println("path not exists ", cameraConfigxmlPathDir)
		err := os.MkdirAll(cameraConfigxmlPathDir, 0711)
		if err != nil {
			log.Println("Error creating directory")
			log.Println(err)
		}
	}
}
