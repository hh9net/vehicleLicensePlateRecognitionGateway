package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"vehicleLicensePlateRecognitionGateway/dto"
	"vehicleLicensePlateRecognitionGateway/service"
	"vehicleLicensePlateRecognitionGateway/types"
)

//@Summary 登录api
//@Tags 登录
//@version 1.0
//@Accept application/json
//@Param req body dto.Reqlogin true "请求参数"
//@Success 200 object dto.Response 成功后返回值
//@Failure 404 object dto.ResponseFailure 查询失败
//@Router /user/login [post]
func Login(c *gin.Context) {
	req := dto.Reqlogin{}
	respFailure := dto.ResponseFailure{}

	if err := c.Bind(&req); err != nil {
		logrus.Println("Login json unmarshal err: %v", err)
		respFailure.Code = -1
		respFailure.Message = fmt.Sprintf("json unmarshal err: %v", err)
	}
	//登录处理
	code, err := service.Login(req)
	if err != nil {
		respFailure.Code = code
		respFailure.Message = fmt.Sprintf("登录时error:%v", err)
		logrus.Println("登录时error：", respFailure.Message)
		c.JSON(http.StatusOK, dto.Response{Code: respFailure.Code, Data: "登录时error", Message: "登录时error"})
		return
	}

	if code == types.StatusPleaseRegister {
		logrus.Println("用户名输入错误,请重新输入")
		c.JSON(http.StatusOK, dto.Response{Code: types.StatusPleaseRegister, Data: types.StatusText(types.StatusPleaseRegister), Message: "用户名输入错误,请重新输入"})
		return
	}

	if code == types.StatusPasswordError {
		logrus.Println("密码错误,请重新输入")
		c.JSON(http.StatusOK, dto.Response{Code: types.StatusPasswordError, Data: types.StatusText(types.StatusPasswordError), Message: "密码错误,请重新输入"})
		return
	}

	if code == types.StatusSuccessfully {
		logrus.Println("用户登录成功")
		c.JSON(http.StatusOK, dto.Response{Code: types.StatusSuccessfully, Data: types.StatusText(types.StatusSuccessfully), Message: "用户登录成功"})
		return
	}
}
