package database

import (
	"fmt"
	"time"
	_ "time/tzdata"

	"jericho-go/settings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySql struct {
	Schema   string
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Charset  string
}

var (
	mysqlConn *gorm.DB
)

// NewMySql 构造函数
func NewMySql() *MySql {
	return &MySql{}
}

// getConn 获取数据库链接
func (receiver *MySql) getConn(dbConnName string) (db *gorm.DB) {
	time.Local = time.FixedZone("CST", 8*3600) // 由于MySQL数据库没有时区,所以时区自动改为0

	config := settings.NewSetting()
	if dbConnName == "" {
		defaultDbConnName := config.DB.Section("db").Key("db-conn-name").MustString("default")
		dbConnName = defaultDbConnName
	}

	receiver.Username = config.DB.Section(dbConnName).Key("username").MustString("root")
	receiver.Password = config.DB.Section(dbConnName).Key("password").MustString("root")
	receiver.Host = config.DB.Section(dbConnName).Key("host").MustString("127.0.0.1")
	receiver.Port = config.DB.Section(dbConnName).Key("port").MustString("3306")
	receiver.Database = config.DB.Section(dbConnName).Key("database").MustString("fix_workshop")
	receiver.Charset = config.DB.Section(dbConnName).Key("charset").MustString("utf8mb4")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		receiver.Username,
		receiver.Password,
		receiver.Host,
		receiver.Port,
		receiver.Database,
		receiver.Charset,
	)

	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize:                          500,
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	db = db.Session(&gorm.Session{
		SkipDefaultTransaction: false, // 开启自动事务
		QueryFields:            true,
		PrepareStmt:            true,
		AllowGlobalUpdate:      false, // 不允许全局修改,必须带有条件
	})

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 空闲复用时长

	return
}

// GetConn 获取数据库链接
func (receiver *MySql) GetConn(dbConnName string) *gorm.DB {
	if mysqlConn == nil {
		mysqlConn = receiver.getConn(dbConnName)
	}

	return mysqlConn
}

// NewConn 获取新数据库链接
func (receiver *MySql) NewConn(dbConnName string) *gorm.DB {
	return receiver.getConn(dbConnName)
}
