package sighting

import (
	"fmt"
	"libs/log"
	"sort"
	"strconv"
	"zhibiaocalcsvr/src/calc/zhibiao/zbbase"

	"zhibiaocalcsvr/src/etc"
	"zhibiaocalcsvr/src/global"
)

const (
	STATE_INIT = 0 + iota
	STATE_SETUP
	STATE_PLUS
	STATE_CLEAR = 4

	TRADDAY_COUNT = 250 // 统计250个交易日内最大收益率
)

var Traddaylist []int32

type BSCalculator struct {
	tradday_limit    int32
	Last_profit_date uint32

	zbbase.BaseModule
}

type sqlBsItem struct {
	stockcode string
	tradday   uint32
	stype     int
	close     int //收盘价
	open      int //开盘价
	reason    string
}

type profit_s struct {
	stockcode string
	i_close   float64 // 建仓价
	c_close   float64 // 清仓价

	initdate  uint32 // 建仓时间
	cleardate uint32 // 清仓时间
}

type profit_item struct {
	stockcode string
	rate      float64 // 收益率
	initdate  uint32  // 建仓时间
	cleardate uint32  // 清仓时间
}

func (self *BSCalculator) Init() {
	log.Release("self *BSCalculator) Init() ----- ")
	self.Filedname = "sighting"
	self.TBname = "cpgj.tb_common_v2"
}

//本程序的数据源(日K)全部来自mongo中的历史K线数据,最新交易日数据在每天收盘后才会入库
//每次启动程序都会从头开始计算,但是只生成上一个交易日之后到最新交易日之间的买卖点数据
//上一个交易日在数据库表tb_common_v2中获取,这个交易日不一定是最新交易日的上一个交易日
//中间可能会有多个交易日,本次将全部计算到最新数据
func (self *BSCalculator) Calc() (error, uint32) {
	log.Release("sighting bs point calculate start ")
	new_last_calc_time := self.Last_calc_date
	new_last_calc_profit_time := self.Last_profit_date
	flags := 0
	var err error

	for stock_code, stock_item := range global.StockData {
		if global.InBlackList(stock_code) {
			continue
		}

		// 指数部分 不参与计算
		if _, ok := global.GServer.GetZhiBiaoMgr().GetIndexCodes()[stock_code]; ok {
			continue
		}

		len_cbs := len(stock_item.Cbs)
		//Cbs是有序的(按交易日期从小到大排序) 日期与索引i对应
		//获取stock_item.Cbs[i].Tradday >= 20160101 的索引 cfgday
		//每次都从头开始算,但买卖点是根据上一个交易日来生成的
		cfgday := sort.Search(len_cbs, func(i int) bool { return stock_item.Cbs[i].Tradday >= uint32(etc.Config.Calc.Sighting_StartDate) })
		if cfgday >= len_cbs-1 {
			continue
		}
		if self.Last_calc_date == stock_item.Cbs[cfgday].Tradday {
			cfgday++
		}
		//-------------------------------------------------------------------------------------------------------------//
		//截取满足日期条件的K线数据
		iCbs := stock_item.Cbs[cfgday:]
		//新上市个股从上市第60个交易日为起点开始计算，不到60个交易日按普通模式显示K线
		//62是方便计算: 前一个交易日的60日均线、前两个交易日的60日均线
		len_icbs := len(iCbs)
		if len_icbs < 62 {
			continue
		}

		//当某只股票有65根K线时 len(ma5) == 65-5+1 == 61
		ma5, _ := MA_func(5, iCbs, func(info *global.CombineItem) float64 {
			return info.FQKLine.Close
		})
		//当某只股票有65根K线时 len(ma10) == 65-10+1 == 56
		ma10, _ := MA_func(10, iCbs, func(info *global.CombineItem) float64 {
			return info.FQKLine.Close
		})
		//当某只股票有65根K线时 len(ma20) == 65-20+1 == 46
		ma20, _ := MA_func(20, iCbs, func(info *global.CombineItem) float64 {
			return info.FQKLine.Close
		})
		//当某只股票有65根K线时 len(ma30) == 65-30+1 == 36
		ma30, _ := MA_func(30, iCbs, func(info *global.CombineItem) float64 {
			return info.FQKLine.Close
		})
		//当某只股票有65根K线时 len(ma60) == 65-60+1== 6
		ma60, _ := MA_func(60, iCbs, func(info *global.CombineItem) float64 {
			return info.FQKLine.Close
		})
		//-----------------------------------------------看多条件计算----------------------------------------------------//
		rc1 := []bool{} //看多条件1(强烈看多)
		rc2 := []bool{} //看多条件2(看多)
		//65根K线时 len(ma60) == 6
		len_ma60 := len(ma60)

		//看多条建1(强烈看多)：n-61
		// 收盘价>60日均线 And 60日均线>前1交易日60日均线 And 前1交易日60日均线>前2交易日60日均线 And
		// 5日均线>60日均线 And 20日均线>60日均线 And 收盘价>30日均线
		for i := 2; i < len_ma60; i++ {
			tmp := iCbs[i+59].FQKLine.Close > ma60[i].Value &&
				ma60[i].Value > ma60[i-1].Value &&
				ma60[i-1].Value > ma60[i-2].Value &&
				ma5[i+55].Value > ma60[i].Value &&
				ma20[i+40].Value > ma60[i].Value &&
				iCbs[i+59].FQKLine.Close > ma30[i+30].Value
			rc1 = append(rc1, tmp)
		}
		//65根K线时 len(ma20) == 46
		len_ma20 := len(ma20)
		//看多条件2(看多)： n-20
		//5日均线>20日均线 And 20日均线>前1交易日20日均线
		for i := 1; i < len_ma20; i++ {
			tmp := ma5[i+15].Value > ma20[i].Value && ma20[i].Value > ma20[i-1].Value
			rc2 = append(rc2, tmp)
		}
		//log.Debug("kjdo条件1 2------------------------------------------------------")
		//-----------------------------------------------建仓条件计算start-----------------------------------------------//
		buy1 := []bool{}    //建仓条件1(买强)
		buy2 := []bool{}    //建仓条件2(买)
		buy3 := []bool{}    //建仓条件3(翻多)
		len_rc1 := len(rc1) //65根K线时 len(rc1) == 4
		len_rc2 := len(rc2) //65根K线时 len(rc2) == 45
		//log.Debug("len_rc1 = %v, len_rc2 = %v", len_rc1, len_rc2)
		//log.Debug("ma5 = %v, ma10 = %v,ma20 = %v, ma30 = %v,ma60 = %v", len(ma5), len(ma10),len(ma20),len(ma30),len(ma60))

		//建仓条件1(买强)：n-62
		//收盘价>5日均线 And 收盘价>10日均线 And 收盘价>20日均线 And 看多条件1成立 And 前1交易日看多条件1不成立
		for i := 1; i < len_rc1; i++ {
			tmp := iCbs[i+61].FQKLine.Close > ma5[i+57].Value &&
				iCbs[i+61].FQKLine.Close > ma10[i+52].Value &&
				iCbs[i+61].FQKLine.Close > ma20[i+42].Value &&
				rc1[i] && !rc1[i-1]
			//65根K线时 len(buy1) == 4-1 == 3
			//需要特别注意buy1 与 rc1的对应关系: buy1[0] <==> rc1[1]
			buy1 = append(buy1, tmp)
		}

		//建仓条件2(买)：n-59
		//收盘价>5日均线 And 收盘价>10日均线 And 收盘价>20日均线 And 看多条件2成立 And 前1交易日看多条件2不成立 And 个股涨跌幅>3%And 收盘价>60日均线
		for i := 39; i < len_rc2; i++ {
			tmp := iCbs[i+20].FQKLine.Close > ma5[i+16].Value &&
				iCbs[i+20].FQKLine.Close > ma10[i+11].Value &&
				iCbs[i+20].FQKLine.Close > ma20[i+1].Value &&
				rc2[i] && !rc2[i-1] &&
				(iCbs[i+20].FQKLine.Close-iCbs[i+20].FQKLine.PreClose)/iCbs[i+20].FQKLine.PreClose > 0.03
			//65根K线时 len(buy2) == 46-1 这行是错的，提交代码的时候会删除
			//65根K线时 len(buy2) == 6-1  5根
			buy2 = append(buy2, tmp)
		}

		//建仓条件3)(翻多)：n-59
		//收盘价>5日均线 and 收盘价>10日均线and 收盘价>20日均线and 收盘价>30日均线 and 收盘价>60日均线and 看多条件=True and 收盘价= 最近10日个交易日最高收盘价 And 收盘价>60日均线
		for i := 39; i < len_rc2; i++ {
			isMaxClose := true
			// 最近10个交易日 收盘价最大值
			max_close := iCbs[i+20].FQKLine.Close
			for index := i + 20; index > i+10; index-- {
				if iCbs[index].FQKLine.Close > max_close {
					//max_close = iCbs[index].FQKLine.Close
					isMaxClose = false
				}
			}

			// 多看条件判断
			rcResult := false
			if i+20 < 61 {
				if rc2[i] {
					rcResult = true
				}
			} else {
				if rc2[i] || rc1[i+20-61] {
					rcResult = true
				}
			}
			tmp := iCbs[i+20].FQKLine.Close > ma5[i+20-4].Value && //收盘价 > 5日均线
				iCbs[i+20].FQKLine.Close > ma10[i+20-9].Value && //收盘价 > 10日均线
				iCbs[i+20].FQKLine.Close > ma20[i+20-19].Value && //收盘价 > 20日均线
				iCbs[i+20].FQKLine.Close > ma30[i+20-29].Value && //收盘价 > 30日均线
				iCbs[i+20].FQKLine.Close > ma60[i+20-59].Value && //收盘价 > 60日均线
				isMaxClose && //收盘价 = 最近10日个交易日最高收盘价
				rcResult             //看多条件=True
			buy3 = append(buy3, tmp) //翻多
		}

		//------------------------------------------------建仓条件计算end------------------------------------------------//
		min_len := len(buy1)            //
		buy2 = buy2[len(buy2)-min_len:] //只取最后 min_len 个建仓条件
		buy3 = buy3[len(buy3)-min_len:]
		len_ma5 := len(ma5)
		index_ma5 := len_ma5 - min_len //只取最后 min_len 个5日均价
		len_ma10 := len(ma10)
		index_ma10 := len_ma10 - min_len //只取最后 min_len 个10日均价
		index_ma20 := len_ma20 - min_len //只取最后 min_len 个20日均价
		len_ma30 := len(ma30)
		index_ma30 := len_ma30 - min_len //只取最后 min_len 个30日均价
		index_ma60 := len_ma60 - min_len //只取最后 min_len 个60日均价
		cbs := iCbs[len_icbs-min_len:]   //截取有效K线数据的最后 min_len 个元素(最近交易日的K线数据)
		prev_state := 0                  //记录上一个买卖点状态
		if len_icbs-min_len-1 >= 0 {
			prev_state = iCbs[len_icbs-min_len-1].C_BsPoint.State
		}
		stype := 0

		sql_items := []sqlBsItem{} //只临时存放某只股票的计算结果
		sql_profit_s := []profit_item{}
		ps_max := profit_item{
			rate: -999,
		}
		var ps profit_s

		for k := 0; k < min_len; k++ { //相当于遍历 buy1，因为buy1长度最小
			if new_last_calc_time < cbs[k].Tradday {
				new_last_calc_time = cbs[k].Tradday
			}
			if new_last_calc_profit_time < cbs[k].Tradday {
				new_last_calc_profit_time = cbs[k].Tradday
			}
			//----------------------------------------------清仓条件1--------------------------------------------------//
			//连续2天收盘价 < 20日线
			close2 := iCbs[len_icbs-min_len-1+k].FQKLine.Close < ma20[index_ma20+k-1].Value &&
				cbs[k].FQKLine.Close < ma20[index_ma20+k].Value
			//连续4天收盘价 < 20日线
			close4 := cbs[k].FQKLine.Close < ma20[index_ma20+k].Value &&
				iCbs[len_icbs-min_len-1+k].FQKLine.Close < ma20[index_ma20+k-1].Value &&
				iCbs[len_icbs-min_len-2+k].FQKLine.Close < ma20[index_ma20+k-2].Value &&
				iCbs[len_icbs-min_len-3+k].FQKLine.Close < ma20[index_ma20+k-3].Value

			//----------------------------------------------清仓条件2--------------------------------------------------//
			//放量大跌1
			srate1 := (cbs[k].FQKLine.Close-cbs[k].FQKLine.PreClose)/cbs[k].FQKLine.PreClose < -0.05
			vrate1 := float64(cbs[k].FQKLine.Volume)/float64(iCbs[len_icbs-min_len-1+k].FQKLine.Volume) > 1.2 ||
				float64(cbs[k].FQKLine.Volume)/float64(iCbs[len_icbs-min_len-2+k].FQKLine.Volume) > 1.2
			//放量大跌2
			srate2 := (cbs[k].FQKLine.Close-cbs[k].FQKLine.PreClose)/cbs[k].FQKLine.PreClose < -0.07
			vrate2 := float64(cbs[k].FQKLine.Volume)/float64(iCbs[len_icbs-min_len-1+k].FQKLine.Volume) > 1 ||
				float64(cbs[k].FQKLine.Volume)/float64(iCbs[len_icbs-min_len-2+k].FQKLine.Volume) > 1

			//----------------------------------------------清仓条件3--------------------------------------------------//
			//跌破所有均线
			fall_average := cbs[k].FQKLine.Close < ma5[index_ma5+k].Value &&
				cbs[k].FQKLine.Close < ma10[index_ma10+k].Value &&
				cbs[k].FQKLine.Close < ma20[index_ma20+k].Value &&
				cbs[k].FQKLine.Close < ma30[index_ma30+k].Value &&
				cbs[k].FQKLine.Close < ma60[index_ma60+k].Value

			//----------------------------------------------清仓条件4--------------------------------------------------//
			//时间止损  获取近11个交易日的最高收盘价
			max_close11 := iCbs[len_icbs-min_len+k].FQKLine.Close
			for j := 1; j < 11; j++ {
				if iCbs[len_icbs-min_len+k-j].FQKLine.Close > max_close11 {
					max_close11 = iCbs[len_icbs-min_len+k-j].FQKLine.Close
				}
			}

			//----------------------------------------------清仓条件5--------------------------------------------------//
			// 建仓部位止损，收盘价/建仓点收盘价<=0.94
			if stock_item.BuyClose <= 0 {
				tradday, lastClose, err := GetLastCloseByStock(stock_code, 1) //获取最后一次建仓点收盘价
				if err != nil && lastClose <= 0 {
					log.Error("[GetLastCloseByStock error] stock_code:%s, err:%s", stock_code, err.Error())
					continue
				}
				stock_item.BuyClose = float64(lastClose) / 1000
				stock_item.BuyDay = uint32(tradday)
			}
			close5 := cbs[k].FQKLine.Close/stock_item.BuyClose <= 0.94

			//---------------------------------------------买卖点开始生成----------------------------------------------//

			switch prev_state {
			case STATE_CLEAR: //初始状态和清仓状态一个操作
				fallthrough
			case STATE_INIT: //建仓
				create := stype == STATE_CLEAR || stype == STATE_INIT
				if buy1[k] || buy2[k] || buy3[k] {
					reason := ""
					if buy1[k] {
						reason = "[建仓1]-->买强"
					} else if buy2[k] {
						reason = "[建仓2]-->买"
					} else if buy3[k] {
						reason = "[建仓3]-->翻多"
					}
					stype = STATE_SETUP
					prev_state = stype
					stock_item.BuyClose = cbs[k].FQKLine.Close //建仓点的收盘价
					stock_item.BuyDay = cbs[k].Tradday
					//只计算最后一个交易日之后的数据
					//若表中最后一个交易日为空,则全部计算重新入库,这种情况下,需要先清空两张表,之后再重启程序
					//必须放在生成买卖点之前是因为要保存当前的买卖点状态
					if cbs[k].Tradday > self.Last_calc_date && create {
						log.Release("[建仓 success] tradday:%v stock_code:%s, reason: %s", cbs[k].Tradday, stock_code, reason)
						sql_items = append(sql_items, sqlBsItem{
							stockcode: stock_code,
							tradday:   cbs[k].Tradday,
							stype:     stype,
							close:     int(cbs[k].KLine.Close * 1000),
							open:      int(cbs[k].KLine.Open * 1000),
							reason:    reason,
						})
					}

					// 记录收益率
					if create {
						ps = profit_s{
							stockcode: stock_code,
							i_close:   cbs[k].FQKLine.Close * 1000,
							initdate:  cbs[k].Tradday,
						}
					}
				}
				cbs[k].C_BsPoint.State = prev_state
			case STATE_SETUP: //建仓部位开始: 先判定清仓, 不满足所有清仓条件时再判定加仓
				if stock_item.BuyClose <= 0 {
					tradday, lastClose, err := GetLastCloseByStock(stock_code, 1) //获取最后一次建仓点收盘价
					if err != nil && lastClose <= 0 {
						log.Error("[GetLastCloseByStock error] stock_code:%s, err:%s", stock_code, err.Error())
						continue
					}
					stock_item.BuyClose = float64(lastClose) / 1000
					stock_item.BuyDay = uint32(tradday)
				}
				clear := stype != STATE_CLEAR && stype != STATE_INIT
				plus := stype == STATE_SETUP
				//需要特别注意buy1 与 rc1的对应关系: buy1[k] <==> rc1[k+1]
				if (!rc1[k+1] && close2) || (rc1[k+1] && close4) {
					reason := "[清仓1]-->有效跌破20日均线"
					stype = STATE_CLEAR
					prev_state = stype
					if cbs[k].Tradday > self.Last_calc_date && clear {
						log.Release("[建仓后清仓1 success] stock_code:%s, reason: %s", stock_code, reason)
						sql_items = append(sql_items, sqlBsItem{
							stockcode: stock_code,
							tradday:   cbs[k].Tradday,
							stype:     stype,
							close:     int(cbs[k].KLine.Close * 1000),
							open:      int(cbs[k].KLine.Open * 1000),
							reason:    reason,
						})
					}
					if clear {
						ps.c_close = cbs[k].FQKLine.Close * 1000
						ps.cleardate = cbs[k].Tradday
						ratestr := fmt.Sprintf("%.4f", ps.c_close/ps.i_close-1)
						rate, _ := strconv.ParseFloat(ratestr, 64)

						sql_profit_s = append(sql_profit_s, profit_item{
							rate:      rate,
							initdate:  ps.initdate,
							cleardate: ps.cleardate,
							stockcode: stock_code,
						})

						if rate > ps_max.rate && int32(ps.initdate) >= self.tradday_limit {
							ps_max.stockcode = stock_code
							ps_max.rate = rate
							ps_max.initdate = ps.initdate
							ps_max.cleardate = ps.cleardate
						}

					}
				} else if ((srate1 && vrate1) || (srate2 && vrate2)) && cbs[k].FQKLine.Close < ma10[index_ma10+k].Value {
					reason := "[清仓2]-->大跌,破10日均线"
					stype = STATE_CLEAR
					prev_state = stype
					if cbs[k].Tradday > self.Last_calc_date && clear {
						log.Release("[建仓后清仓2 success] stock_code:%s, reason: %s", stock_code, reason)
						sql_items = append(sql_items, sqlBsItem{
							stockcode: stock_code,
							tradday:   cbs[k].Tradday,
							stype:     stype,
							close:     int(cbs[k].KLine.Close * 1000),
							open:      int(cbs[k].KLine.Open * 1000),
							reason:    reason,
						})
					}
					if clear {
						ps.c_close = cbs[k].FQKLine.Close * 1000
						ps.cleardate = cbs[k].Tradday

						ratestr := fmt.Sprintf("%.4f", ps.c_close/ps.i_close-1)
						rate, _ := strconv.ParseFloat(ratestr, 64)
						sql_profit_s = append(sql_profit_s, profit_item{
							rate:      rate,
							initdate:  ps.initdate,
							cleardate: ps.cleardate,
							stockcode: stock_code,
						})

						if rate > ps_max.rate && int32(ps.initdate) >= self.tradday_limit {
							ps_max.stockcode = stock_code
							ps_max.rate = rate
							ps_max.initdate = ps.initdate
							ps_max.cleardate = ps.cleardate
						}

					}
				} else if fall_average {
					reason := "[清仓3]-->跌破所有均线"
					stype = STATE_CLEAR
					prev_state = stype
					if cbs[k].Tradday > self.Last_calc_date && clear {
						log.Release("[建仓后清仓3 success] stock_code:%s, reason: %s", stock_code, reason)
						sql_items = append(sql_items, sqlBsItem{
							stockcode: stock_code,
							tradday:   cbs[k].Tradday,
							stype:     stype,
							close:     int(cbs[k].KLine.Close * 1000),
							open:      int(cbs[k].KLine.Open * 1000),
							reason:    reason,
						})
					}
					if clear {
						ps.c_close = cbs[k].FQKLine.Close * 1000
						ps.cleardate = cbs[k].Tradday

						ratestr := fmt.Sprintf("%.4f", ps.c_close/ps.i_close-1)
						rate, _ := strconv.ParseFloat(ratestr, 64)
						sql_profit_s = append(sql_profit_s, profit_item{
							rate:      rate,
							initdate:  ps.initdate,
							cleardate: ps.cleardate,
							stockcode: stock_code,
						})

						if rate > ps_max.rate && int32(ps.initdate) >= self.tradday_limit {
							ps_max.stockcode = stock_code
							ps_max.rate = rate
							ps_max.initdate = ps.initdate
							ps_max.cleardate = ps.cleardate
						}

					}
				} else if IsMoreThanTenDays(stock_code, stock_item.BuyDay, iCbs, k, len(iCbs)-min_len) && stock_item.BuyClose == max_close11 {
					reason := "[清仓4]-->时间止损"
					stype = STATE_CLEAR
					prev_state = stype
					if cbs[k].Tradday > self.Last_calc_date && clear {
						log.Release("[建仓后清仓4 success] stock_code:%s, reason: %s", stock_code, reason)
						sql_items = append(sql_items, sqlBsItem{
							stockcode: stock_code,
							tradday:   cbs[k].Tradday,
							stype:     stype,
							close:     int(cbs[k].KLine.Close * 1000),
							open:      int(cbs[k].KLine.Open * 1000),
							reason:    reason,
						})
					}
					if clear {
						ps.c_close = cbs[k].FQKLine.Close * 1000
						ps.cleardate = cbs[k].Tradday

						ratestr := fmt.Sprintf("%.4f", ps.c_close/ps.i_close-1)
						rate, _ := strconv.ParseFloat(ratestr, 64)
						sql_profit_s = append(sql_profit_s, profit_item{
							rate:      rate,
							initdate:  ps.initdate,
							cleardate: ps.cleardate,
							stockcode: stock_code,
						})

						if rate > ps_max.rate && int32(ps.initdate) >= self.tradday_limit {
							ps_max.stockcode = stock_code
							ps_max.rate = rate
							ps_max.initdate = ps.initdate
							ps_max.cleardate = ps.cleardate
						}

					}
				} else if close5 {
					reason := "[清仓5]-->建仓部位止损"
					stype = STATE_CLEAR
					prev_state = stype
					if cbs[k].Tradday > self.Last_calc_date && clear {
						log.Release("[建仓后清仓5 success] stock_code:%s, reason: %s", stock_code, reason)
						sql_items = append(sql_items, sqlBsItem{
							stockcode: stock_code,
							tradday:   cbs[k].Tradday,
							stype:     stype,
							close:     int(cbs[k].KLine.Close * 1000),
							open:      int(cbs[k].KLine.Open * 1000),
							reason:    reason,
						})
					}
					if clear {
						ps.c_close = cbs[k].FQKLine.Close * 1000
						ps.cleardate = cbs[k].Tradday

						ratestr := fmt.Sprintf("%.4f", ps.c_close/ps.i_close-1)
						rate, _ := strconv.ParseFloat(ratestr, 64)
						sql_profit_s = append(sql_profit_s, profit_item{
							rate:      rate,
							initdate:  ps.initdate,
							cleardate: ps.cleardate,
							stockcode: stock_code,
						})

						if rate > ps_max.rate && int32(ps.initdate) >= self.tradday_limit {
							ps_max.stockcode = stock_code
							ps_max.rate = rate
							ps_max.initdate = ps.initdate
							ps_max.cleardate = ps.cleardate
						}

					}
				} else if cbs[k].FQKLine.Close/stock_item.BuyClose >= 1.09 && cbs[k].FQKLine.Close > ma5[index_ma5+k].Value {
					reason := "[加仓0]-->浮盈加仓"
					stype = STATE_PLUS
					prev_state = stype
					stock_item.AddClose = cbs[k].FQKLine.Close
					stock_item.AddDay = cbs[k].Tradday
					if cbs[k].Tradday > self.Last_calc_date && plus {
						log.Release("[加仓 success] stock_code:%s, reason: %s", stock_code, reason)
						sql_items = append(sql_items, sqlBsItem{
							stockcode: stock_code,
							tradday:   cbs[k].Tradday,
							stype:     stype,
							close:     int(cbs[k].KLine.Close * 1000),
							open:      int(cbs[k].KLine.Open * 1000),
							reason:    reason,
						})
					}
				}
				cbs[k].C_BsPoint.State = prev_state
			case STATE_PLUS: //加仓部位开始:  判定清仓
				if stock_item.BuyClose <= 0 {
					tradday, lastClose, err := GetLastCloseByStock(stock_code, 1) //获取最后一次建仓点收盘价
					if err != nil && lastClose <= 0 {
						log.Error("[GetLastCloseByStock error] stock_code:%s, err:%s", stock_code, err.Error())
						continue
					}
					stock_item.BuyClose = float64(lastClose) / 1000
					stock_item.BuyDay = uint32(tradday)
				}
				if stock_item.AddClose <= 0 {
					plus_tradday, plus_close, err := GetLastCloseByStock(stock_code, 2) //获取最后一次加仓点收盘价
					if err != nil && plus_close <= 0 {
						log.Error("[GetLastCloseByStock STATE_PLUS] stock_code:%s, err:%s", stock_code, err.Error())
						continue
					}
					stock_item.AddClose = float64(plus_close) / 1000
					stock_item.AddDay = uint32(plus_tradday)
				}
				clear := stype != STATE_CLEAR && stype != STATE_INIT
				if (!rc1[k+1] && close2) || (rc1[k+1] && close4) {
					reason := "[清仓1]-->有效跌破20日均线"
					stype = STATE_CLEAR
					prev_state = stype
					if cbs[k].Tradday > self.Last_calc_date && clear {
						log.Release("[加仓后清仓1 success] stock_code:%s, reason: %s", stock_code, reason)
						sql_items = append(sql_items, sqlBsItem{
							stockcode: stock_code,
							tradday:   cbs[k].Tradday,
							stype:     stype,
							close:     int(cbs[k].KLine.Close * 1000),
							open:      int(cbs[k].KLine.Open * 1000),
							reason:    reason,
						})
					}
					if clear {
						ps.c_close = cbs[k].FQKLine.Close * 1000
						ps.cleardate = cbs[k].Tradday

						ratestr := fmt.Sprintf("%.4f", ps.c_close/ps.i_close-1)
						rate, _ := strconv.ParseFloat(ratestr, 64)
						sql_profit_s = append(sql_profit_s, profit_item{
							rate:      rate,
							initdate:  ps.initdate,
							cleardate: ps.cleardate,
							stockcode: stock_code,
						})

						if rate > ps_max.rate && int32(ps.initdate) >= self.tradday_limit {
							ps_max.stockcode = stock_code
							ps_max.rate = rate
							ps_max.initdate = ps.initdate
							ps_max.cleardate = ps.cleardate
						}

					}
				} else if ((srate1 && vrate1) || (srate2 && vrate2)) && cbs[k].FQKLine.Close < ma10[index_ma10+k].Value {
					reason := "[清仓2]-->大跌,破10日均线"
					stype = STATE_CLEAR
					prev_state = stype
					if cbs[k].Tradday > self.Last_calc_date && clear {
						log.Release("[加仓后清仓2 success] stock_code:%s, reason: %s", stock_code, reason)
						sql_items = append(sql_items, sqlBsItem{
							stockcode: stock_code,
							tradday:   cbs[k].Tradday,
							stype:     stype,
							close:     int(cbs[k].KLine.Close * 1000),
							open:      int(cbs[k].KLine.Open * 1000),
							reason:    reason,
						})
					}

					if clear {
						ps.c_close = cbs[k].FQKLine.Close * 1000
						ps.cleardate = cbs[k].Tradday

						ratestr := fmt.Sprintf("%.4f", ps.c_close/ps.i_close-1)
						rate, _ := strconv.ParseFloat(ratestr, 64)
						sql_profit_s = append(sql_profit_s, profit_item{
							rate:      rate,
							initdate:  ps.initdate,
							cleardate: ps.cleardate,
							stockcode: stock_code,
						})

						if rate > ps_max.rate && int32(ps.initdate) >= self.tradday_limit {
							ps_max.stockcode = stock_code
							ps_max.rate = rate
							ps_max.initdate = ps.initdate
							ps_max.cleardate = ps.cleardate
						}

					}
				} else if fall_average {
					reason := "[清仓3]-->跌破所有均线"
					stype = STATE_CLEAR
					prev_state = stype
					if cbs[k].Tradday > self.Last_calc_date && clear {
						log.Release("[加仓后清仓3 success] stock_code:%s, reason: %s", stock_code, reason)
						sql_items = append(sql_items, sqlBsItem{
							stockcode: stock_code,
							tradday:   cbs[k].Tradday,
							stype:     stype,
							close:     int(cbs[k].KLine.Close * 1000),
							open:      int(cbs[k].KLine.Open * 1000),
							reason:    reason,
						})
					}
					if clear {
						ps.c_close = cbs[k].FQKLine.Close * 1000
						ps.cleardate = cbs[k].Tradday

						ratestr := fmt.Sprintf("%.4f", ps.c_close/ps.i_close-1)
						rate, _ := strconv.ParseFloat(ratestr, 64)
						sql_profit_s = append(sql_profit_s, profit_item{
							rate:      rate,
							initdate:  ps.initdate,
							cleardate: ps.cleardate,
							stockcode: stock_code,
						})

						if rate > ps_max.rate && int32(ps.initdate) >= self.tradday_limit {
							ps_max.stockcode = stock_code
							ps_max.rate = rate
							ps_max.initdate = ps.initdate
							ps_max.cleardate = ps.cleardate
						}

					}
				} else if IsMoreThanTenDays(stock_code, stock_item.BuyDay, iCbs, k, len(iCbs)-min_len) && stock_item.BuyClose == max_close11 {
					reason := "[清仓4]-->时间止损"
					stype = STATE_CLEAR
					prev_state = stype
					if cbs[k].Tradday > self.Last_calc_date && clear {
						log.Release("[加仓后清仓4 success] stock_code:%s, reason: %s", stock_code, reason)
						sql_items = append(sql_items, sqlBsItem{
							stockcode: stock_code,
							tradday:   cbs[k].Tradday,
							stype:     stype,
							close:     int(cbs[k].KLine.Close * 1000),
							open:      int(cbs[k].KLine.Open * 1000),
							reason:    reason,
						})
					}

					if clear {
						ps.c_close = cbs[k].FQKLine.Close * 1000
						ps.cleardate = cbs[k].Tradday

						ratestr := fmt.Sprintf("%.4f", ps.c_close/ps.i_close-1)
						rate, _ := strconv.ParseFloat(ratestr, 64)
						sql_profit_s = append(sql_profit_s, profit_item{
							rate:      rate,
							initdate:  ps.initdate,
							cleardate: ps.cleardate,
							stockcode: stock_code,
						})

						if rate > ps_max.rate && int32(ps.initdate) >= self.tradday_limit {
							ps_max.stockcode = stock_code
							ps_max.rate = rate
							ps_max.initdate = ps.initdate
							ps_max.cleardate = ps.cleardate
						}

					}

				} else if close5 {
					reason := "[清仓5]-->建仓部位止损"
					stype = STATE_CLEAR
					prev_state = stype
					if cbs[k].Tradday > self.Last_calc_date && clear {
						log.Release("[加仓后清仓5 success] stock_code:%s, reason: %s", stock_code, reason)
						sql_items = append(sql_items, sqlBsItem{
							stockcode: stock_code,
							tradday:   cbs[k].Tradday,
							stype:     stype,
							close:     int(cbs[k].KLine.Close * 1000),
							open:      int(cbs[k].KLine.Open * 1000),
							reason:    reason,
						})
					}
					if clear {
						ps.c_close = cbs[k].FQKLine.Close * 1000
						ps.cleardate = cbs[k].Tradday

						ratestr := fmt.Sprintf("%.4f", ps.c_close/ps.i_close-1)
						rate, _ := strconv.ParseFloat(ratestr, 64)
						sql_profit_s = append(sql_profit_s, profit_item{
							rate:      rate,
							initdate:  ps.initdate,
							cleardate: ps.cleardate,
							stockcode: stock_code,
						})

						if rate > ps_max.rate && int32(ps.initdate) >= self.tradday_limit {
							ps_max.stockcode = stock_code
							ps_max.rate = rate
							ps_max.initdate = ps.initdate
							ps_max.cleardate = ps.cleardate
						}
					}
				}
				cbs[k].C_BsPoint.State = prev_state
			}
		}
		//操作数据库:某只股票计算完成后入库
		//从2018年开始数据量太大,所以无法将所有股票的计算结果统一入库,
		//只能将每只股票的计算结果单独入库
		if len(sql_items) > 0 {
			flags++
			err = InsertBS(sql_items)
			if err != nil {
				log.Error("[insertBs sql error] err: %v", err)
				return err, self.Last_calc_date
			}
		}

		if len(sql_profit_s) > 0 {
			err = InsertPr(sql_profit_s, ps_max, self.Last_profit_date)
			if err != nil {
				log.Error("[InsertPr sql error] err: %v", err)
				return err, self.Last_profit_date
			}
		}
	}

	if self.Last_profit_date != new_last_calc_profit_time {
		self.Last_profit_date = new_last_calc_profit_time
	}

	if flags > 0 {
		log.Release("sighting bs point sql operation done, A total of %d stocks produced BS point", flags)
		return err, new_last_calc_time
	}
	log.Release("No bs point have been produced since the previous trading day, Last trading date not updated: %d", self.Last_calc_date)
	return err, self.Last_calc_date
}

func (self *BSCalculator) Start() error {
	//回测结果入库：执行一次即可, 执行前需清空表sighting_back_test_v2
	//OpMain()

	// 拉取交易日
	traddaylist, err := GetTraddayList(false)
	if err != nil {
		return err
	}

	// 获取第250个交易日的时间
	Traddaylist = traddaylist
	if len(traddaylist) < TRADDAY_COUNT {
		self.tradday_limit = traddaylist[len(traddaylist)-1]
	} else {
		self.tradday_limit = traddaylist[TRADDAY_COUNT-1]
	}
	log.Release("self.tradday_limit = %v", self.tradday_limit)

	// 计算最后一次的计算时间。
	self.GetLastCalcDate(self.TBname, self.Filedname, self.Name())

	// 开始计算
	err, last_calc_date := self.Calc()
	if err != nil {
		return err
	}
	// 更新入库时间
	err = self.UpdateCalcDateToMySql(self.TBname, self.Filedname, last_calc_date)

	return err
}

// 将计算的时间入库
func (self *BSCalculator) UpdateCalcDateToMySql(tablename, filedname string, last_calc_time uint32) error {
	insert_sql := fmt.Sprintf(`INSERT INTO %s(stype,start_date,last_calc_date) VALUES ('%s',%d,%d) ON DUPLICATE KEY UPDATE start_date=%d,last_calc_date='%d';`,
		tablename, filedname,
		etc.Config.Calc.Sighting_StartDate, last_calc_time,
		etc.Config.Calc.Sighting_StartDate, last_calc_time)
	log.Release("update last calc date:%d, insert_sql = %v", last_calc_time, insert_sql)

	_, err := global.GServer.GetMysqlEngine().Exec(insert_sql)
	if err != nil {
		log.Error("update last_calc_date fail:%v", err)
		return err
	}

	insert_sql_pr := fmt.Sprintf(`INSERT INTO %s(stype,start_date,last_calc_date) VALUES ('%s',%d,%d) ON DUPLICATE KEY UPDATE start_date=%d,last_calc_date='%d';`,
		tablename, "sighting_profit",
		etc.Config.Calc.Sighting_StartDate, self.Last_profit_date,
		etc.Config.Calc.Sighting_StartDate, self.Last_profit_date)
	log.Release("update last calc date:%d, insert_sql_pr = %v", self.Last_profit_date, insert_sql_pr)
	_, err = global.GServer.GetMysqlEngine().Exec(insert_sql_pr)
	if err != nil {
		log.Error("update last_profit_date fail:%v", err)
	}

	return err
}

// 拉取最后计算时间
func (self *BSCalculator) GetLastCalcDate(tbname, filedname, calcname string) {
	if tbname == "" || filedname == "" {
		log.Error("GetLastCalcDate() param is empty, tbname:%v, filedname:%v, calcname:%v", tbname, filedname, calcname)
		self.Last_calc_date = uint32(etc.Config.Calc.Sighting_StartDate)
		return
	}
	defer func() {
		log.Release("GetLastCalcDate() --sighting--Last_calc_date = %v, sighting_profit time = %v", self.Last_calc_date, self.Last_profit_date)
	}()

	filednames := []string{"sighting", "sighting_profit"}
	for k, v := range filednames {
		sql := fmt.Sprintf("SELECT start_date,last_calc_date from %s WHERE stype='%s';", tbname, v)
		log.Release("GetLastCalcDate() --%v-- sql = %v", self.Name(), sql)
		rows, err := global.GServer.GetMysqlEngine().Query(sql)
		defer rows.Close()
		if err != nil {
			log.Error("get last_calc_date fail:%s", err.Error())
			self.Last_calc_date = uint32(etc.Config.Calc.Sighting_StartDate)
			continue
		}
		var start_date uint32
		var last_calc_date uint32
		if rows.Next() {
			rows.Scan(&start_date, &last_calc_date)
		}
		if last_calc_date != 0 {
			if k == 0 {
				self.Last_calc_date = last_calc_date
			} else {
				self.Last_profit_date = last_calc_date
			}
			continue
		}

		if start_date != 0 {
			if k == 0 {
				self.Last_calc_date = last_calc_date
			} else {
				self.Last_profit_date = last_calc_date
			}
			continue
		}
		if etc.Config.Calc.Sighting_StartDate != 0 {
			if k == 0 {
				self.Last_calc_date = last_calc_date
			} else {
				self.Last_profit_date = last_calc_date
			}
			continue
		}

		if k == 0 {
			self.Last_calc_date = last_calc_date
		} else {
			self.Last_profit_date = last_calc_date
		}
	}
}

// 服务时间
func (self *BSCalculator) Name() string {
	return "sighting"
}
