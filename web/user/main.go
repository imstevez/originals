package main

import (
	"net/http"

	"github.com/micro/go-log"
	"github.com/micro/go-web"
)

func main() {
	service := web.NewService(
		web.Name("go.micro.web.user"),
		web.Version("v1.0"),
	)

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	service.Handle("/", http.FileServer(http.Dir("./htmls/build/")))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
