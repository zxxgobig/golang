package global

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
)

const (
	Service_Priority_Calc = 1 + iota
)

type ICalcService interface {
	AlarmError(err_msg string, level int)

	//GetWheelTimer() timer.WheelTimer
	//GetGoServer() *g.Go
	GetMysqlEngine() *sql.DB
	GetMongoEngine() *mgo.Session
	GetRedisHQEngine() *redis.Pool
	GetRedisZXEngine() *redis.Pool
	GetZhiBiaoMgr() IZhiBiao
}

var GServer ICalcService
