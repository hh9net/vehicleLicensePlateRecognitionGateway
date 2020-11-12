package db

import (
	"github.com/sirupsen/logrus"

	"testing"
)

//用户直接注册

//
func TestQueryUsermsg(t *testing.T) {
	Newdb()
	err, resp := QueryUserLoginmsg("1324312admin55645")
	if err != nil {
		logrus.Print("查询用户能否被注册，失败", err)
	}
	logrus.Printf("查询用户已经被注册 %s,%s,%s,%s", resp.FVcGongsid, resp.FVcMingc, resp.FVcZhangh, resp.FVcMim)

}
