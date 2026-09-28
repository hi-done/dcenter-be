// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IUser interface {
		GetInviteCode(ctx context.Context) (string, error)
		// SignIn 用户登录：校验用户名密码与账号状态，签发 access/refresh 双 token，返回新 token 对
		SignIn(ctx context.Context, userName string, password string) (string, string, error)
		// SignOut 用户注销：删除 Redis 中的 refresh token，使当前会话失效
		SignOut(ctx context.Context) error
		// SignUp 用户注册：校验邀请码、加密密码并创建账号，返回新用户ID
		SignUp(ctx context.Context, userName string, password string, inviteCode string) (uint64, error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
