// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DcAuth is the golang structure for table dc_auth.
type DcAuth struct {
	Id        uint64      `json:"id"        orm:"id"         description:"主键ID"`              // 主键ID
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"用户ID，关联dc_user.id"` // 用户ID，关联dc_user.id
	MenuId    uint64      `json:"menuId"    orm:"menu_id"    description:"菜单ID，关联dc_menu.id"` // 菜单ID，关联dc_menu.id
	CreatorId uint64      `json:"creatorId" orm:"creator_id" description:"创建人ID，操作授权的管理员"`    // 创建人ID，操作授权的管理员
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`              // 创建时间
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:"软删除时间"`             // 软删除时间
}
