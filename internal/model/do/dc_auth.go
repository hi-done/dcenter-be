// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DcAuth is the golang structure of table dc_auth for DAO operations like Where/Data.
type DcAuth struct {
	g.Meta    `orm:"table:dc_auth, do:true"`
	Id        any         // 主键ID
	UserId    any         // 用户ID，关联dc_user.id
	MenuId    any         // 菜单ID，关联dc_menu.id
	CreatorId any         // 创建人ID，操作授权的管理员
	CreatedAt *gtime.Time // 创建时间
	DeletedAt *gtime.Time // 软删除时间
}
