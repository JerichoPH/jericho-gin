package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	RbacRoleModel struct {
		MySqlModel
		Name            string                 `gorm:"type:varchar(64);not null;comment:角色名称;" json:"name"`
		Accounts        []*AccountModel        `gorm:"many2many:pivot_rbac_roles__accounts;foreignKey:uuid;joinForeignKey:rbac_role_uuid;references:uuid;joinReferences:accountUuid;" json:"accounts"`
		RbacPermissions []*RbacPermissionModel `gorm:"many2many:pivot_rbac_roles__rbac_permissions;foreignKey:uuid;joinForeignKey:rbac_role_uuid;references:uuid;joinReferences:rbacPermissionUuid;" json:"rbacPermissions"`
		RbacMenus       []*RbacMenuModel       `gorm:"many2many:pivot_rbac_roles__rbac_menus;foreignKey:uuid;joinForeignKey:rbac_role_uuid;references:uuid;joinReferences:rbacMenuUuid;" json:"rbacMenus"`
	}

	RbacPermissionModel struct {
		MySqlModel
		Name        string           `gorm:"type:varchar(64);not null;comment:权限名称;" json:"name"`
		Description *string          `gorm:"type:text;comment:权限描述;" json:"description"`
		Uri         string           `gorm:"type:varchar(255);not null;default:'';comment:权限所属路由;" json:"uri"`
		RbacRoles   []*RbacRoleModel `gorm:"many2many:pivot_rbac_roles__rbac_permissions;foreignKey:uuid;joinForeignKey:rbac_permission_uuid;references:uuid;joinReferences:rbacRoleUuid;" json:"rbacRoles"`
	}

	RbacMenuModel struct {
		MySqlModel
		Name        string           `gorm:"type:varchar(64);not null;comment:菜单名称" json:"name"`
		SubTitle    string           `gorm:"type:varchar(255);not null;default:'';comment:菜单副标题" json:"subTitle"`
		Description *string          `gorm:"type:text;comment:菜单描述" json:"description"`
		Uri         string           `gorm:"type:varchar(255);not null;default:'';comment:菜单所属路由" json:"uri"`
		RbacRoles   []*RbacRoleModel `gorm:"many2many:pivot_rbac_roles__rbac_menus;foreignKey:uuid;joinForeignKey:rbac_menu_uuid;references:uuid;joinReferences:rbacRoleUuid;" json:"rbacRoles"`
	}

	PivotRbacRoleAccountModel struct {
		RbacRoleUuid string `gorm:"type:varchar(36);not null;default:'';comment:角色uuid" json:"rbacRoleUuid"`
		AccountUuid  string `gorm:"type:varchar(36);not null;default:'';comment:用户uuid" json:"accountUuid"`
	}

	PivotRbacRoleRbacPermissionModel struct {
		RbacRoleUuid       string `gorm:"type:varchar(36);not null;default:'';comment:角色uuid" json:"rbacRoleUuid"`
		RbacPermissionUuid string `gorm:"type:varchar(36);not null;default:'';comment:权限uuid" json:"rbacPermissionUuid"`
	}

	PivotRbacRoleRbacMenuModel struct {
		RbacRoleUuid string `gorm:"type:varchar(36);not null;default:'';comment:角色uuid" json:"rbacRoleUuid"`
		RbacMenuUuid string `gorm:"type:varchar(36);not null;default:'';comment:菜单uuid" json:"rbacMenuUuid"`
	}
)

// TableName 角色表名称
func (RbacRoleModel) TableName() string {
	return "rbac_roles"
}

// NewRbacRoleModel 创建一个新的 RBAC 角色模型
func NewRbacRoleModel() *MySqlModel {
	return NewMySqlModel().SetModel(&RbacRoleModel{})
}

// GetListByQuery 根据Query获取角色列表
func (receiver RbacRoleModel) GetListByQuery(ctx *gin.Context) *gorm.DB {
	return NewRbacRoleModel().
		SetWheresExtraHasValue(map[string]func(string, *gorm.DB) *gorm.DB{
			"name": func(value string, db *gorm.DB) *gorm.DB {
				return db.Where(fmt.Sprintf("name like '%%%s%%'", value))
			},
		}).
		SetWheresExtraHasValues(map[string]func([]string, *gorm.DB) *gorm.DB{
			"names[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("name in (?)", values)
			},
		}).
		SetCtx(ctx).
		GetDbUseQuery("").
		Table("rbac_roles as rr")
}

// TableName 权限表名称
func (RbacPermissionModel) TableName() string {
	return "rbac_permissions"
}

// NewRbacPermissionModel 返回一个新的 RbacPermissionModel 模型实例化的指针
func NewRbacPermissionModel() *MySqlModel {
	return NewMySqlModel().SetModel(&RbacPermissionModel{})
}

// GetListByQuery 根据Query获取权限列表
func (receiver RbacPermissionModel) GetListByQuery(ctx *gin.Context) *gorm.DB {
	return NewRbacPermissionModel().
		SetWheresExtraHasValue(map[string]func(string, *gorm.DB) *gorm.DB{
			"name": func(value string, db *gorm.DB) *gorm.DB {
				return db.Where(fmt.Sprintf("rp.name like '%%%s%%'", value))
			},
			"rbac_role_uuid": func(value string, db *gorm.DB) *gorm.DB {
				return db.Where("rr.uuid =?", value)
			},
		}).
		SetWheresExtraHasValues(map[string]func([]string, *gorm.DB) *gorm.DB{
			"names[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("rp.name in (?)", values)
			},
			"rbac_role_uuids[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("rr.uuid in (?)", values)
			},
		}).
		SetCtx(ctx).
		GetDbUseQuery("").
		Table("rbac_permissions as rp").
		Joins("left join pivot_rbac_roles__rbac_permissions prrrp on rp.uuid = prrrp.rbac_permission_uuid").
		Joins("left join rbac_roles rr on prrrp.rbac_role_uuid = rr.uuid").Debug()
}

// TableName 菜单表名称
func (RbacMenuModel) TableName() string {
	return "rbac_menus"
}

// NewRbacMenuModel 返回一个新的 RbacMenuModel 模型实例指针
func NewRbacMenuModel() *MySqlModel {
	return NewMySqlModel().SetModel(&RbacMenuModel{})
}

// TableName 角色与用户对应关系表名称
func (PivotRbacRoleAccountModel) TableName() string {
	return "pivot_rbac_rolse__accounts"
}

// NewPivotRbacRoleAccountModel 返回一个新的 PivotRbacRoleAccountModel 模型实例
func NewPivotRbacRoleAccountModel() *MySqlModel {
	return NewMySqlModel().SetModel(&PivotRbacRoleAccountModel{})
}

// TableName 角色与权限对应关系表名称
func (PivotRbacRoleRbacPermissionModel) TableName() string {
	return "pivot_rbac_roles__rbac_permissions"
}

// NewPivotRbacRoleRbacPermissionModel 返回一个新的 PivotRbacRoleRbacPermissionModel 模型的实例。
func NewPivotRbacRoleRbacPermissionModel() *MySqlModel {
	return NewMySqlModel().SetModel(&PivotRbacRoleRbacPermissionModel{})
}

// TableName 角色与菜单对应关系表名称
func (PivotRbacRoleRbacMenuModel) TableName() string {
	return "pivot_rbac_role__rbac_menus"
}

// NewPivotRbacRoleRbacMenuModel 返回一个新的 PivotRbacRoleRbacMenuModel 模型的实例。
func NewPivotRbacRoleRbacMenuModel() *MySqlModel {
	return NewMySqlModel().SetModel(&PivotRbacRoleRbacMenuModel{})
}
