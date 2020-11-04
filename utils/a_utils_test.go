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
