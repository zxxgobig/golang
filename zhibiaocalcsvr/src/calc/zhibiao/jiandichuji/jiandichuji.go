package jiandichuji

import (
	"common/structs"
	"fmt"
	"libs/log"
	"time"
	"zhibiaocalcsvr/src/calc/zhibiao/zbbase"
	"zhibiaocalcsvr/src/global"
)

const (
	Jdcj_type_short = 1 + iota
	Jdcj_type_mid
	Jdcj_type_long
)

type JianDiChuJi struct {
	zbbase.BaseModule
}

func (m *JianDiChuJi) Init() {
	m.Filedname = "jdcj"
	m.TBname = "db_zbjs_common.tb_lastcalcinfo"
}

/*
【见底出击】

HXLXA= 5日均线
HXLXB= 10日均线
HXLXD= (收盘价-13日均线)/13日均线*100/0.7
HXLXE= (收盘价-6日均线)/6日均线*100*2
HXLXF= (收盘价-11日均线)/11日均线*100/0.3/1.3/1.5

C1= 昨天HXLXF<-12 AND 今天HXLXF>=-12

C2= 昨天HXLXD<-12 今天 HXLXD>=-12
LJD= ((30日均线-最低价)/60日均线)*200

短期趋势=  昨天HXLXE<-9 并且 今天HXLXE>=-9 显示黄色三角
长期趋势=  昨天LJD>30 并且今天 LJD<=30  显示紫色柱子
中期趋势=  C1 和C2同时成立 显示红色三角};
*/

type jdcjItem struct {
	Code        string
	Trade_day   uint32
	Upsert_time uint32
	Type        int
}

func (m *JianDiChuJi) Name() string {
	return zbbase.M_jdcj
}
func (m *JianDiChuJi) Calc() (err error, new_calc_time uint32) {
	bt := time.Now()
	opr_items := []*jdcjItem{}
	total_calc_cnt := 0
	last_del_items := []*jdcjItem{} //最新的实时变化需要执行删除操作
	last_trad_day := uint32(0)

	//短期，中期，长期
	for code, kline_arr := range global.StockData {
		klen := len(kline_arr.Cbs)
		if klen == 0 {
			log.Error("JianDiChuJi::Calc code:%s kline arr == 0", code)
			continue
		}

		if kline_arr.Cbs[klen-1].Upserttime > new_calc_time {
			new_calc_time = kline_arr.Cbs[klen-1].Upserttime
		}

		//计算各类均线
		day_6_avg, _ := zbbase.MA(6, kline_arr.Cbs)   //6日
		day_11_avg, _ := zbbase.MA(11, kline_arr.Cbs) //11日
		day_13_avg, _ := zbbase.MA(13, kline_arr.Cbs) //13日
		day_30_avg, _ := zbbase.MA(30, kline_arr.Cbs) //30日
		day_60_avg, _ := zbbase.MA(60, kline_arr.Cbs) //30日

		var yes_kline *structs.KlineInfo
		for index, kline := range kline_arr.Cbs {
			if kline.Tradday > last_trad_day {
				last_trad_day = kline.Tradday
			}

			total_calc_cnt++
			last_one := index == klen-1
			if yes_kline == nil ||
				(kline.Upserttime < m.Last_calc_date && !last_one) { //保底计算最后一笔
				tmp := structs.KlineInfo{}
				tmp.Upsert_time = kline.Upserttime
				tmp.Update_time = kline.Upserttime
				tmp.Trad_day = kline.Tradday
				tmp.Kline_binfo = kline.KLine
				yes_kline = &tmp
				continue
			}

			//短期趋势=  昨天HXLXE<-9 并且 今天HXLXE>=-9 显示黄色三角
			if v, ok := day_6_avg[yes_kline.Upsert_time]; ok {
				yes_HXLXE := zbbase.HXLXE(yes_kline.Kline_binfo.Close, v.Value)
				tod_HXLXE := zbbase.HXLXE(kline.KLine.Close, day_6_avg[kline.Upserttime].Value)
				if yes_HXLXE < -9 && tod_HXLXE >= -9 {
					opr_items = append(opr_items, &jdcjItem{
						Code:        code,
						Trade_day:   kline.Tradday,
						Upsert_time: kline.Upserttime,
						Type:        Jdcj_type_short,
					})
				} else if last_one {
					last_del_items = append(last_del_items, &jdcjItem{
						Code:        code,
						Trade_day:   kline.Tradday,
						Upsert_time: kline.Upserttime,
						Type:        Jdcj_type_short,
					})
				}
			}

			//长期趋势=  昨天LJD>30 并且今天 LJD<=30  显示紫色柱子
			v60, ok60 := day_60_avg[yes_kline.Upsert_time]
			if ok60 {
				v30, _ := day_30_avg[yes_kline.Upsert_time]
				yes_LJD := zbbase.LJD(v30.Value, yes_kline.Kline_binfo.Low, v60.Value)
				tod_LJD := zbbase.LJD(day_30_avg[kline.Upserttime].Value, kline.KLine.Low, day_60_avg[kline.Upserttime].Value)
				if yes_LJD > 30 && tod_LJD <= 30 {
					opr_items = append(opr_items, &jdcjItem{
						Code:        code,
						Trade_day:   kline.Tradday,
						Upsert_time: kline.Upserttime,
						Type:        Jdcj_type_long,
					})

				} else if last_one {
					last_del_items = append(last_del_items, &jdcjItem{
						Code:        code,
						Trade_day:   kline.Tradday,
						Upsert_time: kline.Upserttime,
						Type:        Jdcj_type_long,
					})
				}
			}

			//中期趋势=  C1 和C2同时成立 显示红色三角;
			v13, ok13 := day_13_avg[yes_kline.Upsert_time]
			if ok13 {
				v11, _ := day_11_avg[yes_kline.Upsert_time]

				yes_HXLXD := zbbase.HXLXD(yes_kline.Kline_binfo.Close, v13.Value)
				yes_HXLXF := zbbase.HXLXF(yes_kline.Kline_binfo.Close, v11.Value)

				tod_HXLXD := zbbase.HXLXD(kline.KLine.Close, day_13_avg[kline.Upserttime].Value)
				tod_HXLXF := zbbase.HXLXF(kline.KLine.Close, day_11_avg[kline.Upserttime].Value)

				if zbbase.C1(yes_HXLXF, tod_HXLXF) && zbbase.C2(yes_HXLXD, tod_HXLXD) {
					opr_items = append(opr_items, &jdcjItem{
						Code:        code,
						Trade_day:   kline.Tradday,
						Upsert_time: kline.Upserttime,
						Type:        Jdcj_type_mid,
					})
				} else if last_one {
					last_del_items = append(last_del_items, &jdcjItem{
						Code:        code,
						Trade_day:   kline.Tradday,
						Upsert_time: kline.Upserttime,
						Type:        Jdcj_type_mid,
					})
				}
			}

			yes_kline.Upsert_time = kline.Upserttime
			yes_kline.Update_time = kline.Upserttime
			yes_kline.Trad_day = kline.Tradday
			yes_kline.Kline_binfo = kline.KLine
		}
	}
	calc_time := time.Now()

	mysql_engine := global.GServer.GetMysqlEngine()

	//批量操作
	index := 0

	for len(last_del_items) > 0 {
		bg_index := index * zbbase.SQL_OPR_BATCH_CNT
		end_index := bg_index + zbbase.SQL_OPR_BATCH_CNT
		if end_index >= len(last_del_items) {
			end_index = len(last_del_items)
		}

		items := last_del_items[bg_index:end_index]
		//先删除
		del_sql := "delete from db_zbjs_hjjj.tb_jdcj where "
		for i, item := range items {
			line_sql := ""
			if i != 0 {
				line_sql += " or "
			}
			line_sql += fmt.Sprintf(`(stockcode="%s" and sdata=%d and upserttime=%d and `+"`rsttype`=%d)",
				item.Code, item.Trade_day, item.Upsert_time, item.Type)
			del_sql += line_sql
		}
		del_sql += ";"
		_, err = mysql_engine.Exec(del_sql)
		if err != nil {
			log.Release("JianDiChuJi::Calc last one del_sql:%s err:%v", del_sql, err)
			return err, new_calc_time
		}

		if end_index == len(last_del_items) {
			break
		}

		index++
	}

	index = 0
	today_cnt := 0
	for len(opr_items) > 0 {
		bg_index := index * zbbase.SQL_OPR_BATCH_CNT
		end_index := bg_index + zbbase.SQL_OPR_BATCH_CNT
		if end_index >= len(opr_items) {
			end_index = len(opr_items)
		}

		items := opr_items[bg_index:end_index]
		//先删除
		del_sql := "delete from db_zbjs_hjjj.tb_jdcj where "
		for i, item := range items {
			line_sql := ""
			if i != 0 {
				line_sql += " or "
			}
			line_sql += fmt.Sprintf(`(stockcode="%s" and sdata = %d and upserttime=%d and `+"`rsttype`=%d)",
				item.Code, item.Trade_day, item.Upsert_time, item.Type)
			del_sql += line_sql
		}
		del_sql += ";"
		_, err = mysql_engine.Exec(del_sql)
		if err != nil {
			log.Release("JianDiChuJi::Calc del_sql:%s err:%v", del_sql, err)
			return err, new_calc_time
		}

		//插入
		insert_sql := "insert into db_zbjs_hjjj.tb_jdcj(stockcode, sdata, upserttime, `rsttype`) values "
		for i, item := range items {
			line_sql := ""
			if i != 0 {
				line_sql += ","
			}
			line_sql += fmt.Sprintf(`("%s",%d,%d,%d)`, item.Code, item.Trade_day, item.Upsert_time, item.Type)
			insert_sql += line_sql
			if item.Trade_day == last_trad_day {
				today_cnt++
			}
		}
		insert_sql += ";"
		_, err = mysql_engine.Exec(insert_sql)
		if err != nil {
			log.Release("JianDiChuJi::Calc insert_sql err:%v", err)
			return err, new_calc_time
		}

		if end_index == len(opr_items) {
			break
		}

		index++
	}

	log.Release("JianDiChuJi Calc stock total time:%v, calc time:%v, calc_cnt:%d, last_del_cnt:%d, suc_cnt:%d, %d:%d",
		time.Now().Sub(bt), calc_time.Sub(bt), total_calc_cnt, len(last_del_items), len(opr_items), last_trad_day, today_cnt)
	return nil, new_calc_time
}

// 开始见底出击指标计算
func (m *JianDiChuJi) Start() error {
	// 获取最后一次计算时间。
	m.GetLastCalcDate(m.TBname, m.Filedname, m.Name())

	// 计算。
	err, timestmap := m.Calc()
	if err != nil {
		return err
	}
	// 更新计算时间。
	err = m.UpdateCalcDateToMySql(m.TBname, m.Filedname, timestmap)

	return err
}
