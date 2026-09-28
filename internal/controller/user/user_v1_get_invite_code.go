package user

import (
	"context"

	v1 "dcenter/api/user/v1"
	"dcenter/internal/service"
)

func (c *ControllerV1) GetInviteCode(ctx context.Context, req *v1.GetInviteCodeReq) (res *v1.GetInviteCodeRes, err error) {
	code, err := service.User().GetInviteCode(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetInviteCodeRes{InviteCode: code}, nil
}
