package global

import (
	"bytes"
	"common/structs"
	"compress/gzip"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"libs/log"
	"sync"
	"sync/atomic"
	"time"
	"zhibiaocalcsvr/src/etc"
)

var (
	StockData  map[string]*StockItem   // 日k线数据
	FuQuanData map[string][]FuQuanInfo // 资讯拉取的复权相关信息

	Last_calc_date uint32

	mutex_kline sync.Mutex
	mutex_ps    sync.Mutex
)

type FuQuanInfo struct {
	Cqcxr             uint32
	Multiply_backward float64
	Add_backward      float64
}

type StockItem struct {
	Cbs      []*CombineItem
	BuyClose float64 //建仓点的收盘价
	BuyDay   uint32  //建仓点日期
	AddClose float64 //加仓点的收盘价
	AddDay   uint32  //加仓点日期

	FirstBuyPrice  float64 //第一次建仓的价格
	FirstBuyDay    uint32  //第一次建仓的日期
	SecondBuyPrice float64
	SecondBuyDay   uint32
}

type CombineItem struct {
	Stockcode    string
	Tradday      uint32
	Upserttime   uint32
	KLine        *structs.KlineBinInfo //不复权的k线
	FQKLine      *structs.KlineBinInfo //后复权的K线
	S_PsPoint    *PsItem
	C_BsPoint    *BsItem
	S_BsPoint    *BsItem
	Amplitude    float64 //10日平均振幅
}

type BsItem struct {
	State int
}

type PsItem struct {
	Tradday    uint32
	Upserttime uint32
	Avgs       []AvgItem
	Highest    []MostItem
	Lowest     []MostItem
	Jump       []JumpItem
}

type AvgItem struct {
	//周期，单位为交易日
	Period int
	//均价
	Price float64
}

type MostItem struct {
	//周期，单位为交易日
	Period int

	//发生最高价、最低价那天的交易日，eg. 20190527
	Tradday uint32

	//tradday当天15:00的unix时间戳
	Upserttime uint32

	//最高价、最低价
	Price float64
}

type JumpItem struct {
	//发生跳断点的日期
	Tradday    uint32
	Upserttime uint32
	//跳断点跳前价格
	Begin_price float64
	//跳断点跳后价格
	Stop_price float64
}

// 拉取的k线包括A股和部分指数的，具体的各个指标计算要根据需求考虑要不要过滤掉指数部分。
func GetHistoryKLine1440(stock_codes []string) error {
	StockData = make(map[string]*StockItem)
	log.Release("zhibiaocalc start get %d stockcodes", len(stock_codes))

	//开启多个goroutines同时读取K线
	codes_per_goroutine := 1000
	goroutines := len(stock_codes)/codes_per_goroutine + 1
	ctx, cancel := context.WithCancel(context.Background())
	var read_stockcodes int32 //已读股票数量，只用于日志输出
	var read_klines int32     //已读k线数量，只用于日志输出
	var wg sync.WaitGroup
	wg.Add(goroutines)
	start_time := time.Now()
	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()

			session := GServer.GetMongoEngine().Copy()
			defer session.Close()

			begin_index := index * codes_per_goroutine
			end_index := begin_index + codes_per_goroutine
			if end_index > len(stock_codes) {
				end_index = len(stock_codes)
			}

			collection := session.DB("HS").C("kline1440")
			for _, stock_code := range stock_codes[begin_index:end_index] {

				//当其他goroutines发生错误，此goroutine也退出
				select {
				case <-ctx.Done():
					return
				default:
				}

				//获取该stock_code在过去X个月的日K
				mgo_klines := []bson.M{}
				query := bson.M{"stockcode": stock_code, "tradday": bson.M{"$gte": etc.Config.Calc.Sighting_StartDate}}
				err := collection.Find(query).Sort("upserttime").All(&mgo_klines)
				if err != nil {
					log.Error("GetHistoryKLine1440 get kline1440 from mongodb fail:%s", err.Error())
					cancel()
					return
				}
				if len(mgo_klines) == 0 {
					continue
				}
				//解析mongo返回的k线数据
				var klines []*structs.KlineInfo
				for _, mgo_kline := range mgo_klines {
					kline_info := &structs.KlineBinInfoInput{}
					err = binary.Read(bytes.NewBuffer(mgo_kline["klineinfo"].([]byte)), binary.LittleEndian, kline_info)
					if err != nil {
						log.Error("GetHistoryKLine1440 read klineinfo fail:%s", err.Error())
						cancel()
						return
					}
					kb := &structs.KlineBinInfo{
						High:     float64(kline_info.High) / 1000.0,
						Open:     float64(kline_info.Open) / 1000.0,
						Low:      float64(kline_info.Low) / 1000.0,
						Close:    float64(kline_info.Close) / 1000.0,
						Avg:      float64(kline_info.Avg) / 1000.0,
						PreClose: float64(kline_info.PreClose) / 1000.0,
						Turnover: float64(kline_info.Turnover) / 1000.0,
						Volume:   kline_info.Volume,
					}
					kline := &structs.KlineInfo{
						Upsert_time: uint32(mgo_kline["upserttime"].(int)),
						Update_time: uint32(mgo_kline["datetime"].(int)),
						Trad_day:    uint32(mgo_kline["tradday"].(int)),
						Kline_binfo: kb,
					}
					klines = append(klines, kline)
				}
				tmp_code_num := atomic.AddInt32(&read_stockcodes, 1)
				tmp_kline_num := atomic.AddInt32(&read_klines, int32(len(klines)))
				var cbs []*CombineItem
				for _, kline := range klines {
					cbs = append(cbs, &CombineItem{
						Stockcode:  stock_code,
						Tradday:    kline.Trad_day,
						Upserttime: kline.Upsert_time,
						KLine:      kline.Kline_binfo,
						S_PsPoint:  &PsItem{},
						C_BsPoint:  &BsItem{},
						S_BsPoint:  &BsItem{},
					})
				}
				mutex_kline.Lock()
				tmp := StockItem{Cbs: cbs}
				StockData[stock_code] = &tmp
				mutex_kline.Unlock()
				if tmp_code_num%100 == 0 {
					log.Release("get %d stockcodes with %d klines", tmp_code_num, tmp_kline_num)
				}
			}
		}(i)
	}

	wg.Wait()
	//cancel()					//这行没有起实际作用，只为了ide不报错
	if ctx.Err() != nil {
		return errors.New("GetHistoryKLine1440 get klines fail")
	}
	end_time := time.Now()
	log.Release("total stockcodes:%d, total klines:%d, used time:%fs", read_stockcodes, read_klines, end_time.Sub(start_time).Seconds())
	return nil
}

//计算复权（后复权）
func ConvertToFuQuan() {
	log.Release("convert kline to fuquan kline")
	for code, stock_item := range StockData {
		var fuquan_info []FuQuanInfo
		if temp, ok := FuQuanData[code]; !ok {
			log.Error("fuquan info not find code:%s", code)
			//continue	//不能continue
		} else {
			fuquan_info = temp
		}
		kline_len := len(stock_item.Cbs)
		fuquan_len := len(fuquan_info)
		for i := kline_len - 1; i >= 0; i = i - 1 {
			flag := false
			for j := fuquan_len - 1; j >= 0; j = j - 1 {
				if stock_item.Cbs[i].Tradday >= fuquan_info[j].Cqcxr {
					mul := float64(float64(fuquan_info[j].Multiply_backward) / 10000.0)
					add := float64(float64(fuquan_info[j].Add_backward) / 10000.0)
					if mul != 0 {
						fqkline := &structs.KlineBinInfo{
							High:     stock_item.Cbs[i].KLine.High*mul + add,
							Open:     stock_item.Cbs[i].KLine.Open*mul + add,
							Low:      stock_item.Cbs[i].KLine.Low*mul + add,
							Close:    stock_item.Cbs[i].KLine.Close*mul + add,
							Avg:      stock_item.Cbs[i].KLine.Avg*mul + add,
							PreClose: stock_item.Cbs[i].KLine.PreClose*mul + add,
							Turnover: stock_item.Cbs[i].KLine.Turnover*mul + add,
							Volume:   stock_item.Cbs[i].KLine.Volume,
						}
						stock_item.Cbs[i].FQKLine = fqkline
						flag = true
						break
					}

				}
			}

			if flag == false {
				stock_item.Cbs[i].FQKLine = stock_item.Cbs[i].KLine
			}
		}
	}
	log.Release("ConvertToFuQuan() convert to fuquan klines done")
}

func GetFuQuanFactor() {
	log.Release("get fuquan factor from redis...")
	FuQuanData = make(map[string][]FuQuanInfo)
	r_conn := GServer.GetRedisZXEngine().Get()
	defer r_conn.Close()
	fuquan_prefix := "dr_detail_"
	for stock_code, _ := range StockData {
		redis_stock_key := fuquan_prefix + stock_code
		sk_data, err := redis.Bytes(r_conn.Do("GET", redis_stock_key))
		if err != nil {
			log.Error("redis get error:%+v, key:%s", err, redis_stock_key)
			continue
		}
		temp, err := gzipDecode(sk_data)
		var info []FuQuanInfo
		if err := json.Unmarshal(temp, &info); err != nil {
			log.Error("json.Unmarshal error:%+v, stockcode:%s, raw_json:%s", err, stock_code, string(temp))
			continue
		}
		log.Release("stockcode:%s, fuquan:%+v", stock_code, info)
		FuQuanData[stock_code] = info
	}
	log.Release("get all fuquan factor done")
}

func gzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}
