package user

import (
	"context"

	v1 "dcenter/api/user/v1"
	"dcenter/internal/service"
)

func (c *ControllerV1) SignOut(ctx context.Context, req *v1.SignOutReq) (res *v1.SignOutRes, err error) {
	if err = service.User().SignOut(ctx); err != nil {
		return nil, err
	}
	return &v1.SignOutRes{}, nil
}
