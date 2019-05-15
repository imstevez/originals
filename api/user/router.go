package main

import (
	"context"
	"fmt"
	"net/http"
	"originals/api/user/handler"
	tokenProto "originals/srv/token/proto"
	userProto "originals/srv/user/proto"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
)

const (
	loginInfoContextKey     = "login_info"
	loginTokenHeaderKey     = "x-login-token"
	loginTokenSecretKey     = "login_token_secret_key"
	loginTokenRefreshPeriod = 30
	registerInfoContextKey  = "register_info"
	registerTokenHeaderKey  = "x-Register-token"
)

// LogCli go-micro服务调用日志中间件
type LogCli struct {
	client.Client
}

// Call go-micro服务调用日志中间件方法
func (lc *LogCli) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	fmt.Printf("[service] service: %s | endpoint: %s\n", req.Service(), req.Endpoint())
	return lc.Client.Call(ctx, req, rsp)
}

var TokenStatusToMsg = map[tokenProto.TokenStatus]string{
	tokenProto.TokenStatus_CANCELED: "register token has been canceled",
	tokenProto.TokenStatus_INVALID:  "register token is invalid",
	tokenProto.TokenStatus_EXPIRED:  "register token is expired",
}

func RegisterAuth() gin.HandlerFunc {
	tokenSrv := tokenProto.NewTokenService("go.micro.srv.token", &LogCli{client.DefaultClient})
	return func(ctx *gin.Context) {
		// 从请求头部提取注册token
		registerToken := ctx.GetHeader(registerTokenHeaderKey)
		if registerToken == "" {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		// 验证注册token
		tokenVerifyRsp, err := tokenSrv.VerifyRegisterToken(context.TODO(), &tokenProto.VerifyRegisterTokenReq{
			Token: registerToken,
		})
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}

		// 注册token验证不通过
		if tokenVerifyRsp.TokenStatus != tokenProto.TokenStatus_OK {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": TokenStatusToMsg[tokenVerifyRsp.TokenStatus]})
			ctx.Abort()
		}

		// 设置注册用户信息数据
		ctx.Set(registerInfoContextKey, map[string]interface{}{
			"email": tokenVerifyRsp.Claims.Email,
		})

		ctx.Next()
		return
	}
}

// LoginAuth gin用户登陆状态验证中间件
func LoginAuth() gin.HandlerFunc {
	tokenSrv := tokenProto.NewTokenService("go.micro.srv.token", &LogCli{client.DefaultClient})
	return func(ctx *gin.Context) {
		// 从请求头部提取登陆token
		loginToken := ctx.GetHeader(loginTokenHeaderKey)
		if loginToken == "" {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		// 验证登陆token
		tokenVerifyReq := tokenProto.VerifyLoginTokenReq{Token: loginToken}
		tokenVerifyRsp, err := tokenSrv.VerifyLoginToken(context.TODO(), &tokenVerifyReq)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}

		// 登陆token验证不通过
		if tokenVerifyRsp.TokenStatus != tokenProto.TokenStatus_OK {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": TokenStatusToMsg[tokenVerifyRsp.TokenStatus]})
			ctx.Abort()
		}

		// 登陆token临期刷新
		if tokenVerifyRsp.Claims.ExpiresAt-time.Now().Unix() < loginTokenRefreshPeriod {
			refreshLoginTokenRsp, err := tokenSrv.RefreshLoginToken(context.TODO(), &tokenProto.RefreshLoginTokenReq{
				Token:     loginToken,
				SecretKey: loginTokenSecretKey,
			})
			if err != nil {
				ctx.Status(http.StatusInternalServerError)
				ctx.Abort()
			} else {
				ctx.Header(loginTokenHeaderKey, refreshLoginTokenRsp.Token)
			}
		}

		// 设置登陆用户信息数据
		ctx.Set(loginInfoContextKey, map[string]interface{}{
			"user_id":  tokenVerifyRsp.Claims.UserId,
			"email":    tokenVerifyRsp.Claims.Email,
			"nickname": tokenVerifyRsp.Claims.Nickname,
			"avatar":   tokenVerifyRsp.Claims.Avatar,
		})

		ctx.Next()
		return
	}
}

// Cors gin跨域中间件
func Cors() gin.HandlerFunc {
	allowHeaders := []string{
		"Content-Type",
		loginTokenHeaderKey,
	}
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Header("Access-Control-Allow-Headers", strings.Join(allowHeaders, ","))
		ctx.Header("Access-Control-Expose-Headers", strings.Join(allowHeaders, ","))
		if ctx.Request.Method == "OPTIONS" {
			ctx.Status(200)
		}
		ctx.Next()
	}
}

// 注册路由
func router() *gin.Engine {
	// 注册依赖服务到handler
	c := client.DefaultClient
	tokenSrv := tokenProto.NewTokenService("go.micro.srv.token", &LogCli{c})
	userSrv := userProto.NewUserService("go.micro.srv.user", &LogCli{c})
	h := &handler.User{
		TokenSrv: tokenSrv,
		UserSrv:  userSrv,
	}

	// 创建gin路由, 注册handler
	r := gin.Default()
	r.Use(Cors())
	userApi := r.Group("/user")
	{
		userApi.Static("/statics", "./static")
		// 用户权限
		authApi := userApi.Group("/auth")
		{
			authApi.POST("/register", h.Register)
			authApi.POST("/Complete", h.Complete).Use(RegisterAuth())
			authApi.POST("/login", h.Login)
			authApi.POST("/logout", h.Logout)
		}

		// 用户信息
		infoApi := userApi.Group("/info")
		infoApi.Use(LoginAuth())
		{
			infoApi.GET("/", h.Profile)
		}
	}

	return r
}
