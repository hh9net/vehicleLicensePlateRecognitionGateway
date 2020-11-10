package utils

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var GormClient *GormDB

var HmdGormClient *HmdGormDB

type GormDB struct {
	dbConfig *DBConfig
	lock     sync.RWMutex // lock
	Client   *gorm.DB     // mysql client
}

type HmdGormDB struct {
	dbConfig  *HmdDBConfig
	lock      sync.RWMutex // lock
	HmdClient *gorm.DB     // mysql client
}

type HSDZGormDB struct {
	dbConfig   *HSDZDBConfig
	lock       sync.RWMutex // lock
	HSDZClient *gorm.DB     // mysql client
}

// 本方法会给GormClient赋值，多次调用GormClient指向最后一次调用的GormDB
func InitGormDB(dbConfig *DBConfig) *GormDB {
	logrus.Infoln("starting db")
	if err := dbConfig.check(); err != nil {
		logrus.WithError(err).Errorln("error db config!")
		return nil
	}
	myDB := &GormDB{
		dbConfig: dbConfig,
	}
	db, err := gorm.Open("mysql", dbConfig.DBAddr)
	if err != nil {
		logrus.Fatalln("db initing fail:", err)
		return nil
	}
	err = db.DB().Ping()
	if err != nil {
		logrus.Fatalln("db ping fail:", err)
		return nil
	}
	logrus.WithField("addr:", dbConfig.DBAddr).Infoln("connecting db success!")
	myDB.Client = db
	myDB.initByDBConfigs()
	myDB.autoCreateTable()
	go myDB.timer()
	GormClient = myDB //gormClient
	return myDB
}

// 本方法会给HmdGormClient赋值，多次调用HmdGormClient指向最后一次调用的GormDB
func HmdInitGormDB(dbConfig *HmdDBConfig) *HmdGormDB {
	logrus.Infoln("starting Hmd db")
	if err := dbConfig.hmdcheck(); err != nil {
		logrus.WithError(err).Errorln("error Hmd db config!")
		return nil
	}
	hmdmyDB := &HmdGormDB{
		dbConfig: dbConfig,
	}
	db, err := gorm.Open("mysql", dbConfig.HmdDBAddr)
	if err != nil {
		logrus.Fatalln("Hmd db initing fail", err)
		return nil
	}
	err = db.DB().Ping()
	if err != nil {
		logrus.Fatalln("Hmd db ping fail", err)
		return nil
	}
	logrus.WithField("addr", dbConfig.HmdDBAddr).Infoln("Hmd connecting db success!")
	hmdmyDB.HmdClient = db
	hmdmyDB.initByDBConfigs()
	hmdmyDB.autoCreateTable()
	go hmdmyDB.timer()
	HmdGormClient = hmdmyDB //hmdgormClient
	return hmdmyDB
}

// 初始化参数
func (p *GormDB) initByDBConfigs() {
	p.Client.DB().SetMaxIdleConns(p.dbConfig.MaxIdleConns)
	p.Client.DB().SetMaxOpenConns(p.dbConfig.MaxOpenConns)
	p.Client.LogMode(p.dbConfig.LogMode)
}

//auto create table
func (p *GormDB) autoCreateTable() {
	if p.dbConfig.AutoCreateTables == nil || len(p.dbConfig.AutoCreateTables) == 0 {
		return
	}
	logrus.WithField("addr", p.dbConfig.DBAddr).Infoln("begin initAutoDB")
	for _, item := range p.dbConfig.AutoCreateTables {
		p.autoCreate(item)
	}
}

// 初始化参数
func (p *HmdGormDB) initByDBConfigs() {
	p.HmdClient.DB().SetMaxIdleConns(p.dbConfig.MaxIdleConns)
	p.HmdClient.DB().SetMaxOpenConns(p.dbConfig.MaxOpenConns)
	p.HmdClient.LogMode(p.dbConfig.LogMode)
}

//auto create table
func (p *HmdGormDB) autoCreateTable() {
	if p.dbConfig.AutoCreateTables == nil || len(p.dbConfig.AutoCreateTables) == 0 {
		return
	}
	logrus.WithField("addr", p.dbConfig.HmdDBAddr).Infoln("begin initAutoDB")
	for _, item := range p.dbConfig.AutoCreateTables {
		p.autoCreate(item)
	}
}

func (p *GormDB) autoCreate(it interface{}) {
	err := p.Client.AutoMigrate(it).Error
	if err != nil {
		logrus.Errorln("auto create ", it, " error", err)
	}
}

func (p *HmdGormDB) autoCreate(it interface{}) {
	err := p.HmdClient.AutoMigrate(it).Error
	if err != nil {
		logrus.Errorln("auto create ", it, " error", err)
	}
}

func (p *GormDB) timer() {
	if p.dbConfig.DetectionInterval < 0 {
		return
	}
	timer1 := time.NewTicker(time.Duration(int64(p.dbConfig.DetectionInterval) * int64(time.Second)))
	for {
		select {
		case <-timer1.C:
			err := p.Client.DB().Ping()
			if err != nil {
				logrus.Errorln("mysql connect fail,err:", err)
				logrus.Infoln("reconnect beginning...")
				p.reConnect()
			}
		}
	}
}

func (p *HmdGormDB) timer() {
	if p.dbConfig.DetectionInterval < 0 {
		return
	}
	timer1 := time.NewTicker(time.Duration(int64(p.dbConfig.DetectionInterval) * int64(time.Second)))
	for {
		select {
		case <-timer1.C:
			err := p.HmdClient.DB().Ping()
			if err != nil {
				logrus.Errorln("mysql connect fail,err:", err)
				logrus.Infoln("reconnect beginning...")
				p.hmdreConnect()
			}
		}
	}
}

//重连接
func (p *GormDB) reConnect() {
	db, err := gorm.Open("mysql", p.dbConfig.DBAddr)
	if err != nil {
		logrus.Fatalln("db reconnect open addr fail", err)
		return
	}
	err = db.DB().Ping()
	if err != nil {
		logrus.Fatalln("db reconnect ping fail", err)
		return
	}
	p.initByDBConfigs()
	logrus.WithField("db addr", p.dbConfig.DBAddr).Infoln("reconnect db success!")
}

//hmd重连接
func (p *HmdGormDB) hmdreConnect() {
	db, err := gorm.Open("mysql", p.dbConfig.HmdDBAddr)
	if err != nil {
		logrus.Fatalln("db reconnect open addr fail", err)
		return
	}
	err = db.DB().Ping()
	if err != nil {
		logrus.Fatalln("db reconnect ping fail", err)
		return
	}
	p.initByDBConfigs()
	logrus.WithField("db addr", p.dbConfig.HmdDBAddr).Infoln("reconnect db success!")
}

type DBConfig struct {
	DBAddr            string
	AutoCreateTables  []interface{} //自动创建的表，不设置则不创建表
	MaxIdleConns      int           //数据库连接池设置—— 最大空闲数，不设置则为10
	MaxOpenConns      int           //数据库连接池设置—— 最大打开的连接数，不设置则为100
	LogMode           bool          //是否打印gorm的日志 配置文件是打印
	DetectionInterval int           //心跳检测间隔，单位为s，默认30s,小于0则不检测
}

//hmd
type HmdDBConfig struct {
	HmdDBAddr         string
	AutoCreateTables  []interface{} //自动创建的表，不设置则不创建表
	MaxIdleConns      int           //数据库连接池设置—— 最大空闲数，不设置则为10
	MaxOpenConns      int           //数据库连接池设置—— 最大打开的连接数，不设置则为100
	LogMode           bool          //是否打印gorm的日志 配置文件是打印
	DetectionInterval int           //心跳检测间隔，单位为s，默认30s,小于0则不检测
}

type HSDZDBConfig struct {
	HSDZDBAddr        string
	AutoCreateTables  []interface{} //自动创建的表，不设置则不创建表
	MaxIdleConns      int           //数据库连接池设置—— 最大空闲数，不设置则为10
	MaxOpenConns      int           //数据库连接池设置—— 最大打开的连接数，不设置则为100
	LogMode           bool          //是否打印gorm的日志 配置文件是打印
	DetectionInterval int           //心跳检测间隔，单位为s，默认30s,小于0则不检测
}

func (p *DBConfig) check() error {
	if p.DBAddr == "" {
		logrus.Println("empty sql addr")
		return fmt.Errorf("empty sql addr")
	}
	if p.MaxIdleConns <= 0 {
		p.MaxIdleConns = 10
	}
	if p.MaxOpenConns <= 0 {
		p.MaxOpenConns = 100
	}
	if p.DetectionInterval == 0 {
		p.DetectionInterval = 30
	}
	return nil
}

//hmd
func (p *HmdDBConfig) hmdcheck() error {
	if p.HmdDBAddr == "" {
		logrus.Println("empty sql addr")
		return fmt.Errorf("empty sql addr")
	}
	if p.MaxIdleConns <= 0 {
		p.MaxIdleConns = 10
	}
	if p.MaxOpenConns <= 0 {
		p.MaxOpenConns = 100
	}
	if p.DetectionInterval == 0 {
		p.DetectionInterval = 30
	}
	return nil
}

// 本方法会给HmdGormClient赋值，多次调用HmdGormClient指向最后一次调用的GormDB
func HSDZInitGormDB(HSDZstr string) *gorm.DB {
	logrus.Infoln("starting hs-dz db")

	db, err := gorm.Open("mysql", HSDZstr)
	if err != nil {
		logrus.Println("hs-dz db initing fail", err)
		return nil
	}
	err = db.DB().Ping()
	if err != nil {
		logrus.Println("hs-dz db ping fail", err)
		return nil
	}

	return db
}

// Param 分页参数
type Param struct {
	DB        *gorm.DB
	PageIndex int
	PageSize  int
	OrderBy   []string
	ShowSQL   bool
}

type Pagination struct {
	CurrentPage int `json:"current_page" form:"current_page"`
	PageSize    int `json:"page_size" form:"page_size"`
	LastPage    int `json:"last_page"`
	Total       int `json:"total" form:"total"`
}

// 分页查询
func Paging(p Param, result interface{}) (Pagination, error) {
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			p.DB = p.DB.Order(o)
		}
	}

	if p.PageIndex == 0 && p.PageSize == 0 {
		if err := p.DB.Find(result).Error; err != nil {
			logrus.Errorf("Paging db get record err: %v", err.Error())
			return Pagination{}, err
		}
		return Pagination{}, nil
	}

	pagination := Pagination{}
	if p.PageIndex <= 0 {
		p.PageIndex = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			p.DB = p.DB.Order(o)
		}
	}

	totalCount := 0
	err := p.DB.Count(&totalCount).Error
	if err != nil {
		logrus.Errorf("Paging db get count err: %v", err.Error())
		return pagination, err
	}
	pagination.Total = totalCount
	pagination.LastPage = totalCount/p.PageSize + 1
	if p.PageIndex > pagination.LastPage {
		p.PageIndex = pagination.LastPage
	}

	if err := p.DB.Limit(p.PageSize).Offset((p.PageIndex - 1) * p.PageSize).Find(result).Error; err != nil {
		logrus.Errorf("Paging db get record err: %v", err.Error())
		return pagination, err
	}

	pagination.CurrentPage = p.PageIndex
	pagination.PageSize = p.PageSize
	return pagination, nil
}
