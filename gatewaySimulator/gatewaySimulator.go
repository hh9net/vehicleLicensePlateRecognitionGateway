package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"strings"
	"vehicleLicensePlateRecognitionGateway/dto"
)

//模拟云端网关与本地网关的交互

func main() {
	//192.168.26.248
	IpAddress := "192.168.26.248:9898"
	logrus.Print("服务端 IpAddress：", IpAddress)
	router := gin.New()
	router.Use(Cors()) //跨域资源共享
	apiV1 := router.Group("/gw/api/v1")
	APIV1Init(apiV1)

	http.Handle("/", router)
	gin.SetMode(gin.ReleaseMode)

	runerr := router.Run(IpAddress)
	if runerr != nil {
		logrus.Print("Run error", runerr)
		return
	}
}
func APIV1Init(route *gin.RouterGroup) {
	AuthAPIInit(route)
}

func AuthAPIInit(route *gin.RouterGroup) {
	//上传抓拍信息
	route.POST("/upload", Upload)
	//获取token
	route.GET("/getToken/:name", GetToken)
	//获取相机列表
	route.GET("/getList/:name", GetList)

}

//以下为cors实现
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*") // 这是允许访问所有域

			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段

			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析

			c.Header("Access-Control-Max-Age", "172800")          // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false") //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")             // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

//1、回复token
func GetToken(c *gin.Context) {
	id := c.Params.ByName("name")
	log.Println("id", id)
	if id == "f9555afb-56c0-4bf8-a067-63b8ca4be538" {
		data := new(dto.GetTokenRespXML)
		data.Token = "a2caedfcb22b21bedafe"
		data.Expire = "2h"
		data.Code = 0
		data.Msg = "请求成功"
		data.Oss.BacketName = "ydcpsbxt"
		data.Oss.ObjectPrefix = "jiangsu/suhuaiyangs"
		c.XML(200, data)
	}
}

//2、回复相机列表
func GetList(c *gin.Context) {

	token := c.Params.ByName("name")
	log.Println("token", token)
	if token == "a2caedfcb22b21bedafe" {
		datas := new(dto.GetCameraList)

		d := `
<ListData>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_1_0</id>
   <name>南区入口(卡口)</name>
   <stationId>9c667aef8bb64a1e99ba328e76cb1a65</stationId>
   <laneType>0</laneType>
   <devCompId>UNIVIEW</devCompId>
   <description>2</description>
   <devIp>10.113.1.37</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>0</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_1_0</id>
   <name>南区入口(卡口)</name>
   <stationId>9c667aef8bb64a1e99ba328e76cb1a65</stationId>
   <laneType>0</laneType>
   <devCompId>UNIVIEW</devCompId>
   <description>2</description>
   <devIp>10.113.1.37</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>0</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_1_0</id>
   <name>南区入口(卡口)</name>
   <stationId>9c667aef8bb64a1e99ba328e76cb1a65</stationId>
   <laneType>0</laneType>
   <devCompId>HIKITS</devCompId>
   <description>2</description>
   <devIp>10.113.1.37</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>1</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_1_1</id>
   <name>南区入口(卡口)</name>
   <stationId>9c667aef8bb64a1e99ba328e76cb1a65</stationId>
   <laneType>0</laneType>
   <devCompId>HIKITS</devCompId>
   <description>2</description>
   <devIp>10.113.1.37</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>2</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_1_2</id>
   <name>南区入口(卡口)</name>
   <stationId>9c667aef8bb64a1e99ba328e76cb1a65</stationId>
   <laneType>0</laneType>
   <devCompId>HIKITS</devCompId>
   <description>2</description>
   <devIp>10.113.1.36</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>2</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_1_0</id>
   <name>南区入口(卡口)</name>
   <stationId>9c667aef8bb64a1e99ba328e76cb1a65</stationId>
   <laneType>0</laneType>
   <devCompId>HIK</devCompId>
   <description>2</description>
   <devIp>10.113.1.37</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>0</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_1_1</id>
   <name>南区出口(卡口)</name>
   <stationId>9c667aef8bb64a1e99ba328e76cb1a65</stationId>
   <laneType>1</laneType>
   <devCompId>HIK</devCompId>
   <description>2</description>
   <devIp>10.113.1.36</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>0</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_2_0</id>
   <name>北区入口(卡口)</name>
   <stationId>20ebf41475174ff7a8ed46fc902aa3a4</stationId>
   <laneType>0</laneType>
   <devCompId>HIK</devCompId>
   <description>2</description>
   <devIp>10.113.1.16</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>0</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
 <Data>
   <id>sxjgl_shygs_321300_G2513_K101_415_3_2_1</id>
   <name>北区出口(卡口)</name>
   <stationId>20ebf41475174ff7a8ed46fc902aa3a4</stationId>
   <laneType>1</laneType>
   <devCompId>DEYA</devCompId>
   <description>2</description>
   <devIp>10.113.1.17</devIp>
   <port>80</port>
   <userName>admin</userName>
   <password>123456</password>
   <channel>0</channel>
   <laneNo></laneNo>
   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
 </Data>
</ListData>
`
		//		d1 := `
		//<ListData>
		// <Data>
		//   <id>sxjgl_shygs_321300_G2513_K101_415_3_1_0</id>
		//   <name>南区入口(卡口)</name>
		//   <stationId>9c667aef8bb64a1e99ba328e76cb1a65</stationId>
		//   <laneType>0</laneType>
		//   <devCompId>UNIVIEW</devCompId>
		//   <description>2</description>
		//   <devIp>10.113.1.37</devIp>
		//   <port>80</port>
		//   <userName>admin</userName>
		//   <password>123456</password>
		//   <channel>0</channel>
		//   <laneNo></laneNo>
		//   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
		// </Data>
		// <Data>
		//   <id>sxjgl_shygs_321300_G2513_K101_415_3_1_1</id>
		//   <name>南区出口(卡口)</name>
		//   <stationId>9c667aef8bb64a1e99ba328e76cb1a65</stationId>
		//   <laneType>1</laneType>
		//   <devCompId>UNIVIEW</devCompId>
		//   <description>2</description>
		//   <devIp>10.113.1.36</devIp>
		//   <port>80</port>
		//   <userName>admin</userName>
		//   <password>123456</password>
		//   <channel>0</channel>
		//   <laneNo></laneNo>
		//   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
		// </Data>
		// <Data>
		//   <id>sxjgl_shygs_321300_G2513_K101_415_3_2_0</id>
		//   <name>北区入口(卡口)</name>
		//   <stationId>20ebf41475174ff7a8ed46fc902aa3a4</stationId>
		//   <laneType>0</laneType>
		//   <devCompId>UNIVIEW</devCompId>
		//   <description>2</description>
		//   <devIp>10.113.1.16</devIp>
		//   <port>80</port>
		//   <userName>admin</userName>
		//   <password>123456</password>
		//   <channel>0</channel>
		//   <laneNo></laneNo>
		//   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
		// </Data>
		// <Data>
		//   <id>sxjgl_shygs_321300_G2513_K101_415_3_2_1</id>
		//   <name>北区出口(卡口)</name>
		//   <stationId>20ebf41475174ff7a8ed46fc902aa3a4</stationId>
		//   <laneType>1</laneType>
		//   <devCompId>UNIVIEW</devCompId>
		//   <description>2</description>
		//   <devIp>10.113.1.17</devIp>
		//   <port>80</port>
		//   <userName>admin</userName>
		//   <password>123456</password>
		//   <channel>0</channel>
		//   <laneNo></laneNo>
		//   <gantryId>157c8013-bcd3-4fde-a548-b8f6473862b2</gantryId>
		// </Data>
		//</ListData>
		//`

		uerr := xml.Unmarshal([]byte(d), datas)
		if uerr != nil {
			log.Println("执行线程4 解析 receivexml文件夹中xml文件内容时，错误信息为：", uerr)
		}

		c.XML(200, datas)
	}

}

//3、上传抓拍信息
func Upload(c *gin.Context) {
	//获取抓拍结果
	req := dto.DateXML{}
	if err := c.Bind(&req); err != nil {
		log.Println(" err: %v", err)
		return
	}
	log.Println("Token:", req.Token, req)
	data := new(dto.ResultRespXML)
	data.Msg = "接收成功"
	data.Code = 0
	c.XML(200, data)
	d, _ := xml.MarshalIndent(req, "  ", "  ")
	log.Println()
	log.Println("+++++++++++++++++++++++++++")
	log.Println(req.LprInfo.VehicleImgPath)
	s := strings.Split(req.LprInfo.VehicleImgPath, "/")
	log.Println(s[3], s[4])
	createxml(s[4], d)
}

/*
<result>
 <code>0</code>
 <msg>接收成功</msg>
</result>
*/

//创建xml文件
func createxml(xmlname string, outputxml []byte) string {

	fw, f_werr := os.Create("./gatewaySimulator/captureResultReceive/" + xmlname + ".xml") //go run main.go
	if f_werr != nil {
		log.Println("Read:", f_werr)
		return ""
	}
	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, outputxml...)

	_, ferr := fw.Write((xmlOutPutData))
	if ferr != nil {
		log.Printf("  Write xml file error: %v\n", ferr)
		return ""
	}

	defer func() {
		_ = fw.Close()
	}()

	return "/captureResultReceive/" + xmlname + ".xml"

}
