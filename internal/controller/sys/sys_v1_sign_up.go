package sys

import (
	"context"

	v1 "dcenter/api/sys/v1"
	"dcenter/internal/service"
)

func (c *ControllerV1) SignUp(ctx context.Context, req *v1.SignUpReq) (res *v1.SignUpRes, err error) {
	userId, err := service.User().SignUp(ctx, req.UserName, req.Password, req.InviteCode)
	if err != nil {
		return nil, err
	}
	return &v1.SignUpRes{UserId: userId}, nil
}
