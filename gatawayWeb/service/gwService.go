package service

import log "github.com/sirupsen/logrus"

var Deviceid string //网关设备id
//1、获取网关设备的token
func GetGatawayToken() {

	Resp, err := GetCameraList(Deviceid)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(Resp.Token, err)
}
