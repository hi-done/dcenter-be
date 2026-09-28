// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DcMenu is the golang structure for table dc_menu.
type DcMenu struct {
	Id          uint64      `json:"id"          orm:"id"           description:"菜单ID"`                            // 菜单ID
	Title       string      `json:"title"       orm:"title"        description:"菜单显示名称"`                          // 菜单显示名称
	Icon        string      `json:"icon"        orm:"icon"         description:"菜单图标"`                            // 菜单图标
	Sort        int         `json:"sort"        orm:"sort"         description:"菜单排序，数值越小越靠前"`                    // 菜单排序，数值越小越靠前
	ParentId    uint64      `json:"parentId"    orm:"parent_id"    description:"父菜单ID，0代表一级菜单"`                   // 父菜单ID，0代表一级菜单
	VisibleRank int         `json:"visibleRank" orm:"visible_rank" description:"菜单可见等级，数值越高权限要求越高；用于判断普通用户是否可查看"` // 菜单可见等级，数值越高权限要求越高；用于判断普通用户是否可查看
	Status      int         `json:"status"      orm:"status"       description:"菜单状态：1启用，0停用"`                    // 菜单状态：1启用，0停用
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"`                            // 创建时间
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:"软删除时间"`                           // 软删除时间
}
