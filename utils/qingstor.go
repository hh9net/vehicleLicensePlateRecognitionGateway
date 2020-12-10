package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingstor-sdk-go/config"
	qs "github.com/yunify/qingstor-sdk-go/service"
	"io"
	"os"
	"time"
)

//http://gdkjetcpark.sh1a.qingstor.com/

//可以在不同区域创建 存储空间(Bucket)
const (
	Zone         = "sh1a" //地区
	BucketName   = "gdkjetcpark"
	AccessKey    = "QGYIKYOMPKJLFWVURGBG"
	AccessSecret = "QLPAaRF1legVvjbA8nfz2bN2EiuKRvD9f8HKZISX"
)

//chep
//AccessKey=FCAWLYKENEJOPADCNONN
//SecretKey=ryb5a7yZdHSX0rS8ceaLi3VeaCCwxpw0mU7I179m
//QsZone=sh1a
//#QsBucketname=cpsbxt
//QsBucketname=ydcpsbxt
const (
	CHEPZone         = "sh1a" //地区
	CHEPBucketName   = "ydcpsbxt"
	CHEPAccessKey    = "FCAWLYKENEJOPADCNONN"
	CHEPAccessSecret = "ryb5a7yZdHSX0rS8ceaLi3VeaCCwxpw0mU7I179m"
	UPloadOK         = 201
)

var BacketName string

func QingStorUpload(path, fname, prefix string) (int, int64, string) {

	//发起请求前首先建立需要初始化服务:
	//1、初始化了一个 QingStor Service
	//configuration, _ := config.New("ACCESS_KEY_ID", "SECRET_ACCESS_KEY")
	configuration, _ := config.New(CHEPAccessKey, CHEPAccessSecret)
	qsService, _ := qs.Init(configuration)

	//2、初始化并创建 Bucket, 需要指定 Bucket[桶] 名称和所在 Zone:
	log.Println("qsService.Bucket (BacketName , CHEPZone)", BacketName, CHEPZone) //此处打印
	bucket, _ := qsService.Bucket(BacketName, CHEPZone)

	//2、创建一个 Object 例如上传一张屏幕截图:
	// Open file
	f, err := os.Open(path)
	if err != nil {
		log.Print("上传oss 创建一个 Object error：", err)
		return 0, 0, ""
	}
	log.Println("os.Open fname:", f.Name()) //prefix:/jiangsu/suhuaiyangs/

	defer func() {
		_ = f.Close()
	}()

	// Put object          &service: 包名称
	//Output, err := bucket.PutObject(fname, &service.PutObjectInput{Body: f})
	log.Printf("PutObject:prefix+/+fname:%s", prefix+"/"+fname)
	Output, PutObjecterr := bucket.PutObject(prefix+"/"+fname, &qs.PutObjectInput{Body: f})
	if PutObjecterr != nil {
		// 所有 >= 400 的 HTTP 返回码都被视作错误 Example: QingStor Error: StatusCode 403, Code "permission_denied"...
		log.Println("上传结果有错误:", PutObjecterr)
	} else {
		// Print the HTTP status code. Example: 201
		log.Println("http://" + BacketName + "." + CHEPZone + ".qingstor.com/" + prefix + "/" + fname)
		log.Println("上传结果:", qs.IntValue(Output.StatusCode))
	}

	return qs.IntValue(Output.StatusCode), time.Now().Unix(), prefix + "/" + fname
}

func QingStorGetFile(fname, fm string) {
	//发起请求前首先建立需要初始化服务:
	//1、初始化了一个 QingStor Service
	//configuration, _ := config.New("ACCESS_KEY_ID", "SECRET_ACCESS_KEY")
	configuration, _ := config.New(CHEPAccessKey, CHEPAccessSecret)
	qsService, _ := qs.Init(configuration)

	//2、初始化并创建 Bucket, 需要指定 Bucket[桶] 名称和所在 Zone:
	//bucket, _ := qsService.Bucket("test-bucket", "pek3a")
	bucket, _ := qsService.Bucket(CHEPBucketName, CHEPZone)
	//putBucketOutput, _ := bucket.Put()

	getOutput, err := bucket.GetObject(
		fname,
		&qs.GetObjectInput{},
	)

	if err != nil {
		// Example: QingStor Error: StatusCode 404, Code "object_not_exists"...
		log.Println(err)
		//if qsErr, ok := err.(*qsErrors.QingStorError); ok {
		//	println(qsErr.StatusCode, qsErr.Code)
		//}
	} else {
		defer func() {
			_ = getOutput.Close() // 一定记得关闭GetObjectOutput, 否则容易造成链接泄漏
		}()
		log.Println("fm:", fm)
		f, err := os.OpenFile(fm, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Println(err)
		}
		defer func() {
			_ = f.Close()
		}()

		// Download with 32k temporary buffer
		_, err = io.CopyBuffer(f, getOutput.Body, make([]byte, 32*1024))
		if err != nil {
			log.Println(err)
		}
	}
}

func QingStorDeleteFile(fname string) {
	//发起请求前首先建立需要初始化服务:
	//1、初始化了一个 QingStor Service
	//configuration, _ := config.New("ACCESS_KEY_ID", "SECRET_ACCESS_KEY")
	configuration, _ := config.New(CHEPAccessKey, CHEPAccessSecret)
	qsService, _ := qs.Init(configuration)

	//2、初始化并创建 Bucket, 需要指定 Bucket[桶] 名称和所在 Zone:
	//bucket, _ := qsService.Bucket("test-bucket", "pek3a")
	bucket, _ := qsService.Bucket(CHEPBucketName, CHEPZone)
	//putBucketOutput, _ := bucket.Put()

	Output, _ := bucket.DeleteObject(fname)

	// Print the HTTP status code.
	// Example: 204[delete ok]
	log.Println("delete :[+++++++++++++", qs.IntValue(Output.StatusCode), "+++++++++++++]", fname)
}

//bucket_not_exists	当访问的 bucket 不存在时，返回此错误	404
//object_not_exists	当访问的 object 不存在时，返回此错误	404
