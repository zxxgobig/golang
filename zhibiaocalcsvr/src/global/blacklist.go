package global

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"io"
	"libs/log"
	"os"
	"strings"

	"common/pbmessage"
	"zhibiaocalcsvr/src/etc"
)

var (
	black_list map[string]int
)

const (
	RiskStock = 2 //风险股
)

func InBlackList(stock_code string) bool {
	_, ok := black_list[stock_code]
	return ok
}

func InitBlackList() error {
	black_list = make(map[string]int)
	err := loadBlacklistFile()
	if err != nil {
		return err
	}
	err = loadRiskStock()
	if err != nil {
		return err
	}
	return nil
}

func loadRiskStock() error {
	r_conn := GServer.GetRedisZXEngine().Get()
	defer r_conn.Close()
	redis_stock_key := "stockcodes"
	sk_data, err := redis.Bytes(r_conn.Do("get", redis_stock_key))
	if err != nil {
		log.Error("load stockcodes from redis fail:%v", err)
		return err
	}
	sks := &QuoteProto.QuoteBaseStockInfo{}
	err = proto.Unmarshal(sk_data, sks)
	if err != nil {
		log.Error("unmarshal stockcode fail:%v", err)
		return err
	}
	if sks.CodesTable == nil || len(sks.CodesTable) == 0 {
		err = errors.New("redis stockcodes empty")
		log.Error(err.Error())
		return err
	}

	for _, stock := range sks.CodesTable {
		if *stock.DealStatus == RiskStock {
			black_list[*stock.StockCode] = 1
			log.Debug("risk stock  %s", *stock.StockCode)
		}
	}
	return nil
}

func loadBlacklistFile() error {
	file, err := os.Open(etc.Config.Calc.Blacklist)
	defer file.Close()
	if err != nil {
		log.Error("open file %s fail:%s", etc.Config.Calc.Blacklist, err.Error())
		return err
	}
	file_reader := bufio.NewReader(file)
	for {
		line, err := file_reader.ReadString('\n')
		line = strings.Trim(line, " \n\t")
		if line != "" {
			black_list[line] = 1
			log.Debug("blacklist %s", line)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			err = errors.New(fmt.Sprintf("read file %s error:%v", etc.Config.Calc.Blacklist, err))
			log.Error(err.Error())
			return err
		}
	}
	return nil
}
