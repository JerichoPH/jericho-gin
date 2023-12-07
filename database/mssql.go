package database

import (
	"fmt"
	"jericho-gin/settings"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type MsSql struct {
	Schema   string
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// NewMsSql 构造函数
func NewMsSql() *MsSql {
	return &MsSql{}
}

var msSqlConn *gorm.DB

func (receiver *MsSql) getConn(dbConnName string) (db *gorm.DB) {
	config := settings.NewSetting()

	receiver.Username = config.DB.Section(dbConnName).Key("username").MustString("sa")
	receiver.Password = config.DB.Section(dbConnName).Key("password").MustString("JW087073yjz..")
	receiver.Host = config.DB.Section(dbConnName).Key("host").MustString("127.0.0.1")
	receiver.Port = config.DB.Section(dbConnName).Key("port").MustString("1433")
	receiver.Database = config.DB.Section(dbConnName).Key("databases").MustString("Dwqcgl")

	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%s?database=%s",
		receiver.Schema,
		receiver.Username,
		receiver.Password,
		receiver.Host,
		receiver.Port,
		receiver.Database,
	)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 空闲复用时长

	return
}

// GetConn 获取数据库链接
func (receiver *MsSql) GetConn(dbConnName string) *gorm.DB {
	if msSqlConn == nil {
		msSqlConn = receiver.getConn(dbConnName)
	}
	return msSqlConn
}

// NewConn 获取新数据库链接
func (receiver *MsSql) NewConn(dbConnName string) *gorm.DB {
	return receiver.getConn(dbConnName)
}
