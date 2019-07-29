package common

import (
	"github.com/gomodule/redigo/redis"
	"time"
)


func OpenRedis() error {
	var err error
	RedisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(Config.Redis.Network, Config.Redis.Address, redis.DialDatabase(Config.Redis.Database), redis.DialPassword(Config.Redis.Password))
			if Config.Redis.ShowCommand {
				conn = redis.NewLoggingConn(conn, Log.info, "[redis]")
			}
			return conn, err
		},
	}
	return err
}
