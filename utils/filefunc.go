package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func MoveFile(src string, des string) error {
	//err := os.Rename("./a", "/tmp/a")
	err := os.Rename(src, des)
	if err != nil {
		log.Println("移动文件错误:", err)
		return err
	}
	log.Printf("移动文件：%s to： %s 成功", src, des)
	return nil
}

//删除文件
func DelFile(src string) {
	//"./1.txt"
	del := os.Remove(src)
	if del != nil {
		log.Println("删除失败:", del)
		return
	}
	log.Println("删除图片ok", src)
}

//删除指定path下的所有文件
func DelAllFile(src string) {
	//"./dir"
	delDir := os.RemoveAll(src)
	if delDir != nil {
		log.Println(delDir)
	}
}
