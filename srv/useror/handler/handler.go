package handler

import (
	"context"
	"errors"
	"fmt"
	"originals/conf"
	"originals/email"
	"originals/jwt"
	"originals/srv/user/model"
	"originals/srv/user/proto"
	"originals/utils"
	"sync"
	"time"

	"github.com/micro/go-log"
)

type UserSrvHandler struct {
	freshTokenMu sync.Mutex
	Model        *model.UserSrvModel
}

const signUpEmailBody = `
<!doctype html>
<html>
	<body>
		<header><h3>Originals Beta v1.0<h3><hr></header>
		<article>
			<p>Hi there,</p>
			<p>Before use the <b>originals</b>, please take a few minutes to set your account. This link will take you to the page:<br>
			<a href="/user/signup?token=%s"><i>Account Setting</i></a></p>
			<p><b>Thanks</b></p>
		</article>
		<footer>
			<p>-------------------------------<br>O-P-T</p>
		<footer>
	</body>
</html>
`

const (
	// Reg patterns
	RegEmail    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	RegPassword = `^[a-zA-Z0-9]{4,16}$`
	RegNickname = `^[a-zA-Z0-9_]{3,16}$`
)

var (
	// Errors
	ErrInvalidEmail    = errors.New("invalid email")
	ErrEmailExist      = errors.New("email has been sign up")
	ErrEmptySecret     = errors.New("jwt secret key is empty")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidNickname = errors.New("invalid nickname")
	ErrUserNotExist    = errors.New("user not exist")
	ErrPasswordWrong   = errors.New("wrong password")
	ErrTokenCanceled   = errors.New("token has been canceled")
)

// InviteUser
func (hdlr *UserSrvHandler) InviteUser(ctx context.Context, req *proto.InviteUserReq, rsp *proto.InviteUserRsp) error {
	log.Log("Received User.StartEmail request")

	// Verify email format
	if match := utils.RegMatch(req.Email, RegEmail); !match {
		return ErrInvalidEmail
	}

	// Check if email has been sign up
	if emailCount, err := hdlr.Model.CountUserEmail(req.Email); err != nil {
		return err
	} else if emailCount > 0 {
		return ErrEmailExist
	}

	// Create sign up token
	secret := conf.SrvConf["user"].Extra["jwt_secret_invite"]
	if secret == "" {
		return ErrEmptySecret
	}
	claims := jwt.EmailClaims{Email: req.Email}
	claims.StandardClaims.ExpiresAt = time.Now().Add(30 * time.Minute).Unix()
	tokenStr, err := jwt.CreateToken(claims, secret)
	if err != nil {
		return err
	}

	// Send sign up email
	mailBody := fmt.Sprintf(signUpEmailBody, tokenStr)
	signUpMail := &email.Email{
		Recivers: []string{req.Email},
		Subject:  "Originals-起源-Beta v1.0 注册测试邮件",
		Body:     mailBody,
	}
	if err := email.SendMail(signUpMail); err != nil {
		return err
	}

	rsp.InviteToken = tokenStr
	return nil
}

// ParseInviteToken
func (hdlr *UserSrvHandler) ParseInviteToken(ctx context.Context, req *proto.ParseInviteTokenReq, rsp *proto.ParseInviteTokenRsp) error {
	secret := conf.SrvConf["user"].Extra["jwt_secret_invite"]
	if secret == "" {
		return ErrEmptySecret
	}
	claims := &jwt.EmailClaims{}
	if err := jwt.ParseToken(req.InviteToken, secret, claims); err != nil {
		return err
	}
	rsp.Email = claims.Email
	return nil
}

// InsertUser
func (hdlr *UserSrvHandler) InsertUser(ctx context.Context, req *proto.InsertUserReq, rsp *proto.InsertUserRsp) error {
	if err := verifyInsertUser(req); err != nil {
		return err
	}
	password, salt := utils.Password(req.Password)
	user := &proto.User{
		Email:        req.Email,
		Password:     password,
		PasswordSalt: salt,
		Mobile:       req.Mobile,
		Nickname:     req.Nickname,
		ImageUrl:     req.ImageUrl,
	}
	id, err := hdlr.Model.InsertUser(user)
	if err != nil {
		return err
	}
	if id == 0 {
		return ErrEmailExist
	}
	rsp.UserId = id
	return nil
}

// GetAuthToken
func (hdlr *UserSrvHandler) GetAuthToken(ctx context.Context, req *proto.GetAuthTokenReq, rsp *proto.GetAuthTokenRsp) error {
	user, err := hdlr.Model.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotExist
	}
	if user.Password != utils.Hash(req.Password, user.PasswordSalt) {
		return ErrPasswordWrong
	}
	secret := conf.SrvConf["user"].Extra["jwt_secret_auth"]
	if secret == "" {
		return ErrEmptySecret
	}
	claims := jwt.UserClaims{
		UserId:   user.Id,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		ImageUrl: user.ImageUrl,
	}
	claims.ExpiresAt = time.Now().Add(5 * time.Second).Unix()
	tokenStr, err := jwt.CreateToken(claims, secret)
	if err != nil {
		return err
	}
	if err := hdlr.Model.UpdateLastLoginDate(user.Id); err != nil {
		return err
	}
	rsp.AuthToken = tokenStr
	return nil
}

// VerityAuth
func (hdlr *UserSrvHandler) VerifyAuthToken(ctx context.Context, req *proto.VerifyAuthTokenReq, rsp *proto.VerifyAuthTokenRsp) error {
	if canceled, err := hdlr.Model.IsTokenCanceled(req.AuthToken); err != nil {
		return err
	} else if canceled {
		return ErrTokenCanceled
	}

	secret := conf.SrvConf["user"].Extra["jwt_secret_auth"]
	if secret == "" {
		return ErrEmptySecret
	}
	var claims = &jwt.UserClaims{}
	if err := jwt.ParseToken(req.AuthToken, secret, claims); err != nil {
		if err != jwt.ErrTokenExpired {
			return err
		}

		refreshDeadLine := time.Unix(claims.ExpiresAt, 0).Add(5 * time.Minute)
		if time.Now().After(refreshDeadLine) {
			return jwt.ErrTokenExpired
		}

		hdlr.freshTokenMu.Lock()
		defer hdlr.freshTokenMu.Unlock()
		if freshToken, err := hdlr.Model.GetFreshToken(req.AuthToken); err == nil {
			rsp.FreshToken = freshToken
		} else {
			if err != model.ErrKeyNotExist {
				return err
			}
			claims.ExpiresAt = time.Now().Add(5 * time.Second).Unix()
			if freshToken, err := jwt.CreateToken(claims, secret); err != nil {
				return err
			} else {
				if err := hdlr.Model.SetFreshToken(req.AuthToken, freshToken); err != nil {
					return err
				}
				rsp.FreshToken = freshToken
			}
		}
	}

	rsp.UserId = claims.UserId
	rsp.Email = claims.Email
	rsp.Mobile = claims.Mobile
	rsp.Nickname = claims.Nickname
	rsp.ImageUrl = claims.ImageUrl

	return nil
}

// CancelAuthToken
func (hdlr *UserSrvHandler) CancelAuthToken(ctx context.Context, req *proto.CancelAuthTokenReq, rsp *proto.CancelAuthTokenRsp) error {
	secret := conf.SrvConf["user"].Extra["jwt_secret_auth"]
	if secret == "" {
		return ErrEmptySecret
	}
	var claims = &jwt.UserClaims{}
	if err := jwt.ParseToken(req.AuthToken, secret, claims); err != nil {
		if err != jwt.ErrTokenInvalid {
			return err
		}
	}
	if err := hdlr.Model.CancelToken(req.AuthToken, time.Unix(claims.ExpiresAt, 0)); err != nil {
		return err
	}
	return nil
}

// verifyInsertUser
func verifyInsertUser(req *proto.InsertUserReq) (err error) {
	switch {
	case !utils.RegMatch(req.Email, RegEmail):
		err = ErrInvalidEmail
	case !utils.RegMatch(req.Password, RegPassword):
		err = ErrInvalidPassword
	case !utils.RegMatch(req.Nickname, RegNickname):
		err = ErrInvalidNickname
	default:
		err = nil
	}
	return
}
