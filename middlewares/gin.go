package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	tokenProto "originals/srv/token/proto"

	"github.com/gin-gonic/gin"
	"github.com/go-log/log"
	"github.com/micro/go-grpc/client"
)

const (
	// 登陆配置
	loginInfoContextKey     = "login_info"
	loginTokenHeaderKey     = "x-login-token"
	loginTokenSecretKey     = "login_token_secret_key"
	loginTokenRefreshPeriod = 30

	// 注册配置
	registerInfoContextKey = "register_info"
	registerTokenHeaderKey = "x-Register-token"
	registerTokenSecretKey = "register_token_secret_key"
)

var TokenStatusToMsg = map[tokenProto.TokenStatus]string{
	tokenProto.TokenStatus_CANCELED: "token canceled",
	tokenProto.TokenStatus_INVALID:  "token invalid",
	tokenProto.TokenStatus_EXPIRED:  "token expired",
}

// 错误日志格式
func logErr(srv, mtd string, err error) {
	log.Logf("service: %s | method: %-24s | error: %v\n", srv, mtd, err)
}

// Cors gin跨域中间件
func Cors() gin.HandlerFunc {
	allowHeaders := []string{
		"Content-Type",
		loginTokenHeaderKey,
		registerTokenHeaderKey,
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

// RegisterAuth 注册完成权限验证中间件
func RegisterAuth() gin.HandlerFunc {
	tokenSrv := tokenProto.NewTokenService("go.micro.srv.token", &LogCli{client.DefaultClient})
	return func(ctx *gin.Context) {
		// 从请求头部提取注册token
		registerToken := ctx.GetHeader(registerTokenHeaderKey)
		if registerToken == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 验证注册token
		tokenVerifyRsp, err := tokenSrv.VerifyRegisterToken(context.TODO(), &tokenProto.VerifyRegisterTokenReq{
			Token:     registerToken,
			SecretKey: registerTokenSecretKey,
		})
		if err != nil {
			logErr("srv.token", "VerifyRegisterToken", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 注册token验证不通过
		if tokenVerifyRsp.TokenStatus != tokenProto.TokenStatus_OK {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": TokenStatusToMsg[tokenVerifyRsp.TokenStatus]})
			ctx.Abort()
			return
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
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 验证登陆token
		tokenVerifyRsp, err := tokenSrv.VerifyLoginToken(context.TODO(), &tokenProto.VerifyLoginTokenReq{
			Token:     loginToken,
			SecretKey: loginTokenSecretKey,
		})
		if err != nil {
			logErr("srv.token", "VerifyLoginToken", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 登陆token验证不通过
		if tokenVerifyRsp.TokenStatus != tokenProto.TokenStatus_OK {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": TokenStatusToMsg[tokenVerifyRsp.TokenStatus]})
			ctx.Abort()
			return
		}

		// 登陆token临期刷新
		if tokenVerifyRsp.Claims.ExpiresAt-time.Now().Unix() < loginTokenRefreshPeriod {
			refreshLoginTokenRsp, err := tokenSrv.RefreshLoginToken(context.TODO(), &tokenProto.RefreshLoginTokenReq{
				Token:     loginToken,
				SecretKey: loginTokenSecretKey,
			})
			if err != nil {
				logErr("srv.token", "RefreshLoginToken", err)
				ctx.AbortWithStatus(http.StatusInternalServerError)
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
