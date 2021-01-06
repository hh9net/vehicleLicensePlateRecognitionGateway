package dto

//json

type Resp struct {
	Code int    `json:"code"` //201
	Msg  string `json:"msg"`  //接收成功
	//Data interface{} `json:"data"`
}

//2.5.	版本查询接口响应值
type VersionInfo struct {
	GatewayId string `json:"gatewayId"` //1	gatewayId		网关id
	Newver    string `json:"newver"`    //2	newver	v1.0.22_20201222gateway	版本号
	NewverNum string `json:"newverNum"` //3	newverNum	22	新版本号大于当前版本号，则自动执行更新；
	ResTime   string `json:"resTime"`   //4	resTime	2020-02-04 15:01:04	版本请求时间
	Osspath   string `json:"osspath"`   //5	osspath	http://xxxx.xxx/xxx.zip //oss版本文件路径
	Prosize   int    `json:"prosize"`   //6	prosize	111235	byte
	Checksume string `json:"checksume"` //7	checksume	EA12BF5CA127E	文件md5 校验值
}
