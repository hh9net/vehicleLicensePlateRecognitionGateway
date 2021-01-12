package gatewayWeb

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"net/http"
	"vehicleLicensePlateRecognitionGateway/dto"
)

//网关基本数据
func GatewayBasicDataQuery(c *gin.Context) {
	//失败响应
	respFailure := ResponseFailure{}
	//登录处理
	code, err, data := QueryGatewayBasicData()
	if err != nil {
		respFailure.Code = code
		respFailure.Message = fmt.Sprintf("%v", err)
		c.JSON(http.StatusOK, respFailure)
		return
	}

	if code == StatusSuccessfully {
		log.Println("查询网关基本数据ok")
		c.JSON(http.StatusOK, dto.Response{Code: StatusSuccessfully, Data: *data, Message: "查询网关基本数据ok"})
		return
	}
}

//网关动态数据
func GatewayDynamicDataQuery(c *gin.Context) {
	//失败响应
	respFailure := ResponseFailure{}
	//登录处理
	code, err, data := QueryGatewayDynamicData()
	if err != nil {
		respFailure.Code = code
		respFailure.Message = fmt.Sprintf("%v", err)
		c.JSON(http.StatusOK, respFailure)
		return
	}

	//
	if code == StatusSuccessfully {
		log.Println("查询网关动态数据OK")
		c.JSON(http.StatusOK, dto.Response{Code: StatusSuccessfully, Data: *data, Message: "查询网关动态数据OK"})
		return
	}
}

// 摄像头基本信息列表查询
func CameraInfoDataQuery(c *gin.Context) {
	//失败响应
	respFailure := ResponseFailure{}
	//登录处理
	code, err, data := QueryCameraInfoData()
	if err != nil {
		respFailure.Code = code
		respFailure.Message = fmt.Sprintf("%v", err)
		c.JSON(http.StatusOK, respFailure)
		return
	}

	if code == StatusSuccessfully {
		log.Println("查询摄像头基本信息列表ok")
		c.JSON(http.StatusOK, dto.Response{Code: StatusSuccessfully, Data: *data, Message: "查询摄像头基本信息列表ok"})
		return
	}
}
