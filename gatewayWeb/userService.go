package gatewayWeb

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"net/http"
	"vehicleLicensePlateRecognitionGateway/dto"
)

var (
	UserName string
	Password string
)

func Login(c *gin.Context) {
	req := Reqlogin{}
	respFailure := ResponseFailure{}

	if err := c.Bind(&req); err != nil {
		log.Errorf("Login json unmarshal err: %v", err)
		respFailure.Code = -1
		respFailure.Message = fmt.Sprintf("json unmarshal err: %v", err)
		c.JSON(http.StatusOK, respFailure)
		return
	}
	//登录处理
	code, err := Loginserver(req)
	if err != nil {
		respFailure.Code = code
		respFailure.Message = fmt.Sprintf("%v", err)
		//c.JSON(http.StatusOK, respFailure)
		//return
	}

	if code == StatusPleaseRegister {
		log.Println("用户名输入错误")
		c.JSON(http.StatusOK, dto.Response{Code: StatusPleaseRegister, Data: StatusText(StatusPleaseRegister), Message: "用户名输入错误,请重新输入"})
		return
	}

	if code == StatusPasswordError {
		log.Println("密码错误,请重新输入")
		c.JSON(http.StatusOK, dto.Response{Code: StatusPasswordError, Data: StatusText(StatusPasswordError), Message: "密码错误,请重新输入"})
		return
	}

	if code == StatusNoVerificationcode {
		log.Println("验证码错误,请重新输入")
		c.JSON(http.StatusOK, dto.Response{Code: StatusNoVerificationcode, Data: StatusText(StatusNoVerificationcode), Message: "验证码错误,请重新输入"})
		return
	}

	if code == StatusSuccessfully {
		log.Println("用户登录成功")
		c.JSON(http.StatusOK, dto.Response{Code: StatusSuccessfully, Data: "Successful", Message: "用户登录成功"})
		return
	}
}

//登录
func Loginserver(req Reqlogin) (int, error) {
	log.Print("登录请求参数：", req)
	//获取用户数据

	if UserName != req.UserName {
		log.Println("用户名不正确错误")
		return StatusPleaseRegister, nil
	} else if Password != req.Password {
		log.Println("密码错误")
		return StatusPasswordError, nil
	}
	log.Println("密码正确")
	//返回数据
	return StatusSuccessfully, nil
}
