// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DcUser is the golang structure for table dc_user.
type DcUser struct {
	Id          uint64      `json:"id"          orm:"id"            description:"用户ID"`                           // 用户ID
	UserName    string      `json:"userName"    orm:"user_name"     description:"用户名"`                            // 用户名
	Email       string      `json:"email"       orm:"email"         description:"邮箱"`                             // 邮箱
	Phone       string      `json:"phone"       orm:"phone"         description:"手机号"`                            // 手机号
	Password    string      `json:"password"    orm:"password"      description:"密码哈希(bcrypt)"`                   // 密码哈希(bcrypt)
	InviterId   uint64      `json:"inviterId"   orm:"inviter_id"    description:"直接邀请人ID，关联dc_user.id，无邀请人为null"` // 直接邀请人ID，关联dc_user.id，无邀请人为null
	UserRank    int         `json:"userRank"    orm:"user_rank"     description:"用户权限等级"`                         // 用户权限等级
	Status      int         `json:"status"      orm:"status"        description:"账号状态：1正常 0禁用 2锁定"`               // 账号状态：1正常 0禁用 2锁定
	LastLoginAt *gtime.Time `json:"lastLoginAt" orm:"last_login_at" description:"最近登录时间"`                         // 最近登录时间
	LastLoginIp string      `json:"lastLoginIp" orm:"last_login_ip" description:"最近登录IP"`                         // 最近登录IP
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    description:"创建时间"`                           // 创建时间
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"    description:"软删除时间，NULL=未删除"`                 // 软删除时间，NULL=未删除
}
