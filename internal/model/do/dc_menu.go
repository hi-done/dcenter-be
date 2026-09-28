// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DcMenu is the golang structure of table dc_menu for DAO operations like Where/Data.
type DcMenu struct {
	g.Meta      `orm:"table:dc_menu, do:true"`
	Id          any         // 菜单ID
	Title       any         // 菜单显示名称
	Icon        any         // 菜单图标
	Sort        any         // 菜单排序，数值越小越靠前
	ParentId    any         // 父菜单ID，0代表一级菜单
	VisibleRank any         // 菜单可见等级，数值越高权限要求越高；用于判断普通用户是否可查看
	Status      any         // 菜单状态：1启用，0停用
	CreatedAt   *gtime.Time // 创建时间
	DeletedAt   *gtime.Time // 软删除时间
}
