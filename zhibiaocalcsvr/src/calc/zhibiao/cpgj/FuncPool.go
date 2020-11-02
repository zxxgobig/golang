package cpgj

import (
	"fmt"
	"libs/log"
	"zhibiaocalcsvr/src/global"
)

type MathItem struct {
	Value float64
}

// cycle日均线
func MA_func(cycle uint, src_block []*global.CombineItem, fn func(info *global.CombineItem) float64) (dst_block []*MathItem, err error) {
	if src_block == nil {
		err = fmt.Errorf("ma src block nil")
		return
	}
	if int(cycle) > len(src_block) {
		return
	}
	ret_len := len(src_block) - int(cycle) + 1
	sum := float64(0)
	for i := 0; i < ret_len; i++ {
		//假设 cycle == 60
		if i == 0 {
			//前 60 日收盘价总和
			for k := 0; k < int(cycle); k++ {
				sum += fn(src_block[k])
			}
		} else {
			//前i个交易日的 60日收盘价 - 前i-1个交易日的 60日收盘价  总和
			sum += fn(src_block[int(cycle)+i-1]) - fn(src_block[i-1])
		}
		//一个周期60天,假设某只股票有65根K线(k[0]--k[64] 时间戳从小到大排序)：
		//则 ret_len = 65 - 60 + 1 == 6个交易日的60日均价
		//k[0]--k[59]的收盘价均值是k[59]的60日均价
		//k[1]--k[60]的收盘价均值是k[60]的60日均价
		//k[2]--k[61]的收盘价均值是k[61]的60日均价
		//k[3]--k[62]的收盘价均值是k[62]的60日均价
		//k[4]--k[63]的收盘价均值是k[63]的60日均价
		//k[5]--k[64]的收盘价均值是k[64]的60日均价
		dst_block = append(dst_block, &MathItem{sum / float64(cycle)})
	}
	return
}

//买卖点结果入主库
func InsertBS(sql_items []sqlBsItem) (err error) {
	i := 0
	insert_sql := "INSERT INTO tb_cpgj_v2 (stockcode,tradday,stype,open,close,virtual,reason) VALUES"
	for _, item := range sql_items {
		insert_sql += fmt.Sprintf("('%s',%d,%d,%d,%d,%d,'%s'),", item.stockcode, item.tradday, item.stype, item.open, item.close, item.virtual, item.reason)
		i++
	}
	if i > 0 {
		insert_sql = insert_sql[:len(insert_sql)-1]
		insert_sql += ";"
		log.Release("[insertBs sql info] sql:[%s]", insert_sql)
		_, err = global.GServer.GetMysqlEngine().Exec(insert_sql)
	}
	return
}

//从机中 获取某只股票最后一次的买卖点状态和虚拟仓位
func GetLastBSTypeAndVirtual(stock_code string) (stype, virtual int, err error) {
	sql_tt := fmt.Sprintf("SELECT tradday, stype, virtual FROM tb_cpgj_v2 WHERE stockcode='%s' ORDER BY tradday desc limit 1;", stock_code)
	rows, err := global.GServer.GetMysqlEngine().Query(sql_tt)
	if err != nil {
		return
	}
	defer func() {
		rows.Close()
	}()

	if !rows.Next() { //这只股票没有任何历史买卖点
		log.Release("[GetLastBSTypeAndVirtual is NULL] stockcode:%s no history BS", stock_code)
		return
	}
	var tradday int
	rows.Scan(&tradday, &stype, &virtual)
	return
}

//获取某只股票最后一次【建仓/加仓】点的收盘价
//stype: 1--建仓点   2--加仓点
func GetLastCloseByStock(stock_code string, stype int) (tradday, close int, err error) {
	sql_tt := fmt.Sprintf("SELECT tradday, stype, close FROM tb_cpgj_v2 WHERE stockcode='%s' and stype = %d ORDER BY tradday desc limit 1;", stock_code, stype)
	rows, err := global.GServer.GetMysqlEngine().Query(sql_tt)
	if err != nil {
		return
	}
	defer func() {
		rows.Close()
	}()

	if !rows.Next() { //这只股票没有任何历史买卖点
		log.Release("[GetLastCloseByStock is NULL] stockcode:%s no history BS", stock_code)
		return
	}
	rows.Scan(&tradday, &stype, &close)
	return
}

// 判断当前交易日距离上一个时间点(建仓点或者是加仓点)，是否超过10个交易日。
// cbs: 原始未截取过k线
// buyTradday: 上一个时间点
// index: 截取过的k线下标索引
// distance: k线截取的长度
// 返回值 true:超过10个交易日. false:没超过10个交易日。
func IsMoreThanTenDays(stockcode string, buyTradday uint32, cbs []*global.CombineItem, index int, distance int) bool {
	// 还原k线截取之前的下表索引
	cur_index := index + distance

	//if "sh600519" == stockcode {
	//	log.Debug("PPPPPP buyTradday = %v, tradday = %v <--1--> ", buyTradday, cbs[cur_index].Tradday)
	//}
	if cur_index < 10 {
		return false
	}

	res := true
	for i := cur_index; i > cur_index-10; i-- {
		if cbs[i].Tradday < buyTradday {
			res = false
			break
		}
	}
	//if "sh600519" == stockcode {
	//	log.Debug("PPPPPP buyTradday = %v,tradday = %v <--2--> , res = %v", buyTradday, cbs[cur_index].Tradday, res)
	//}
	return res
}
