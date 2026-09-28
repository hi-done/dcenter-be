package v1

import "github.com/gogf/gf/v2/frame/g"

type SignInReq struct {
	g.Meta   `path:"/sys/signIn" method:"post" summary:"用户登录" tags:"系统"`
	UserName string `json:"userName" v:"required" dc:"用户名"`
	Password string `json:"password" v:"required" dc:"密码"`
}

type SignInRes struct {
	AccessToken  string `json:"accessToken"  dc:"访问令牌"`
	RefreshToken string `json:"refreshToken" dc:"刷新令牌"`
}
