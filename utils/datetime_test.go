package utils

import (
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestStrTimeTotime(t *testing.T) {
	StrTimeTotime("1212-12-12 12:12:12")
}

func TestStrTimeToNowtime(t *testing.T) {
	StrTimeToNowtime()
}

//KuaizhaoTimeNowFormat
func TestKuaizhaoTimeNowFormat(t *testing.T) {
	logrus.Print(KuaizhaoTimeNowFormat())
}

//Yesterdaydate()
func TestYesterdaydate(t *testing.T) {
	logrus.Print(Yesterdaydate())
}

//OldData
func TestOldData(t *testing.T) {
	logrus.Print(OldData(7))
}

//StrdateToNowdate()
func TestDateToNowdate(t *testing.T) {
	logrus.Println(DateToNowdate())
}

// DateFormatTimeToTime
func TestDateFormatTimeToTime(t *testing.T) {
	logrus.Println(DateFormatTimeTostrdate(time.Now()))
}

func TestGetTimestamp(t *testing.T) {
	logrus.Println(GetTimestamp())

}

//
func TestTimestampToFormat(t *testing.T) {
	logrus.Println(TimestampToFormat(1600678120))
	logrus.Println(StrTimeTotime(TimestampToFormat(1600678120)))
	//"2020-09-22 15:20:01"
}

//
func TestStrTimeTimestamp(t *testing.T) {
	logrus.Println(StrTimeToTimestamp("2020-09-22 15:20:01"))
	//1600759201
}
func TestSecondsToTime(t *testing.T) {
	//logrus.Println(SecondsToTime(36069))//"10小时1分9秒"
	//logrus.Println(SecondsToTime(86469))//"1天0小时1分9秒"
	logrus.Println(SecondsToTime(96469))  //"1天2小时47分49秒"
	logrus.Println(SecondsToTime(964690)) //"11天3小时58分10秒"
}

//GetSomeTimestamp
func TestGetSomeTimestamp(t *testing.T) {
	logrus.Println(GetSomeTimestamp())
	//1600759201
}

func TestTimeDifference(t *testing.T) {
	t1 := time.Now()
	time.Sleep(time.Second * 5)
	t2 := time.Now()
	logrus.Println(TimeDifference(t1, t2))

}
