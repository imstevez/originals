package main

import (
	"fmt"
	"originals/conf"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
)

func main() {
	// Load Configs
	if err := conf.InitConf(); err != nil {
		log.Fatal(err)
	}

	// New service
	srvConf := conf.SrvConf["user"]
	srvName := fmt.Sprintf("%s.%s", conf.SrvNameSpace, srvConf.Name)
	service := micro.NewService(
		initMysqlDB,
		initRedis,
		registerHandler,
		micro.Name(srvName),
		micro.Version(srvConf.Version),
		sayBye,
	)

	// Initialise service
	service.Init()

	// Run services
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
