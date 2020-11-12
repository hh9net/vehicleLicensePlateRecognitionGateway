package service

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"syscall"
	"vehicleLicensePlateRecognitionGateway/dto"
)

var Deviceid string //网关设备id Token

var Token string

//1、启动进程
func Runmain() error {
	// 打印当前进程号
	fmt.Println("当前进程id：", syscall.Getpid())
	//cmd := exec.Command("../grpcSimulator/grpc_main", "test_file")
	cmd := exec.Command("./grpcSimulator/grpc_main", "-addr", ":8099")

	buf, err := cmd.Output()
	fmt.Printf("output: %s\n", string(buf))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	//执行Cmd中包含的命令，阻塞直到命令执行完成
	Runerr := cmd.Run()
	if Runerr != nil {
		log.Println("++++++ Execute Command failed. ++", "+++++++ Runerr:", Runerr.Error())
		return Runerr
	}

	log.Println("Execute Command finished.")
	return nil
}

// 进程管理
func ProcessManagementService() {

	//1、获取网关设备的token
	token, getTokenerr := GetGatawayToken()
	if getTokenerr != nil {
		log.Println("获取网关设备的token 失败")
	}
	Token = token
	//2、根据token获取camera列表
	CameraList, listerr := GetGatawayCameraList()
	if listerr != nil {
		log.Println("获取相机列表错误", listerr)
	}

	log.Println(" 相机列表数据 ：", CameraList)

	for _, cmera := range CameraList.Data {
		cmera.Description = ""

	}

A:
	//1、进程启动
	if err := Runmain(); err != nil {
		log.Println("重启")

		var a int
		//2、进程重启
		Rerr := Runmain()
		a = a + 1
		if Rerr != nil {
			log.Println("重启 error!", Rerr)

			goto A
		}
	}

}

//1、获取网关设备的token
func GetGatawayToken() (string, error) {

	Resp, err := GetToken(Deviceid)
	if err != nil {
		log.Println("GetToken error:", err)
		return "", err
	}
	log.Println(Resp.Token, err)
	return Resp.Token, nil
}

//2、根据token获取camera列表
func GetGatawayCameraList() (*dto.GetCameraList, error) {

	Resp, err := GetCameraList(Token)
	if err != nil {
		log.Println("GetToken error:", err)
		return nil, err
	}
	log.Println(Resp.Data[0].Description, err)
	return Resp, nil
}
