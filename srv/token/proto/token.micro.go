// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: token.proto

package token

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Token service

type TokenService interface {
	GetInviteToken(ctx context.Context, in *GetInviteTokenReq, opts ...client.CallOption) (*GetInviteTokenRsp, error)
	VerifyInviteToken(ctx context.Context, in *VerifyInviteTokenReq, opts ...client.CallOption) (*VerifyInviteTokenRsp, error)
	GetAuthToken(ctx context.Context, in *GetAuthTokenReq, opts ...client.CallOption) (*GetAuthTokenRsp, error)
	VerifyAuthToken(ctx context.Context, in *VerifyAuthTokenReq, opts ...client.CallOption) (*VerifyAuthTokenRsp, error)
	CancelToken(ctx context.Context, in *CancelTokenReq, opts ...client.CallOption) (*CancelTokenRsp, error)
}

type tokenService struct {
	c    client.Client
	name string
}

func NewTokenService(name string, c client.Client) TokenService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "token"
	}
	return &tokenService{
		c:    c,
		name: name,
	}
}

func (c *tokenService) GetInviteToken(ctx context.Context, in *GetInviteTokenReq, opts ...client.CallOption) (*GetInviteTokenRsp, error) {
	req := c.c.NewRequest(c.name, "Token.GetInviteToken", in)
	out := new(GetInviteTokenRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenService) VerifyInviteToken(ctx context.Context, in *VerifyInviteTokenReq, opts ...client.CallOption) (*VerifyInviteTokenRsp, error) {
	req := c.c.NewRequest(c.name, "Token.VerifyInviteToken", in)
	out := new(VerifyInviteTokenRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenService) GetAuthToken(ctx context.Context, in *GetAuthTokenReq, opts ...client.CallOption) (*GetAuthTokenRsp, error) {
	req := c.c.NewRequest(c.name, "Token.GetAuthToken", in)
	out := new(GetAuthTokenRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenService) VerifyAuthToken(ctx context.Context, in *VerifyAuthTokenReq, opts ...client.CallOption) (*VerifyAuthTokenRsp, error) {
	req := c.c.NewRequest(c.name, "Token.VerifyAuthToken", in)
	out := new(VerifyAuthTokenRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenService) CancelToken(ctx context.Context, in *CancelTokenReq, opts ...client.CallOption) (*CancelTokenRsp, error) {
	req := c.c.NewRequest(c.name, "Token.CancelToken", in)
	out := new(CancelTokenRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Token service

type TokenHandler interface {
	GetInviteToken(context.Context, *GetInviteTokenReq, *GetInviteTokenRsp) error
	VerifyInviteToken(context.Context, *VerifyInviteTokenReq, *VerifyInviteTokenRsp) error
	GetAuthToken(context.Context, *GetAuthTokenReq, *GetAuthTokenRsp) error
	VerifyAuthToken(context.Context, *VerifyAuthTokenReq, *VerifyAuthTokenRsp) error
	CancelToken(context.Context, *CancelTokenReq, *CancelTokenRsp) error
}

func RegisterTokenHandler(s server.Server, hdlr TokenHandler, opts ...server.HandlerOption) error {
	type token interface {
		GetInviteToken(ctx context.Context, in *GetInviteTokenReq, out *GetInviteTokenRsp) error
		VerifyInviteToken(ctx context.Context, in *VerifyInviteTokenReq, out *VerifyInviteTokenRsp) error
		GetAuthToken(ctx context.Context, in *GetAuthTokenReq, out *GetAuthTokenRsp) error
		VerifyAuthToken(ctx context.Context, in *VerifyAuthTokenReq, out *VerifyAuthTokenRsp) error
		CancelToken(ctx context.Context, in *CancelTokenReq, out *CancelTokenRsp) error
	}
	type Token struct {
		token
	}
	h := &tokenHandler{hdlr}
	return s.Handle(s.NewHandler(&Token{h}, opts...))
}

type tokenHandler struct {
	TokenHandler
}

func (h *tokenHandler) GetInviteToken(ctx context.Context, in *GetInviteTokenReq, out *GetInviteTokenRsp) error {
	return h.TokenHandler.GetInviteToken(ctx, in, out)
}

func (h *tokenHandler) VerifyInviteToken(ctx context.Context, in *VerifyInviteTokenReq, out *VerifyInviteTokenRsp) error {
	return h.TokenHandler.VerifyInviteToken(ctx, in, out)
}

func (h *tokenHandler) GetAuthToken(ctx context.Context, in *GetAuthTokenReq, out *GetAuthTokenRsp) error {
	return h.TokenHandler.GetAuthToken(ctx, in, out)
}

func (h *tokenHandler) VerifyAuthToken(ctx context.Context, in *VerifyAuthTokenReq, out *VerifyAuthTokenRsp) error {
	return h.TokenHandler.VerifyAuthToken(ctx, in, out)
}

func (h *tokenHandler) CancelToken(ctx context.Context, in *CancelTokenReq, out *CancelTokenRsp) error {
	return h.TokenHandler.CancelToken(ctx, in, out)
}
