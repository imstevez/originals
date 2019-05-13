package handler

import (
	"context"
	"fmt"
	"mycms/utils"
	"originals/email"
	tokenProto "originals/srv/token/proto"
	"originals/srv/user/model"
	proto "originals/srv/user/proto"

	"errors"
)

type User struct {
	Model    *model.UserModel
	TokenCli tokenProto.TokenService
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
			<p>If the link above doesn't work, please copy this link to your browser: http://localhost:3000/complete/%s</p>
			<p><b>Thanks</b></p>
		</article>
		<footer>
			<p>-------------------------------<br>O-P-T</p>
		<footer>
	</body>
</html>
`

// IsEmailRegistered
func (u *User) IsEmailRegistered(ctx context.Context, req *proto.IsEmailRegisteredReq, rsp *proto.IsEmailRegisteredRsp) (err error) {
	rsp.Registered, err = u.Model.IsEmailExist(req.Email)
	return
}

// CreateUser
func (u *User) CreateUser(ctx context.Context, req *proto.CreateUserReq, rsp *proto.CreateUserRsp) (err error) {
	password, salt := utils.Password(req.Password)
	rsp.UserId, err = u.Model.InsertUser(&model.InsertUserObj{
		Email:        req.Email,
		Password:     password,
		PasswordSalt: salt,
		Avatar:       req.Avatar,
		Nickname:     req.Nickname,
	})
	return
}

// VerifyUser
func (u *User) VerifyUser(ctx context.Context, req *proto.VerifyUserReq, rsp *proto.VerifyUserRsp) (err error) {
	var sUser *model.SecretUser
	if sUser, err = u.Model.GetUserSecret(req.Email); err != nil {
		return
	}
	if sUser.UserId == 0 {
		rsp.Result = proto.VerifyUserRsp_NOT_EXIST
		return
	}
	if sUser.Password != utils.Hash(req.Password, sUser.PasswordSalt) {
		rsp.Result = proto.VerifyUserRsp_PWD_WRONG
		return
	}
	rsp.Result = proto.VerifyUserRsp_OK
	rsp.UserInfo = &proto.BaseUserInfo{
		UserId:   sUser.UserId,
		Email:    sUser.Email,
		Nickname: sUser.NickName,
		Avatar:   sUser.Avatar,
	}
	return
}

// UpdateUserLoginDate
func (u *User) UpdateUserLoginDate(ctx context.Context, req *proto.UpdateUserLoginDateReq, rsp *proto.UpdateUserLoginDateRsp) (err error) {
	return u.Model.UpdateLastLoginDate(req.UserId)
}

// Invite
func (u *User) Invite(ctx context.Context, req *proto.InviteReq, rsp *proto.InviteRsp) error {
	if req.Email == "" {
		rsp.Status = proto.Status_ParamInvalid
		return nil
	}
	if exist, err := u.Model.IsEmailExist(req.Email); err != nil {
		return err
	} else if exist {
		rsp.Status = proto.Status_EmailRegistered
		return nil
	}
	tokenReq := &tokenProto.GetInviteTokenReq{
		Claims: &tokenProto.InviteClaims{
			Email: req.Email,
		},
	}
	tokenRsp, err := u.TokenCli.GetInviteToken(ctx, tokenReq)
	if err != nil {
		return err
	}
	if tokenRsp.Status != tokenProto.Status_OK || tokenRsp.Token == "" {
		return errors.New("get invite token failed")
	}
	rsp.InviteToken = tokenRsp.Token

	mailBody := fmt.Sprintf(registerEmailBody, tokenRsp.Token, tokenRsp.Token)
	registerMail := &email.Email{
		Recivers: []string{req.Email},
		Subject:  "Originals-起源-Beta v1.0 注册测试邮件",
		Body:     mailBody,
	}
	if err := email.SendMail(registerMail); err != nil {
		rsp.Status = proto.Status_EmailSendFailed
		return nil
	}
	rsp.Status = proto.Status_OK
	return nil
}

// Register
func (u *User) Register(ctx context.Context, req *proto.RegisterReq, rsp *proto.RegisterRsp) error {
	if req.Email == "" {
		rsp.Status = proto.Status_ParamInvalid
		return nil
	}
	password, salt := utils.Password(req.Password)
	id, err := u.Model.InsertUser(&model.InsertUserObj{
		Email:        req.Email,
		Password:     password,
		PasswordSalt: salt,
		Nickname:     req.Nickname,
		Avatar:       req.ImageUrl,
	})
	if err != nil {
		return err
	}
	if id == 0 {
		rsp.Status = proto.Status_EmailRegistered
		return nil
	}
	rsp.Status = proto.Status_OK
	rsp.UserId = id
	return nil
}

// Login
func (u *User) Login(ctx context.Context, req *proto.LoginReq, rsp *proto.LoginRsp) error {
	if req.Email == "" {
		rsp.Status = proto.Status_ParamInvalid
		return nil
	}
	sUser, err := u.Model.GetUserSecret(req.Email)
	if err != nil {
		return err
	}
	if sUser == nil || sUser.UserId == 0 {
		rsp.Status = proto.Status_UserNotExist
		return nil
	}
	if sUser.Password != utils.Hash(req.Password, sUser.PasswordSalt) {
		rsp.Status = proto.Status_PasswordWrong
		return nil
	}
	tokenReq := &tokenProto.GetAuthTokenReq{
		Claims: &tokenProto.AuthClaims{
			UserId:   sUser.UserId,
			Email:    sUser.Email,
			Mobile:   sUser.Mobile,
			Nickname: sUser.NickName,
			ImageUrl: sUser.ImageUrl,
		},
	}
	tokenRsp, err := u.TokenCli.GetAuthToken(ctx, tokenReq)
	if err != nil {
		return err
	}
	if tokenRsp.Status != tokenProto.Status_OK || tokenRsp.Token == "" {
		return errors.New("get auth token failed")
	}
	if err := u.Model.UpdateLastLoginDate(sUser.UserId); err != nil {
		return err
	}

	rsp.AuthToken = tokenRsp.Token
	rsp.Status = proto.Status_OK
	return nil
}

// Logout
func (u *User) Logout(ctx context.Context, req *proto.LogoutReq, rsp *proto.LogoutRsp) error {
	if req.AuthToken == "" {
		rsp.Status = proto.Status_ParamInvalid
		return nil
	}
	tokenReq := &tokenProto.CancelTokenReq{
		Token: req.AuthToken,
	}
	tokenRsp, err := u.TokenCli.CancelToken(ctx, tokenReq)
	if err != nil {
		return err
	}
	if tokenRsp.Status != tokenProto.Status_OK {
		return errors.New("cancel auth token failed")
	}
	rsp.Status = proto.Status_OK
	return nil
}
