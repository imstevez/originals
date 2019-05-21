package main

import (
	"net/http"
	"originals/api/user/handler"
	mid "originals/middlewares"
	tokenProto "originals/srv/token/proto"
	userProto "originals/srv/user/proto"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
)

// 注册路由
func router() *gin.Engine {
	// 注册依赖服务到handler
	c := client.DefaultClient
	tokenSrv := tokenProto.NewTokenService("go.micro.srv.token", &mid.LogCli{c})
	userSrv := userProto.NewUserService("go.micro.srv.user", &mid.LogCli{c})
	h := &handler.User{
		TokenSrv: tokenSrv,
		UserSrv:  userSrv,
	}

	// 创建gin路由, 注册handler
	r := gin.Default()

	r.Use(mid.Cors())
	userApi := r.Group("/user")
	{
		// Index
		userApi.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "Hello, there is originals user api.",
			})
		})

		// 静态服务
		userApi.Static("/statics", "/statics")

		// 用户权限
		authApi := userApi.Group("/auth")
		{
			authApi.POST("/register", h.Register)
			authApi.POST("/login", h.Login)
			authApi.POST("/logout", h.Logout)
			authApi.Use(mid.RegisterAuth())
			{
				authApi.POST("/complete", h.Complete)
			}
		}

		// 用户信息
		infoApi := userApi.Group("/info").Use(mid.LoginAuth())
		{
			infoApi.GET("/", h.Profile)
		}
	}

	return r
}
