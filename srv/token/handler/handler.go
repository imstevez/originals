package handler

import (
	"context"
	"originals/srv/token/model"
	proto "originals/srv/token/proto"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	freshTokenMu sync.Mutex
	Model        *model.TokenModel
}

const (
	inviteTokenKey = "invite_token_key"
	authTokenKey   = "auth_token_key"
)

var (
	inviteTokenLive    = 10 * time.Minute
	authTokenLive      = 5 * time.Minute
	authTokenFreshLive = 5 * time.Minute
)

type inviteTokenClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type authTokenClaims struct {
	UserId   int64  `json:"user_id"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	NickName string `json:"nick_name"`
	ImageUrl string `json:"image_url"`
	jwt.StandardClaims
}

// GetInviteToken
func (t *Token) GetInviteToken(ctx context.Context, req *proto.GetInviteTokenReq, rsp *proto.GetInviteTokenRsp) error {
	if req.Claims == nil {
		rsp.Status = proto.Status_ParamInvalid
		return nil
	}
	claims := inviteTokenClaims{
		Email: req.Claims.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(inviteTokenLive).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(inviteTokenKey))
	if err != nil {
		return err
	}
	rsp.Status = proto.Status_OK
	rsp.Token = tokenStr
	return nil
}

// VerifyInviteToken
func (t *Token) VerifyInviteToken(ctx context.Context, req *proto.VerifyInviteTokenReq, rsp *proto.VerifyInviteTokenRsp) error {
	claims := inviteTokenClaims{}
	token, err := jwt.ParseWithClaims(req.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(inviteTokenKey), nil
	})
	if token != nil {
		rsp.Claims = &proto.InviteClaims{
			Email: claims.Email,
		}
	}
	if err != nil {
		if vErr := err.(*jwt.ValidationError); vErr.Errors != jwt.ValidationErrorExpired {
			rsp.Status = proto.Status_TokenInvalid
			return nil
		}
		rsp.Status = proto.Status_TokenExpired
		return nil
	}
	rsp.Status = proto.Status_OK
	return nil
}

// GetAuthToken
func (t *Token) GetAuthToken(ctx context.Context, req *proto.GetAuthTokenReq, rsp *proto.GetAuthTokenRsp) error {
	if req.Claims == nil {
		rsp.Status = proto.Status_ParamInvalid
		return nil
	}
	claims := authTokenClaims{
		UserId:   req.Claims.UserId,
		Email:    req.Claims.Email,
		Mobile:   req.Claims.Mobile,
		NickName: req.Claims.Nickname,
		ImageUrl: req.Claims.ImageUrl,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(authTokenLive).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(authTokenKey))
	if err != nil {
		return err
	}
	rsp.Status = proto.Status_OK
	rsp.Token = tokenStr
	return nil
}

// VerifyAuthToken
func (t *Token) VerifyAuthToken(ctx context.Context, req *proto.VerifyAuthTokenReq, rsp *proto.VerifyAuthTokenRsp) error {
	if canceled, err := t.Model.IsTokenCanceled(req.Token); err != nil {
		return err
	} else if canceled {
		rsp.Status = proto.Status_TokenCanceled
		return nil
	}

	claims := authTokenClaims{}
	token, err := jwt.ParseWithClaims(req.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(authTokenKey), nil
	})
	if token != nil {
		rsp.Claims = &proto.AuthClaims{
			UserId:   claims.UserId,
			Email:    claims.Email,
			Mobile:   claims.Mobile,
			Nickname: claims.NickName,
			ImageUrl: claims.ImageUrl,
		}
	}

	if err == nil {
		rsp.Status = proto.Status_OK
		return nil
	}

	if vErr := err.(*jwt.ValidationError); vErr.Errors != jwt.ValidationErrorExpired {
		rsp.Status = proto.Status_TokenInvalid
		return nil
	}

	freshDeadLine := time.Unix(claims.ExpiresAt, 0).Add(authTokenFreshLive)
	if time.Now().After(freshDeadLine) {
		rsp.Status = proto.Status_TokenExpired
		return nil
	}

	t.freshTokenMu.Lock()
	defer t.freshTokenMu.Unlock()
	if freshToken, err := t.Model.GetFreshToken(req.Token); err == nil {
		rsp.Status = proto.Status_TokenRefreshed
		rsp.FreshToken = freshToken
		return nil
	} else if err == model.ErrKeyNotExist {
		claims.IssuedAt = time.Now().Unix()
		claims.ExpiresAt = time.Now().Add(authTokenLive).Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		freshToken, err := token.SignedString([]byte(authTokenKey))
		if err != nil {
			return err
		}

		err = t.Model.SetFreshToken(req.Token, freshToken, authTokenFreshLive)
		if err != nil {
			return err
		}

		rsp.Status = proto.Status_TokenRefreshed
		rsp.FreshToken = freshToken
		return nil
	} else {
		return err
	}
}

// CancelAuthToken
func (t *Token) CancelToken(ctx context.Context, req *proto.CancelTokenReq, rsp *proto.CancelTokenRsp) error {
	claims := authTokenClaims{}
	token, _ := jwt.ParseWithClaims(req.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(authTokenKey), nil
	})
	if token == nil || !token.Valid {
		rsp.Status = proto.Status_TokenInvalid
		return nil
	}

	err := t.Model.CancelToken(req.Token, time.Unix(claims.ExpiresAt, 0).Add(authTokenFreshLive))
	if err != nil {
		return err
	}
	rsp.Status = proto.Status_OK
	return nil
}
