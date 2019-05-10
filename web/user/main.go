package main

import (
	"net/http"

	"github.com/micro/go-log"
	"github.com/micro/go-web"
)

func main() {

	// New service
	service := web.NewService(
		web.Name("go.micro.web.user"),
		web.Version("v1.0"),
	)

	// Initialise service
	service.Init()

	service.Handle("/", http.FileServer(http.Dir("originals-user/build")))

	// Run services
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
