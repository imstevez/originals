package conf

import (
	"github.com/micro/go-config"
)

var (
	SqlConf   database
	EmailConf email
)

func InitConf() error {
	// Config file
	if err := config.LoadFile("../../conf.json"); err != nil {
		return err
	}

	// Sql database
	if err := config.Get("sql").Scan(&SqlConf); err != nil {
		return err
	}

	// Email
	if err := config.Get("email").Scan(&EmailConf); err != nil {
		return err
	}

	return nil
}
