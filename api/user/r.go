package main

import (
	"context"
	"fmt"
	"originals/api/user/handler"
	tokenSrvProto "originals/srv/token/proto"
	userSrvProto "originals/srv/user/proto"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
)

type logCli struct {
	client.Client
}

func (lc *logCli) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	fmt.Printf("[wrapper] client request to service: %s endpoint: %s\n", req.Service(), req.Endpoint())
	return lc.Client.Call(ctx, req, rsp)
}

func initRouter() *gin.Engine {
	c := client.DefaultClient
	userSrv := userSrvProto.NewUserService("go.micro.srv.user", &logCli{c})
	tokenSrv := tokenSrvProto.NewTokenService("go.micro.srv.token", &logCli{c})
	h := &handler.User{
		UserSrv:  userSrv,
		TokenSrv: tokenSrv,
	}
	r := gin.Default()
	r.GET("/user", h.Index)
	r.GET("/user/invite", h.Invite)
	r.POST("/user/register", h.Register)
	return r
}
