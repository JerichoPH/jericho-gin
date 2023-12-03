package database

import (
	"fmt"
	"jericho-gin/settings"
	"jericho-gin/wrongs"
	"os"

	"gorm.io/gorm"
)

// GormLauncher gorm启动器
type GormLauncher struct {
	DbDriver string
}

// NewGormLauncher 初始化数据库启动器
func NewGormLauncher() *GormLauncher {
	return &GormLauncher{}
}

// GetDbDriver 获取数据库设置
func (receiver *GormLauncher) GetDbDriver() string {
	return receiver.DbDriver
}

// SetDbDriver 设置数据库
func (receiver *GormLauncher) SetDbDriver(dbDriver string) *GormLauncher {
	receiver.DbDriver = dbDriver
	return receiver
}

// GetConn 获取当前数据库链接
func (receiver *GormLauncher) GetConn(dbConnName string) (dbSession *gorm.DB) {
	var (
		dbDriver string
		config   *settings.Setting
	)
	config = settings.NewSetting()

	if receiver.DbDriver != "" {
		dbDriver = receiver.DbDriver
	} else {
		dbDriver = config.DB.Section("db").Key("db-driver").MustString("")
	}

	if dbConnName == "" {
		dbConnName = config.DB.Section("db").Key("db-conn-name").MustString("default")
	}

	switch dbDriver {
	case "postgresql":
		dbSession = NewPostgreSql().GetConn(dbConnName)
	case "mysql":
		dbSession = NewMySql().GetConn(dbConnName)
	case "mssql":
		dbSession = NewMsSql().GetConn(dbConnName)
	default:
		wrongs.ThrowForbidden("没有配置数据库")
	}
	return
}

// BootAutoMigrate 初始化数据库迁移
func (receiver *GormLauncher) BootAutoMigrate(dst []any) {
	if errAutoMigrate := receiver.GetConn("").AutoMigrate(dst...); errAutoMigrate != nil {
		fmt.Println("数据库迁移错误:", errAutoMigrate.Error())
		os.Exit(1)
	}
}

// NewConn 获取新数据库链接
func (receiver *GormLauncher) NewConn(dbConnName string) (dbSession *gorm.DB) {
	var (
		dbDriver string
		config   *settings.Setting
	)
	config = settings.NewSetting()

	if receiver.DbDriver != "" {
		dbDriver = receiver.DbDriver
	} else {
		dbDriver = config.DB.Section("db").Key("db-driver").MustString("")
	}

	if dbConnName == "" {
		dbConnName = config.DB.Section("db").Key("db-conn-name").MustString("mysql")
	}

	switch dbDriver {
	case "postgresql":
		dbSession = NewPostgreSql().NewConn(dbConnName)
	case "mysql":
		dbSession = NewMySql().NewConn(dbConnName)
	case "mssql":
		dbSession = NewMsSql().NewConn(dbConnName)
	default:
		wrongs.ThrowForbidden("没有配置数据库")
	}
	return
}
