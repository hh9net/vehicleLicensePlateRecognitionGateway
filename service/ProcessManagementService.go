package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"syscall"
)

//1、启动进程
func Runmain() error {
	// 打印当前进程号
	fmt.Println("当前进程id", syscall.Getpid())
	//cmd := exec.Command("../grpcSimulator/grpc_main", "test_file")
	cmd := exec.Command("./grpcSimulator/grpc_main", "-addr", ":8099")

	buf, err := cmd.Output()
	fmt.Printf("output: %s\n", string(buf))
	fmt.Printf("err: %v", err)

	//执行Cmd中包含的命令，阻塞直到命令执行完成
	err = cmd.Run()
	if err != nil {
		log.Println("Execute Command failed++++++++:", "err:", err.Error())
		return err
	}

	log.Println("Execute Command finished.")
	return nil
}
