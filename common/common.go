package common

import (
	"order-backend/model"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
)

var DB *xorm.EngineGroup
var Log *Logger
var Config *ServerConfig
var RedisPool *redis.Pool
var App *AppConfig

//全局APP信息,在auth中间件中付值
var APP model.App
var OrderS OrderService
func LoadConfig() error {
	Config = &ServerConfig{}
	return Config.load()
}
func SetupLogger() error {
	var err error
	Log, err = NewLogger(Config.Log.LogFile, Config.Log.TraceLevel)
	return err
}

