package calc

import (
	"database/sql"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"libs/chanrpc"
	"libs/go"
	"libs/log"
	"libs/utils"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"zhibiaocalcsvr/src/calc/zhibiao"
	"zhibiaocalcsvr/src/calc/zhibiao/cpgj"
	"zhibiaocalcsvr/src/etc"
	"zhibiaocalcsvr/src/global"
)

type CalcService struct {
	mysql_     *sql.DB
	dev_mysql_ *sql.DB
	mongo_     *mgo.Session
	hqredis_   *redis.Pool
	zxredis_   *redis.Pool
	zhibiao_   global.IZhiBiao
	utils.Pooller

	bs_calc *cpgj.BSCalculator //买卖点计算器
}

func (ser *CalcService) Start() error {
	log.Release("zhibiaocalc start")
	global.GServer = ser

	var err error
	ser.GoServer = g.New(etc.Go_Len)
	ser.ChanServer = chanrpc.NewServer(etc.Chan_Server_Len)

	// 初始化redis
	ser.hqredis_, ser.zxredis_ = global.InitRedis()
	log.Release("redis  初始化成功。----")

	// 初始化mysql
	err, ser.mysql_ = global.InitMysqlMaster()
	if err != nil {
		log.Error("calc service start failed, mysql init failed, err = %v", err)
		return err
	} else {
		log.Release("mysql  初始化成功。----")
	}

	err, ser.dev_mysql_ = global.InitMysqlDev()
	if err != nil {
		log.Error("calc service start failed, mysql dev init failed, err = %v", err)
		return err
	} else {
		log.Release("mysql-dev 初始化成功。----")
	}

	// 初始化mongo
	err, ser.mongo_ = global.InitMongo()
	if err != nil {
		log.Error("calc service start failed, mongo init failed, err = %v", err)
		return err
	} else {
		log.Release("mongo  初始化成功。----")
	}

	//启动指标计算模块
	ser.zhibiao_ = zhibiao.ZhiBiaoMgr
	err = ser.zhibiao_.Start()
	if err != nil {
		global.GServer.AlarmError(fmt.Sprintf("zhibiaocalcsvr start failed,err:%v",err.Error()),1)
		return err
	}

	return nil
}

// Pool应为Poll
func (self *CalcService) Pool(cs chan bool) {
	time_up := time.NewTicker(time.Duration(etc.Config.Calc.PrintStatus) * time.Second)
	for {
		select {
		case <-cs:
			self.Close()
			return
		case cc := <-self.GoServer.ChanCb:
			self.GoServer.Cb(cc)
		case <-time_up.C:
			log.Release("server is running")
		}
	}
}

func (self *CalcService) GetName() string {
	return "指标定时计算服务"
}

func (self *CalcService) Close() {
	log.Release("Zhibiaocalc close")
	self.GoServer.Close()
	self.zhibiao_.Stop()
	self.mysql_.Close()
	self.hqredis_.Close()
	self.zxredis_.Close()
	self.mongo_.Close()
}

func (self *CalcService) GetPriority() int {
	return global.PRIORITY_CALC
}

func (ser *CalcService) GetGoServer() *g.Go {
	return ser.GoServer
}

func (ser *CalcService) GetMysqlEngine() *sql.DB {
	return ser.mysql_
}

func (ser *CalcService) GetMongoEngine() *mgo.Session {
	return ser.mongo_
}

func (ser *CalcService) GetRedisHQEngine() *redis.Pool {
	return ser.hqredis_
}

func (ser *CalcService) GetRedisZXEngine() *redis.Pool {
	return ser.zxredis_
}

func (ser *CalcService) GetZhiBiaoMgr() global.IZhiBiao {
	return ser.zhibiao_
}

func (ser *CalcService) AlarmError(err_msg string, level int) {
	ser.GoServer.Go(func() {
		cli := &http.Client{
			Timeout: time.Second * 2,
		}
		resp, err := cli.PostForm(etc.Config.Alarm_url, url.Values{
			"ServerName": []string{fmt.Sprintf("zhibiaocalcsvr_%d", etc.Config.ID)},
			"ServerIp":   []string{""},
			"ServerPort": []string{""},
			"ErrMsg":     []string{err_msg},
			"ErrLv":      []string{strconv.Itoa(level)},
		})

		if err != nil {
			log.Error("AlarmError %s post error:%v", err_msg, err)
			return
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			log.Release("AlarmError %s resp:%s", err_msg, string(data))
		} else {
			log.Error("AlarmError %s read resp err:%v", err_msg, err)
		}
	}, func() {})
}

func (self *CalcService) clearMemoryData() {
	global.StockData = nil
	global.FuQuanData = nil
}
