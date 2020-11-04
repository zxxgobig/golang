package global

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"zhibiaocalcsvr/src/etc"
)

//读写分离：主库只写
func InitMysqlMaster() (err error, master_mysql *sql.DB) {
	data_source_name := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		etc.Config.Mysql.Master.User, etc.Config.Mysql.Master.Pwd,
		etc.Config.Mysql.Master.IP, etc.Config.Mysql.Master.Port,
		etc.Config.Mysql.Master.DB)
	mysql, err := connectMysql(data_source_name, etc.Config.Mysql.Master.MaxIdle, etc.Config.Mysql.Master.MaxOpen)
	if err != nil {
		return err, nil
	}
	master_mysql = mysql

	return nil, master_mysql
}

func connectMysql(data_source_name string, max_idle int, max_open int) (*sql.DB, error) {
	mysql, err := sql.Open("mysql", data_source_name)
	if err != nil {
		return nil, err
	}
	mysql.SetMaxIdleConns(max_idle)
	mysql.SetMaxOpenConns(max_open)
	mysql.SetConnMaxLifetime(time.Hour * 4)

	return mysql, err
}

//读写分离：主库只写
func InitMysqlDev() (err error, dev_mysql *sql.DB) {
	data_source_name := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		etc.Config.Mysql.Master.User, "1234560.",
		"10.216.251.99", 3306,
		etc.Config.Mysql.Master.DB)
	mysql, err := connectMysql(data_source_name, etc.Config.Mysql.Master.MaxIdle, etc.Config.Mysql.Master.MaxOpen)
	if err != nil {
		return err, nil
	}
	dev_mysql = mysql

	return nil, dev_mysql
}
