package caopantixing

import (
	"common/utils"
	"fmt"
	"libs/log"
	"time"
	"zhibiaocalcsvr/src/calc/zhibiao/zbbase"
	"zhibiaocalcsvr/src/global"
)

const (
	BO_DUAN  = 0
	FAN_TAN  = 1
	CHAO_DIE = 2
	RED      = 0
	GREEN    = 1
)

type columnStruct struct {
	Code        string
	Trad_day    uint32
	Upsert_time uint32
	High        int
	Low         int
	Color       int
}

type dotStruct struct {
	Code        string
	Trad_day    uint32
	Upsert_time uint32
	Rsttype     int
}

type lastCalcInfo struct {
	upserttime uint32
}

type CaoPanTiXing struct {
	zbbase.BaseModule
	last_calc_info_ map[string]lastCalcInfo
}

func (m *CaoPanTiXing) Init() {
	m.last_calc_info_ = make(map[string]lastCalcInfo)

	m.Filedname = "cptx"
	m.TBname = "db_zbjs_common.tb_lastcalcinfo"
}

func (m *CaoPanTiXing) Name() string {
	return zbbase.M_cptx
}

/*
Z:=MA(C,120);																						//c close 收盘价 ， ma 均线计算公式 ，z 120日均线
VAR3:=(MA(C,5)-Z)/Z;
VAR4:=MA((CLOSE-LLV(LOW,10))/(HHV(HIGH,10)-LLV(LOW,10))*100,3);
波段进场:IF(CLOSE>Z AND REF(VAR4,1)<30 AND VAR4>REF(VAR4,1) AND REF(VAR4,1)<REF(VAR4,2),80,50);		// ref 前n个对应值， eg：ref 1 昨天  ， ref 2 前天， REF(VAR4,1)  昨天var4的结果
反弹进场:IF(REF(VAR4,1)<5 AND VAR4>REF(VAR4,1) AND REF(VAR4,1)<REF(VAR4,2) AND VAR3<-0.3,80,50);
超跌进场:IF(CROSS(VAR4,5) AND VAR3<-0.4,80,50),COLORMAGENTA;
STICKLINE(C>=MA(C,10),VAR4,REF(VAR4,1),8,0),COLORRED;									// STICKLINE 画柱线  第一个参数时间， 最后一个是颜色， 中间四个参数应该是坐标
STICKLINE(C<MA(C,10),VAR4,REF(VAR4,1),8,0),COLORGREEN;


{公式解释
Z=120均线
VAR3=（5日均线-120均线）/120均线
VAR4= (收盘价-10日内最低价的最低值)/(10日内最高价的最高值-10日内最低价的最低值)*100  取3日简单移动平均
波段进场=  收盘价>Z  并且 昨天的VAR4<30 并且VAR4上涨 昨天VAR4下降
反弹进场= 昨天VAR4《5 并且 VAR4上涨 昨天VAR4下降 并且 VAR3<-0.3
超跌进场= VAR4上穿数值5 并且 VAR3<-0


收盘价大等于10日均线 画红色柱子
收盘价小于10日均线 画绿色柱子}
*/

func (m *CaoPanTiXing) Calc() (error, uint32) {
	log.Release("CaoPanTiXing Calc start")
	insert_columns := []columnStruct{}
	delete_columns := []columnStruct{}
	insert_dots := []dotStruct{}
	delete_dots := []dotStruct{}
	begin_time := time.Now()
	var new_calc_time uint32
	var err error

	for code, klines := range global.StockData {
		//需要120日均线的3日移动平均
		klen := len(klines.Cbs)
		if klen < 123 {
			continue
		}

		// 指数部分 不参与计算
		if _, ok := global.GServer.GetZhiBiaoMgr().GetIndexCodes()[code]; ok {
			continue
		}

		klinedata := make([]*global.CombineItem, klen)
		copy(klinedata, klines.Cbs)

		if klinedata[klen-1].Upserttime > new_calc_time {
			new_calc_time = klinedata[klen-1].Upserttime
		}

		//公式：Z:=MA(C,120);
		ma120, err := zbbase.MA_func(120, klinedata, func(kline *global.CombineItem) float64 {
			return kline.KLine.Close
		})
		if err != nil {
			log.Error("code:%s, %v", code, err)
			return err, new_calc_time
		}
		ma5, err := zbbase.MA_func(5, klinedata, func(kline *global.CombineItem) float64 {
			return kline.KLine.Close
		})
		if err != nil {
			log.Error("code:%s, %v", code, err)
			return err, new_calc_time
		}
		len_ma := len(ma120)
		ma5 = ma5[len(ma5)-len_ma:]
		var3s := []*zbbase.MathItem{}
		//公式：VAR3:=(MA(C,5)-Z)/Z;
		for i := 0; i < len_ma; i++ {
			var3s = append(var3s, &zbbase.MathItem{(ma5[i].Value - ma120[i].Value) / ma120[i].Value})
		}
		llv_low10 := zbbase.LLV_func(klinedata, 10, func(info *global.CombineItem) float64 {
			return info.KLine.Low
		})
		hhv_high10 := zbbase.HHV_func(klinedata, 10, func(info *global.CombineItem) float64 {
			return info.KLine.High
		})
		len_var4s_temp := len(llv_low10)
		var4s_temp := []*zbbase.MathItem{}
		for i := 0; i < len_var4s_temp; i++ {
			var4s_temp = append(var4s_temp, &zbbase.MathItem{
				(klinedata[i+9].KLine.Close - llv_low10[i].Value) / (hhv_high10[i].Value - llv_low10[i].Value) * 100})
		}
		//公式：VAR4:=MA((CLOSE-LLV(LOW,10))/(HHV(HIGH,10)-LLV(LOW,10))*100,3);
		var4s := zbbase.MA_arr(3, var4s_temp)
		len_var4s := len(var4s)

		for i := 0; i < len_ma; i++ {
			//如果计算过这天的，就跳过
			if klinedata[klen-len_ma+i].Upserttime < m.Last_calc_date {
				continue
			}
			//如果这天是最后一次计算的时间，就先删除（大于是考虑重启间隔时间大于一天）
			if klinedata[klen-len_ma+i].Upserttime >= m.Last_calc_date && m.Last_calc_date != 0 {
				delete_dots = append(delete_dots, dotStruct{
					Code:        code,
					Trad_day:    klinedata[klen-len_ma+i].Tradday,
					Upsert_time: klinedata[klen-len_ma+i].Upserttime,
				})
			}
			if klinedata[klen-len_ma+i].KLine.Close > ma120[i].Value &&
				var4s[len_var4s-len_ma+i-1].Value < 30 &&
				var4s[len_var4s-len_ma+i].Value > var4s[len_var4s-len_ma+i-1].Value &&
				var4s[len_var4s-len_ma+i-1].Value < var4s[len_var4s-len_ma+i-2].Value {
				//波段进场
				insert_dots = append(insert_dots, dotStruct{
					Code:        code,
					Trad_day:    klinedata[klen-len_ma+i].Tradday,
					Upsert_time: klinedata[klen-len_ma+i].Upserttime,
					Rsttype:     BO_DUAN,
				})
			}
			if var4s[len_var4s-len_ma+i-1].Value < 5 &&
				var4s[len_var4s-len_ma+i].Value > var4s[len_var4s-len_ma+i-1].Value &&
				var4s[len_var4s-len_ma+i-1].Value < var4s[len_var4s-len_ma+i-2].Value &&
				var3s[i].Value < -0.3 {
				//反弹进场
				insert_dots = append(insert_dots, dotStruct{
					Code:        code,
					Trad_day:    klinedata[klen-len_ma+i].Tradday,
					Upsert_time: klinedata[klen-len_ma+i].Upserttime,
					Rsttype:     FAN_TAN,
				})
			}
			if var4s[len_var4s-len_ma+i-1].Value < 5 &&
				var4s[len_var4s-len_ma+i].Value > 5 &&
				var3s[i].Value < -0.4 {
				//超跌进场
				insert_dots = append(insert_dots, dotStruct{
					Code:        code,
					Trad_day:    klinedata[klen-len_ma+i].Tradday,
					Upsert_time: klinedata[klen-len_ma+i].Upserttime,
					Rsttype:     CHAO_DIE,
				})
			}
		}
		ma10, err := zbbase.MA_func(10, klinedata, func(info *global.CombineItem) float64 {
			return info.KLine.Close
		})
		len_ma10 := len(ma10)
		var low, high int
		for i := 1; i < len_var4s; i++ {
			if var4s[i].Value > var4s[i-1].Value {
				high = int(var4s[i].Value * 1000)
				low = int(var4s[i-1].Value * 1000)
			} else {
				low = int(var4s[i].Value * 1000)
				high = int(var4s[i-1].Value * 1000)
			}
			//如果这天算过，就跳过
			if klinedata[klen-len_var4s+i].Upserttime < m.Last_calc_date {
				continue
			}
			if klinedata[klen-len_var4s+i].Upserttime >= m.Last_calc_date && m.Last_calc_date > 0 {
				delete_columns = append(delete_columns, columnStruct{
					Code:        code,
					Trad_day:    klinedata[klen-len_var4s+i].Tradday,
					Upsert_time: klinedata[klen-len_var4s+i].Upserttime,
				})
			}
			if klinedata[klen-len_var4s+i].KLine.Close >= ma10[len_ma10-len_var4s+i].Value {
				//红色柱子
				insert_columns = append(insert_columns, columnStruct{
					Code:        code,
					Trad_day:    klinedata[klen-len_var4s+i].Tradday,
					Upsert_time: klinedata[klen-len_var4s+i].Upserttime,
					High:        high,
					Low:         low,
					Color:       RED,
				})
			} else {
				//绿色柱子
				insert_columns = append(insert_columns, columnStruct{
					Code:        code,
					Trad_day:    klinedata[klen-len_var4s+i].Tradday,
					Upsert_time: klinedata[klen-len_var4s+i].Upserttime,
					High:        high,
					Low:         low,
					Color:       GREEN,
				})
			}
		}

	}
	calc_time := time.Now()
	log.Release("CaoPanTiXing calc done, insert_columns:%d, delete_columns:%d, insert_dots:%d, delete_dots:%d, calc used time:%dms",
		len(insert_columns), len(delete_columns), len(insert_dots), len(delete_dots), utils.Millisecond(calc_time.Sub(begin_time)))
	deleteRecord(delete_columns, delete_dots)
	insertRecord(insert_columns, insert_dots)
	sql_time := time.Now()
	log.Release("CaoPanTiXing sql operation done, used time:%dms", utils.Millisecond(sql_time.Sub(calc_time)))
	return err, new_calc_time
}

func deleteRecord(del_cols []columnStruct, del_dots []dotStruct) (err error) {
	mysql_engine := global.GServer.GetMysqlEngine()

	index := 0
	for len(del_cols) > 0 {
		bg_index := index * zbbase.SQL_OPR_BATCH_CNT
		end_index := (index + 1) * zbbase.SQL_OPR_BATCH_CNT
		if end_index > len(del_cols) {
			end_index = len(del_cols)
		}
		items := del_cols[bg_index:end_index]
		delete_sql := `delete from db_zbjs_dcxg.tb_cptx_column where `
		for i, item := range items {
			var del_sql string
			if i != 0 {
				del_sql += " or "
			}
			del_sql += fmt.Sprintf(`(stockcode="%s" and sdata=%d and upserttime=%d and rsttype=%d)`,
				item.Code, item.Trad_day, item.Upsert_time, RED)
			del_sql += fmt.Sprintf(` or (stockcode="%s" and sdata=%d and upserttime=%d and rsttype=%d)`,
				item.Code, item.Trad_day, item.Upsert_time, GREEN)
			delete_sql += del_sql;
		}
		delete_sql += ";"
		_, err = mysql_engine.Exec(delete_sql)
		if err != nil {
			log.Error("CaoPanTiXing::deleteRecord delete_sql err:%v", err)
			return
		}
		if end_index == len(del_cols) {
			break
		}
		index++
	}

	index = 0
	for len(del_dots) > 0 {
		bg_index := index * zbbase.SQL_OPR_BATCH_CNT
		end_index := (index + 1) * zbbase.SQL_OPR_BATCH_CNT
		if end_index > len(del_dots) {
			end_index = len(del_dots)
		}
		items := del_dots[bg_index:end_index]
		delete_sql := `delete from db_zbjs_dcxg.tb_cptx_dot where `
		for i, item := range items {
			var del_sql string
			if i != 0 {
				del_sql += " or "
			}
			del_sql += fmt.Sprintf(`(stockcode="%s" and sdata=%d and upserttime=%d and rsttype=%d)`,
				item.Code, item.Trad_day, item.Upsert_time, BO_DUAN)
			del_sql += fmt.Sprintf(` or (stockcode="%s" and sdata=%d and upserttime=%d and rsttype=%d)`,
				item.Code, item.Trad_day, item.Upsert_time, FAN_TAN)
			del_sql += fmt.Sprintf(` or (stockcode="%s" and sdata=%d and upserttime=%d and rsttype=%d)`,
				item.Code, item.Trad_day, item.Upsert_time, CHAO_DIE)
			delete_sql += del_sql
		}
		delete_sql += ";"
		_, err = mysql_engine.Exec(delete_sql)
		if err != nil {
			log.Error("CaoPanTiXing::deleteRecord delete_sql err:%v", err)
			return
		}
		if end_index == len(del_dots) {
			break
		}
		index++
	}
	return
}

func insertRecord(columns []columnStruct, dots []dotStruct) (err error) {
	mysql_engine := global.GServer.GetMysqlEngine()
	index := 0
	for len(columns) > 0 {
		begin_index := index * zbbase.SQL_OPR_BATCH_CNT
		end_index := (index + 1) * zbbase.SQL_OPR_BATCH_CNT
		if end_index > len(columns) {
			end_index = len(columns)
		}

		items := columns[begin_index:end_index]
		insert_sql := "insert into db_zbjs_dcxg.tb_cptx_column(stockcode, sdata, upserttime, rsttype, low, high) values "
		for i, item := range items {
			line_sql := ""
			if i != 0 {
				line_sql += ","
			}
			line_sql += fmt.Sprintf(`("%s",%d,%d,%d,%d,%d)`,
				item.Code, item.Trad_day, item.Upsert_time, item.Color, item.Low, item.High)
			insert_sql += line_sql
		}
		insert_sql += ";"
		_, err = mysql_engine.Exec(insert_sql)
		if err != nil {
			log.Error("CaoPanTiXing::insertRecord sql err:%v", err)
			return
		}
		if end_index == len(columns) {
			break
		}
		index++
	}
	index = 0
	for len(dots) > 0 {
		begin_index := index * zbbase.SQL_OPR_BATCH_CNT
		end_index := (index + 1) * zbbase.SQL_OPR_BATCH_CNT
		if end_index > len(dots) {
			end_index = len(dots)
		}
		items := dots[begin_index:end_index]
		insert_sql := "insert into db_zbjs_dcxg.tb_cptx_dot(stockcode, sdata, upserttime, rsttype) values "
		for i, item := range items {
			line_sql := ""
			if i != 0 {
				line_sql += ","
			}
			line_sql += fmt.Sprintf(`("%s",%d,%d,%d)`,
				item.Code, item.Trad_day, item.Upsert_time, item.Rsttype)
			insert_sql += line_sql
		}
		insert_sql += ";"
		_, err = mysql_engine.Exec(insert_sql)
		if err != nil {
			log.Error("CaoPanTiXing::insertRecord sql err:%v", err)
			return
		}
		if end_index == len(dots) {
			break
		}
		index++
	}
	return
}

// 开始指标计算
func (m *CaoPanTiXing) Start() error {
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
