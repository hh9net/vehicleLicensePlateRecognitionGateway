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
	QingStorGetFile("222.zip")
}

func TestQingStorUpload(t *testing.T) {
	QingStorUpload("222.zip")
}

func TestQingStorDeleteFile(t *testing.T) {
	QingStorDeleteFile("222.zip")
}
