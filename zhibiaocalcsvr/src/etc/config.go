package etc

import (
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	cfg_filepath = "../etc/zhibiaocalcsvr.xml"
)

var (
	Go_Len          = 1000
	Chan_Server_Len = 1000
	Timer_Len       = 100
	Timer_duration  = time.Millisecond * 50
)

var (
	Config ZBConfig

	BinaryEndian          binary.ByteOrder = binary.BigEndian //消息格式为大端解析
	LogicMainTimeDuration                  = time.Millisecond * 50

	SendChanSize         = 100        //发送队列容量
	GoChanSize           = 10         //异步处理队列容量
	LogSize       int64  = 2048000000 //log 日志大小
	LenMsgLen            = 4          //消息头位数
	MaxRecvMsgLen uint32 = 65535      //接收消息大小
	MaxsendMsgLen uint32 = 65535      //发送消息大小

	CloseRate       = 0.02
	MinKLineLength  = 251
	AmplitudePeriod = 10
)

type ZBConfig struct {
	RunCalc   bool   `xml:"run_calc"`
	Timeout   int    `xml:"timeout"`
	Alarm_url string `xml:"alarm_url"`
	ID        int    `xml:"id"`

	Log struct {
		File  string `xml:"file"`
		Level string `xml:"level"`
		Size  int64  `xml:"size"`
	} `xml:"log"`

	Calc struct {
		Common_StartDate   int    `xml:"common_start_date"`
		Sighting_StartDate int    `xml:"sighting_start_date"`
		PrintStatus        int    `xml:"print_status"`
		RawAvgs            string `xml:"avgs"`
		Averages           []int
		RawLowHigh         string `xml:"lowest_highest"`
		LowHigh            []int
		Blacklist          string `xml:"black_list"`
	} `xml:"calc"`

	Mysql struct {
		Master mysqlConfig `xml:"master"`
		Slave  mysqlConfig `xml:"slave"`
	} `xml:"mysql"`
	Redis struct {
		HangQing redisConfig `xml:"hangqing"`
		ZiXun    redisConfig `xml:"zixun"`
	} `xml:"redis"`
	Mongo struct {
		IP     string `xml:"ip"`
		User   string `xml:"user"`
		Passwd string `xml:"passwd"`
	} `xml:"mongo"`
}

type mysqlConfig struct {
	User    string `xml:"user"`
	Pwd     string `xml:"passwd"`
	IP      string `xml:"ip"`
	Port    int    `xml:"port"`
	MaxIdle int    `xml:"max_idle"`
	MaxOpen int    `xml:"max_open"`
	DB      string `xml:"db"`
}

type redisConfig struct {
	IP   string `xml:"ip"`
	Pwd  string `xml:"passwd"`
	Port int    `xml:"port"`
}

func LoadConfig() error {
	fdata, err := ioutil.ReadFile(cfg_filepath)
	if err != nil {
		panic(err)
	}
	err = xml.Unmarshal(fdata, &Config)
	if err != nil {
		panic(err)
	}

	Config.Calc.Averages = []int{}
	tmp := strings.Split(Config.Calc.RawAvgs, ",")
	fmt.Println("avg", tmp)
	for _, it := range tmp {
		num, err := strconv.Atoi(it)
		if err != nil {
			return err
		}
		Config.Calc.Averages = append(Config.Calc.Averages, num)
	}

	Config.Calc.LowHigh = []int{}
	tmp = strings.Split(Config.Calc.RawLowHigh, ",")
	fmt.Println("lowhigh", tmp)
	for _, it := range tmp {
		num, err := strconv.Atoi(it)
		if err != nil {
			return err
		}
		Config.Calc.LowHigh = append(Config.Calc.LowHigh, num)
	}

	if Config.Timeout == 0 {
		Config.Timeout = 20
	}
	return nil
}
