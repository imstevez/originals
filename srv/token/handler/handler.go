package handler

import (
	"context"
	"originals/srv/token/model"
	"originals/srv/token/proto"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"shendu.com/encoding/json"
)

type Token struct {
	reFreshTokenMu sync.Mutex
	Model          *model.TokenModel
}

type registerTokenClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// GetRegisterToken 获取生成用户注册token
func (t *Token) GetRegisterToken(ctx context.Context, req *proto.GetRegisterTokenReq, rsp *proto.GetRegisterTokenRsp) (err error) {
	claims := registerTokenClaims{
		Email: req.Claims.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: req.Claims.ExpiresAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rsp.Token, err = token.SignedString([]byte(req.SecretKey))
	return
}

// VerifyRegisterToken 解析，验证用户注册token
func (t *Token) VerifyRegisterToken(ctx context.Context, req *proto.VerifyRegisterTokenReq, rsp *proto.VerifyRegisterTokenRsp) (err error) {
	// Token是否已被取消
	canceled, err := t.Model.IsTokenCanceled(req.Token)
	if err != nil {
		return err
	}
	if canceled {
		rsp.TokenStatus = proto.TokenStatus_CANCELED
		return
	}

	// 解析token
	claims := registerTokenClaims{}
	token, err := jwt.ParseWithClaims(req.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(req.SecretKey), nil
	})
	if token != nil {
		rsp.Claims = &proto.RegisterTokenClaims{
			Email:     claims.Email,
			ExpiresAt: claims.ExpiresAt,
		}
	}

	// 验证token, 过期设置TokenStatus_INVALID状态, 其它设置TokenStatus_INVALID状态
	if err != nil {
		if vErr := err.(*jwt.ValidationError); vErr.Errors != jwt.ValidationErrorExpired {
			rsp.TokenStatus = proto.TokenStatus_INVALID
			return nil
		}
		rsp.TokenStatus = proto.TokenStatus_EXPIRED
		return nil
	}

	// 合法token, 返回OK状态
	rsp.TokenStatus = proto.TokenStatus_OK
	return
}

type loginTokenClaims struct {
	UserId   int64  `json:"user_id"`
	Email    string `json:"email"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
	jwt.StandardClaims
}

// GetLoginToken 获取生成用户登陆token
func (t *Token) GetLoginToken(ctx context.Context, req *proto.GetLoginTokenReq, rsp *proto.GetLoginTokenRsp) (err error) {
	claims := loginTokenClaims{
		UserId:   req.Claims.UserId,
		Email:    req.Claims.Email,
		NickName: req.Claims.Nickname,
		Avatar:   req.Claims.Avatar,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: req.Claims.ExpiresAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rsp.Token, err = token.SignedString([]byte(req.SecretKey))
	return
}

// VerifyLoginToken 解析, 验证用户登陆token
func (t *Token) VerifyLoginToken(ctx context.Context, req *proto.VerifyLoginTokenReq, rsp *proto.VerifyLoginTokenRsp) (err error) {
	// Token是否已被取消
	canceled, err := t.Model.IsTokenCanceled(req.Token)
	if err != nil {
		return
	}
	if canceled {
		rsp.TokenStatus = proto.TokenStatus_CANCELED
		return
	}

	// 解析token
	claims := loginTokenClaims{}
	token, err := jwt.ParseWithClaims(req.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(req.SecretKey), nil
	})
	if token != nil {
		rsp.Claims = &proto.LoginTokenClaims{
			UserId:    claims.UserId,
			Email:     claims.Email,
			Nickname:  claims.NickName,
			Avatar:    claims.Avatar,
			ExpiresAt: claims.ExpiresAt,
		}
	}

	// 验证token, 过期设置TokenStatus_INVALID状态, 其它设置TokenStatus_INVALID状态
	if err != nil {
		if vErr := err.(*jwt.ValidationError); vErr.Errors != jwt.ValidationErrorExpired {
			rsp.TokenStatus = proto.TokenStatus_INVALID
			return nil
		}
		rsp.TokenStatus = proto.TokenStatus_EXPIRED
		return nil
	}

	// 合法token, 返回OK状态
	rsp.TokenStatus = proto.TokenStatus_OK
	return
}

// RefreshLoginToken 刷新用户登陆token
//
// 用于token临期等场景下的token更新, 方法将在原Claims的基础上, 更新
// IssuedAt和ExpiresAt属性, 并生成新的登陆token, 为了防止同一个token
// 并发刷新时生成多个新token, 新生成的token将被置入缓存, 当token已被
// 刷新时, 将直接返回缓存的刷新token
//
// 缓存刷新token将以原token为键新token为值, 缓存到原token过期失效
//
func (t *Token) RefreshLoginToken(ctx context.Context, req *proto.RefreshLoginTokenReq, rsp *proto.RefreshLoginTokenRsp) (err error) {
	// 加锁, 防止并发刷新
	t.reFreshTokenMu.Lock()
	defer t.reFreshTokenMu.Unlock()

	// 查询Token刷新缓存
	rsp.Token, err = t.Model.GetFreshToken(req.Token)
	if err == nil || err != model.ErrKeyNotExist {
		return
	}

	// 解析Payload
	decodeBytes, err := jwt.DecodeSegment(strings.Split(req.Token, ".")[1])
	if err != nil {
		return err
	}
	var claims loginTokenClaims
	err = json.Unmarshal(decodeBytes, &claims)
	if err != nil {
		return
	}

	// 生成新token
	cacheLive := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt += req.AddSeconds
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	rsp.Token, err = token.SignedString([]byte(req.SecretKey))

	// 缓存刷新token
	err = t.Model.SetRefreshToken(req.Token, rsp.Token, cacheLive)
	if err != nil {
		return err
	}
	return
}

// CancelToken 取消token
func (t *Token) CancelToken(ctx context.Context, req *proto.CancelTokenReq, rsp *proto.CancelTokenRsp) (err error) {
	decodeBytes, err := jwt.DecodeSegment(strings.Split(req.Token, ".")[1])
	if err != nil {
		return err
	}
	var claims loginTokenClaims
	err = json.Unmarshal(decodeBytes, &claims)
	if err != nil {
		return
	}

	expiration := time.Until(time.Unix(claims.ExpiresAt, 0))
	err = t.Model.CancelToken(req.Token, expiration)
	if err != nil {
		return err
	}
	return
}
