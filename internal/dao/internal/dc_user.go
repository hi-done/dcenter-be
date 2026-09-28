// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DcUserDao is the data access object for the table dc_user.
type DcUserDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DcUserColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DcUserColumns defines and stores column names for the table dc_user.
type DcUserColumns struct {
	Id          string // 用户ID
	UserName    string // 用户名
	Email       string // 邮箱
	Phone       string // 手机号
	Password    string // 密码哈希(bcrypt)
	InviterId   string // 直接邀请人ID，关联dc_user.id，无邀请人为null
	UserRank    string // 用户权限等级
	Status      string // 账号状态：1正常 0禁用 2锁定
	LastLoginAt string // 最近登录时间
	LastLoginIp string // 最近登录IP
	CreatedAt   string // 创建时间
	DeletedAt   string // 软删除时间，NULL=未删除
}

// dcUserColumns holds the columns for the table dc_user.
var dcUserColumns = DcUserColumns{
	Id:          "id",
	UserName:    "user_name",
	Email:       "email",
	Phone:       "phone",
	Password:    "password",
	InviterId:   "inviter_id",
	UserRank:    "user_rank",
	Status:      "status",
	LastLoginAt: "last_login_at",
	LastLoginIp: "last_login_ip",
	CreatedAt:   "created_at",
	DeletedAt:   "deleted_at",
}

// NewDcUserDao creates and returns a new DAO object for table data access.
func NewDcUserDao(handlers ...gdb.ModelHandler) *DcUserDao {
	return &DcUserDao{
		group:    "default",
		table:    "dc_user",
		columns:  dcUserColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DcUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DcUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DcUserDao) Columns() DcUserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DcUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DcUserDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *DcUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
