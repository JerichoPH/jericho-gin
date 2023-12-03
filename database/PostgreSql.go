package database

import (
	"fmt"
	"jericho-gin/settings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSql struct {
	Host     string
	Port     string
	Username string
	Database string
	Password string
	SSLMode  string
}

// NewPostgreSql 构造函数
func NewPostgreSql() *PostgreSql {
	return &PostgreSql{}
}

var postgresqlConn *gorm.DB

// getConn 获取链接
//
//	@receiver receiver
//	@return db
func (receiver *PostgreSql) getConn(dbConnName string) (db *gorm.DB) {
	config := settings.NewSetting()

	receiver.Host = config.DB.Section(dbConnName).Key("host").MustString("127.0.0.1")
	receiver.Port = config.DB.Section(dbConnName).Key("port").MustString("5432")
	receiver.Username = config.DB.Section(dbConnName).Key("username").MustString("postgres")
	receiver.Password = config.DB.Section(dbConnName).Key("password").MustString("zces@1234")
	receiver.Database = config.DB.Section(dbConnName).Key("database").MustString("postgres")
	receiver.SSLMode = config.DB.Section(dbConnName).Key("ssl_mode").MustString("disable")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		receiver.Host,
		receiver.Port,
		receiver.Username,
		receiver.Database,
		receiver.Password,
		receiver.SSLMode,
	)

	mySqlConn, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize:                          1000,
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	db = mySqlConn.Session(&gorm.Session{
		SkipDefaultTransaction: true,
		QueryFields:            true,
		PrepareStmt:            true,
	})

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 空闲复用时长

	return
}

// GetConn 获取数据库链接
func (receiver *PostgreSql) GetConn(dbConnName string) *gorm.DB {
	if postgresqlConn == nil {
		postgresqlConn = receiver.getConn(dbConnName)
	}
	return postgresqlConn
}

// NewConn 获取新数据库链接
func (receiver *PostgreSql) NewConn(dbConnName string) *gorm.DB {
	return receiver.getConn(dbConnName)
}
