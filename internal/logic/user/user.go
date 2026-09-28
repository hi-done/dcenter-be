package user

import "dcenter/internal/service"

type sUser struct {
}

func init() {
	service.RegisterUser(New())
}

func New() service.IUser {
	return &sUser{}
}
