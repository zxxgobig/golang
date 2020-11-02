package zbbase

import (
	"zhibiaocalcsvr/src/global"
)

// func JJ(kline *structs.KlineInfo) MathItem {
// 	return (kline.Close + kline.Low + kline.High) / 3
// }

func JJ(klines []*global.CombineItem) []*MathItem {
	var results []*MathItem
	for _, kline := range klines {
		results = append(results,
			&MathItem{(kline.FQKLine.Close + kline.FQKLine.Low + kline.FQKLine.High) / 3})
	}
	return results
}
