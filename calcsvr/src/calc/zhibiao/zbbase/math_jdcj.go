package zbbase

import (
	"zhibiaocalcsvr/src/global"
)

//HXLXA= 5日均线
func HXLAXA(src_block []*global.CombineItem) (dst_block map[uint32]*MathMAInfo, err error) {
	dst_block, err = MA(5, src_block)
	return
}

//HXLXB= 10日均线
func HXLXB(src_block []*global.CombineItem) (dst_block map[uint32]*MathMAInfo, err error) {
	dst_block, err = MA(10, src_block)
	return
}

//HXLXD= (收盘价-13日均线)/13日均线*100/0.7
func HXLXD(cur_close float64, ma13 float64) (value float64) {
	value = (cur_close - ma13) / ma13 * float64(100) / float64(0.7)
	return
}

//HXLXE= (收盘价-6日均线)/6日均线*100*2
func HXLXE(cur_close float64, ma6 float64) (value float64) {
	value = (cur_close - ma6) / ma6 * float64(100) * float64(2)
	return
}

//HXLXF= (收盘价-11日均线)/11日均线*100/0.3/1.3/1.5
func HXLXF(cur_close float64, ma11 float64) (value float64) {
	value = (cur_close - ma11) / ma11 * float64(100) / float64(0.3) / float64(1.3) / float64(1.5)
	return
}

//C1= 昨天HXLXF<-12 AND 今天HXLXF>=-12
func C1(hxlxf_yestoday, hxlxf_today float64) bool {
	return hxlxf_yestoday < float64(-12) && hxlxf_today >= float64(-12)
}

//昨天HXLXD<-12 今天 HXLXD>=-12
func C2(hxlxd_yestoday, hxlxd_today float64) bool {
	return hxlxd_yestoday < float64(-12) && hxlxd_today >= float64(-12)
}

//((30日均线-最低价)/60日均线)*200
func LJD(ma30 float64, low float64, ma60 float64) float64 {
	return (ma30 - low) / ma60 * float64(200)
}
