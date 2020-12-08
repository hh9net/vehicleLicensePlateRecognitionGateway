package utils

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func DateTimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05") //后面的参数是固定的 否则将无输出
}

func DateFormatTimeToTime(data time.Time) time.Time {
	datestr := data.Format("2006-01-02 00:00:00")
	const Layout = "2006-01-02 15:04:05" //时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")
	time1, _ := time.ParseInLocation(Layout, datestr /*需要转换的时间类型字符串*/, loc)
	return time1
}

func DateFormatTimeTostrdate(data time.Time) string {
	datestr := data.Format("2006-01-02 00:00:00")
	b := []byte(datestr)
	return string(b[0:10])
}

//2006-01-02 15:04:05
func DateTimeNowFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//2006-01-02 15:04:05
func KuaizhaoTimeNowFormat() string {
	return time.Now().Format("2006-01-02 15:04:00")
}

func DateNowFormat() string {
	return time.Now().Format("2006-01-02")
}

//处理时间字符串转时间
func StrTimeTotime(strTime string) time.Time {
	const Layout = "2006-01-02 15:04:05" //时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tim, _ := time.ParseInLocation(Layout, strTime /*需要转换的时间类型字符串*/, loc)
	log.Println(tim)
	return tim
}

//处理时间字符串转时间
func StrTimeToNowtime() time.Time {
	strTime := time.Now().Format("2006-01-02 15:04:05")
	const Layout = "2006-01-02 15:04:05" //时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tim, _ := time.ParseInLocation(Layout, strTime /*需要转换的时间类型字符串*/, loc)
	log.Println(tim)
	return tim
}

//处理时间字符串转时间
func DateToNowdate() time.Time {
	strTime := time.Now().Format("2006-01-02")
	const Layout = "2006-01-02" //时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")

	tim, _ := time.ParseInLocation(Layout, strTime /*需要转换的时间类型字符串*/, loc)
	log.Println(tim)
	return tim
}

//处理时间字符串转时间
func StrTimeTodefaultdate() time.Time {
	strTime := "2020-01-01 00:00:00"
	const Layout = "2006-01-02 15:04:05" //时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")

	tim, _ := time.ParseInLocation(Layout, strTime /*需要转换的时间类型字符串*/, loc)
	log.Println(tim)
	return tim
}

//处理时间字符串转时间
func StrTimeTodefaultdatetimestr() string {
	return "2020-01-01 00:00:00"
}

//获取昨天的日期
func Yesterdaydate() string {
	nTime := time.Now()
	yesTime := nTime.AddDate(0, 0, -1)
	return yesTime.Format("2006-01-02")
}

//currentTime := time.Now()
//oldTime := currentTime.AddDate(0, 0, -2)
//获取前num天的日期
func OldData(num int) []string {
	days := num
	nTime := time.Now()
	//switch num {
	//case 1:
	//	days = -1
	//case 7:
	//	days = -7
	//case 14:
	//	days = -14
	//case 30:
	//	days = -7
	//default: //default case
	//	log.Println("  number  error")
	//	return nil
	//}
	daystrs := make([]string, 0)
	for i := 0; i < num; i++ {
		yesTime := nTime.AddDate(0, 0, -days)
		daystrs = append(daystrs, yesTime.Format("20060102"))
		days--
	}
	//log.Println("daystrs:", daystrs)
	return daystrs
}

//获取当前时间戳
func GetTimestamp() int64 {
	return time.Now().Unix()
}

//获取当前时间戳
func GetSomeTimesstamp(t time.Time) int64 {
	return t.Unix()
}

//获取十分钟前的时间戳
func GetSomeTimestamp() int64 {
	now := time.Now()               //获取当前时间
	t := now.Add(time.Minute * -10) // 获取上分钟时间
	//10分钟前的时间戳，
	//lastM := t.Unix() - int64(now.Second()) // 上-分钟时间-上-分钟秒数
	//now.Unix()

	return t.Unix()
}

//时间戳转时间字符串
func TimestampToFormat(timeUnix int64) string {

	log.Println(time.Unix(timeUnix, 0))
	// 之后可以用Format 比如
	strTime := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	log.Println(strTime)
	return strTime
}

//时间字符串转时间戳 "2020-09-22 15:20:01"
func StrTimeToTimestamp(strtime string) int64 {
	s := strings.Split(strtime, " ")
	s1 := strings.Split(s[0], "-")
	y, _ := strconv.Atoi((s1[0]))
	m, _ := strconv.Atoi((s1[1]))
	d, _ := strconv.Atoi((s1[2]))

	t := strings.Split(s[1], ":")
	h, _ := strconv.Atoi((t[0]))
	min, _ := strconv.Atoi((t[1]))
	sec, _ := strconv.Atoi((t[2]))

	tstamp := time.Date(y, time.Month(m), d, h, min, sec, 0, time.Local).Unix()
	log.Println(tstamp)
	return tstamp
}

//秒转
func SecondsToTurn(intime time.Time, outtime time.Time) string {

	var day string
	in := intime.Format("2006-01-02 15:04:05")
	out := outtime.Format("2006-01-02 15:04:05")
	inT, _ := time.Parse("2006-01-02 15:04:05", in)
	outT, _ := time.Parse("2006-01-02 15:04:05", out)
	t := outT.Sub(inT).String()
	log.Println("时间差", t)

	if strings.Contains(t, "h") {
		//h
		hourstr := strings.Split(t, "h")
		hour, _ := strconv.Atoi(hourstr[0])
		//m eg:59m59s
		m := strings.Split(hourstr[1], "m")
		//s eg:59s
		s := strings.Split(m[1], "s")
		//h
		if hour >= 24 {
			//小时转天
			d := hour / 24
			h := hour % 24
			dstr := strconv.Itoa(d)
			hstr := strconv.Itoa(h)
			day = dstr + "天" + hstr + "时" + m[0] + "分" + s[0] + "秒"
			return day
		}

		if 24 > hour && hour > 0 {
			day = hourstr[0] + "时" + m[0] + "分" + s[0] + "秒"
			return day
		}
	} else if strings.Contains(t, "m") {
		//m eg:59m59s
		m := strings.Split(t, "m")
		//s eg:59s
		s := strings.Split(m[1], "s")

		day = m[0] + "分" + s[0] + "秒"
		return day
	} else if strings.Contains(t, "s") {
		//s eg:59s
		s := strings.Split(t, "s")
		day = s[0] + "秒"
		return day
	}
	return ""
}

//时差得秒
func TimeDifference(intime time.Time, outtime time.Time) string {
	in := intime.Format("2006-01-02 15:04:05")
	out := outtime.Format("2006-01-02 15:04:05")
	inT, _ := time.Parse("2006-01-02 15:04:05", in)
	outT, _ := time.Parse("2006-01-02 15:04:05", out)
	t := outT.Sub(inT).String()
	log.Println("时间差", t)
	return t
}

/*
时间常量
*/
const (
	//定义每分钟的秒数
	SecondsPerMinute = 60
	//定义每小时的秒数
	SecondsPerHour = SecondsPerMinute * 60
	//定义每天的秒数
	SecondsPerDay = SecondsPerHour * 24
)

/*
时间转换函数
*/
func resolveTime(seconds int) (day int, hour int, minute int, sec int) {
	//天数
	day = seconds / SecondsPerDay
	//天余（小时秒数） ：小时数
	hour = (seconds % SecondsPerDay) / SecondsPerHour
	//小时秒数余（分钟秒数）：分钟数
	minute = ((seconds % SecondsPerDay) % (SecondsPerHour)) / SecondsPerMinute
	//
	sec = ((seconds % SecondsPerDay) % (SecondsPerHour)) % SecondsPerMinute //% SecondsPerMinute

	return
}

//秒数转时间 "11天3小时58分10秒"
func SecondsToTime(seconds int) string {
	//打印返回参数
	day, hour, minute, second := resolveTime(seconds)
	log.Println(day, hour, minute, second)
	d := strconv.Itoa(day)
	h := strconv.Itoa(hour)
	m := strconv.Itoa(minute)
	s := strconv.Itoa(second)
	str := d + "天" + h + "时" + m + "分" + s + "秒"

	str1 := strings.Split(str, "天")
	if str1[0] == "0" {
		str2 := strings.Split(str1[1], "时")
		if str2[0] == "0" {
			str3 := strings.Split(str2[1], "分")
			if str3[0] == "0" {
				//"9秒"
				return str3[1]
			}
			//"1分9秒"
			return str2[1]
		}
		//"1时1分9秒"
		return str1[1]
	} else {
		//"1天2时47分49秒"
		return str
	}

}
