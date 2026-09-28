package v1

import "github.com/gogf/gf/v2/frame/g"

type SignOutReq struct {
	g.Meta `path:"/user/signOut" method:"post" summary:"用户注销"`
}

type SignOutRes struct {
}
