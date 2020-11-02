package cpgj

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"libs/log"
	"sort"

	"zhibiaocalcsvr/src/etc"
	"zhibiaocalcsvr/src/global"
)

/*--------------------------------------------------回测结果入库 start--------------------------------------------------*/
//回测 数据源结构
type BackTestData struct {
	stockcode string
	tradday   int32
	stype     int32
	close     int32
	virtual   int32
}

//回测 结果
type BackTestResult struct {
	stockcode string
	buyday    int32
	clearday  int32
	holddays  int32
	bclose    float64
	cclose    float64
	rate      float64
	pmrate    float64
}

//回测 临时结构,存放某只股票某个周期内的所有收盘价和虚拟仓位
type TmpBackTest struct {
	close1  float64 //建仓收盘价
	close2  float64 //加仓收盘价
	close3  float64 //减仓收盘价
	close4  float64 //清仓收盘价
	virtual float64 //本周期内:当前买卖点到清仓点时的仓位变化
}

var (
	back_test_data   map[string][]BackTestData   //回测 数据源
	back_test_result map[string][]BackTestResult //回测 结果
	tradday_list     []int32                     //上证指数交易日列表(>= 20180101)
)

func OpMain() {
	GetTraddayList() //获取上证指数的交易日列表
	err := GetAllBackTestData()
	if err != nil {
		return
	}
	SortAndCalc()
	for code, vbtr := range back_test_result {
		if err = InsertBTR(vbtr); err != nil { //回测结果入库
			log.Error("OpMain::InsertBTR error, stockcode:%s, error:[%v]", code, err)
			return
		}
	}
	log.Release("OpMain back test result output done !!")
	return
}

//从mongo中获取上证指数20180101之后的交易日列表
func GetTraddayList() {
	session := global.GServer.GetMongoEngine().Clone()
	defer session.Close()
	table := session.DB("HS").C("kline1440")
	klines := []bson.M{}
	query := bson.M{"stockcode": "sh000001", "tradday": bson.M{"$gte": etc.Config.Calc.Common_StartDate}}
	_ = table.Find(query).Sort("tradday").All(&klines)
	for _, vk := range klines {
		tradday_list = append(tradday_list, int32(vk["tradday"].(int)))
	}
	log.Release("from mongo GetTraddayList tradday_list:%v", tradday_list)
	return
}

func GetAllBackTestData() (err error) {
	sql_tt := fmt.Sprintf("select stockcode, tradday, stype, close, virtual FROM tb_cpgj_v2 where tradday > 20180101;")
	rows, err := global.GServer.GetMysqlEngine().Query(sql_tt)
	if err != nil {
		log.Error("GetAllBackTestData error:[%v], sql:[%s]", err, sql_tt)
		return
	}
	defer func() {
		rows.Close()
	}()
	btd := BackTestData{}
	back_test_data = make(map[string][]BackTestData)
	for rows.Next() { //这只股票没有任何历史买卖点
		rows.Scan(&btd.stockcode, &btd.tradday, &btd.stype, &btd.close, &btd.virtual)
		back_test_data[btd.stockcode] = append(back_test_data[btd.stockcode], btd)
	}
	log.Release("GetAllBackTestData sql:[%s]", sql_tt)
	return
}

func SortAndCalc() {
	back_test_result = make(map[string][]BackTestResult)
	for code, vbtd := range back_test_data {
		tmp := TmpBackTest{}
		btr := BackTestResult{}
		stype := int32(0)             //上一个买卖点类型
		virtual := float64(50)        //仓位变化: 建--50  加--50  减--50  清--50或100
		out := AscSortByTradday(vbtd) //排序
		for _, vo := range out { //计算
			if (stype == 0 || stype == 4) && vo.stype == 1 {
				stype = vo.stype
				btr.buyday = vo.tradday
				tmp.virtual = float64(50)
				tmp.close1 = float64(vo.close) / 1000
				btr.bclose = tmp.close1
				btr.stockcode = code
			} else if stype == 1 && vo.stype == 2 {
				stype = vo.stype
				tmp.close2 = float64(vo.close) / 1000
				tmp.virtual = float64(100)
			} else if stype == 2 && vo.stype == 3 {
				stype = vo.stype
				tmp.close3 = float64(vo.close) / 1000
				tmp.virtual = float64(50)
				if tmp.close2 != 0 { //加减仓收益率 只统计有减仓的情况
					btr.pmrate = tmp.close3/tmp.close2 - 1
				}
			} else if stype != 0 && stype != 4 && vo.stype == 4 {
				stype = vo.stype
				btr.clearday = vo.tradday
				len_list := len(tradday_list)
				i1 := sort.Search(len_list, func(i int) bool { return tradday_list[i] >= btr.buyday })
				i2 := sort.Search(len_list, func(i int) bool { return tradday_list[i] >= btr.clearday })
				btr.holddays = int32(i2 - i1 + 1) //某只股票一个周期内的持有的天数(交易日)
				tmp.close4 = float64(vo.close) / 1000
				btr.cclose = tmp.close4
				close := tmp.close1 + tmp.close2 //建仓加仓收盘价总和
				if close != 0 { //收益率
					btr.rate = (tmp.close3*virtual+tmp.close4*tmp.virtual)/(close*virtual) - 1
				}
				back_test_result[code] = append(back_test_result[code], btr)
				//清空当前周期的缓存, 进行下个周期的计算
				tmp = TmpBackTest{}
				btr = BackTestResult{}
			}
		}
	}
}

func InsertBTR(btr []BackTestResult) (err error) {
	if len(btr) == 0 {
		log.Release("[InsertBTR sql data is NULL] len(btr) == 0")
		return
	}
	i := 0
	insert_sql := "insert into cpgjv2_back_test_v2 (stockcode, buyday, clearday, holddays, bclose, cclose, rate, pmrate) VALUES"
	for _, vb := range btr {
		insert_sql += fmt.Sprintf("('%s',%d,%d,%d,%f,%f,%f,%f),", vb.stockcode, vb.buyday, vb.clearday, vb.holddays, vb.bclose, vb.cclose, vb.rate, vb.pmrate)
		i++
	}
	if i > 0 {
		insert_sql = insert_sql[:len(insert_sql)-1]
		insert_sql += ";"
		_, err = global.GServer.GetMysqlEngine().Exec(insert_sql)
	}
	return
}

//每只股票的k线数据 根据日期从小到大排序
func AscSortByTradday(data []BackTestData) (out []BackTestData) {
	l := len(data)
	if l <= 0 {
		return
	}
	k := 0                //记录每次循环最小值的索引
	min := BackTestData{} //记录每次循环的最小值
	for i := 0; i < l-1; i++ { //只需排出最小的前n个值，不需要全部排序
		min = data[i]
		k = i
		for j := i + 1; j < l; j++ {
			if data[j].tradday < min.tradday {
				min = data[j] //记录本次循环的最小值
				k = j         //记录本次循环最小值的索引
			}
		}
		//与最小的做交换
		data[k] = data[i]
		data[i] = min
	}
	out = make([]BackTestData, l)
	copy(out, data) //采用深拷贝
	return
}

/*---------------------------------------------------回测结果入库 end---------------------------------------------------*/
