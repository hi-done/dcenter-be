package v1

import "github.com/gogf/gf/v2/frame/g"

type SignUpReq struct {
	g.Meta     `path:"/sys/signUp" method:"post" summary:"用户注册" tags:"系统"`
	UserName   string `json:"userName" v:"required|length:3,20|regex:^[a-zA-Z0-9_]+$" dc:"用户名（英文，3-20位字母/数字/下划线）"`
	Password   string `json:"password" v:"required|length:6,32" dc:"密码（6-32位）"`
	InviteCode string `json:"inviteCode" v:"required" dc:"邀请码"`
}

type SignUpRes struct {
	UserId uint64 `json:"userId" dc:"新用户ID"`
}
