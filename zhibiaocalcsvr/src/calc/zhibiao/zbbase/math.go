package zbbase

import (
	"common/structs"
	"fmt"
	"zhibiaocalcsvr/src/global"
)

type MathMAInfo struct {
	Value float64
}

type MathItem struct {
	Value float64
}

/**
 * @brief MA(X,N),返回X的N日移动平均值 算法：(X1+X2+X3+...+Xn)/N
 */
func MA(cycle uint, src_block []*global.CombineItem) (dst_block map[uint32]*MathMAInfo, err error) {
	if src_block == nil {
		err = fmt.Errorf("ma src block nil")
		return
	}

	dst_block = make(map[uint32]*MathMAInfo)
	if int(cycle) > len(src_block) {
		//err = fmt.Errorf("src block len:%d cycle:%d too max", len(src_block), cycle)
		return
	}

	ret_len := len(src_block) - int(cycle) + 1
	sum := float64(0)
	for i := 0; i < ret_len; i++ {
		if i == 0 {
			for k := 0; k < int(cycle); k++ {
				sum += src_block[k].KLine.Close
			}
		} else {
			sum += (src_block[int(cycle)+i-1].KLine.Close - src_block[i-1].KLine.Close)
		}

		dst_block[src_block[int(cycle)+i-1].Upserttime] = &MathMAInfo{sum / float64(cycle)}
	}

	return
}

/*
 同MA，输入输出格式不一样
 */
func MA_arr(cycle uint, src_block []*MathItem) (dst_block []*MathItem) {
	if src_block == nil {
		return
	}

	if int(cycle) > len(src_block) {
		return
	}

	dst_block = []*MathItem{}
	ret_len := len(src_block) - int(cycle) + 1
	sum := float64(0)
	for i := 0; i < ret_len; i++ {
		if i == 0 {
			for k := 0; k < int(cycle); k++ {
				sum += src_block[k].Value
			}
		} else {
			sum += (src_block[int(cycle)+i-1].Value - src_block[i-1].Value)
		}

		dst_block = append(dst_block, &MathItem{sum / float64(cycle)})
	}

	return
}

/*
	同MA，输入输出格式不一样
 */
func MA_func(cycle uint, src_block []*global.CombineItem, fn func(info *global.CombineItem) float64) (dst_block []*MathItem, err error) {
	if src_block == nil {
		err = fmt.Errorf("ma src block nil")
		return
	}

	if int(cycle) > len(src_block) {
		//err = fmt.Errorf("src block len:%d cycle:%d too max", len(src_block), cycle)
		return
	}

	ret_len := len(src_block) - int(cycle) + 1
	sum := float64(0)
	for i := 0; i < ret_len; i++ {
		if i == 0 {
			for k := 0; k < int(cycle); k++ {
				sum += fn(src_block[k])
			}
		} else {
			//sum += (src_block[int(cycle) + i - 1].Kline_binfo.Close - src_block[i - 1].Kline_binfo.Close)
			sum += fn(src_block[int(cycle)+i-1]) - fn(src_block[i-1])
		}

		//dst_block[src_block[int(cycle) + i - 1].Upsert_time] = &MathMAInfo{sum / float64(cycle)}
		dst_block = append(dst_block, &MathItem{sum / float64(cycle)})
	}

	return
}

/*
* @breif *求移动平均值
*        Y=SMA(X,N,M) 算法：Y=[M*X+(N-M)*Y']/N,其中Y'表示上一周期Y值,N必须大于M
*	     SMA(CLOSE,30,1)表示求30日移动平均价 1为权重
*/
func SMA(n, m int, src_block []*MathItem, last_sma_value *MathItem) (dst_block []*MathItem, err error) {
	if src_block == nil {
		err = fmt.Errorf("sma src block nil")
		return
	}

	src_len := len(src_block)
	if m > src_len || n > src_len {
		err = fmt.Errorf("sma m or n param too long")
		return
	}

	ret_len := src_len - n - 1
	if ret_len > 0 {
		dst_block = make([]*MathItem, ret_len)
		last_ret := src_block[n].Value
		if last_sma_value != nil {
			last_ret = last_sma_value.Value
		}

		j := 0
		for i := n + 1; i < src_len; i++ {
			last_ret = (float64(m)*src_block[i].Value + float64(n-m)*last_ret) / float64(n)
			dst_block[j] = &MathItem{last_ret}
			j++
		}
	}

	return
}

//前一个版本的函数不适配顶部搭建的算法。。。
func SMA2(n, m int, src_block []*MathItem, last_sma_value *MathItem) (dst_block []*MathItem, err error) {
	if src_block == nil || len(src_block) <= n {
		return
	}

	ret_len := len(src_block) - n
	if ret_len > 0 {
		if last_sma_value != nil {
			dst_block = append(dst_block, last_sma_value)
		} else {
			dst_block = append(dst_block, src_block[n])
		}
		val := dst_block[0].Value
		for i := 1; i < ret_len; i++ {
			val = (float64(m)*src_block[n+i].Value + float64(n-m)*val) / float64(n)
			dst_block = append(dst_block, &MathItem{val})
		}
	}
	return
}

/**
* @breif 前N个对应值(REF)
*/
func REF(src_block []*MathItem, n int) (dst_block []*MathItem) {
	if src_block == nil {
		return
	}

	src_len := len(src_block)
	if src_len > n {
		ret_len := src_len - n
		dst_block = make([]*MathItem, ret_len)
		for i := 0; i < ret_len; i++ {
			dst_block[i] = &MathItem{src_block[i].Value}
		}
	}

	return
}

func SUM(src_block []*MathItem, n int) (dst_block []*MathItem) {
	if src_block == nil {
		return
	}

	src_len := len(src_block)
	if src_len >= n {
		ret_len := src_len - n + 1
		sum := float64(0)
		dst_block = make([]*MathItem, ret_len)
		for i := 0; i < ret_len; i++ {
			if i == 0 {
				for k := i; k < i+n; k++ {
					sum += src_block[k].Value
				}
			} else {
				sum += (src_block[n+i-1].Value - src_block[i-1].Value)
			}

			dst_block[i] = &MathItem{sum}
		}
	}
	return
}

/*
	以n为周期的最高值
 */
func HHV(src_block []*MathItem, n int) (dst_block []*MathItem) {
	if src_block == nil {
		return
	}

	src_len := len(src_block)
	if src_len >= n {
		ret_len := src_len - n + 1
		dst_block = make([]*MathItem, ret_len)
		for i := 0; i < ret_len; i++ {
			if i == 0 {
				dst_block[i] = &MathItem{src_block[i].Value}
				for k := i; k < n; k++ {
					if dst_block[i].Value < src_block[k].Value {
						dst_block[i].Value = src_block[k].Value
					}
				}
			} else {
				if dst_block[i-1].Value > src_block[i-1].Value {
					dst_block[i] = &MathItem{dst_block[i-1].Value}
					if dst_block[i-1].Value < src_block[i+n-1].Value {
						dst_block[i].Value = src_block[i+n-1].Value
					}
				} else {
					dst_block[i] = &MathItem{src_block[i].Value}
					for j := i; j < i+n; j++ {
						if dst_block[i].Value < src_block[j].Value {
							dst_block[i].Value = src_block[j].Value
						}
					}
				}
			}
		}
	}

	return
}

/*
	同HHV，入参不一样
 */
//func HHV_func(src_block []*structs.KlineInfo, n int, fn func(info *structs.KlineInfo) float64) (dst_block []*MathItem) {
func HHV_func(src_block []*global.CombineItem, n int, fn func(info *global.CombineItem) float64) (dst_block []*MathItem) {
	if src_block == nil {
		return
	}

	src_len := len(src_block)
	if src_len >= n {
		ret_len := src_len - n + 1
		dst_block = make([]*MathItem, ret_len)
		for i := 0; i < ret_len; i++ {
			if i == 0 {
				dst_block[i] = &MathItem{fn(src_block[i])}
				for k := i; k < n; k++ {
					if dst_block[i].Value < fn(src_block[k]) {
						dst_block[i].Value = fn(src_block[k])
					}
				}
			} else {
				if dst_block[i-1].Value > fn(src_block[i-1]) {
					dst_block[i] = &MathItem{dst_block[i-1].Value}
					if dst_block[i-1].Value < fn(src_block[i+n-1]) {
						dst_block[i].Value = fn(src_block[i+n-1])
					}
				} else {
					dst_block[i] = &MathItem{fn(src_block[i])}
					for j := i; j < i+n; j++ {
						if dst_block[i].Value < fn(src_block[j]) {
							dst_block[i].Value = fn(src_block[j])
						}
					}
				}
			}
		}
	}

	return
}

/*
	以n为周期的最低值
 */
func LLV(src_block []*MathItem, n int) (dst_block []*MathItem) {
	if src_block == nil {
		return
	}

	src_len := len(src_block)
	if src_len >= n {
		ret_len := src_len - n + 1
		dst_block = make([]*MathItem, ret_len)
		for i := 0; i < ret_len; i++ {
			if i == 0 {
				dst_block[i] = &MathItem{src_block[i].Value}
				for k := i; k < n; k++ {
					if dst_block[i].Value > src_block[k].Value {
						dst_block[i].Value = src_block[k].Value
					}
				}
			} else {
				if dst_block[i-1].Value < src_block[i-1].Value {
					dst_block[i] = &MathItem{dst_block[i-1].Value}
					if dst_block[i-1].Value > src_block[i+n-1].Value {
						dst_block[i].Value = src_block[i+n-1].Value
					}
				} else {
					dst_block[i] = &MathItem{src_block[i].Value}
					for j := i; j < i+n; j++ {
						if dst_block[i].Value > src_block[j].Value {
							dst_block[i].Value = src_block[j].Value
						}
					}
				}
			}
		}
	}

	return
}

/*
	同LLV，入参变了
 */
//llv_low10 := zbbase.LLV_func(klines, 10, func(info *structs.KlineInfo) float64 {

//func LLV_func(src_block []*structs.KlineInfo, n int, fn func(info *structs.KlineInfo) float64) (dst_block []*MathItem) {
func LLV_func(src_block []*global.CombineItem, n int, fn func(info *global.CombineItem) float64) (dst_block []*MathItem) {
	if src_block == nil {
		return
	}

	src_len := len(src_block)
	if src_len >= n {
		ret_len := src_len - n + 1
		dst_block = make([]*MathItem, ret_len)
		for i := 0; i < ret_len; i++ {
			if i == 0 {
				dst_block[i] = &MathItem{fn(src_block[i])}
				for k := i; k < n; k++ {
					if dst_block[i].Value > fn(src_block[k]) {
						dst_block[i].Value = fn(src_block[k])
					}
				}
			} else {
				if dst_block[i-1].Value < fn(src_block[i-1]) {
					dst_block[i] = &MathItem{dst_block[i-1].Value}
					if dst_block[i-1].Value > fn(src_block[i+n-1]) {
						dst_block[i].Value = fn(src_block[i+n-1])
					}
				} else {
					dst_block[i] = &MathItem{fn(src_block[i])}
					for j := i; j < i+n; j++ {
						if dst_block[i].Value > fn(src_block[j]) {
							dst_block[i].Value = fn(src_block[j])
						}
					}
				}
			}
		}
	}

	return
}

/**
* @breif 若求X的N日指数平滑移动平均，则表达式为：EMA（X，N），算法是：
		 若Y=EMA(X，N)，则Y=[2*X+(N-1)*Y’]/(N+1)，其中Y’表示上一周期的Y值。
*/
func EMA(cycle int, src_block []*MathItem, last_ema_value *MathItem) (dst_block []*MathItem) {
	if src_block == nil || len(src_block) <= cycle {
		return
	}
	ret_len := len(src_block) - cycle
	if ret_len > 0 {
		if last_ema_value != nil {
			dst_block = append(dst_block, last_ema_value)
		} else {
			dst_block = append(dst_block, src_block[cycle])
		}
		val := dst_block[0].Value
		for i := 1; i < ret_len; i++ {
			val = (2.0*src_block[cycle+i].Value + float64(cycle-1)*val) / float64(cycle+1);
			dst_block = append(dst_block, &MathItem{val})
		}
	}
	return
}

func EMA_func(cycle int, src_block []*structs.KlineInfo, last_ema_value *MathItem, fn func(info *structs.KlineInfo) float64) (dst_block []*MathItem) {
	if src_block == nil || len(src_block) <= cycle {
		return
	}
	ret_len := len(src_block) - cycle
	if ret_len > 0 {
		if last_ema_value != nil {
			dst_block = append(dst_block, last_ema_value)
		} else {
			dst_block = append(dst_block, &MathItem{fn(src_block[cycle])})
		}
		val := dst_block[0].Value
		for i := 1; i < ret_len; i++ {
			val = (2.0*fn(src_block[cycle+i]) + float64(cycle-1)*val) / float64(cycle+1);
			dst_block = append(dst_block, &MathItem{val})
		}
	}
	return
}

/**
* @brief 两条线交叉

* 两条线交叉
  nType 0表示取左值,1表示取右值
* 用法:
*   CROSS(L,R);表示当L从下方向上穿过R时返回L或R的单值(通过nType确定)，否则返回无效值
*/
func CROSS(ema_as []*MathItem, ema_bs []*MathItem, klines []*structs.KlineInfo) []uint32 {
	if (ema_as == nil || ema_bs == nil) {
		return nil;
	}
	min := func(a int, b int) int {
		if a > b {
			return b
		}
		return a
	}
	size := min(len(klines),
		min(len(ema_as), len(ema_bs)))
	klines = klines[len(klines)-1-size:]
	ema_as = ema_as[len(ema_as)-1-size:]
	ema_bs = ema_bs[len(ema_bs)-1-size:]

	var out_dots []uint32
	pre_left := ema_as[0]
	pre_right := ema_bs[0]
	var now_left, now_right *MathItem
	for i := 0; i < size; i++ {
		now_left = ema_as[i]
		now_right = ema_bs[i]
		if now_left.Value >= now_right.Value && pre_left.Value < pre_right.Value {
			out_dots = append(out_dots, klines[i].Upsert_time)
		}
		pre_left = now_left
		pre_right = now_right
	}
	return out_dots
}

/**
 * 线性回归斜率
*/
func SLOPE_func(cycle int, src_block []*structs.KlineInfo, last_slope_value *MathItem, fn func(info *structs.KlineInfo) float64) (dst_block []*MathItem) {
	if src_block == nil || len(src_block) <= cycle {
		return
	}

	stride := float64(cycle)
	for i := 0; i+cycle < len(src_block); i++ {
		var dx, dy, dxy, dx2, dbase float64
		for j := i; j < i+cycle; j++ {
			tempX := float64(j)
			tempY := fn(src_block[j])
			dx += tempX;
			dy += tempY;
			dx2 += tempX * tempX;
			dxy += tempX * tempY;
		}
		dbase = stride*dx2 - dx*dx
		if dbase != 0 {
			dst_block = append(dst_block, &MathItem{(stride*dxy - dx*dy) / dbase})
		}
	}
	return
}
