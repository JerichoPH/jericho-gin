package database

import (
	"context"
	"fmt"
	"jericho-gin/settings"
	"jericho-gin/wrongs"
	"net/url"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Schema   string
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

var mongoConn *mongo.Client

// NewMongo 构造函数
func NewMongo() *Mongo {
	return &Mongo{}
}

// 获取链接
func (receiver *Mongo) getConn(dbConnName string) *mongo.Client {
	config := settings.NewSetting()

	receiver.Username = config.DB.Section(dbConnName).Key("username").MustString("root")
	receiver.Password = config.DB.Section(dbConnName).Key("password").MustString("zces@1234")
	receiver.Host = config.DB.Section(dbConnName).Key("host").MustString("127.0.0.1")
	receiver.Port = config.DB.Section(dbConnName).Key("port").MustInt(27017)
	receiver.Database = config.DB.Section(dbConnName).Key("database").MustString("admin")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb://" + url.QueryEscape("127.0.0.1:27017")).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		wrongs.ThrowForbidden(fmt.Sprintf("数据库链接失败：%s", err.Error()))
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			wrongs.ThrowForbidden(fmt.Sprintf("数据库关闭失败：%s", err.Error()))
		}
	}()

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		wrongs.ThrowForbidden(err.Error())
	}

	return client
}

// GetConn 获取数据库链接
func (receiver *Mongo) GetConn(dbConnName string) *mongo.Client {
	if mongoConn == nil {
		mongoConn = receiver.getConn(dbConnName)
	}
	return mongoConn
}

// NewConn 获取新数据库链接
func (receiver *Mongo) NewConn(dbConnName string) *mongo.Client {
	return receiver.getConn(dbConnName)
}
