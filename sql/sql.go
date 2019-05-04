package sql

import (
	"database/sql"
	"fmt"
	"originals/conf"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var MysqlDB *sql.DB

func InitMysqlDB() error {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
		conf.SqlConf.UserName,
		conf.SqlConf.Password,
		conf.SqlConf.Protocol,
		conf.SqlConf.Host,
		conf.SqlConf.Port,
		conf.SqlConf.Database,
	)
	var err error
	if MysqlDB, err = sql.Open("mysql", dsn); err != nil {
		return err
	}
	MysqlDB.SetConnMaxLifetime(time.Duration(conf.SqlConf.MaxLifeTime) * time.Second)
	MysqlDB.SetMaxOpenConns(conf.SqlConf.MaxOpenConns)
	MysqlDB.SetMaxIdleConns(conf.SqlConf.MaxIdleConns)
	return nil
}
