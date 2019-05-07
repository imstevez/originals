package handler

import (
	"context"
	"errors"
	"fmt"
	"originals/conf"
	"originals/email"
	"github.com/dgrijalva/jwt-go"
	"originals/srv/token/model"
	"originals/srv/token/proto"
	"originals/utils"
	"sync"
	"time"

	"github.com/micro/go-log"
)

type TokenHandler struct {
	freshTokenMu sync.Mutex
	Model        *model.TokenModel
}


const (
	inviteTokenKey = "invite_token_key"
	authTokenKey = "auth_token_key"
	inviteTokenExp = 30
	authTokenExp = 5
	authTokenRefreshLimit = 5
	authTokenRefreshLive = 1
)

type inviteClaims struct {
	Email string
	jwt.StandardClaims
}

// GetInviteToken
func (hdlr *TokenHandler) GetInviteToken(ctx context.Context, req *proto.GetInviteTokenReq, rsp *proto.GetInviteTokenRsp) error {
	var claims = inviteClaims{
		Email: req.Claims.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(inviteTokenKey))
	if err != nil {
		return err
	}
	rsp.Token = tokenStr
	return nil
}

// VerifyInvite
func (hdlr *TokenHandler) VerifyInvite(ctx context.Context, req *proto.VerifyInviteReq, rsp *proto.VerifyInviteRsp) error {
	claims := &inviteClaims{}
	token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(inviteTokenKey), nil
	})
	if err != nil {
		valiErr := err.(*jwt.ValidationError)
		if valiErr.Errors == jwt.ValidationErrorExpired {
			rsp.Status = proto.Status_ErrTokenExpired
		} else {
			rsp.Status = proto.Status_ErrTokenInvalid
		}
		return nil
	}
	if token == nil || !token.Valid {
		rsp.Status = proto.Status_ErrTokenInvalid
		return nil
	}
	rsp.Claims.Email = claims.Email
	return nil
}

// ParseInviteToken
func (hdlr *UserSrvHandler) ParseInviteToken(ctx context.Context, req *proto.ParseInviteTokenReq, rsp *proto.ParseInviteTokenRsp) error {
	secret := conf.SrvConf["user"].Extra["jwt_secret_invite"]
	if secret == "" {
		return ErrEmptySecret
	}
	claims := &jwt.EmailClaims{}

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
	claims.ExpiresAt = time.Now().Add(5 * time.Minute).Unix()
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
			claims.ExpiresAt = time.Now().Add(5 * time.Minute).Unix()
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
