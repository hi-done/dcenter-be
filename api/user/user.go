// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	"dcenter/api/user/v1"
)

type IUserV1 interface {
	GetInviteCode(ctx context.Context, req *v1.GetInviteCodeReq) (res *v1.GetInviteCodeRes, err error)
	SignOut(ctx context.Context, req *v1.SignOutReq) (res *v1.SignOutRes, err error)
}
