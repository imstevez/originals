package main

import (
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
	service := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("v1.0"),
	)

	// Initialise service
	service.Init(
		initMysqlDB,
		registerHandler,
	)

	// Run services
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
