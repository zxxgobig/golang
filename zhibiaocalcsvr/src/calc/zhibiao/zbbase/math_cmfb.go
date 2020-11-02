package zbbase

import (
	"math"
	"zhibiaocalcsvr/src/global"
)

/*
LC:=REF(CLOSE,1); 												// lc 昨日收盘价
RSI5:=SMA(MAX(CLOSE-LC,0),7,1)/SMA(ABS(CLOSE-LC),7,1)*100; 		// MAX(CLOSE-LC,0) 差值大于0取差值，不然取 0，  ABS(CLOSE-LC)取绝对值
																// SMA(CLOSE,30,1)表示求30日移动平均价 1为权重   公式中有
																// SMA(X,N,M)，求X的N日移动平均，M为权重。算法：若Y=SMA(X,N,M) 则 Y=(M*X+(N-M)*Y')/N，其中Y'表示上一周期Y值，N必须大于M。(百度说明)
TR1:=SUM(MAX(MAX(HIGH-LOW,ABS(HIGH-REF(CLOSE,1))),ABS(LOW-REF(CLOSE,1))),7); // ll = MAX(HIGH-LOW,ABS(HIGH-REF(CLOSE,1)) 当日最高价-最低价 和 当日最高价-昨收价的绝对值 的高值
																		// ll 和 最低价-昨收价的绝对值的高值
																		// 取7天+和

HD:=HIGH-REF(HIGH,1);
LD:=REF(LOW,1)-LOW;
DMP:=SUM(IF(HD>0 AND HD>LD,HD,0),7); 								// 满足HD>0 AND HD>LD条件 取 hd 不然取0 ， 结果七天+和
DMM:=SUM(IF(LD>0 AND LD>HD,LD,0),7);
PDI:=DMP*100/TR1;
MDI:=DMM*100/TR1;
ADX:=MA(ABS(MDI-PDI)/(MDI+PDI)*100,5);

WR10:=100*(HHV(HIGH,7)-CLOSE)/(HHV(HIGH,7)-LLV(LOW,7));
MM:=-(WR10*2+ADX-RSI5);
NN:=MA(MM,5);
DRAWTEXT(REF(MM,1)>0 AND MM<REF(MM,1),H*1.01,'分'),COLORGREEN;
DRAWTEXT(REF(MM,1)<-200 AND MM>REF(MM,1),L*0.99,'集'),COLORYELLOW;
*/

func MM(src_kline []*global.CombineItem, max_sma_info, abs_sma_info *MathItem) ([]*MathItem, []*MathItem, []*MathItem) {
	if src_kline == nil || len(src_kline) < 9 { //SMA(V,7,1)导致
		return nil, nil, nil
	}

	src_kline_len := len(src_kline)
	base_arr := make([]*MathItem, src_kline_len)
	for i, item := range src_kline {
		base_arr[i] = &MathItem{item.KLine.Close}
	}
	lc_arr := REF(base_arr, 1)

	for i, item := range src_kline {
		base_arr[i] = &MathItem{item.KLine.High}
	}
	hd_arr := REF(base_arr, 1)

	//HHV(HIGH, 7)
	hhv_arr := HHV(base_arr, 7)

	for i, item := range src_kline {
		base_arr[i] = &MathItem{item.KLine.Low}
	}
	ld_arr := REF(base_arr, 1)

	//LLV(LOW,7)
	llv_arr := LLV(base_arr, 7)

	/*
		HD:=HIGH-REF(HIGH,1);
		LD:=REF(LOW,1)-LOW;
	*/
	ref_len := len(lc_arr)
	for i := 0; i < ref_len; i++ {
		cur_k := src_kline[src_kline_len-1-i]
		ld_arr[ref_len-1-i].Value = ld_arr[ref_len-1-i].Value - cur_k.KLine.Low
		hd_arr[ref_len-1-i].Value = cur_k.KLine.High - hd_arr[ref_len-1-i].Value
	}

	max_arr := make([]*MathItem, ref_len)
	abs_arr := make([]*MathItem, ref_len)
	for i, _ := range lc_arr {
		v := math.Max(src_kline[src_kline_len-1-i].KLine.Close-lc_arr[ref_len-1-i].Value, 0)
		max_arr[ref_len-1-i] = &MathItem{v}

		v = math.Abs(src_kline[src_kline_len-1-i].KLine.Close - lc_arr[ref_len-1-i].Value)
		abs_arr[ref_len-1-i] = &MathItem{v}
	}

	max_sma, _ := SMA(7, 1, max_arr, max_sma_info)
	abs_sma, _ := SMA(7, 1, abs_arr, abs_sma_info)
	sma_len := len(max_sma)

	//RSI5:=SMA(MAX(CLOSE-LC,0),7,1)/SMA(ABS(CLOSE-LC),7,1)*100
	rsi5_arr := make([]*MathItem, sma_len)
	for i, max_item := range max_sma {
		rsi5_arr[i] = &MathItem{
			Value: max_item.Value / abs_sma[i].Value * float64(100),
		}
	}

	//TR1:=SUM(MAX(MAX(HIGH-LOW,ABS(HIGH-REF(CLOSE,1))),ABS(LOW-REF(CLOSE,1))),7)
	tr1_src_arr := make([]*MathItem, ref_len)
	for i := 0; i < ref_len; i++ {
		cur_k := src_kline[src_kline_len-1-i]
		v := math.Max(cur_k.KLine.High-cur_k.KLine.Low, math.Abs(cur_k.KLine.High-lc_arr[ref_len-1-i].Value))
		v = math.Max(v, math.Abs(cur_k.KLine.Low-lc_arr[ref_len-1-i].Value))
		tr1_src_arr[ref_len-1-i] = &MathItem{v}
	}
	tr1_arr := SUM(tr1_src_arr, 7)

	//WR10:=100*(HHV(HIGH,7)-CLOSE)/(HHV(HIGH,7)-LLV(LOW,7));
	hhv_len := len(hhv_arr)
	wr10_arr := make([]*MathItem, hhv_len)
	for i := 0; i < hhv_len; i++ {
		cur_k := src_kline[src_kline_len-1-i]
		hl_v := hhv_arr[hhv_len-1-i].Value - llv_arr[hhv_len-1-i].Value
		if hl_v != 0 {
			hl_v = float64(100) * (hhv_arr[hhv_len-1-i].Value - cur_k.KLine.Close) / hl_v
		}
		wr10_arr[hhv_len-1-i] = &MathItem{
			Value: hl_v,
		}
	}

	/*
	DMP:=SUM(IF(HD>0 AND HD>LD,HD,0),7); 								// 满足HD>0 AND HD>LD条件 取 hd 不然取0 ， 结果七天+和
	DMM:=SUM(IF(LD>0 AND LD>HD,LD,0),7);
	*/
	hd_len := len(hd_arr)
	dmp_src_arr := make([]*MathItem, hd_len)
	dmm_src_arr := make([]*MathItem, hd_len)
	for i := 0; i < hd_len; i++ {
		v := float64(0)
		hd_v := hd_arr[hd_len-1-i].Value
		ld_v := ld_arr[hd_len-1-i].Value
		if hd_v > 0 && hd_v > ld_v {
			v = hd_v
		}
		dmp_src_arr[hd_len-1-i] = &MathItem{v}

		v = float64(0)
		if ld_v > 0 && ld_v > hd_v {
			v = ld_v
		}
		dmm_src_arr[hd_len-1-i] = &MathItem{v}
	}
	dmp_arr := SUM(dmp_src_arr, 7)
	dmm_arr := SUM(dmm_src_arr, 7)

	/*
	PDI:=DMP*100/TR1;
	MDI:=DMM*100/TR1;
	ADX:=MA(ABS(MDI-PDI)/(MDI+PDI)*100,5);
	*/
	dmp_len := len(dmp_arr)
	tr1_len := len(tr1_arr)
	ma_arr := make([]*MathItem, dmp_len)
	for i := 0; i < dmp_len; i++ {
		tr_v := tr1_arr[tr1_len-1-i].Value
		if tr_v == 0 {
			ma_arr[dmp_len-1-i] = &MathItem{0}
		} else {
			pdi_v := dmp_arr[dmp_len-1-i].Value * float64(100) / tr_v
			mdi_v := dmm_arr[dmp_len-1-i].Value * float64(100) / tr_v
			v := math.Abs(mdi_v-pdi_v) / (mdi_v + pdi_v) * float64(100)
			ma_arr[dmp_len-1-i] = &MathItem{v}
		}
	}

	adx_arr := MA_arr(5, ma_arr)

	rsi5_len := len(rsi5_arr)
	wr10_len := len(wr10_arr)
	adx_len := len(adx_arr)
	min_len := rsi5_len
	if wr10_len < min_len {
		min_len = wr10_len
	}
	if adx_len < min_len {
		min_len = adx_len
	}
	//MM:=-(WR10*2+ADX-RSI5);
	mm_arr := make([]*MathItem, min_len)
	for i := 0; i < min_len; i++ {
		v := -(wr10_arr[wr10_len-1-i].Value*float64(2) + adx_arr[adx_len-1-i].Value - rsi5_arr[rsi5_len-1-i].Value)
		mm_arr[min_len-1-i] = &MathItem{v}
	}

	return mm_arr, max_sma, abs_sma
}
