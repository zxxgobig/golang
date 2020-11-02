package choumafenbu

import (
	"fmt"
	"libs/log"
	"time"
	"zhibiaocalcsvr/src/calc/zhibiao/zbbase"
	"zhibiaocalcsvr/src/global"
)

const (
	Cmfb_type_ji  = iota //集
	Cmfb_type_fen        //分
)

const (
	before_kline_len = 20
)

type cmfbItem struct {
	Code        string
	Trade_day   uint32
	Upsert_time uint32
	Type        int
}

type lastCalcInfo struct {
	abs_sma    *zbbase.MathItem
	max_sma    *zbbase.MathItem
	upserttime uint32
}

type ChouMaFenBu struct {
	zbbase.BaseModule

	last_calc_info_ map[string]*lastCalcInfo
}

func (m *ChouMaFenBu) Init() {
	m.last_calc_info_ = make(map[string]*lastCalcInfo)
	m.Filedname = "cmfb"
	m.TBname = "db_zbjs_common.tb_lastcalcinfo"
}

func (m *ChouMaFenBu) Name() string {
	return zbbase.M_cmfb
}

func (m *ChouMaFenBu) Calc() (error, uint32) {
	bt := time.Now()
	opr_items := []*cmfbItem{}
	del_items := []*cmfbItem{}
	last_trad_day := uint32(0)
	var new_calc_time uint32
	var err error

	//短期，中期，长期
	for code, kline_arr := range global.StockData {
		klen := len(kline_arr.Cbs)
		if klen == 0 {
			log.Error("ChouMaFenBu Calc code:%s kline arr == 0", code)
			continue
		}

		// 指数部分 不参与计算
		if _, ok := global.GServer.GetZhiBiaoMgr().GetIndexCodes()[code]; ok {
			continue
		}

		klinedata := make([]*global.CombineItem, klen)
		copy(klinedata, kline_arr.Cbs)

		if klinedata[klen-1].Upserttime > new_calc_time {
			new_calc_time = klinedata[klen-1].Upserttime
		}

		var mm_arr, max_sma, abs_sma []*zbbase.MathItem
		if last_calc_info, ok := m.last_calc_info_[code]; ok {
			var last_index int
			for last_index = klen - 1; last_index >= 0; last_index-- {
				if klinedata[last_index].Upserttime == last_calc_info.upserttime {
					break
				}
			}

			klinedata = klinedata[last_index:]
			klen = len(klinedata)
			mm_arr, max_sma, abs_sma = zbbase.MM(klinedata, last_calc_info.max_sma, last_calc_info.abs_sma)
		} else {
			mm_arr, max_sma, abs_sma = zbbase.MM(klinedata, nil, nil)
		}

		//计算MM值
		mm_len := len(mm_arr)
		sma_len := len(abs_sma)
		/*
			昨天的MM大于0 并且  MM小于昨天的MM  写上 ‘分’ 字
			昨天的MM小于-200 并且 MM大于昨天的MM  写上'集' 字}
		*/

		if sma_len > before_kline_len {
			m.last_calc_info_[code] = &lastCalcInfo{
				upserttime: klinedata[sma_len-before_kline_len].Upserttime, //ref值 1
				abs_sma:    abs_sma[sma_len-1-before_kline_len],
				max_sma:    max_sma[sma_len-1-before_kline_len],
			}
		}

		for i := 0; i < mm_len; i++ {
			last_one := (i == 0)
			cur_index := mm_len - 1 - i
			last_index := cur_index - 1
			if last_index < 0 {
				break
			}

			cur_kline := klinedata[klen-1-i]
			if cur_kline.Tradday > last_trad_day {
				last_trad_day = cur_kline.Tradday
			}

			if cur_kline.Upserttime < m.Last_calc_date && !last_one { //计算的起点
				break
			}

			if m.Last_calc_date != 0 {
				del_items = append(del_items, &cmfbItem{
					Code:        code,
					Trade_day:   cur_kline.Tradday,
					Upsert_time: cur_kline.Upserttime,
					Type:        Cmfb_type_fen,
				})
				del_items = append(del_items, &cmfbItem{
					Code:        code,
					Trade_day:   cur_kline.Tradday,
					Upsert_time: cur_kline.Upserttime,
					Type:        Cmfb_type_ji,
				})
			}

			if mm_arr[last_index].Value > 0 && mm_arr[cur_index].Value < mm_arr[last_index].Value { //分
				opr_items = append(opr_items, &cmfbItem{
					Code:        code,
					Trade_day:   cur_kline.Tradday,
					Upsert_time: cur_kline.Upserttime,
					Type:        Cmfb_type_fen,
				})
			} else if mm_arr[last_index].Value < float64(-200) && mm_arr[cur_index].Value > mm_arr[last_index].Value { //集
				opr_items = append(opr_items, &cmfbItem{
					Code:        code,
					Trade_day:   cur_kline.Tradday,
					Upsert_time: cur_kline.Upserttime,
					Type:        Cmfb_type_ji,
				})
				// } else if last_one {
				// 	last_del_items = append(last_del_items, &cmfbItem{
				// 		Code: code,
				// 		Trade_day: cur_kline.Trad_day,
				// 		Upsert_time: cur_kline.Upsert_time,
				// 		Type: Cmfb_type_fen,
				// 	})

				// 	last_del_items = append(last_del_items, &cmfbItem{
				// 		Code: code,
				// 		Trade_day: cur_kline.Trad_day,
				// 		Upsert_time: cur_kline.Upsert_time,
				// 		Type: Cmfb_type_ji,
				// 	})
			}
		}
	}

	calc_time := time.Now()
	mysql_engine := global.GServer.GetMysqlEngine()
	index := 0

	for len(del_items) > 0 {
		bg_index := index * zbbase.SQL_OPR_BATCH_CNT
		end_index := bg_index + zbbase.SQL_OPR_BATCH_CNT
		if end_index >= len(del_items) {
			end_index = len(del_items)
		}

		items := del_items[bg_index:end_index]
		//先删除
		del_sql := "delete from db_zbjs_hjjj.tb_cmfb where "
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
			log.Release("ChouMaFenBu::Calc del_sql:%s err:%v", del_sql, err)
			return err, new_calc_time
		}

		if end_index == len(del_items) {
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
		//插入
		insert_sql := "insert into db_zbjs_hjjj.tb_cmfb(stockcode, sdata, upserttime, `rsttype`) values "
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
			log.Release("ChouMaFenBu::Calc insert_sql err:%v", err)
			return err, new_calc_time
		}

		if end_index == len(opr_items) {
			break
		}

		index++
	}

	log.Release("ChouMaFenBu Calc stock len:%d total time:%v calc time:%v, del cnt:%d, suc cnt:%d %d:%d",
		len(global.StockData), time.Now().Sub(bt), calc_time.Sub(bt), len(del_items), len(opr_items), last_trad_day, today_cnt)

	return err, new_calc_time
}

// 开始筹码分布指标计算
func (m *ChouMaFenBu) Start() error {
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
