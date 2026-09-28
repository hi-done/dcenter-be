package sys

import (
	"context"

	v1 "dcenter/api/sys/v1"
	"dcenter/internal/service"
)

func (c *ControllerV1) SignIn(ctx context.Context, req *v1.SignInReq) (res *v1.SignInRes, err error) {
	accessToken, refreshToken, err := service.User().SignIn(ctx, req.UserName, req.Password)
	if err != nil {
		return nil, err
	}
	return &v1.SignInRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
