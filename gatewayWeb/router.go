package gatewayWeb

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func RouteInit(IpAddress string) {
	logrus.Print("gwWeb服务端 IpAddress：", IpAddress)
	router := gin.New()
	router.Use(Cors()) //跨域资源共享
	apiV1 := router.Group("/vpGateWay/api/v1")
	APIV1Init(apiV1)

	gwApi := router.Group("/vpGateWay/api/v1/gw")
	GwApiV1Init(gwApi)
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

//网关组路由
func GwApiV1Init(route *gin.RouterGroup) {
	GWAuthAPIInit(route)
}

func GWAuthAPIInit(route *gin.RouterGroup) {
	//Gataway
	route.GET("/gatewaybasicdata", GatewayBasicDataQuery)
	route.GET("/gatewaydynamicdata", GatewayDynamicDataQuery)
	route.GET("/camerainfodata", CameraInfoDataQuery)
}

func AuthAPIInit(route *gin.RouterGroup) {
	//用户注册
	//route.POST("/user/register", controller.Register)
	//用户登录
	//route.GET("/user/imagecaptcha", controller.Imagecaptcha)
	route.POST("/user/login", Login)
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
