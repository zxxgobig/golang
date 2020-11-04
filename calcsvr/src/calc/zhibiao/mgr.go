package zhibiao

import (
	"bufio"
	"common/pbmessage"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"io"
	"libs/log"
	"os"
	"sync"
	"time"
	"zhibiaocalcsvr/src/calc/zhibiao/caopantixing"
	"zhibiaocalcsvr/src/calc/zhibiao/choumafenbu"
	"zhibiaocalcsvr/src/calc/zhibiao/cpgj"
	"zhibiaocalcsvr/src/calc/zhibiao/huanglanqujian"
	"zhibiaocalcsvr/src/calc/zhibiao/jiandichuji"
	"zhibiaocalcsvr/src/calc/zhibiao/sighting"
	"zhibiaocalcsvr/src/calc/zhibiao/zbbase"
	"zhibiaocalcsvr/src/global"
)

var ZhiBiaoMgr global.IZhiBiao

type zhibiaoMgr struct {
	kline_lock_       sync.RWMutex
	calc_stock_codes_ map[string]struct{}
	calc_index_codes_ map[string]struct{}
}

func init() {
	ZhiBiaoMgr = &zhibiaoMgr{
		calc_stock_codes_: make(map[string]struct{}),
		calc_index_codes_: make(map[string]struct{}),
	}
}

func (mgr *zhibiaoMgr) Start() error {
	// 加载部分指数code
	err := mgr.loadIndexCodes()
	if err != nil {
		log.Error("loadIndexCodes failed ,err = %v", err)
		return err
	}

	// 读取码表
	err, stockcodes := mgr.LoadStockCodes()
	if err != nil {
		log.Error("LoadStockCodes failed, err = %v", err)
		return err
	} else {
		log.Debug("LoadStockCodes success")
	}

	// 拉取日k
	err = global.GetHistoryKLine1440(stockcodes)
	if err != nil {
		log.Error("get kline1440 failed, err = %v", err)
		return err
	} else {
		log.Debug("GetHistoryKLine1440 success")
	}

	//获取复权因子
	global.GetFuQuanFactor()
	//将K线转换成后复权
	global.ConvertToFuQuan()

	modules := []zbbase.ZhibiaoModule{}
	modules = append(modules, &cpgj.BSCalculator{})             // 操盘管家
	modules = append(modules, &huanglanqujian.HuangLanQuJian{}) // bs点
	modules = append(modules, &jiandichuji.JianDiChuJi{})       // 见底出击
	modules = append(modules, &choumafenbu.ChouMaFenBu{})       // 筹码分布
	modules = append(modules, &caopantixing.CaoPanTiXing{})     // 操盘提醒
	modules = append(modules, &sighting.BSCalculator{})         // 狙击镜

	wg := sync.WaitGroup{}
	wg.Add(1)
	// 单独开启协程计算 各个指标
	go func() {
		var errmsg string
		for _, v := range modules {
			v.Init()
			err := v.Start()
			if err == nil {
				log.Release("%v 计算成功", v.Name())
			} else {
				msg:= fmt.Sprintf("%v 计算失败 ,", v.Name())
				errmsg += msg
				log.Release(msg)
			}
			time.Sleep(time.Second)
		}

		//  如果那个模块计算失败， 发送报警。
		if errmsg != "" {
			global.GServer.AlarmError(errmsg, 1)
		}

		wg.Done()
	}()
	wg.Wait()

	log.Release("指标定时计算服务计算完成，10秒后结束程序.")
	go func() {
		for i := 9; i >= 0; i-- {
			log.Release("倒计时 %v", i)
			time.Sleep(time.Second)
		}
		log.Release("指标定时计算服务 退出。")
		time.Sleep(time.Second)
		os.Exit(0)
	}()

	return nil
}

func (mgr *zhibiaoMgr) Stop() {
	log.Release("zhibiaoMgr stop")
}

// 加载部分指数(这部分指数目前仅供bs点计算，其它指标计算暂时用不到)
func (mgr *zhibiaoMgr) loadIndexCodes() error {
	file, err := os.Open("../etc/calc_index_codes.txt")
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == nil {
			mgr.calc_index_codes_[string(line)] = struct{}{}
			log.Release("Load index %s", string(line))
		} else if err == io.EOF {
			return nil
		} else {
			return err
		}
	}
	return err
}

// 获取指数code
func (mgr *zhibiaoMgr) GetIndexCodes() map[string]struct{} {
	return mgr.calc_index_codes_
}

//读取码表
func (mgr *zhibiaoMgr) LoadStockCodes() (error, []string) {
	r_conn := global.GServer.GetRedisHQEngine().Get()
	defer r_conn.Close()
	redis_stock_key := "stockcodes"
	sk_data, err := redis.Bytes(r_conn.Do("get", redis_stock_key))
	if err != nil {
		log.Error("redishq get stockcodes failed")
		return err, nil
	}

	sks := &QuoteProto.QuoteBaseStockInfo{}
	err = proto.Unmarshal(sk_data, sks)
	if err != nil {
		log.Error("unmarshal stockcodes from redis failed")
		return err, nil
	}

	if sks.CodesTable == nil {
		err = fmt.Errorf("redis stockcodes empty")
		return err, nil
	}

	stockcodes := []string{}
	mgr.kline_lock_.Lock()
	mgr.calc_stock_codes_ = make(map[string]struct{})
	for _, stock_item := range sks.CodesTable {
		//A股和部分指数需要计算(以日k数据计算的服务)
		if *stock_item.StockType == "A" {
			stockcodes = append(stockcodes, *stock_item.StockCode)
			mgr.calc_stock_codes_[*stock_item.StockCode] = struct{}{}
			continue
		}
		if _, ok := mgr.calc_index_codes_[*stock_item.StockCode]; ok {
			stockcodes = append(stockcodes, *stock_item.StockCode)
			mgr.calc_stock_codes_[*stock_item.StockCode] = struct{}{}
		}
	}
	mgr.kline_lock_.Unlock()

	return nil, stockcodes
}
