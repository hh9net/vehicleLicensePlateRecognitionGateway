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
	QingStorGetFile("/jiangsu/suhuaiyangs/sxjgl_yzjtd_320200_G2_K1071_2_0_004_20201124143417_000031.jpg", "20201124", "sxjgl_yzjtd_320200_G2_K1071_2_0_004_20201124143417_000031.jpg")
}

func TestQingStorUpload(t *testing.T) {
	//	QingStorUpload("222.zip")
}

func TestQingStorDeleteFile(t *testing.T) {
	QingStorDeleteFile("sxjgl_yzjtd_320200_G2_K1071_2_0_004_20201124143417_000031.jpg")
}
