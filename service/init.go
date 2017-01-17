package service

import (
	"github.com/jinzhu/gorm"
	"github.com/zwh8800/66ana/conf"
	"gopkg.in/redis.v5"
)

var dbConn *gorm.DB
var redisClient *redis.Client

func init() {
	var err error
	dbConn, err = gorm.Open(conf.Conf.DB.Driver, conf.Conf.DB.Dsn)
	if err != nil {
		panic(err)
	}
	dbConn.DB().SetMaxOpenConns(conf.Conf.DB.MaxConnection)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.Redis.Addr,
		Password: conf.Conf.Redis.Password,
		DB:       conf.Conf.Redis.DB,
	})
	err = redisClient.Ping().Err()
	if err != nil {
		panic(err)
	}
}
