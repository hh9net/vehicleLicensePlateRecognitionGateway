package db

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"operationPlatform/dto"
	"testing"
)

//HandleDayTasks()
func TestHandleDayTasks(t *testing.T) {
	Newdb()
	HandleDayTasks()
	HandleHourTasks()
	HandleMinutesTasks()
	HandleSecondTasks()
}

func TestGatewayDataPostWithJson1(t *testing.T) {
	Newdb()
	GatewayDataPostWithJson()
}
func TestGatewayDataPostWithJson(t *testing.T) {
	Newdb()
	b := []byte(`{"code": 0,"data": [{"msghead": {"parkid": "2002009998","companyid": "3202999999","terminal_id": "CE4C37043A520C93","msgtype": "3"},"programe_runtime": "348","deviceid": "CE4C37043A520C93","gatewayip": "192.168.200.215","getway_version": "build2020-06-29 08:56:41|ver1","lastversion_updatedatetime": "2020-09-02 16:50:00","antenna_infos": [{"laneid": "1101","rsuip": "192.168.200.248","rsuport": "21003","power": "12","waittime": "12","allow_repeattime": "0","isregister": "0","antenna_status": "0","antenna_status_updatetime": ""}]}],"msg": "SUCCESS"}`)
	Resp := new(dto.GatewayDeviceMsgResp)
	unmerr := json.Unmarshal(b, Resp)
	if unmerr != nil {
		log.Println("json.Unmarshal error")
	}
	log.Printf("Post request with json result:%v\n", Resp)
}

func TestErrorDataPostWithJson(t *testing.T) {

	Newdb()
	var startetime, endtime int64
	startetime = 1600087800
	endtime = 1600341810
	ErrorDataPostWithJson(startetime, endtime)
}

func TestRestartDataPostWithJson(t *testing.T) {
	Newdb()
	Restart_address = "http://180.76.177.104:8080/etcpark/v1/metric/list"
	var startetime, endtime int64
	startetime = 1600087800
	endtime = 1900978029
	RestartDataPostWithJson(startetime, endtime)
}

func TestMetricDataPostWithJson(t *testing.T) {
	Newdb()
	m1 := "gateway.park.gateway.cpupercent"
	m2 := "gateway.park.gateway.mempercent"
	log.Println(MetricDataPostWithJson(m1))

	log.Println(MetricDataPostWithJson(m2))

}
