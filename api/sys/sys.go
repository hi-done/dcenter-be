// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"context"

	"dcenter/api/sys/v1"
)

type ISysV1 interface {
	SignIn(ctx context.Context, req *v1.SignInReq) (res *v1.SignInRes, err error)
	SignUp(ctx context.Context, req *v1.SignUpReq) (res *v1.SignUpRes, err error)
}
