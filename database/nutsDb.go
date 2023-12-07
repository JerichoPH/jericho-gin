package database

import (
	"github.com/nutsdb/nutsdb"
	"log"
)

// NutsDb nuts-db
type NutsDb struct {
}

var nutsDbConn *nutsdb.DB

// NewNutsDb 初始化nuts-db对象
func NewNutsDb() *NutsDb {
	return &NutsDb{}
}

// 获取nuts-db数据库链接
func (receiver NutsDb) getConn() *nutsdb.DB {
	var err error
	db, err := nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir("./nuts-db"),
	)
	if err != nil {
		log.Fatalf("nuts-db 启动失败：%s", err.Error())
	}
	defer func(db *nutsdb.DB) {
		err = db.Close()
		if err != nil {
			log.Fatalf("nuts-db 关闭失败：%s", err.Error())
		}
	}(db)

	return db
}

// GetConn 获取nuts-db数据库链接
func (receiver NutsDb) GetConn() *nutsdb.DB {
	if nutsDbConn == nil {
		nutsDbConn = receiver.getConn()
	}

	return nutsDbConn
}

// NewConn 获取全新的nuts-db数据库链接
func (receiver NutsDb) NewConn() *nutsdb.DB {
	return receiver.getConn()
}
