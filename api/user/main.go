package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-web"
)

func main() {

	// New service
	service := web.NewService(
		web.Name("go.micro.api.user"),
		web.Version("v1.0"),
	)

	// Initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// Register handler
	service.Handle("/", initRouter())

	// Run services
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
