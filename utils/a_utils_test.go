package utils

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestNumberToChinese(t *testing.T) {
	logrus.Print(NumberToChinese(3091746200))
}

func TestByteToMB(t *testing.T) {
	logrus.Print(ByteToMB(1173741834))

	logrus.Print(ByteToGB(4029083648))
}

func TestQingStorGetFile(t *testing.T) {
	//http://ydcpsbxt.sh1a.qingstor.com/cloud_lpr/jiangsu/ggzx/sxjgl_ggzx_320600_G40_K212_2_1_1102_20201209113306_000412.jpg
	QingStorGetFile("sxjgl_ggzx_320600_G40_K212_2_1_1102_20201209113306_000412.jpg", "cloud_lpr/jiangsu/ggzx/sxjgl_ggzx_320600_G40_K212_2_1_1102_20201209113306_000412.jpg")
}

func TestQingStorUpload(t *testing.T) {
	//	QingStorUpload("222.zip")
}

func TestQingStorDeleteFile(t *testing.T) {
	QingStorDeleteFile("sxjgl_yzjtd_320200_G2_K1071_2_0_004_20201124143417_000031.jpg")
}

//
func TestMainMutexFile(t *testing.T) {
	MainMutexFileCreate()
}

//
func TestProcessMutexBegin(t *testing.T) {
	ProcessMutexBegin()
}
