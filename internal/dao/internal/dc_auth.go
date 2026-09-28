// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DcAuthDao is the data access object for the table dc_auth.
type DcAuthDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DcAuthColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DcAuthColumns defines and stores column names for the table dc_auth.
type DcAuthColumns struct {
	Id        string // 主键ID
	UserId    string // 用户ID，关联dc_user.id
	MenuId    string // 菜单ID，关联dc_menu.id
	CreatorId string // 创建人ID，操作授权的管理员
	CreatedAt string // 创建时间
	DeletedAt string // 软删除时间
}

// dcAuthColumns holds the columns for the table dc_auth.
var dcAuthColumns = DcAuthColumns{
	Id:        "id",
	UserId:    "user_id",
	MenuId:    "menu_id",
	CreatorId: "creator_id",
	CreatedAt: "created_at",
	DeletedAt: "deleted_at",
}

// NewDcAuthDao creates and returns a new DAO object for table data access.
func NewDcAuthDao(handlers ...gdb.ModelHandler) *DcAuthDao {
	return &DcAuthDao{
		group:    "default",
		table:    "dc_auth",
		columns:  dcAuthColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DcAuthDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DcAuthDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DcAuthDao) Columns() DcAuthColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DcAuthDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DcAuthDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DcAuthDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
