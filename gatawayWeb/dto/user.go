package dto

//成功响应
type Response struct {
	Code    int         `json:"code"  example:"200"`
	Data    interface{} `json:"data"`
	Message string      `json:"message" example:"响应成功信息"`
}

//失败响应
type ResponseFailure struct {
	Code    int         `json:"code"  example:"404"`
	Data    interface{} `json:"data"`
	Message string      `json:"message" example:"响应失败信息"`
}

//注册请求
type ReqRegister struct {
	UserName string `json:"username" example:"abc"`
	Password string `json:"password" example:"abc123"`
	Name     string `json:"name" example:"王小二"`
	Num      string `json:"num" example:"13913661234"`
	Email    string `json:"email" example:"abc123@163.com"`
}

//登录请求
type Reqlogin struct {
	UserName         string `json:"username" example:"abc"`
	Password         string `json:"password" example:"abc123"`
	Verificationcode string `json:"verificationcode" example:"abc123"`
}
