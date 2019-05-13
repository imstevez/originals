package handler

import (
	"fmt"
	"originals/email"
	tokenSrvProto "originals/srv/token/proto"
	userSrvProto "originals/srv/user/proto"
	"regexp"
	"time"

	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserSrv  userSrvProto.UserService
	TokenSrv tokenSrvProto.TokenService
}

type Rsp struct {
	Code    int64                  `json:"code"`
	Message string                 `json:"message"`
	Result  map[string]interface{} `json:"result,omitempty"`
}

const (
	passwordReg = `^[a-zA-Z0-9]{4,16}$`
	emailReg    = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	nickNameReg = `[a-zA-Z0-9]{4,16}`
)

type RegisterReq struct {
	Email string `json:"email"`
}

const registerEmailBody = `
<!doctype html>
<html>
	<body>
		<header><h3>Originals Beta v1.0<h3><hr></header>
		<article>
			<p>Hi there,</p>
			<p>Before use the <b>originals</b>, please take a few minutes to complete your account. This link will take you to the page:<br>
			<a href="http://localhost:3000/complete/%s"><i>Account Setting</i></a></p>
			<p>If the link above doesn't work, please copy this link to your browser: <b>http://localhost:3000/complete/%s</b></p>
			<p><b>Thanks</b></p>
		</article>
		<footer>
			<p>-------------------------------<br>Originals-Team</p>
		<footer>
	</body>
</html>
`

// Register
func (u *User) Register(ctx *gin.Context) {
	var (
		req RegisterReq
		rsp Rsp
	)
	if err := ctx.Bind(&req); err != nil {
		rsp.Code = 301
		rsp.Message = "param error: " + err.Error()
		ctx.JSON(200, rsp)
		return
	}
	if !regVerify(req.Email, emailReg) {
		rsp.Code = 302
		rsp.Message = "email invalid"
		ctx.JSON(200, rsp)
		return
	}
	isEmailRegisteredRsp, err := u.UserSrv.IsEmailRegistered(context.TODO(), &userSrvProto.IsEmailRegisteredReq{
		Email: req.Email,
	})
	if err != nil {
		rsp.Code = 500
		rsp.Message = "Internal error: " + err.Error()
		ctx.JSON(200, rsp)
		return
	}
	if isEmailRegisteredRsp.Registered {
		rsp.Code = 303
		rsp.Message = "email has been registered"
		ctx.JSON(200, rsp)
	}

	getInviteTokenRsp, err := u.TokenSrv.GetInviteToken(ctx, &tokenSrvProto.GetInviteTokenReq{
		Claims: &tokenSrvProto.InviteClaims{
			Email: req.Email,
		},
	})
	if err != nil {
		rsp.Code = 500
		rsp.Message = "Internal error: " + err.Error()
		ctx.JSON(200, rsp)
		return
	}

	mailBody := fmt.Sprintf(registerEmailBody, getInviteTokenRsp.Token, getInviteTokenRsp.Token)
	registerMail := &email.Email{
		Recivers: []string{req.Email},
		Subject:  "Originals-起源-Beta v1.0 注册测试邮件",
		Body:     mailBody,
	}
	if err := email.SendMail(registerMail); err != nil {
		rsp.Code = 304
		rsp.Message = "email send failed"
		ctx.JSON(200, rsp)
		return
	}

	rsp.Code = 200
	rsp.Message = "success"
	rsp.Result = map[string]interface{}{
		"invite_token": getInviteTokenRsp.Token,
	}

	ctx.JSON(200, rsp)
	return
}

// Complete
func (u *User) Complete(ctx *gin.Context) {
	var (
		createUserReq userSrvProto.CreateUserReq
		ok            bool
		rsp           Rsp
	)
	token := ctx.GetHeader("x-originals-token")
	if token == "" {
		rsp.Code = 301
		rsp.Message = "token is empty"
		ctx.JSON(200, rsp)
		return
	}
	verifyInviteTokenRsp, err := u.TokenSrv.VerifyInviteToken(context.TODO(), &tokenSrvProto.VerifyInviteTokenReq{
		Token: token,
	})
	if err != nil {
		rsp.Code = 500
		rsp.Message = "internal error: " + err.Error()
		ctx.JSON(200, rsp)
		return
	}
	createUserReq.Email = verifyInviteTokenRsp.Claims.Email

	isEmailRegisteredRsp, err := u.UserSrv.IsEmailRegistered(context.TODO(), &userSrvProto.IsEmailRegisteredReq{
		Email: createUserReq.Email,
	})
	if err != nil {
		rsp.Code = 500
		rsp.Message = "Internal error: " + err.Error()
		ctx.JSON(200, rsp)
		return
	}
	if isEmailRegisteredRsp.Registered {
		rsp.Code = 302
		rsp.Message = "email has been registered"
		ctx.JSON(200, rsp)
	}

	avatarFile, err := ctx.FormFile("avatar")
	if err != nil {
		rsp.Code = 303
		rsp.Message = "avatar image error"
		ctx.JSON(200, rsp)
	}
	if avatarFile != nil {
		dst := fmt.Sprintf("./file/avatar/%s_%s", time.Now().Unix(), avatarFile.Filename)
		if err := ctx.SaveUploadedFile(avatarFile, dst); err != nil {
			rsp.Code = 500
			rsp.Message = "internal error: " + err.Error()
			ctx.JSON(200, rsp)
			return
		}
	}

	if inviteToken, ok = ctx.GetPostForm("token"); !ok {
		rsp.Code = 301
		rsp.Message = "invite_token empty"
		ctx.JSON(200, rsp)
		return
	}

	if registerReq.Password, ok = ctx.GetPostForm("password"); !ok {
		rsp.Code = 301
		rsp.Message = "password undefined"
		ctx.JSON(200, rsp)
		return
	}
	if !regVerify(registerReq.Password, passwordReg) {
		rsp.Code = 301
		rsp.Message = "password invalid"
		ctx.JSON(200, rsp)
		return
	}
	if registerReq.Nickname, ok = ctx.GetPostForm("nickname"); !ok {
		rsp.Code = 301
		rsp.Message = "nickname undefined"
		ctx.JSON(200, rsp)
		return
	}
	if !regVerify(registerReq.Nickname, nickNameReg) {
		rsp.Code = 301
		rsp.Message = "nickname invalid"
		ctx.JSON(200, rsp)
		return
	}
	registerReq.Mobile, _ = ctx.GetPostForm("mobile")
	registerReq.ImageUrl, _ = ctx.GetPostForm("image_url")
	tokenSrvRsp, err := u.TokenSrv.VerifyInviteToken(context.TODO(), &tokenSrvProto.VerifyInviteTokenReq{
		Token: inviteToken,
	})
	if err != nil {
		rsp.Code = 500
		rsp.Message = "internal error: " + err.Error()
		ctx.JSON(200, rsp)
		return
	}
	switch tokenSrvRsp.Status {
	case tokenSrvProto.Status_OK:
		registerReq.Email = tokenSrvRsp.Claims.Email
	case tokenSrvProto.Status_TokenInvalid:
		rsp.Code = 401
		rsp.Message = "invite_token invalid"
		ctx.JSON(200, rsp)
		return
	case tokenSrvProto.Status_TokenExpired:
		rsp.Code = 402
		rsp.Message = "invite_token expired"
		ctx.JSON(200, rsp)
		return
	default:
		rsp.Code = 500
		rsp.Message = "internal error"
		ctx.JSON(200, rsp)
		return
	}
	userSrvRsp, err := u.UserSrv.Register(context.TODO(), &registerReq)
	if err != nil {
		rsp.Code = 500
		rsp.Message = "internal error: " + err.Error()
		ctx.JSON(200, rsp)
		return
	}
	switch userSrvRsp.Status {
	case userSrvProto.Status_OK:
		rsp.Code = 200
		rsp.Message = "success"
		rsp.Result = map[string]interface{}{
			"user_id": userSrvRsp.UserId,
		}
	case userSrvProto.Status_EmailRegistered:
		rsp.Code = 302
		rsp.Message = "email registered"
	default:
		rsp.Code = 500
		rsp.Message = "internal error"
	}
	ctx.JSON(200, rsp)
	return
}

// Login
func (u *User) Login(ctx *gin.Context) {
	var (
		loginReq userSrvProto.LoginReq
		rsp      Rsp
		ok       bool
	)
	if loginReq.Email, ok = ctx.GetPostForm("email"); !ok {
		rsp.Code = 301
		rsp.Message = "email empty"
		ctx.JSON(200, rsp)
		return
	}
	if loginReq.Password, ok = ctx.GetPostForm("password"); !ok {
		rsp.Code = 301
		rsp.Message = "password empty"
		ctx.JSON(200, rsp)
		return
	}

	loginRsp, err := u.UserSrv.Login(context.TODO(), &loginReq)
	if err != nil {
		rsp.Code = 500
		rsp.Message = "internal error: " + err.Error()
		ctx.JSON(200, rsp)
		return
	}
	switch loginRsp.Status {
	case userSrvProto.Status_OK:
		rsp.Code = 200
		rsp.Message = "success"
		rsp.Result = map[string]interface{}{
			"token": loginRsp.AuthToken,
		}
	case userSrvProto.Status_UserNotExist:
		rsp.Code = 401
		rsp.Message = "email not exist"
	case userSrvProto.Status_PasswordWrong:
		rsp.Code = 402
		rsp.Message = "password wrong"
	default:
		rsp.Code = 500
		rsp.Message = "internal error"
	}
	ctx.JSON(200, rsp)
	return
}

// Logout
func (u *User) Logout(ctx *gin.Context) {
	var (
		logoutReq userSrvProto.LogoutReq
		rsp       Rsp
	)
	logoutReq.AuthToken = ctx.GetHeader("x-originals-token")
	if logoutReq.AuthToken == "" {
		rsp.Code = 301
		rsp.Message = "token empty"
		ctx.JSON(200, rsp)
		return
	}
	logoutRsp, err := u.UserSrv.Logout(context.TODO(), &logoutReq)
	if err != nil {
		rsp.Code = 500
		rsp.Message = "internal error: " + err.Error()
		ctx.JSON(200, rsp)
		return
	}
	switch logoutRsp.Status {
	case userSrvProto.Status_OK:
		rsp.Code = 200
		rsp.Message = "success"
		ctx.Header("x-originals-token", "")
	default:
		rsp.Code = 500
		rsp.Message = "internal error"
	}
	ctx.JSON(200, rsp)
	return
}

// List
func (u *User) Profile(ctx *gin.Context) {
	user := make(map[string]interface{})
	user["user_id"] = ctx.MustGet("user_id").(int64)
	user["email"] = ctx.MustGet("email").(string)
	user["mobile"] = ctx.MustGet("mobile").(string)
	user["nickname"] = ctx.MustGet("nickname").(string)
	user["image_url"] = ctx.MustGet("image_url").(string)
	ctx.JSON(200, gin.H{"code": 200, "result": user})
}

// regVerfy
func regVerify(str, pattern string) bool {
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(str)
}
