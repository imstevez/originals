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

func Auth() gin.HandlerFunc {
	c := client.DefaultClient
	tokenSrv := tokenSrvProto.NewTokenService("go.micro.srv.token", &logCli{c})
	return func(ctx *gin.Context) {
		authToken := ctx.GetHeader("x-originals-token")
		if authToken == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
			ctx.Abort()
			return
		}
		tokenVerifyReq := tokenSrvProto.VerifyAuthTokenReq{Token: authToken}
		tokenVerifyRsp, err := tokenSrv.VerifyAuthToken(context.TODO(), &tokenVerifyReq)
		if err != nil {
			ctx.JSON(500, gin.H{"code": 500, "message": "internal error"})
			ctx.Abort()
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
		case tokenSrvProto.Status_TokenRefreshed:
			ctx.Set("user_id", tokenVerifyRsp.Claims.UserId)
			ctx.Set("email", tokenVerifyRsp.Claims.Email)
			ctx.Set("mobile", tokenVerifyRsp.Claims.Mobile)
			ctx.Set("nickname", tokenVerifyRsp.Claims.Nickname)
			ctx.Set("image_url", tokenVerifyRsp.Claims.ImageUrl)
			ctx.Next()
			ctx.Header("x-originals-token", tokenVerifyRsp.FreshToken)
		case tokenSrvProto.Status_TokenInvalid:
			ctx.JSON(401, gin.H{"code": 401, "message": "auth token invalid"})
			ctx.Abort()
		case tokenSrvProto.Status_TokenCanceled:
			ctx.JSON(401, gin.H{"code": 402, "message": "auth token canceled"})
			ctx.Abort()
		case tokenSrvProto.Status_TokenExpired:
			ctx.JSON(401, gin.H{"code": 403, "message": "auth token expired"})
			ctx.Abort()
		default:
			ctx.JSON(500, gin.H{"code": 500, "message": "internal error"})
			ctx.Abort()
		}
		return
	}
}

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, x-originals-token")
		if ctx.Request.Method == "OPTIONS" {
			ctx.Status(200)
		}
		ctx.Next()
	}
}

func initRouter() *gin.Engine {
	c := client.DefaultClient
	userSrv := userSrvProto.NewUserService("go.micro.srv.user", &logCli{c})
	tokenSrv := tokenSrvProto.NewTokenService("go.micro.srv.token", &logCli{c})

	h := &handler.User{UserSrv: userSrv, TokenSrv: tokenSrv}
	r := gin.Default()
	r.Use(Cors())

	user := r.Group("/user")
	{
		user.GET("/invite", h.Invite)
		user.POST("/register", h.Register)
		user.POST("/login", h.Login)
		user.GET("/logout", h.Logout)
		profile := user.Group("/profile")
		profile.Use(Auth())
		{
			profile.GET("/", h.Profile)
		}
	}

	return r
}
