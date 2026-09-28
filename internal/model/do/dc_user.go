// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DcUser is the golang structure of table dc_user for DAO operations like Where/Data.
type DcUser struct {
	g.Meta      `orm:"table:dc_user, do:true"`
	Id          any         // 用户ID
	UserName    any         // 用户名
	Email       any         // 邮箱
	Phone       any         // 手机号
	Password    any         // 密码哈希(bcrypt)
	InviterId   any         // 直接邀请人ID，关联dc_user.id，无邀请人为null
	UserRank    any         // 用户权限等级
	Status      any         // 账号状态：1正常 0禁用 2锁定
	LastLoginAt *gtime.Time // 最近登录时间
	LastLoginIp any         // 最近登录IP
	CreatedAt   *gtime.Time // 创建时间
	DeletedAt   *gtime.Time // 软删除时间，NULL=未删除
}
