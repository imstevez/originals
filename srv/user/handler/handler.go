package handler

import (
	"context"
	"originals/srv/user/model"
	"originals/srv/user/proto"
	"originals/utils"
)

type User struct {
	Model *model.UserModel
}

// IsEmailRegistered 验证用户邮箱是否已注册
func (u *User) IsEmailRegistered(ctx context.Context, req *proto.IsEmailRegisteredReq, rsp *proto.IsEmailRegisteredRsp) (err error) {
	rsp.Registered, err = u.Model.IsEmailExist(req.Email)
	return
}

// CreateNewUser 创建新用户
func (u *User) CreateNewUser(ctx context.Context, req *proto.CreateNewUserReq, rsp *proto.CreateNewUserRsp) (err error) {
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

// VerifyUser 验证用户邮箱及密码
//
// 邮箱不存在设置UserVerifyStatus_NOT_EXIST状态, 密码错误设置
// UserVerifyStatus_PWD_ERROR状态, 正常设置UserVerifyStatus_OK,
// 并返回用户基本信息数据
//
func (u *User) VerifyUser(ctx context.Context, req *proto.VerifyUserReq, rsp *proto.VerifyUserRsp) (err error) {
	var sUser *model.SecretUser
	if sUser, err = u.Model.GetUserSecret(req.Email); err != nil {
		return
	}
	if sUser.UserId == 0 {
		rsp.VerifyStatus = proto.UserVerifyStatus_NOT_EXIST
		return
	}
	if sUser.Password != utils.Hash(req.Password, sUser.PasswordSalt) {
		rsp.VerifyStatus = proto.UserVerifyStatus_PWD_ERROR
		return
	}
	rsp.VerifyStatus = proto.UserVerifyStatus_OK
	rsp.UserInfo = &proto.UserInfo{
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
