package conf

import (
	"github.com/micro/go-config"
)

var (
	SrvNameSpace string
	ApiNameSpace string
	WebNameSpace string
	SqlConf      database
	EmailConf    email
	SrvConf      map[string]service
)

func InitConf() error {
	// Config file
	if err := config.LoadFile("/Users/stevez/gocode/src/originals/conf.json"); err != nil {
		return err
	}

	// Name spaces
	SrvNameSpace = config.Get("name_space", "srv").String("go.micro.srv")
	ApiNameSpace = config.Get("name_space", "api").String("go.micro.api")
	WebNameSpace = config.Get("name_space", "web").String("go.micro.web")

	// Sql database
	if err := config.Get("sql").Scan(&SqlConf); err != nil {
		return err
	}

	// Email
	if err := config.Get("email").Scan(&EmailConf); err != nil {
		return err
	}

	// Services
	if err := config.Get("srv").Scan(&SrvConf); err != nil {
		return nil
	}
	return nil
}
