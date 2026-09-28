package v1

import "github.com/gogf/gf/v2/frame/g"

type GetInviteCodeReq struct {
	g.Meta `path:"/user/getInviteCode" method:"get" summary:"获取邀请码"`
}

type GetInviteCodeRes struct {
	InviteCode string `json:"inviteCode"`
}
