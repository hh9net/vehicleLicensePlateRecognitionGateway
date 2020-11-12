package types

import "time"

//23 用户表
//CREATE TABLE `b_sys_yongh` (
type BSysYongh struct {
	FVcId          string    `gorm:"column:F_VC_ID"`         //`F_VC_ID` varchar(32) NOT NULL COMMENT 'ID 由系统生成的唯一标识；',
	FVcZhangh      string    `gorm:"column:F_VC_ZHANGH"`     //`F_VC_ZHANGH` varchar(32) NOT NULL COMMENT '账号 idx-',
	FVcMingc       string    `gorm:"column:F_VC_MINGC"`      //`F_VC_MINGC` varchar(32) NOT NULL COMMENT '姓名',
	FVcNic         string    `gorm:"column:F_VC_NIC"`        //`F_VC_NIC` varchar(32) DEFAULT NULL COMMENT '昵称',
	FNbNiann       int       `gorm:"column:F_NB_NIANN"`      //`F_NB_NIANN` int(11) DEFAULT NULL COMMENT '年龄',
	FNbXingb       int       `gorm:"column:F_NB_XINGB"`      //`F_NB_XINGB` int(11) DEFAULT NULL COMMENT '性别 0：男性；1：女性；',
	FVcDianh       string    `gorm:"column:F_VC_DIANH"`      //`F_VC_DIANH` varchar(32) NOT NULL COMMENT '电话',
	FVcYouj        string    `gorm:"column:F_VC_YOUJ"`       //`F_VC_YOUJ` varchar(32) DEFAULT NULL COMMENT '邮件',
	FVcToux        string    `gorm:"column:F_VC_TOUX"`       //`F_VC_TOUX` varchar(1024) DEFAULT NULL COMMENT '头像',
	FVcMim         string    `gorm:"column:F_VC_MIM"`        //`F_VC_MIM` varchar(32) NOT NULL COMMENT '密码 md5加密',
	FVcGongsid     string    `gorm:"column:F_VC_GONGSID"`    //`F_VC_GONGSID` varchar(32) DEFAULT NULL COMMENT '公司ID',
	FVcZuzid       string    `gorm:"column:F_VC_ZUZID"`      //`F_VC_ZUZID` varchar(32) DEFAULT 'root' COMMENT '组织ID 为root则表示处于公司根组织',
	FNbZhuangt     int       `gorm:"column:F_NB_ZHUANGT"`    //`F_NB_ZHUANGT` int(11) DEFAULT '1' COMMENT '状态 1：正常；2：已禁用；',
	FNbShenfzyyzzt int       `gorm:"column:F_NB_SHENFZYZZT"` //`F_NB_SHENFZYZZT` int(11) DEFAULT '0' COMMENT '身份证验证状态 0：待提交；1：待审核；2：审核通过；3：审核驳回，需修改信息；4：审核拒绝；',
	FVcShenfzhm    string    `gorm:"column:F_VC_SHENFZHM"`   //`F_VC_SHENFZHM` varchar(32) DEFAULT NULL COMMENT '身份证号码',
	FVcShenfzzp    string    `gorm:"column:F_VC_SHENFZZP"`   //`F_VC_SHENFZZP` varchar(1024) DEFAULT NULL COMMENT '身份证照片',
	FDtChuangjsj   time.Time `gorm:"column:F_DT_CHUANGJSJ"`  //`F_DT_CHUANGJSJ` datetime NOT NULL COMMENT '创建时间',
	FDtDenglsj     time.Time `gorm:"column:F_DT_DENGLSJ"`    //`F_DT_DENGLSJ` datetime DEFAULT NULL COMMENT '登录时间',
	FVcDenglip     string    `gorm:"column:F_VC_DENGLIP"`    //`F_VC_DENGLIP` varchar(32) DEFAULT NULL COMMENT '登录IP',
	FDtGengxsj     time.Time `gorm:"column:F_DT_GENGXSJ"`    //`F_DT_GENGXSJ` datetime DEFAULT NULL COMMENT '更新时间',
	FVcWeixid      string    `gorm:"column:F_VC_WEIXID"`     //`F_VC_WEIXID` varchar(32) DEFAULT NULL COMMENT '微信ID',
	FVcQqid        string    `gorm:"column:F_VC_QQID"`       //`F_VC_QQID` varchar(32) DEFAULT NULL COMMENT 'QQ ID',
	FNbGuanlylx    int       `gorm:"column:F_NB_GUANLYLX"`   //`F_NB_GUANLYLX` int(11) DEFAULT '0' COMMENT ' 0：普通用户；1：系统超级管理员；10：单点公司管理员；11：单点停车场管理员；20：总对总公司管理员；21: 总对总停车场管理员31：服务商管理员；',
	FNbYindzt      int       `gorm:"column:F_NB_YINDZT"`     //`F_NB_YINDZT` int(11) DEFAULT '0' COMMENT '引导状态 0：未引导； 1：已引导',
}
