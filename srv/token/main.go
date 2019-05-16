package main

import (
	_ "originals/conf"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
)

func main() {
	// 创建服务
	service := micro.NewService(
		micro.Name("go.micro.srv.token"),
		micro.Version("v1.0"),
	)

	// 初始化服务
	service.Init(
		initRedis,
		registerHandler,
	)

	// 运行服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
