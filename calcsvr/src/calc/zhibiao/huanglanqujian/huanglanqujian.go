package huanglanqujian

import (
	"common/utils"
	"fmt"
	"libs/log"
	"time"
	"zhibiaocalcsvr/src/calc/zhibiao/zbbase"
	"zhibiaocalcsvr/src/global"
)

/*
******************************需求的算法描述（3.7.0）****************************,部分算法已经更新至430
看多条件1：
收盘价>60日均线 And 60日均线>前1交易日60日均线 And 前1交易日60日均线>前2交易日60日均线 And 5日均线>60日均线 And 20日均线>60日均线 And 收盘价>30日均线
看多条件2：
5日均线>20日均线 And 20日均线>前1交易日20日均线

看多条件：
看多条件1 Or 看多条件2

当看多条件成立时，显示黄柱，否则显示蓝柱。（tzn注：这里并没有给出柱子上下坐标！！！！！！！！！！！！！）

//430需求改动，对B点条件1、2 都有改动，并增加B点条件3
B点条件1(买强)：
收盘价>20日均线 And收盘价>10日均线 And收盘价>5日均线 And 看多条件1成立 And 前1交易日看多条件1不成立
B点条件2(买)：
收盘价>20日均线 And收盘价>5日均线 And收盘价>10日均线 And 看多条件2成立And 前1日看多条件2不成立And 个股涨跌幅>3%/指数涨跌幅>1% And 收盘价>60日均线
B点条件3(翻多)：
收盘价>5日均线 and 收盘价>10日均线 and 收盘价>20日均线 and 收盘价>30日均线 and 收盘价>60日均线and 看多条件=True  and 收盘价= 最近10日个交易日最高收盘价 And 收盘价>60日均线
B点条件：
B点条件1成立 Or B点条件2成立 or B点条件3成立

B点显示条件：
第1个B点显示方法为从上市第60个交易日开始，满足B点条件。同时将虚拟仓位设置为100。
第2个B点及以后I点显示方法为满足I点条件且虚拟仓位=0。同时将虚拟仓位设置为100。
显示位置和式样上一版本保持不变。

S点条件1：
(看多条件1不成立 And 连续2天收盘价<20日线) Or (看多条件1成立 And 连续4天收盘价<20日线)
S点条件2：
放量大跌1定义，个股涨跌幅<-5%/指数涨跌幅<-2% And (成交额/前1交易日成交额>1.2 OR 成交额/前2交易日成交额>1.2)
放量大跌2定义，个股涨跌幅<-7%/指数涨跌幅<-3.5% And (成交额/前1交易日成交额>1 OR 成交额/前2交易日成交额>1)
放量大跌定义，放量大跌1 Or放量大跌2
S点条件2定义，放量大跌 And 收盘价<10日均线
S点条件3：
收盘价<5日线 And 收盘价<10日线And 收盘价<20日线And 收盘价<30日线And 收盘价<60日线
S点补充条件1：
前1个已显示的I点到现在已经历10个交易日 And 该I点当日收盘价=近11日个交易日最高收盘价

//430需求改动，对S点增加补充条件2
S点补充条件2(空间止损)：
收盘价/B点收盘价<=0.94

S点显示条件：
虚拟仓位=100 And (S点条件1 Or S点条件2 Or S点条件3 Or S点补充条件1 or S点补充条件2)。同时将虚拟仓位设置为0。


******************************用于参考的通达信代码********************************
A5:=MA(CLOSE,5);
A10:=MA(CLOSE,10);
A20:=MA(CLOSE,20);
A30:=MA(CLOSE,30);
A60:=MA(CLOSE,60);

中短期:=A20;
长期:=A60;

走牛:=C>长期 AND 长期>REF(长期,1) AND REF(长期,1)>=REF(长期,2);
RC1:= IF(走牛 AND A5>长期 AND A20>长期 AND C>A30,1,0);
RC2:= IF((A5>中短期 AND 中短期>REF(中短期,1)),1,0);
RC:= RC1 OR RC2;

B11:=RC1 AND REF(RC1,1)=0 AND C>A20;
B22:=REF(RC2,1)=0 AND C/REF(C,1)>1.03 AND C>中短期 AND RC2;
B1:= B11 OR B22;

放量大跌:=(C/REF(C,1)<0.95 AND (V/REF(V,1)>1.2 OR V/REF(V,2)>1.2)) OR (C/REF(C,1)<0.93 AND (V/REF(V,1)>1 OR V/REF(V,2)>1));
S11:= (COUNT(C<A20,2)=2 AND RC1=0) OR (COUNT(C<A20,4)=4 AND RC1=1);
S22:= 放量大跌 AND C<A10;
S33:= BARSLAST(B1)=10 AND REF(C,10)=HHV(C,11);
S44:= C<A5 AND C<A10 AND C<A20 AND C<A30 AND C<A60;

S0:= S11 OR S22 OR S33 OR S44;
S1:= S0 AND (BARSLAST(REF(B1,1))<BARSLAST(REF(S0,1)) OR (COUNT(B1,0)>=1 AND COUNT(S0,0)=1));


STICKLINE(RC=1,OPEN,CLOSE,3,0), COLORRED;
STICKLINE(RC=1,LOW,OPEN,0,0), COLORRED;
STICKLINE(RC=1,CLOSE,HIGH,0,0),COLORRED;
STICKLINE(RC=0,CLOSE,OPEN,3,0), COLORGREEN;
STICKLINE(RC=0,OPEN,LOW,0,0), COLORGREEN;
STICKLINE(RC=0,HIGH,CLOSE,0,0),COLORGREEN;

M5:MA(CLOSE,5);
M10:MA(CLOSE,10);
M20:MA(CLOSE,20);
M30:MA(CLOSE,30);
M60:MA(CLOSE,60);
DRAWTEXT(B11,LOW*0.97,'强'),COLORYELLOW;
DRAWTEXT(B22,LOW*0.97,'弱'),COLORYELLOW;
DRAWTEXT(S1,HIGH*1.03,'卖'),COLORGREEN;
*/

const (
	YELLOW           = 0
	BLUE             = 1
	STYPE_BUY        = 0
	STYPE_SELL       = 1
	BEFORE_KLIEN_LEN = 20 //TODO:考虑新股的情况，这里应该设多少
	Q                = 0
	I                = 1
)

var (
	STOCK_FACTOR = factorStruct{
		B2:  0.03,
		FD1: -0.05,
		FD2: -0.07,
	}

	INDEX_FACTOR = factorStruct{
		B2:  0.01,
		FD1: -0.02,
		FD2: -0.035,
	}
)

//存放个股和指数系数的结构体
type factorStruct struct {
	B2  float64 //B点条件2涨幅
	FD1 float64 //放量大跌1涨幅
	FD2 float64 //放量大跌2涨幅
}

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

	last_ema    *zbbase.MathItem //最后
	last_column *columnStruct
}

type HuangLanQuJian struct {
	zbbase.BaseModule
	last_calc_info_ map[string]*lastCalcInfo
}

func (m *HuangLanQuJian) Init() {
	log.Release("(m *HuangLanQuJian) Init() ----- ")
	m.Filedname = "hlqj"
	m.TBname = "db_zbjs_common.tb_lastcalcinfo"
	m.last_calc_info_ = make(map[string]*lastCalcInfo)
}

func (m *HuangLanQuJian) Name() string {

	return zbbase.M_hlqj
}

//虽然370开始柱子上下坐标已经没用，但是为了兼容已上线客户端，黄蓝柱算法虽然变了，但是其上下坐标仍沿用老算法
//即：s
//  柱子 = 老算法计算出上下坐标 + 新算法计算出黄蓝色（贼恶心）
//  BS点 = 新算法
//func (m *HuangLanQuJian) Calc(data map[string][]*structs.KlineInfo) (new_calc_time uint32, zb_value string, err error) {
func (m *HuangLanQuJian) Calc() (error, uint32) {
	log.Release("HuangLanQuJian start Calc")

	insert_columns := []columnStruct{}
	delete_columns := []columnStruct{}
	insert_dots := []dotStruct{}
	delete_dots := []dotStruct{}
	begin_time := time.Now()
	var new_calc_time uint32
	var err error

	for code, klines := range global.StockData {
		klen := len(klines.Cbs)
		if klen < 62 {
			continue
		}
		klinedata := make([]*global.CombineItem, klen)
		copy(klinedata, klines.Cbs)

		timeList := GetBuyTimeByStock(code, 1)
		//选择系数因子，个股和指数不一样（从app370开始计算指数BS点）
		factor := STOCK_FACTOR
		if _, ok := global.GServer.GetZhiBiaoMgr().GetIndexCodes()[code]; ok {
			factor = INDEX_FACTOR
		}

		ma5, _ := zbbase.MA_func(5, klinedata, func(info *global.CombineItem) float64 {
			return info.FQKLine.Close
		})
		ma10, _ := zbbase.MA_func(10, klinedata, func(info *global.CombineItem) float64 {
			return info.FQKLine.Close
		})
		ma20, _ := zbbase.MA_func(20, klinedata, func(info *global.CombineItem) float64 {
			return info.FQKLine.Close
		})
		ma30, _ := zbbase.MA_func(30, klinedata, func(info *global.CombineItem) float64 {
			return info.FQKLine.Close
		})
		ma60, _ := zbbase.MA_func(60, klinedata, func(info *global.CombineItem) float64 {
			return info.FQKLine.Close
		})
		zouniu := []bool{}
		rc1 := []bool{} //多看条件1
		rc2 := []bool{} //多看条件2
		// 看多条件1: 收盘价>60日均线 And 60日均线>前1交易日60日均线 And 前1交易日60日均线>前2交易日60日均线 And 5日均线>60日均线 And 20日均线>60日均线 And 收盘价>30日均线
		for i := 2; i < len(ma60); i++ {
			tmp := klinedata[i+59].FQKLine.Close > ma60[i].Value && ma60[i].Value > ma60[i-1].Value && ma60[i-1].Value > ma60[i-2].Value
			zouniu = append(zouniu, tmp)
		}
		for i := 0; i < len(zouniu); i++ {
			tmp := zouniu[i] && ma5[i+57].Value > ma60[i+2].Value && ma20[i+42].Value > ma60[i+2].Value && klinedata[i+61].FQKLine.Close > ma30[i+32].Value
			rc1 = append(rc1, tmp)
		}
		// 看多条件2: 5日均线>20日均线 And 20日均线>前1交易日20日均线
		for i := 1; i < len(ma20); i++ {
			tmp := ma5[i+15].Value > ma20[i].Value && ma20[i].Value > ma20[i-1].Value
			rc2 = append(rc2, tmp)
		}

		i1 := []bool{}
		i2 := []bool{}
		i3 := []bool{}
		//B点条件1(买强):收盘价>20日均线 And收盘价>10日均线 And收盘价>5日均线 And 看多条件1成立 And 前1交易日看多条件1不成立
		for i := 1; i < len(rc1); i++ {
			tmp := klinedata[i+61].FQKLine.Close > ma20[i+42].Value && klinedata[i+61].FQKLine.Close > ma5[i+57].Value && klinedata[i+61].FQKLine.Close > ma10[i+52].Value && rc1[i] && !rc1[i-1]
			i1 = append(i1, tmp)
		}

		//B点条件2(买):收盘价>20日均线 And 收盘价>10日均线 And收盘价>5日均线 And 看多条件2成立And 前1日看多条件2不成立And 个股涨跌幅>3%/指数涨跌幅>1% And 收盘价>60日均线
		/////////////////////////////////////////////////////////////////////old//////////////////////////////////////////////////////////////////////
		//for i := 1; i < len(rc2); i++ {
		//	tmp := klines[i+20].Kline_binfo.Close > ma20[i+1].Value && klines[i+20].Kline_binfo.Close > ma10[11+i].Value && klines[i+20].Kline_binfo.Close > ma5[15+i].Value &&
		//		rc2[i] && !rc2[i-1] && (klines[i+20].Kline_binfo.Close-klines[i+20].Kline_binfo.PreClose)/klines[i+20].Kline_binfo.PreClose > factor.B2
		//	i2 = append(i2, tmp)
		//}
		//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		for i := 29; i < len(rc2)-10; i++ {
			tmp := klinedata[i+30].FQKLine.Close > ma20[i+11].Value &&
				klinedata[i+30].FQKLine.Close > ma10[i+21].Value &&
				klinedata[i+30].FQKLine.Close > ma5[i+26].Value &&
				(klinedata[i+30].FQKLine.Close-klinedata[i+30].FQKLine.PreClose)/klinedata[i+30].FQKLine.PreClose > factor.B2 &&
				klinedata[i+30].FQKLine.Close > ma60[i-29].Value &&
				rc2[i+10] && !rc2[i+10-1]
			i2 = append(i2, tmp)
		}
		///////////////////////////////////////////////////////////////////new////////////////////////////////////////////////////////////////////////

		//////////////////////////////////////////////////////////////////新增加//////////////////////////////////////////////////////////////////////
		// B点条件3(翻多)：收盘价>5日均线 and 收盘价>10日均线 and 收盘价>20日均线 and 收盘价>30日均线 and 收盘价>60日均线and 看多条件=True  and 收盘价= 最近10日个交易日最高收盘价 And 收盘价>60日均线
		for i := 0; i < len(ma60); i++ {
			if i >= 2 {
				break
			}
			maxClose := klinedata[i+59].FQKLine.Close
			isMax := true
			for n := i + 59 - 9; n <= i+59; n++ {
				if klinedata[n].FQKLine.Close > maxClose {
					isMax = false
					break
				}
			}

			tmp := klinedata[i+59].FQKLine.Close > ma5[54].Value && klinedata[i+59].FQKLine.Close > ma10[i+50].Value &&
				klinedata[i+59].FQKLine.Close > ma20[i+40].Value && klinedata[i+59].FQKLine.Close > ma30[i+30].Value &&
				klinedata[i+59].FQKLine.Close > ma60[i].Value && rc2[39+i] && isMax
			i3 = append(i3, tmp)
		}
		for i := 0; i < len(rc1); i++ {
			maxClose := klinedata[i+61].FQKLine.Close
			isMax := true
			for n := i + 61 - 9; n <= i+61; n++ {
				if klinedata[n].FQKLine.Close > maxClose {
					isMax = false
					break
				}
			}
			tmp := klinedata[i+61].FQKLine.Close > ma5[56].Value && klinedata[i+61].FQKLine.Close > ma10[i+52].Value &&
				klinedata[i+61].FQKLine.Close > ma20[i+42].Value && klinedata[i+61].FQKLine.Close > ma30[i+32].Value &&
				klinedata[i+61].FQKLine.Close > ma60[i+2].Value && (rc2[41+i] || rc1[0+i]) && isMax
			i3 = append(i3, tmp)
		}
		//////////////////////////////////////////////////////////////////新增加////////////////////////////////////////////////////////////////////////

		q11 := []bool{}
		q12 := []bool{}
		q1 := []bool{}
		//S点条件1(有效跌破20日均线)：
		//(看多条件1不成立 And 连续2天收盘价<20日线) Or (看多条件1成立 And 连续4天收盘价<20日线)
		for i := 0; i < len(rc1); i++ {
			tmp := !rc1[i] && klinedata[i+61].FQKLine.Close < ma20[i+42].Value && klinedata[i+60].FQKLine.Close < ma20[i+41].Value
			q11 = append(q11, tmp)
		}
		for i := 0; i < len(rc1); i++ {
			tmp := rc1[i] && klinedata[i+61].FQKLine.Close < ma20[i+42].Value && klinedata[i+60].FQKLine.Close < ma20[i+41].Value &&
				klinedata[i+59].FQKLine.Close < ma20[i+40].Value && klinedata[i+58].FQKLine.Close < ma20[i+39].Value
			q12 = append(q12, tmp)
		}
		for i := 0; i < len(rc1); i++ {
			q1 = append(q1, q11[i] || q12[i])
		}

		//S点条件2（大跌破10日均线）：
		//放量大跌1定义，个股涨跌幅<-5%/指数涨跌幅<-2% And (成交额/前1交易日成交额>1.2 OR 成交额/前2交易日成交额>1.2)
		//放量大跌2定义，个股涨跌幅<-7%/指数涨跌幅<-3.5% And (成交额/前1交易日成交额>1 OR 成交额/前2交易日成交额>1)
		//放量大跌定义，放量大跌1 Or放量大跌2
		//S点条件2定义，放量大跌 And 收盘价<10日均线
		q2 := []bool{}
		for i := 10; i < len(klinedata); i++ {
			rate := (klinedata[i].FQKLine.Close - klinedata[i].FQKLine.PreClose) / klinedata[i].FQKLine.PreClose
			fl1 := rate < factor.FD1 &&
				(klinedata[i].FQKLine.Turnover/klinedata[i-1].FQKLine.Turnover > 1.2 || klinedata[i].FQKLine.Turnover/klinedata[i-2].FQKLine.Turnover > 1.2)
			fl2 := rate < factor.FD2 &&
				(klinedata[i].FQKLine.Turnover/klinedata[i-1].FQKLine.Turnover > 1 || klinedata[i].FQKLine.Turnover/klinedata[i-2].FQKLine.Turnover > 1)
			fl := fl1 || fl2
			q2 = append(q2, fl && klinedata[i].FQKLine.Close < ma10[i-10].Value)
		}

		//S点条件3(跌破所有均线)：
		//收盘价<5日线 And 收盘价<10日线And 收盘价<20日线And 收盘价<30日线And 收盘价<60日线
		q3 := []bool{}
		for i := 0; i < len(ma60); i++ {
			close := klinedata[i+59].FQKLine.Close
			tmp := close < ma5[i+55].Value && close < ma10[i+50].Value &&
				close < ma20[i+40].Value && close < ma30[i+30].Value && close < ma60[i].Value
			q3 = append(q3, tmp)
		}

		//前1个已显示的I点到现在已经历10个交易日 And 该I点当日收盘价=近11日个交易日最高收盘价(I点也就是B点)
		q4 := []bool{}
		for i := 0; i < len(klinedata); i++ {
			date, ret := GetLastBuyTime(timeList, int(klinedata[i].Tradday))
			if ret {
				res := IsHighestClose(date, klinedata, i)
				q4 = append(q4, res)
			} else {
				q4 = append(q4, false)
			}
		}

		//收盘价/B点收盘价<=0.94
		q5 := []bool{}
		for i := 0; i < len(klinedata); i++ {
			date, ret := GetLastBuyTime(timeList, int(klinedata[i].Tradday))
			if ret {
				result := GetLastBuyClose(date, klinedata, i)
				if result <= 0 {
					q5 = append(q5, false)
				} else {
					if klinedata[i].FQKLine.Close/result <= 0.94 {
						q5 = append(q5, true)
					} else {
						q5 = append(q5, false)
					}
				}
			} else {
				q5 = append(q5, false)
			}
		}

		min_len := len(i1)
		i2 = i2[len(i2)-min_len:]
		i3 = i3[len(i3)-min_len:]
		iqklines := klinedata[len(klinedata)-min_len:]
		q1 = q1[len(q1)-min_len:]
		q2 = q2[len(q2)-min_len:]
		q3 = q3[len(q3)-min_len:]
		q4 = q4[len(q4)-min_len:]
		q5 = q5[len(q5)-min_len:]

		prev := Q
		iday := 0
		//向数据库插数据要区分两种情况
		//1. 数据库是空的，插入全量数据
		//2. 数据库不是空的
		for i := 0; i < min_len; i++ {
			iday++
			if prev == Q && (i1[i] || i2[i] || i3[i]) {
				prev = I
				iday = 0
				if m.Last_calc_date == 0 {
					insert_dots = append(insert_dots, dotStruct{
						Code:        code,
						Trad_day:    iqklines[i].Tradday,
						Upsert_time: iqklines[i].Upserttime,
						Rsttype:     STYPE_BUY,
					})
				} else if m.Last_calc_date <= iqklines[i].Upserttime {
					insert_dots = append(insert_dots, dotStruct{
						Code:        code,
						Trad_day:    iqklines[i].Tradday,
						Upsert_time: iqklines[i].Upserttime,
						Rsttype:     STYPE_BUY,
					})
					delete_dots = append(delete_dots, dotStruct{
						Code:        code,
						Trad_day:    iqklines[i].Tradday,
						Upsert_time: iqklines[i].Upserttime,
					})
				}
			} else if prev == I && (q1[i] || q2[i] || q3[i] || q4[i] || q5[i]) {
				prev = Q
				if m.Last_calc_date == 0 {
					insert_dots = append(insert_dots, dotStruct{
						Code:        code,
						Trad_day:    iqklines[i].Tradday,
						Upsert_time: iqklines[i].Upserttime,
						Rsttype:     STYPE_SELL,
					})
				} else if m.Last_calc_date <= iqklines[i].Upserttime {
					insert_dots = append(insert_dots, dotStruct{
						Code:        code,
						Trad_day:    iqklines[i].Tradday,
						Upsert_time: iqklines[i].Upserttime,
						Rsttype:     STYPE_SELL,
					})
					delete_dots = append(delete_dots, dotStruct{
						Code:        code,
						Trad_day:    iqklines[i].Tradday,
						Upsert_time: iqklines[i].Upserttime,
					})
				}
			}
		}

		//下面原本是最老的“黄蓝柱”算法，当时要求实时计算，所以使用增量计算，增加了代码复杂度
		//新算法只是“黄蓝”算法，“柱”部分需沿用旧算法，因此保留下面这段代码
		//即：
		//  旧算法算出“黄蓝柱”，新算法将柱子重新染色
		if klinedata[klen-1].Upserttime > new_calc_time {
			new_calc_time = klinedata[klen-1].Upserttime
		}
		var ema_as, ema_bs []*zbbase.MathItem
		var startup_first_calc bool
		if last_calc_info, ok := m.last_calc_info_[code]; ok {
			var last_index int
			startup_first_calc = false
			for last_index = klen - 1; last_index >= 0; last_index-- {
				if klinedata[last_index].Upserttime == last_calc_info.upserttime {
					break
				}
			}
			klinedata = klinedata[last_index-10:]
			jjs := zbbase.JJ(klinedata)
			ema_as = zbbase.EMA(10, jjs, last_calc_info.last_ema)
		} else {
			startup_first_calc = true
			jjs := zbbase.JJ(klinedata)
			ema_as = zbbase.EMA(10, jjs, nil)
		}

		ema_bs = ema_as[0 : len(ema_as)-1] //REF
		ema_as = ema_as[1:]                //截成和ema_bs一样长
		elen := len(ema_as)
		klinedata = klinedata[len(klinedata)-elen:] //截成一样长方便计算
		//yes_color := m.last_calc_info_[code].last_column.Color
		rc1_delta := elen - len(rc1)
		rc2_delta := elen - len(rc2)
		for i := 0; i < elen; i++ {
			//如果计算过这天，或者新股上市60天内，就跳过 
			if klinedata[i].Upserttime < m.Last_calc_date || (startup_first_calc && i < 50 /*48*/) { //48才是精确的，但是会crash
				continue
			}
			if klinedata[i].Upserttime >= m.Last_calc_date && m.Last_calc_date != 0 {
				delete_columns = append(delete_columns, columnStruct{
					Code:        code,
					Trad_day:    klinedata[i].Tradday,
					Upsert_time: klinedata[i].Upserttime,
				})
			}
			today := ema_as[i].Value
			yesterday := ema_bs[i].Value
			color := BLUE
			if rc1[i-rc1_delta] || rc2[i-rc2_delta] {
				color = YELLOW
			}
			if today > yesterday {
				insert_columns = append(insert_columns, columnStruct{
					Code:        code,
					Trad_day:    klinedata[i].Tradday,
					Upsert_time: klinedata[i].Upserttime,
					High:        int(today * 1000),
					Low:         int(yesterday * 1000),
					Color:       color,
				})
			} else {
				insert_columns = append(insert_columns, columnStruct{
					Code:        code,
					Trad_day:    klinedata[i].Tradday,
					Upsert_time: klinedata[i].Upserttime,
					High:        int(yesterday * 1000),
					Low:         int(today * 1000),
					Color:       color,
				})
			}
		}
		if elen > BEFORE_KLIEN_LEN {
			m.last_calc_info_[code] = &lastCalcInfo{}
			m.last_calc_info_[code].last_ema = &zbbase.MathItem{
				ema_as[elen-BEFORE_KLIEN_LEN].Value,
			}
			m.last_calc_info_[code].upserttime = klinedata[elen-BEFORE_KLIEN_LEN].Upserttime
		}
	}
	calc_time := time.Now()
	log.Release("HuangLanQuJian calc done, insert_columns:%d, delete_columns:%d, insert_dots:%d, delete_dots:%d, calc used time:%dms",
		len(insert_columns), len(delete_columns), len(insert_dots), len(delete_dots), utils.Millisecond(calc_time.Sub(begin_time)))

	//操作数据库
	deleteRecord(delete_columns, delete_dots)
	insertRecord(insert_columns, insert_dots)
	sql_time := time.Now()
	log.Release("HuangLanQuJian sql operation done, used time:%dms", utils.Millisecond(sql_time.Sub(calc_time)))

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
		delete_sql := `delete from db_zbjs_dcxg.tb_hlqj_column where `
		for i, item := range items {
			var del_sql string
			if i != 0 {
				del_sql += " or "
			}
			//根据联合主键删除比较快，所以YELLOW BLUE都写进删除语句
			del_sql += fmt.Sprintf(`(stockcode="%s" and sdata=%d and upserttime=%d and rsttype=%d)`,
				item.Code, item.Trad_day, item.Upsert_time, YELLOW)
			del_sql += fmt.Sprintf(` or (stockcode="%s" and sdata=%d and upserttime=%d and rsttype=%d)`,
				item.Code, item.Trad_day, item.Upsert_time, BLUE)
			delete_sql += del_sql
		}
		delete_sql += ";"
		_, err = mysql_engine.Exec(delete_sql)
		if err != nil {
			log.Error("HuangLanQuJian::deleteRecord delete_sql err:%v", err)
			return
		}
		if end_index == len(del_cols) {
			break
		}
		index++
	}

	//操作点数据
	index = 0
	for len(del_dots) > 0 {
		bg_index := index * zbbase.SQL_OPR_BATCH_CNT
		end_index := (index + 1) * zbbase.SQL_OPR_BATCH_CNT
		if end_index > len(del_dots) {
			end_index = len(del_dots)
		}
		items := del_dots[bg_index:end_index]
		delete_sql := `delete from db_zbjs_dcxg.tb_hlqj_dot where `
		for i, item := range items {
			var del_sql string
			if i != 0 {
				del_sql += " or "
			}
			del_sql += fmt.Sprintf(`(stockcode="%s" and sdata=%d and upserttime=%d and rsttype=%d)`,
				item.Code, item.Trad_day, item.Upsert_time, STYPE_BUY)
			del_sql += fmt.Sprintf(` or (stockcode="%s" and sdata=%d and upserttime=%d and rsttype=%d)`,
				item.Code, item.Trad_day, item.Upsert_time, STYPE_SELL)
			delete_sql += del_sql
		}
		delete_sql += ";"
		_, err = mysql_engine.Exec(delete_sql)
		if err != nil {
			log.Error("HuangLanQuJian::deleteRecord delete_sql err:%v", err)
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
		insert_sql := "insert into db_zbjs_dcxg.tb_hlqj_column(stockcode, sdata, upserttime, rsttype, low, high) values "
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
			log.Error("HuangLanQuJian::insertRecord sql err:%v", err)
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
		insert_sql := "insert into db_zbjs_dcxg.tb_hlqj_dot(stockcode, sdata, upserttime, rsttype) values "
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
			log.Error("HuangLanQuJian::insertRecord sql err:%v", err)
			return
		}
		if end_index == len(dots) {
			break
		}
		index++
	}
	return
}

// 获取某一只股票 从2018年1月1日起，到现在为止买点的时间列表。
func GetBuyTimeByStock(stock_code string, stype int) []int {
	timelist := []int{}

	sql_tt := fmt.Sprintf("SELECT sdata FROM db_zbjs_dcxg.tb_hlqj_dot WHERE stockcode='%s' and rsttype = %d and sdata >=20180101 ORDER BY sdata asc ;", stock_code, stype)
	mysql_engine := global.GServer.GetMysqlEngine()

	rows, err := mysql_engine.Query(sql_tt)
	if err != nil {
		return timelist
	}
	defer func() {
		rows.Close()
	}()

	if rows.Next() { //这只股票没有任何历史买卖点
		var timevalue int
		rows.Scan(&timevalue)
		timelist = append(timelist, timevalue)
	}
	return timelist
}

// 获取最近一个买点的时间
func GetLastBuyTime(timelist []int, tradday int) (date int, b bool) {
	tlen := len(timelist)
	if tlen == 1 {
		if timelist[0] < tradday {
			return timelist[0], true
		}
	} else {
		for i, v := range timelist {
			if v > tradday {
				return timelist[i-1], true
			}
		}
	}
	return -1, false
}

// 获取最近一次的收盘价
func GetLastBuyClose(lastBTime int, klines []*global.CombineItem, kindex int) float64 {
	var close float64
	for i := kindex; i > 0; i-- {
		if int(klines[i].Tradday) == lastBTime {
			close = klines[i].FQKLine.Close
			break
		}
	}

	return close
}

// 判断是否符合 "前1个已显示的I点到现在已经历10个交易日 And 该I点当日收盘价=近11日个交易日最高收盘价"
func IsHighestClose(lastBTime int, klines []*global.CombineItem, kindex int) bool {
	var maxClose float64
	cnt := 0 // 统计当前交易日 与 前一个B点之间 相差多少个交易日
	for i := kindex - 1; i > 0; i-- {
		if int(klines[i].Tradday) == lastBTime {
			maxClose = klines[i].FQKLine.Close
			break
		}
		cnt++
	}

	// 没有找到当前日期之前最近的一个B点数据
	// 或者距离最近的一个B点时间差 不足10个交易日
	if maxClose == 0 || cnt <= 10 {
		return false
	}

	// 前一个B点的收盘价是否是  近10个交易日的最大值。
	isMax := true
	for i := kindex - 1; i >= kindex-10; i-- {
		if klines[i].FQKLine.Close > maxClose {
			isMax = false
			break
		}
	}

	return isMax
}

// 开始指标计算
func (m *HuangLanQuJian) Start() error {
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
