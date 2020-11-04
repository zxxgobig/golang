package main

import (
	"libs/log"
	"libs/utils"
	syslog "log"
	"strings"

	"zhibiaocalcsvr/src/calc"
	"zhibiaocalcsvr/src/etc"
	"zhibiaocalcsvr/src/global"
)

func main() {
	err := etc.LoadConfig()
	if err != nil {
		panic(err)
	}
	logger, err := log.New(log.CreateFileLog(etc.Config.Log.File, etc.Config.Log.Size), //日志输出定向到文件
		strings.ToLower(etc.Config.Log.Level),                                          //日志等级
		syslog.LstdFlags|syslog.Lshortfile,
		false)
	if err != nil {
		panic(err)
	}
	log.Export(logger)

	err = global.InitBlackList()

	services := []utils.IService{}
	services = append(services, &calc.CalcService{})

	utils.RunMutli(services...)
}
