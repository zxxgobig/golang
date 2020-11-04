package global

import (
	"gopkg.in/mgo.v2"
	"time"
	"zhibiaocalcsvr/src/etc"
)

func InitMongo() (err error, mgosession *mgo.Session) {
	dial_info := &mgo.DialInfo{
		Addrs:    []string{etc.Config.Mongo.IP},
		Timeout:  time.Second * 3,
		Username: etc.Config.Mongo.User,
		Password: etc.Config.Mongo.Passwd,
	}

	msession, err := mgo.DialWithInfo(dial_info)
	return err, msession
}
