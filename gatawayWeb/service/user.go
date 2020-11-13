package service

import (
	"github.com/sirupsen/logrus"
	"vehicleLicensePlateRecognitionGateway/gatawayWeb/dto"
	"vehicleLicensePlateRecognitionGateway/gatawayWeb/types"
)

var (
	UserName string
	Password string
)

//登录
func Login(req dto.Reqlogin) (int, error) {
	logrus.Print("登录请求参数：", req)
	//获取用户数据

	if UserName != req.UserName {
		logrus.Println("用户名不正确错误")
		return types.StatusPleaseRegister, nil
	} else if Password != req.Password {
		logrus.Println("密码错误")
		return types.StatusPasswordError, nil
	}
	logrus.Println("密码正确")
	//返回数据
	return types.StatusSuccessfully, nil
}
