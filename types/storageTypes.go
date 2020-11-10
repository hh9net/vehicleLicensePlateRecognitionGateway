package types

import "time"

var KafkaIpa string
var KafkaIpb string
var KafkaIpc string

var DdkafkaTopic string
var ZdzkafkaTopic string

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

//CREATE TABLE `b_dm_chongq` 重启(
type BDmChongq struct {
	FNbWeiyjlid int `gorm:"column:F_NB_WEIYJLID; AUTO_INCREMENT ;primary_key"` //`F_NB_WEIYJLID` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一记录ID',

	FVcWanggbh      string    `gorm:"column:F_VC_WANGGBH"`      //`F_VC_WANGGBH` varchar(32) NOT NULL COMMENT '网关编号',
	FDtChongqsj     time.Time `gorm:"column:F_DT_CHONGQSJ"`     //`F_DT_CHONGQSJ` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '重启时间',
	FNbChongqlxgzsc int       `gorm:"column:F_NB_CHONGQLXGZSC"` //`F_NB_CHONGQLXGZSC` bigint(20) DEFAULT NULL COMMENT '重启前连续工作时长 单位：秒',
	FNbChongqlx     int       `gorm:"column:F_NB_CHONGQLX"`     //`F_NB_CHONGQLX` int(11) NOT NULL DEFAULT '0' COMMENT '重启类型 0：自动、1：手动',
	FVcChongqrid    string    `gorm:"column:F_VC_CHONGQRID"`    //`F_VC_CHONGQRID` varchar(32) DEFAULT NULL COMMENT '重启人ID',
	FVcChongqrxm    string    `gorm:"column:F_VC_CHONGQRXM"`    //`F_VC_CHONGQRXM` varchar(32) DEFAULT NULL COMMENT '重启人姓名',
	//PRIMARY KEY (`F_NB_WEIYJLID`),
	//KEY `IDX_WANGGBH` (`F_VC_WANGGBH`),
	//KEY `IDX_CHONGQSJ` (`F_DT_CHONGQSJ`),
	//KEY `IDX_CHONGQLX` (`F_NB_CHONGQLX`)
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='重启 ';
}

//CREATE TABLE `b_dm_gaoj`告警  (
type BDmGaoj struct {
	FNbWeiyjlid int `gorm:"column:F_NB_WEIYJLID; AUTO_INCREMENT ;primary_key"` //	`F_NB_WEIYJLID` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一记录ID',

	FVcWanggbh     string    `gorm:"column:F_VC_WANGGBH"` //	`F_VC_WANGGBH` varchar(32) NOT NULL COMMENT '网关编号',
	FDtGaojsj      time.Time `gorm:"column:F_DT_GAOJSJ"`  //	`F_DT_GAOJSJ` datetime DEFAULT NULL COMMENT '告警时间',
	FVcGaojms      string    `gorm:"column:F_VC_GAOJMS"`  //	`F_VC_GAOJMS` varchar(1024) DEFAULT NULL COMMENT '告警描述',
	FNbChulZhuangt int       `gorm:"column:F_NB_ZHUANGT"` //	`F_NB_ZHUANGT` int(11) NOT NULL DEFAULT '0' COMMENT '状态 0：未处理、1：已处理',
	FVcChulrid     string    `gorm:"column:F_VC_CHULRID"` //	`F_VC_CHULRID` varchar(32) DEFAULT NULL COMMENT '处理人ID',
	FVcChulrxm     string    `gorm:"column:F_VC_CHULRXM"` //	`F_VC_CHULRXM` varchar(32) DEFAULT NULL COMMENT '处理人姓名',
	FDtChulsj      time.Time `gorm:"column:F_DT_CHULSJ"`  //	`F_DT_CHULSJ` datetime DEFAULT NULL COMMENT '处理时间',
	//	PRIMARY KEY (`F_NB_WEIYJLID`),
	//	KEY `IDX_WANGGBH` (`F_VC_WANGGBH`),
	//	KEY `IDX_GAOJSJ` (`F_DT_GAOJSJ`),
	//	KEY `IDX_ZHUANGT` (`F_NB_ZHUANGT`),
	//	KEY `IDX_CHULSJ` (`F_DT_CHULSJ`),
	//	KEY `IDX_CHULRID` (`F_VC_CHULRID`)
	//) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '告警 '
}

//CREATE TABLE `b_dm_ruanjbb`软件版本 ' (
type BDmRuanjbb struct {
	FNbWeiyjlid int `gorm:"column:F_NB_WEIYJLID; AUTO_INCREMENT ;primary_key"` //	`F_NB_WEIYJLID` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一记录ID',

	FVcRuanjbbh  string    `gorm:"column:F_VC_RUANJBBH"`  //	`F_VC_RUANJBBH` varchar(512) NOT NULL COMMENT '软件版本号',
	FVcBanbgxnr  string    `gorm:"column:F_VC_BANBGXNR"`  //	`F_VC_BANBGXNR` varchar(1024) DEFAULT NULL COMMENT '版本更新内容',
	FDtShangcsj  time.Time `gorm:"column:F_DT_SHANGCSJ"`  //	`F_DT_SHANGCSJ` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
	FVcShangczid string    `gorm:"column:F_VC_SHANGCZID"` //	`F_VC_SHANGCZID` varchar(32) NOT NULL COMMENT '上传者ID',
	FVcShangczxm string    `gorm:"column:F_VC_SHANGCZXM"` //	`F_VC_SHANGCZXM` varchar(32) DEFAULT NULL COMMENT '上传者姓名',
	FVcWenjlj    string    `gorm:"column:F_VC_WENJLJ"`    //	`F_VC_WENJLJ` varchar(512) DEFAULT NULL COMMENT '文件路径',
	FNbZhuangt   int       `gorm:"column:F_NB_ZHUANGT"`   //	`F_NB_ZHUANGT` int(11) NOT NULL DEFAULT '0' COMMENT '状态 0：正常、1：已删除',
	//	PRIMARY KEY (`F_NB_WEIYJLID`),
	//	KEY `IDX_SHANGCSJ` (`F_DT_SHANGCSJ`),
	//	KEY `IDX_SHANGCZID` (`F_VC_SHANGCZID`),
	//	KEY `IDX_ZHUANGT` (`F_NB_ZHUANGT`)
	//) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '软件版本 '
}

//CREATE TABLE `b_dm_ruanjgxzx`软件更新执行  (
type BDmRuanjgxzx struct {
	FNbWeiyjlid int `gorm:"column:F_NB_WEIYJLID; AUTO_INCREMENT ;primary_key"` //	`F_NB_WEIYJLID` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一记录ID',

	FVcWanggbh   string    `gorm:"column:F_VC_WANGGBH"`   //	`F_VC_WANGGBH` varchar(32) NOT NULL COMMENT '网关编号',
	FNbBanbid    int       `gorm:"column:F_NB_BANBENID"`  //	`F_NB_BANBENID` int(11) NOT NULL COMMENT '版本ID',
	FVcRuanjbbh  string    `gorm:"column:F_VC_RUANJBBH"`  //	`F_VC_RUANJBBH` varchar(512) NOT NULL COMMENT '软件版本号',
	FNbJihgxcl   int       `gorm:"column:F_NB_JIHGXCL"`   //	`F_NB_JIHGXCL` int(11) NOT NULL DEFAULT '0' COMMENT '计划更新策略 0：立即更新、1：定时更新',
	FDtJihgxsj   time.Time `gorm:"column:F_DT_JIHGXSJ"`   //	`F_DT_JIHGXSJ` datetime DEFAULT NULL COMMENT '计划更新时间',
	FNbZhuangt   int       `gorm:"column:F_NB_ZHUANGT"`   //	 '状态 0：未完成、1：已完成'，2更新中,
	FDtGengxwcsj time.Time `gorm:"column:F_DT_GENGXWCSJ"` //	`F_DT_GENGXWCSJ` datetime DEFAULT NULL COMMENT '更新完成时间',
	//	PRIMARY KEY (`F_NB_WEIYJLID`),
	//	KEY `IDX_WANGGBH` (`F_VC_WANGGBH`),
	//	KEY `IDX_ZHUANGT` (`F_NB_ZHUANGT`),
	//	KEY `IDX_GENGXWCSJ` (`F_DT_GENGXWCSJ`)
	//) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '软件更新执行 '
}

//CREATE TABLE `b_dm_tianxxx` 天线信息 (
type BDmTianxxx struct {
	FNbWeiyjlid int `gorm:"column:F_NB_WEIYJLID; AUTO_INCREMENT ;primary_key"` //	`F_NB_WEIYJLID` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一记录ID',

	FVcWanggbh     string    `gorm:"column:F_VC_WANGGBH"`     //	`F_VC_WANGGBH` varchar(32) NOT NULL COMMENT '网关编号',
	FVcChedwyid    string    `gorm:"column:F_VC_CHEDWYID"`    //	`F_VC_CHEDWYID` varchar(32) NOT NULL COMMENT '车道唯一ID',
	FVcIpdz        string    `gorm:"column:F_VC_IPDZ"`        //	`F_VC_IPDZ` varchar(32) DEFAULT NULL COMMENT 'IP地址',
	FDtShangcqdsj  time.Time `gorm:"column:F_DT_SHANGCQDSJ"`  //	`F_DT_SHANGCQDSJ` datetime DEFAULT NULL COMMENT '上次启动时间',
	FNbLianxgzsc   int       `gorm:"column:F_NB_LIANXGZSC"`   //	`F_NB_LIANXGZSC` bigint(20) DEFAULT NULL COMMENT '连续工作时长 单位：秒',
	FVcZhuczt      string    `gorm:"column:F_VC_ZHUCZT"`      //	`F_VC_WANGGBH` varchar(32) NOT NULL COMMENT '注册状态',
	FVcTianxzt     string    `gorm:"column:F_VC_TIANXZT"`     //	`F_VC_CHEDWYID` varchar(32) NOT NULL COMMENT '天线状态',
	FVcTianxztgxsj string    `gorm:"column:F_VC_TIANXZTGXSJ"` //	`F_VC_IPDZ` varchar(32) DEFAULT NULL COMMENT '天线状态更新时间',
	//`F_VC_ZHUCZT` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '注册状态 1：已注册，其他：未注册',
	//`F_VC_TIANXZT` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '天线状态 1：正常，其他：异常',
	//`F_VC_TIANXZTGXSJ` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '天线状态更新时间 yyyy-MM-dd HH:mm:ss',
	//	PRIMARY KEY (`F_NB_WEIYJLID`),
	//	KEY `IDX_WANGGBH` (`F_VC_WANGGBH`),
	//	KEY `IDX_SHANGCQDSJ` (`F_DT_SHANGCQDSJ`),
	//	KEY `IDX_LIANXGZSC` (`F_NB_LIANXGZSC`)
	//) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '天线信息 '
}

//CREATE TABLE `b_dm_wanggjcxx`网关基础信息  (
type BDmWanggjcxx struct {
	FVcWanggbh    string    `gorm:"column:F_VC_WANGGBH"`    //	`F_VC_WANGGBH` varchar(32) NOT NULL COMMENT '网关编号',
	FVcGongsID    string    `gorm:"column:F_VC_GONGSID"`    //	`F_VC_GONGSID` varchar(32) NOT NULL COMMENT '公司ID',
	FVcTingccbh   string    `gorm:"column:F_VC_TINGCCBH"`   //	`F_VC_TINGCCBH` varchar(32) NOT NULL COMMENT '停车场编号',
	FNbZhuangt    int       `gorm:"column:F_NB_ZHUANGT"`    //	`F_NB_ZHUANGT` int(11) NOT NULL DEFAULT '0' COMMENT '状态 0：离线、1：在线',
	FVcIpdz       string    `gorm:"column:F_VC_IPDZ"`       //	`F_VC_IPDZ` varchar(32) DEFAULT NULL COMMENT 'IP地址',
	FNbCPUsyl     float64   `gorm:"column:F_NB_CPUSYL"`     //	`F_NB_CPUSYL` decimal(32, 10) DEFAULT NULL COMMENT 'CPU使用率 百分比',
	FNbNeicsyl    float64   `gorm:"column:F_NB_NEICSYL"`    //	`F_NB_NEICSYL` decimal(32, 10) DEFAULT NULL COMMENT '内存使用率 百分比',
	FNbYsyncdx    float64   `gorm:"column:F_NB_YISYNCDX"`   //	`F_NB_YISYNCDX` decimal(32, 10) DEFAULT NULL COMMENT '已使用内存大小 单位：MB',
	FNbZongncdx   float64   `gorm:"column:F_NB_ZONGNCDX"`   //	`F_NB_ZONGNCDX` decimal(32, 10) DEFAULT NULL COMMENT '总内存大小 单位：MB',
	FNbYingpsyl   float64   `gorm:"column:F_NB_YINGPSYL"`   //	`F_NB_YINGPSYL` decimal(32, 10) DEFAULT NULL COMMENT '硬盘使用率 百分比',
	FNbYisyypdx   float64   `gorm:"column:F_NB_YISYYPDX"`   //	`F_NB_YISYYPDX` decimal(32, 10) DEFAULT NULL COMMENT '已使用硬盘大小 单位：GB',
	FNbZongypdx   float64   `gorm:"column:F_NB_ZONGYPDX"`   //	`F_NB_ZONGYPDX` decimal(32, 10) DEFAULT NULL COMMENT '总硬盘大小 单位：GB',
	FNbGaojzs     int       `gorm:"column:F_NB_GAOJZS"`     //	`F_NB_GAOJZS` int(11) NOT NULL DEFAULT '0' COMMENT '告警总数',
	FNbWeiclgjs   int       `gorm:"column:F_NB_WEICLGJS"`   //	`F_NB_WEICLGJS` int(11) NOT NULL DEFAULT '0' COMMENT '未处理告警数',
	FNbChongqcs   int       `gorm:"column:F_NB_CHONGQCS"`   //	`F_NB_CHONGQCS` int(11) NOT NULL DEFAULT '0' COMMENT '重启次数',
	FVcDangqbbh   string    `gorm:"column:F_VC_DANGQBBH"`   //	`F_VC_DANGQBBH` varchar(512) DEFAULT NULL COMMENT '当前版本号',
	FDtZuijgxbbsj time.Time `gorm:"column:F_DT_ZUIJGXBBSJ"` //	`F_DT_ZUIJGXBBSJ` datetime DEFAULT NULL COMMENT '最近更新版本时间',
	FNbTianxsl    int       `gorm:"column:F_NB_TIANXSL"`    //	`F_NB_TIANXSL` int(11) DEFAULT NULL COMMENT '天线数量',
	FNbWanglyc    int       `gorm:"column:F_NB_WANGLYC"`    //	`F_NB_WANGLYC` bigint(20) DEFAULT NULL COMMENT '网络延迟 单位：ms',
	FDtChuangjsj  time.Time `gorm:"column:F_DT_CHUANGJSJ"`  //	`F_DT_CHUANGJSJ` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	FDtZuihgxsj   time.Time `gorm:"column:F_DT_ZUIHGXSJ"`   //	`F_DT_ZUIHGXSJ` datetime DEFAULT NULL COMMENT '网关最后更新数据时间',
	FNbYunxsc     int       `gorm:"column:F_NB_YUNXSC"`     //`F_NB_YUNXSC` int DEFAULT NULL COMMENT '运行时长 单位：s',网关设备
	//	PRIMARY KEY (`F_VC_WANGGBH`),
	//	KEY `IDX_TINGCCBH` (`F_VC_TINGCCBH`),
	//	KEY `IDX_ZHUANGT` (`F_NB_ZHUANGT`),
	//	KEY `IDX_BANBH` (`F_VC_DANGQBBH`),
	//	KEY `IDX_GENGXBBSJ` (`F_DT_ZUIJGXBBSJ`),
	//	KEY `IDX_GENGXSJSJ` (`F_DT_ZUIHGXSJ`)
	//) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '网关基础信息 '
}

//b_tcc_tingcc 停车场表-U
type BTccTingcc struct {
	FVcTingccbh     string `gorm:"column:F_VC_TINGCCBH"`      //`F_VC_TINGCCBH` varchar(32) NOT NULL COMMENT '停车场编号',
	FVcGongsbh      string `gorm:"column:F_VC_GONGSBH"`       //`F_VC_GONGSBH` varchar(32) DEFAULT NULL COMMENT '公司/集团编号 idx-',
	FVcQudbh        string `gorm:"column:F_VC_QUDBH"`         //`F_VC_QUDBH` varchar(32) DEFAULT NULL COMMENT '渠道编号',
	FVcTingccwlbh   string `gorm:"column:F_VC_TINGCCWLBH"`    //`F_VC_TINGCCWLBH` varchar(32) DEFAULT NULL COMMENT '停车场网络编号 由于前期要与旧平台同步，改字段请用数字表示',
	FNbTingcclx     int64  `gorm:"column:F_NB_TINGCCLX"`      //`F_NB_TINGCCLX` int(11) NOT NULL DEFAULT '1' COMMENT '停车场类型 1：单点；2：总对总；',
	FVcMingc        string `gorm:"column:F_VC_MINGC"`         //`F_VC_MINGC` varchar(32) DEFAULT NULL COMMENT '名称-NEW',
	FVcDiz          string `gorm:"column:F_VC_DIZ"`           //`F_VC_DIZ` varchar(512) DEFAULT NULL COMMENT '地址',
	FVcJingd        string `gorm:"column:F_VC_JINGD"`         //`F_VC_JINGD` decimal(32,10) DEFAULT NULL COMMENT '经度',
	FVcWeid         string `gorm:"column:F_VC_WEID"`          //`F_VC_WEID` decimal(32,10) DEFAULT NULL COMMENT '维度',
	FVcGuanlyid     string `gorm:"column:F_VC_GUANLYID"`      //`F_VC_GUANLYID` varchar(32) NOT NULL COMMENT '管理员ID-NEW',
	FDtChuangjsj    string `gorm:"column:F_DT_CHUANGJSJ"`     //`F_DT_CHUANGJSJ` datetime DEFAULT NULL COMMENT '创建时间',
	FVcChuangjr     string `gorm:"column:F_VC_CHUANGJR"`      //`F_VC_CHUANGJR` varchar(32) DEFAULT NULL COMMENT '创建人',
	FNbZhuangt      int    `gorm:"column:F_NB_ZHUANGT"`       //`F_NB_ZHUANGT` int(11) DEFAULT '1' COMMENT '状态-U 1：正常；2：待审核；3：停用；',
	FVcVerifyStatus int    `gorm:"column:F_VC_VERIFY_STATUS"` //`F_VC_VERIFY_STATUS` int(11) DEFAULT NULL COMMENT '审核结果-NEW 1：审核通过；2：待审核；3：审核驳回，需修改信息；4：审核拒绝；',
	FVcFuzrdh       string `gorm:"column:F_VC_FUZRDH"`        //`F_VC_FUZRDH` varchar(32) DEFAULT NULL COMMENT '负责人电话-D',
	FVcFuzrxm       string `gorm:"column:F_VC_FUZRXM"`        //`F_VC_FUZRXM` varchar(32) DEFAULT NULL COMMENT '负责人姓名-D',
	FVcShengdm      string `gorm:"column:F_VC_SHENGDM"`       //`F_VC_SHENGDM` varchar(32) DEFAULT NULL COMMENT '省代码',
	FVcShengmc      string `gorm:"column:F_VC_SHENGMC"`       //`F_VC_SHENGMC` varchar(32) DEFAULT NULL COMMENT '省名称',
	FVcShidm        string `gorm:"column:F_VC_SHIDM"`         //`F_VC_SHIDM` varchar(32) DEFAULT NULL COMMENT '市代码',
	FVcShimc        string `gorm:"column:F_VC_SHIMC"`         //`F_VC_SHIMC` varchar(32) DEFAULT NULL COMMENT '市名称',
	FVcQudm         string `gorm:"column:F_VC_QUDM"`          //`F_VC_QUDM` varchar(32) DEFAULT NULL COMMENT '区代码',
	FVcQumc         string `gorm:"column:F_VC_QUMC"`          //`F_VC_QUMC` varchar(32) DEFAULT NULL COMMENT '区名称',
	FNbFeil         int    `gorm:"column:F_NB_FEIL"`          //`F_NB_FEIL` int(11) DEFAULT NULL COMMENT '费率 万分比',
}

// CREATE TABLE `b_tcc_ched`
type BTccChed struct {
	FVcChedwyid  string `gorm:"column:F_VC_CHEDWYID"`  //`F_VC_CHEDWYID` varchar(32) NOT NULL COMMENT '车道唯一ID-NEW',
	FVcGongsid   string `gorm:"column:F_VC_GONGSID"`   //`F_VC_GONGSID` varchar(32) DEFAULT NULL COMMENT '公司ID-U idx-联合唯一索引',
	FVcChedbh    string `gorm:"column:F_VC_CHEDBH"`    //`F_VC_CHEDBH` varchar(32) DEFAULT NULL COMMENT '车道编号-U idx-联合唯一索引',
	FVcTingccbh  string `gorm:"column:F_VC_TINGCCBH"`  //`F_VC_TINGCCBH` varchar(32) DEFAULT NULL COMMENT '停车场编号',
	FNbTingcclx  int    `gorm:"column:F_NB_TINGCCLX"`  //`F_NB_TINGCCLX` int(11) DEFAULT '1' COMMENT '停车场类型 1：单点；2：总对总；',
	FNbChedlx    int    `gorm:"column:F_NB_CHEDLX"`    //`F_NB_CHEDLX` int(11) NOT NULL DEFAULT '1' COMMENT '车道类型 1、入口，2、出口',
	FVcChedmc    string `gorm:"column:F_VC_CHEDMC"`    //`F_VC_CHEDMC` varchar(32) DEFAULT NULL COMMENT '车道名称-NEW',
	FVcChuangjz  string `gorm:"column:F_VC_CHUANGJZ"`  //`F_VC_CHUANGJZ` varchar(32) DEFAULT NULL COMMENT '创建者',
	FDtChuangjsj string `gorm:"column:F_DT_CHUANGJSJ"` //`F_DT_CHUANGJSJ` datetime DEFAULT NULL COMMENT '创建时间',
	FNbZhuangt   int    `gorm:"column:F_NB_ZHUANGT"`   //`F_NB_ZHUANGT` int(11) NOT NULL DEFAULT '1' COMMENT '状态 1：正常；2：停用',
	FVcPsamid    string `gorm:"column:F_VC_PSAMID"`    //`F_VC_PSAMID` varchar(64) NOT NULL COMMENT 'PSAM卡ID-NEW',
	FVcMiyaolj   string `gorm:"column:F_VC_MIYAOLJ"`   //`F_VC_MIYAOLJ` varchar(128) DEFAULT NULL COMMENT '密钥路径 下载密钥数据解密',
	FVcMiy       string `gorm:"column:F_VC_MIY"`       //`F_VC_MIY` varchar(32) NOT NULL COMMENT '秘钥-MOVE id密码一致才通过验证',
	FVcZuijljsj  string `gorm:"column:F_VC_ZUIJLJSJ"`  //`F_VC_ZUIJLJSJ` varchar(32) DEFAULT NULL COMMENT '车道最近连接时间',
	FVcChedyz    string `gorm:"column:F_VC_CHEDYZ"`    //`F_VC_CHEDYZ` int(11) DEFAULT NULL COMMENT '车道验证 1：已验证；2：未验证；3：已开通；',
	FVcYunxzt    string `gorm:"column:F_VC_YUNXZT"`    //`F_VC_YUNXZT` int(11) DEFAULT NULL COMMENT '车道运行状态； 0：正常运行；1：关闭；',
	FVcChenxbb   string `gorm:"column:F_VC_CHENXBB"`   //`F_VC_CHENXBB` varchar(32) DEFAULT NULL COMMENT '车道程序版本',
	FVcXinxcjsj  string `gorm:"column:F_VC_XINXCJSJ"`  //`F_VC_XINXCJSJ` varchar(32) DEFAULT NULL COMMENT '信息采集时间',
	FVcZhongdid  string `gorm:"column:F_VC_ZHONGDID"`  //`F_VC_ZHONGDID` varchar(32) NOT NULL COMMENT '终端ID-MOVE',
	FVcChangsid  string `gorm:"column:F_VC_CHANGSID"`  //`F_VC_CHANGSID` varchar(32) DEFAULT NULL COMMENT '厂商编码ID',
	FVcRsuid     string `gorm:"column:F_VC_RSUID"`     //`F_VC_RSUID` varchar(32) DEFAULT NULL COMMENT 'RSU设备ID-END',
	//PRIMARY KEY (`F_VC_CHEDWYID`)
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='停车场车道表-U ';
}
