// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DcMenuDao is the data access object for the table dc_menu.
type DcMenuDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DcMenuColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DcMenuColumns defines and stores column names for the table dc_menu.
type DcMenuColumns struct {
	Id          string // 菜单ID
	Title       string // 菜单显示名称
	Icon        string // 菜单图标
	Sort        string // 菜单排序，数值越小越靠前
	ParentId    string // 父菜单ID，0代表一级菜单
	VisibleRank string // 菜单可见等级，数值越高权限要求越高；用于判断普通用户是否可查看
	Status      string // 菜单状态：1启用，0停用
	CreatedAt   string // 创建时间
	DeletedAt   string // 软删除时间
}

// dcMenuColumns holds the columns for the table dc_menu.
var dcMenuColumns = DcMenuColumns{
	Id:          "id",
	Title:       "title",
	Icon:        "icon",
	Sort:        "sort",
	ParentId:    "parent_id",
	VisibleRank: "visible_rank",
	Status:      "status",
	CreatedAt:   "created_at",
	DeletedAt:   "deleted_at",
}

// NewDcMenuDao creates and returns a new DAO object for table data access.
func NewDcMenuDao(handlers ...gdb.ModelHandler) *DcMenuDao {
	return &DcMenuDao{
		group:    "default",
		table:    "dc_menu",
		columns:  dcMenuColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DcMenuDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DcMenuDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DcMenuDao) Columns() DcMenuColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DcMenuDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DcMenuDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DcMenuDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
