package global

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"libs/log"
	"time"
	"zhibiaocalcsvr/src/etc"
)

func InitRedis() (hqredis, zxredis *redis.Pool) {
	zxredis = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			opts := []redis.DialOption{}
			if etc.Config.Redis.ZiXun.Pwd != "" {
				opts = append(opts, redis.DialPassword(etc.Config.Redis.ZiXun.Pwd))
			}

			conn, err := redis.Dial("tcp",
				fmt.Sprintf("%s:%d", etc.Config.Redis.ZiXun.IP, etc.Config.Redis.ZiXun.Port),
				opts...)
			if err != nil {
				log.Error("zixun redis start failed ,err = %v", err)
				return nil, err
			}
			return conn, err
		},
		MaxIdle:     100,
		MaxActive:   200,
		IdleTimeout: time.Minute,
	}

	hqredis = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			opts := []redis.DialOption{}
			if etc.Config.Redis.HangQing.Pwd != "" {
				opts = append(opts, redis.DialPassword(etc.Config.Redis.HangQing.Pwd))
			}

			conn, err := redis.Dial("tcp",
				fmt.Sprintf("%s:%d", etc.Config.Redis.HangQing.IP, etc.Config.Redis.HangQing.Port),
				opts...)
			if err != nil {
				log.Error("hangqing redis start failed ,err = %v", err)
				return nil, err
			}
			return conn, err
		},
		MaxIdle:     100,
		MaxActive:   200,
		IdleTimeout: time.Minute,
	}
	return
}
