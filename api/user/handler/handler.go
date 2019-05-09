package handler

import (
	userSrvProto "originals/srv/user/proto"
	"regexp"

	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserSrv userSrvProto.UserService
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

func (u *User) Index(ctx *gin.Context) {
	ctx.JSON(200, Rsp{
		Code:    200,
		Message: "Hello, this is user api.",
	})
}

// Invite send an invite email to user
func (u *User) Invite(ctx *gin.Context) {
	email := ctx.Query("email")
	rsp := Rsp{}
	if !regVerify(email, emailReg) {
		rsp.Code = 301
		rsp.Message = "invalid param: email"
		ctx.JSON(200, rsp)
		return
	}
	userSrvRsp, err := u.UserSrv.Invite(context.TODO(), &userSrvProto.InviteReq{Email: email})
	if err != nil {
		rsp.Code = 500
		rsp.Message = "Internal error: " + err.Error()
		ctx.JSON(200, rsp)
		return
	}
	switch userSrvRsp.Status {
	case userSrvProto.Status_OK:
		rsp.Code = 200
		rsp.Message = "success"
		rsp.Result = map[string]interface{}{
			"invite_token": userSrvRsp.InviteToken,
		}
	case userSrvProto.Status_EmailRegistered:
		rsp.Code = 302
		rsp.Message = "email has been registered"
	case userSrvProto.Status_EmailSendFailed:
		rsp.Code = 303
		rsp.Message = "email send failed"
	default:
		rsp.Code = 500
		rsp.Message = "internal error"
	}
	ctx.JSON(200, rsp)
	return
}

// Register register a user
func (u *User) Register(ctx *gin.Context) {
	var (
		rsp         Rsp
		ok          bool
		inviteToken string
		password    string
		mobile      string
		nickname    string
		imageUrl    string
	)
	if inviteToken, ok = ctx.GetPostForm("param 'invite_token' undefined"); !ok {
		rsp.Code = 301
		rsp.Message = "invite_token is empty"
		ctx.JSON(200, rsp)
		return
	}
	if password, ok = ctx.GetPostForm("password"); !ok {
		rsp.Code = 301
		rsp.Message = "param 'password' undefined"
		ctx.JSON(200, rsp)
		return
	}
	if !regVerify(password, passwordReg) {
		rsp.Code = 301
		rsp.Message = "param 'password' invalid"
		ctx.JSON(200, rsp)
		return
	}
	if nickname, ok = ctx.GetPostForm("nickname"); !ok {
		rsp.Code = 301
		rsp.Message = "param 'nickname' undefined"
		ctx.JSON(200, rsp)
		return
	}
	if !regVerify(nickname, nickNameReg) {
		rsp.Code = 301
		rsp.Message = "param 'nickname' invalid"
		ctx.JSON(200, rsp)
		return
	}
	mobile, _ = ctx.GetPostForm("mobile")
	imageUrl, _ = ctx.GetPostForm("image_url")

}

func regVerify(str, pattern string) bool {
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(str)
}
