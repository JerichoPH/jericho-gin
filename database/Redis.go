package database

import (
	"context"
	"fmt"
	"jericho-gin/settings"
	"jericho-gin/types"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Host     string
	Port     string
	Database uint64
	Password string
	Client   *redis.Client
}

func NewRedis(database uint64) *Redis {
	return &Redis{}
}

// GetConn 获取redis链接
func (receiver Redis) GetConn(database uint64) *redis.Client {
	var err error
	config := settings.NewSetting()

	receiver.Host = config.DB.Section("redis").Key("host").MustString("127.0.0.1")
	receiver.Port = config.DB.Section("redis").Key("port").MustString("6379")
	receiver.Database = config.DB.Section("redis").Key("database").MustUint64(0)
	receiver.Password = config.DB.Section("redis").Key("password").String()
	if database == 0 {
		database = receiver.Database
	}

	receiver.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", receiver.Host, receiver.Port),
		Password: receiver.Password, // no password set
		DB:       int(database),     // use default GetDb
	})

	defer func(db *redis.Client) {
		err = db.Close()
		if err != nil {
			log.Printf("redis 关闭失败：%s", err.Error())
		}
	}(receiver.Client)

	return receiver.Client
}

// SetValue 设置值
func (receiver Redis) SetValue(key string, value interface{}, exp time.Duration) (any, error) {
	val, err := receiver.Client.Set(context.Background(), key, value, exp).Result()
	return val, err
}

// GetValue 获取值
func (receiver Redis) GetValue(key string) (any, error) {
	value, err := receiver.Client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return "", err
	}
	return value, nil
}

// Sort 排序获取值
func (receiver Redis) Sort(key string, offset, count int64, order string) ([]string, error) {
	values, err := receiver.Client.Sort(context.Background(), key, &redis.Sort{Offset: offset, Count: count, Order: order}).Result()
	return values, err
}

// ZRange 获取字典排序
func (receiver Redis) ZRange(key, min, max string, offset, count int64) ([]redis.Z, error) {
	values, err := receiver.Client.ZRangeByScoreWithScores(context.Background(), key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Result()

	return values, err
}

// ZInterStore
func (receiver Redis) ZInterStore(keys []string, weights types.ListFloat64) (int64, error) {
	values, err := receiver.Client.ZInterStore(context.Background(), "out", &redis.ZStore{
		Keys:    keys,
		Weights: weights,
	}).Result()

	return values, err
}

// Do 执行命令
func (receiver Redis) Do(command, key string, value any) (any, error) {
	res, err := receiver.Client.Do(context.Background(), command, key, value).Result()
	return res, err
}
