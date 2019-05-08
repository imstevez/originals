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
		initRedis,
		registerHandler,
		micro.Name("go.micro.srv.token"),
		micro.Version("v1.0"),
		sayBye,
	)

	// Initialise service
	service.Init()

	// Run services
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
