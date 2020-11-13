package db

import (
	"github.com/sirupsen/logrus"
	"vehicleLicensePlateRecognitionGateway/gatawayWeb/types"
	"vehicleLicensePlateRecognitionGateway/utils"
)

//用户注册 查询用户信息是否以及存在
func QueryUserLoginmsg(username string) (error, *types.BSysYongh) {
	db := utils.GormClient.Client
	user := new(types.BSysYongh)
	if err := db.Table("b_sys_yongh").Where("F_VC_ZHANGH = ?", username).First(user).Error; err != nil {
		logrus.Println("查询用户登录信息失败！")
		return err, nil
	}
	logrus.Println("查询用户登录信息 ok:", user.FVcMingc, user.FVcZhangh)

	return nil, user

}

//用户注册 插入数据库
func UserInsert(data *types.BSysYongh) error {
	db := utils.GormClient.Client
	if err := db.Table("b_sys_yongh").Create(&data).Error; err != nil {
		// 错误处理...
		logrus.Println("Insert b_sys_yongh error", err)
		return err
	}
	logrus.Println("用户表插入成功！")
	return nil
}
