package zbbase

import (
	"errors"
	"fmt"
	"libs/log"
	"time"
	"zhibiaocalcsvr/src/etc"
	"zhibiaocalcsvr/src/global"
)

var (
	M_cmfb = "cmfb" //筹码分布
	M_hlqj = "hlqj" //黄蓝区间
	M_jdcj = "jdcj" //见底出击
	M_cptx = "cptx" //见底出击
)

const (
	SQL_OPR_BATCH_CNT = 1000
)

type ZhibiaoModule interface {
	Calc() (error, uint32)
	Start() error
	Init()
	Name() string
	GetLastCalcDate(string, string, string)
	UpdateCalcDateToMySql(string, string, uint32) error
}

type BaseModule struct {
	cm_            ZhibiaoModule
	TBname         string // 表名
	Filedname      string // 字段名
	Last_calc_date uint32 // 最近的一次計算的時間。
}

// 获取最后一次计算时间
func (bm *BaseModule) GetLastCalcDate(tablename, filedname, calcname string) {
	defer func() {
		log.Release("GetLastCalcDate() --%v--Last_calc_date = %v", calcname, bm.Last_calc_date)
	}()

	var timestamp uint32
	location, _ := time.LoadLocation("Asia/Shanghai") //设置时区
	ts, err := time.ParseInLocation("20060102", fmt.Sprintf("%v", etc.Config.Calc.Common_StartDate), location)
	if err != nil {
		timestamp = 0
		log.Error("时间转换错误,err = %v", err)
		return
	}
	timestamp = uint32(ts.Unix())

	if tablename == "" || filedname == "" {
		log.Error("GetLastCalcDate() param is empty, tablename:%v, filedname:%v", tablename, filedname)
		bm.Last_calc_date = 0
		return
	}

	query_sql := fmt.Sprintf(`select calc_time from %v where zb_type="%s";`, tablename, filedname)
	log.Release("GetLastCalcDate() --%v-- sql = %v", calcname, query_sql)
	rows, err := global.GServer.GetMysqlEngine().Query(query_sql)
	if err != nil {
		if filedname == "cpgj" || filedname == "sighting" {
			bm.Last_calc_date = uint32(etc.Config.Calc.Common_StartDate)
		} else {
			bm.Last_calc_date = timestamp
		}
		return
	}

	last_calc_time := uint32(0)
	for rows.Next() {
		err = rows.Scan(&last_calc_time)
		if err != nil {
			//判断是否nil
			var uptime_inf interface{}
			err = rows.Scan(&uptime_inf)
			if err != nil {
				if filedname == "cpgj" || filedname == "sighting" {
					bm.Last_calc_date = uint32(etc.Config.Calc.Common_StartDate)
				} else {
					bm.Last_calc_date = timestamp
				}
				return
			}
		} else {
			bm.Last_calc_date = last_calc_time
		}
	}

	if last_calc_time == 0 {
		if filedname == "cpgj" || filedname == "sighting" {
			bm.Last_calc_date = uint32(etc.Config.Calc.Common_StartDate)
		} else {
			bm.Last_calc_date = timestamp
		}
	}

	rows.Close()
}

func (bm *BaseModule) UpdateCalcDateToMySql(tablename, filedname string, timestamp uint32) (err error) {
	if tablename == "" || filedname == "" || timestamp <= 0 {
		errnew := "updateCalcDateToMySql() param is empty"
		log.Error(errnew)
		return errors.New(errnew)
	}
	update_time_sql := fmt.Sprintf(`INSERT INTO %s(zb_type, calc_time, zb_value)
		VALUES ('%s', %d, '%s') ON DUPLICATE KEY UPDATE calc_time=%d, zb_value='%s';`, tablename, filedname, timestamp, "",
		timestamp, "")
	log.Release("update last calc date:%d, insert_sql = %v", timestamp, update_time_sql)

	_, err = global.GServer.GetMysqlEngine().Exec(update_time_sql)
	if err != nil {
		log.Error("%v calc donetime update failed, err:%v", err)
	}

	return err
}
