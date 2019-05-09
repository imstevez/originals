package main

import (
	"originals/api/user/handler"
	userSrvProto "originals/srv/user/proto"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
)

func initRouter() *gin.Engine {
	c := client.DefaultClient
	userSrv := userSrvProto.NewUserService("go.micro.srv.user", c)
	h := &handler.User{
		UserSrv: userSrv,
	}
	r := gin.Default()
	r.GET("/user", h.Index)
	r.GET("/user/invite", h.Invite)
	r.POST("/user/register", h.Register)
	return r
}
