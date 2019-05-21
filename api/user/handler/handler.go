package handler

import (
	"fmt"
	"net/http"
	"originals/email"
	tokenProto "originals/srv/token/proto"
	userProto "originals/srv/user/proto"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-log"
	"golang.org/x/net/context"
)

const (
	registerTokenLive      = 30
	registerTokenSecretKey = "register_token_secret_key"
	registerInfoContextKey = "register_info"
	registerEmailBody      = `
<!doctype html>
<html>
	<body>
		<header><h3>Originals Beta v1.0<h3><hr></header>
		<article>
			<p>Hi, 你好</p>
			<p>欢迎注册<b>云记</b>, 在使用应用前，您还需要您点击本邮件中的链接完成密码等一些必要的账户设置:<br>
			<a href="http://www.koogo.net/complete/%s"><i>http://www.koogo.net/complete/%s</i></a></p>
			<p>如果以上链接无法打开, 您可以直接复制链接地址到浏览器中打开, 链接30分钟内有效.</p>
			<p><b>谢谢使用</b></p>
		</article>
		<footer>
			<p>-------------------------------<br>云记 团队</p>
		<footer>
	</body>
</html>
`

	loginTokenLive      = 30
	loginTokenSecretKey = "login_token_secret_key"
	loginTokenHeaderKey = "x-login-token"
	loginInfoContextKey = "login_info"
	avatarUri           = "http://www.koogo.net:8080/user/statics/avatar/"

	// Response code
	codeSuccess         = 200
	codeParamErr        = 300
	codeEmailRegistered = 301
	codeEmailSendFailed = 302
	codeUserNotExist    = 303
	codePasswordError   = 304
)

func logErr(srv, mtd string, err error) {
	log.Logf("service: %s | method: %-24s | error: %v\n", srv, mtd, err)
}

type User struct {
	TokenSrv tokenProto.TokenService
	UserSrv  userProto.UserService
}

type RegisterReq struct {
	Email string `json:"email" binding:"email,required"`
}

// Register 验证用户邮箱, 并发送注册邮件
func (u *User) Register(ctx *gin.Context) {
	// 绑定请求参数
	var req RegisterReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": codeParamErr,
			"message": err.Error(),
		})
		return
	}

	// 验证邮箱是否已注册
	isEmailRegisteredRsp, err := u.UserSrv.IsEmailRegistered(context.TODO(), &userProto.IsEmailRegisteredReq{
		Email: req.Email,
	})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if isEmailRegisteredRsp.Registered {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    codeEmailRegistered,
			"message": "email has been registered",
		})
		return
	}

	// 获取注册token
	getRegisterTokenRsp, err := u.TokenSrv.GetRegisterToken(ctx, &tokenProto.GetRegisterTokenReq{
		Claims: &tokenProto.RegisterTokenClaims{
			Email:     req.Email,
			ExpiresAt: time.Now().Add(registerTokenLive * time.Minute).Unix(),
		},
		SecretKey: registerTokenSecretKey,
	})
	if err != nil {
		logErr("srv.token", "GetRegisterToken", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// 发送注册邮件
	mailBody := fmt.Sprintf(registerEmailBody, getRegisterTokenRsp.Token, getRegisterTokenRsp.Token)
	registerMail := &email.Email{
		Receivers: []string{req.Email},
		Subject:   "云记-Beta v1.0 注册测试邮件",
		Body:      mailBody,
	}
	if err := email.SendMail(registerMail); err != nil {
		logErr("srv.email", "SendMail", err)
		ctx.JSON(http.StatusOK, gin.H{
			"code":    codeEmailSendFailed,
			"message": "register email send failed",
		})
		return
	}

	// 注册成功
	ctx.JSON(200, gin.H{
		"code":    codeSuccess,
		"message": "register success",
	})
	return
}

// Complete 创建用户, 完成用户注册
//
// TODO: 表单信息较验
//
func (u *User) Complete(ctx *gin.Context) {
	// 提取用户注册权限验证上下文
	registerInfo := ctx.MustGet(registerInfoContextKey).(map[string]interface{})
	regEmail := registerInfo["email"].(string)
	avatar := ""

	// 解析请求参数
	password := ctx.PostForm("password")
	nickname := ctx.PostForm("nickname")
	avatarFile, err := ctx.FormFile("avatar")
	if err == nil {
		// 保存用户头像图片
		avatar = avatarUri + "default.png"
		if avatarFile != nil {
			fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), avatarFile.Filename)
			dst := fmt.Sprintf("./statics/avatar/%s", fileName)
			if err := ctx.SaveUploadedFile(avatarFile, dst); err != nil {
				logErr("srv.file", "SaveUploadedFile", err)
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			avatar = avatarUri + fileName
		}
	}

	// 创建用户
	createNewUserRsp, err := u.UserSrv.CreateNewUser(context.TODO(), &userProto.CreateNewUserReq{
		Email:    regEmail,
		Password: password,
		Nickname: nickname,
		Avatar:   avatar,
	})
	if err != nil {
		logErr("srv.user", "CreateNewUser", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 用户已完成注册
	if createNewUserRsp.UserId == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    codeEmailRegistered,
			"message": "email register has been completed",
		})
		return
	}

	// 完成注册成功
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "complete success",
		"user_id": createNewUserRsp.UserId,
	})
	return
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登陆接口
func (u *User) Login(ctx *gin.Context) {
	// 绑定请求参数
	var req loginReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    codeParamErr,
			"message": err.Error(),
		})
		return
	}

	// 用户验证
	verifyUserRsp, err := u.UserSrv.VerifyUser(context.TODO(), &userProto.VerifyUserReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logErr("srv.user", "VerifyUser", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 登陆信息验证失败
	if verifyUserRsp.VerifyStatus != userProto.UserVerifyStatus_OK {
		userVerifyStatusToRsp := map[userProto.UserVerifyStatus]gin.H{
			userProto.UserVerifyStatus_NOT_EXIST: gin.H{
				"code":    codeUserNotExist,
				"message": "user is not exist",
			},
			userProto.UserVerifyStatus_PWD_ERROR: gin.H{
				"code":    codePasswordError,
				"message": "password is wrong",
			},
		}
		ctx.JSON(http.StatusOK, userVerifyStatusToRsp[verifyUserRsp.VerifyStatus])
		return
	}

	// 获取登陆token
	getLoginTokenRsp, err := u.TokenSrv.GetLoginToken(context.TODO(), &tokenProto.GetLoginTokenReq{
		Claims: &tokenProto.LoginTokenClaims{
			UserId:    verifyUserRsp.UserInfo.UserId,
			Email:     verifyUserRsp.UserInfo.Email,
			Nickname:  verifyUserRsp.UserInfo.Nickname,
			Avatar:    verifyUserRsp.UserInfo.Avatar,
			ExpiresAt: time.Now().Add(loginTokenLive * time.Minute).Unix(),
		},
		SecretKey: loginTokenSecretKey,
	})
	if err != nil {
		logErr("srv.token", "GetLoginToken", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 登陆成功
	ctx.Header(loginTokenHeaderKey, getLoginTokenRsp.Token)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    codeSuccess,
		"message": "login success",
	})
	return
}

// Logout 退出登陆
func (u *User) Logout(ctx *gin.Context) {
	// 获取授权token
	token := ctx.GetHeader(loginTokenHeaderKey)
	if token == "" {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	// 注销token
	if _, err := u.TokenSrv.CancelToken(context.TODO(), &tokenProto.CancelTokenReq{
		Token: token,
	}); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 退出登陆成功
	ctx.JSON(200, gin.H{
		"code":    codeSuccess,
		"message": "logout success",
	})
	return
}

// Profile 获取用户基础信息数据
func (u *User) Profile(ctx *gin.Context) {
	userInfo := ctx.MustGet(loginInfoContextKey).(map[string]interface{})
	ctx.JSON(200, gin.H{
		"code":    codeSuccess,
		"message": "success",
		"data":    userInfo,
	})
}
