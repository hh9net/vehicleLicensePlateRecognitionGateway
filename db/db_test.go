package db

import (
	"log"
	"operationPlatform/types"
	"testing"
	"time"
)

func TestQueryTable(t *testing.T) {
	//数据库生成的表名是结构体名称的复数形式
	s1 := []string{"b_jsjk_jiesjkptyhb",
		"b_jsjk_zhuanjssjjk",
		"b_jsjk_yuqsjtj",
		"b_jsjk_yicsjtj",
		"b_jsjk_yicsjtcctj",
		"b_jsjk_tingccjssjtj",
		"b_jsjk_shujtbjk",
		"b_jsjk_shujbjk",
		"b_jsjk_shengwtccjsqs",
		"b_jsjk_shengwqftj",
		"b_jsjk_shengwjszysjtj",
		"b_jsjk_shengwjssjfl",
		"b_jsjk_shengwjsqs",
		"b_jsjk_shengnyfssjtj",
		"b_jsjk_shengntccjsqs",
		"b_jsjk_shengnsssjjk",
		"b_jsjk_shengnqktj",
		"b_jsjk_shengnjssjfl",
		"b_jsjk_shengnjsqs",
		"b_jsjk_heimdjk",
		"b_jsjk_shengnjfsjtj",
		"b_jsjk_qingfhd",
		"b_jsjk_jiestj",
		"JieSuanWssjs",
	}

	tablenames := make([]string, 0)
	tablenames = append(tablenames, s1...)
	log.Println(len(tablenames), tablenames)
	Newdb()
	for i, tablename := range tablenames {
		log.Println(1+i, tablename)
		QueryTable(tablename)
	}
}

func TestGatewayInsert(t *testing.T) {
	Newdb()
	gwxx := new(types.BDmWanggjcxx)
	gwxx.FVcWanggbh = "abc1231"
	gwxx.FDtChuangjsj = time.Now()
	gwxx.FNbCPUsyl = 10.9
	ierr := GatewayInsert(gwxx)
	if ierr != nil {
		log.Println("error:", ierr)
	}
}

func TestQueryGatewaydata(t *testing.T) {
	Newdb()
	qerr, gwxx := QueryGatewaydata("abc1231")
	if qerr != nil {
		log.Println("error:", qerr, gwxx)
	}
	log.Println("gwxx:", qerr, gwxx)

}

func TestQueryGatewayALLdata(t *testing.T) {
	Newdb()
	//qerr, gwxxs := QueryGatewayALLdata()
	//if qerr != nil {
	//	log.Println("error:", qerr, gwxxs)
	//}
	//log.Println("gwxx:", qerr, gwxxs)

}

//

func TestQueryRestartOnedata(t *testing.T) {
	Newdb()
	qerr, gwxxs := QueryRestartOnedata("gw1111")
	if qerr != nil {
		log.Println("error:", qerr, gwxxs)
	}
	log.Println("gwxx:", qerr, gwxxs)

}

//QueryAlarm()
func TestQueryAlarm(t *testing.T) {
	Newdb()
	qerr, gjxxs := QueryAlarm()
	if qerr != nil {
		log.Println("error:", qerr, gjxxs)
	}
	log.Println("gjxx:", gjxxs.FDtGaojsj)

}

func TestQueryChedaoMC(t *testing.T) {
	Newdb()
	qerr, gjxxs := QueryChedaoMC("1102")
	if qerr != nil {
		log.Println("error:", qerr, gjxxs)
	}
	log.Println("Chedmc:", gjxxs.FVcChedmc)

}
