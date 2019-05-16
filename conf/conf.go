package conf

import (
	"github.com/micro/go-config"
	"github.com/micro/go-log"
)

var (
	SqlConf   database
	EmailConf email
)

func init() {
	log.Log("initializing global configs")
	// Config file
	if err := config.LoadFile("../../conf.json"); err != nil {
		log.Fatal(err)
	}

	// Sql database
	if err := config.Get("sql").Scan(&SqlConf); err != nil {
		log.Fatal(err)
	}

	// Email
	if err := config.Get("email").Scan(&EmailConf); err != nil {
		log.Fatal(err)
	}
}
