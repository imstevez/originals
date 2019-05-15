package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-web"
)

func main() {
	// 新建服务
	service := web.NewService(
		web.Name("go.micro.api.user"),
		web.Version("v1.0"),
	)

	// 初始化服务
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// 注册路由handler
	service.Handle("/", router())

	// 运行服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
