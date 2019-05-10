package main

import (
	"context"
	"fmt"
	"net/http"
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
	fmt.Printf("[service] client request to service: %s endpoint: %s\n", req.Service(), req.Endpoint())
	return lc.Client.Call(ctx, req, rsp)
}

func initRouter() *gin.Engine {
	c := client.DefaultClient
	userSrv := userSrvProto.NewUserService("go.micro.srv.user", &logCli{c})
	tokenSrv := tokenSrvProto.NewTokenService("go.micro.srv.token", &logCli{c})
	authVerify := func(ctx *gin.Context) {
		authToken := ctx.GetHeader("x-originals-token")
		if authToken == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
			return
		}
		tokenVerifyReq := tokenSrvProto.VerifyAuthTokenReq{Token: authToken}
		tokenVerifyRsp, err := tokenSrv.VerifyAuthToken(context.TODO(), &tokenVerifyReq)
		if err != nil {
			ctx.JSON(500, gin.H{"code": 500, "message": "internal error"})
			return
		}
		switch tokenVerifyRsp.Status {
		case tokenSrvProto.Status_OK:
			ctx.Set("user_id", tokenVerifyRsp.Claims.UserId)
			ctx.Set("email", tokenVerifyRsp.Claims.Email)
			ctx.Set("mobile", tokenVerifyRsp.Claims.Mobile)
			ctx.Set("nickname", tokenVerifyRsp.Claims.Nickname)
			ctx.Set("image_url", tokenVerifyRsp.Claims.ImageUrl)
			ctx.Next()
		case tokenSrvProto.Status_TokenInvalid:
			ctx.JSON(401, gin.H{"code": 401, "message": "auth token invalid"})
		case tokenSrvProto.Status_TokenCanceled:
			ctx.JSON(401, gin.H{"code": 402, "message": "auth token canceled"})
		case tokenSrvProto.Status_TokenExpired:
			ctx.JSON(401, gin.H{"code": 403, "message": "auth token expired"})
		default:
			ctx.JSON(500, gin.H{"code": 500, "message": "internal error"})
		}
		return
	}
	h := &handler.User{UserSrv: userSrv, TokenSrv: tokenSrv}
	router := gin.Default()
	v1 := router.Group("/user/v1")
	{
		v1.GET("/", h.Index)
		v1.GET("/invite", h.Invite)
		v1.POST("/register", h.Register)
		v1.POST("/login", h.Login)
		v1.GET("/logout", h.Logout)

		admin := v1.Group("/admin")
		admin.Use(authVerify)
		{
			admin.GET("/", h.Profile)
		}
	}

	return router
}
