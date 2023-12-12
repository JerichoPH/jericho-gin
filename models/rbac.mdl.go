package models

import (
	"fmt"
	"jericho-gin/database"
	"jericho-gin/wrongs"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	// RbacRoleModel 角色模型
	RbacRoleModel struct {
		MySqlModel
		Name            string                 `gorm:"type:varchar(64);not null;comment:角色名称;" json:"name"`
		Accounts        []*AccountModel        `gorm:"many2many:pivot_rbac_roles__accounts;foreignKey:uuid;joinForeignKey:rbac_role_uuid;references:uuid;joinReferences:accountUuid;" json:"accounts"`
		RbacPermissions []*RbacPermissionModel `gorm:"many2many:pivot_rbac_roles__rbac_permissions;foreignKey:uuid;joinForeignKey:rbac_role_uuid;references:uuid;joinReferences:rbacPermissionUuid;" json:"rbac_permissions"`
		RbacMenus       []*RbacMenuModel       `gorm:"many2many:pivot_rbac_roles__rbac_menus;foreignKey:uuid;joinForeignKey:rbac_role_uuid;references:uuid;joinReferences:rbacMenuUuid;" json:"rbac_menus"`
	}

	// RbacPermissionModel 权限模型
	RbacPermissionModel struct {
		MySqlModel
		Name                    string           `gorm:"type:varchar(64);not null;comment:权限名称;" json:"name"`
		Description             *string          `gorm:"type:text;comment:权限描述;" json:"description"`
		Uri                     string           `gorm:"type:varchar(255);not null;default:'';comment:权限所属路由;" json:"uri"`
		Method                  string           `gorm:"type:varchar(32);not null;comment:请求方法;" json:"method"`
		RbacPermissionGroupUuid string           `gorm:"type:varchar(36);not null;default:'';comment:所属权限组uuid;" json:"rbac_permission_group_uuid"`
		RbacRoles               []*RbacRoleModel `gorm:"many2many:pivot_rbac_roles__rbac_permissions;foreignKey:uuid;joinForeignKey:rbac_permission_uuid;references:uuid;joinReferences:rbacRoleUuid;" json:"rbac_roles"`
	}

	// RbacMenuModel 菜单模型
	RbacMenuModel struct {
		MySqlModel
		Name        string           `gorm:"type:varchar(64);not null;comment:菜单名称" json:"name"`
		SubTitle    string           `gorm:"type:varchar(128);not null;default:'';comment:菜单副标题" json:"sub_title"`
		Description *string          `gorm:"type:text;comment:菜单描述" json:"description"`
		Uri         string           `gorm:"type:varchar(128);not null;default:'';comment:菜单所属路由" json:"uri"`
		Icon        string           `gorm:"type:varchar(64);not null;default:'';comment:菜单图标" json:"icon"`
		RbacRoles   []*RbacRoleModel `gorm:"many2many:pivot_rbac_roles__rbac_menus;foreignKey:uuid;joinForeignKey:rbac_menu_uuid;references:uuid;joinReferences:rbacRoleUuid;" json:"rbac_roles"`
		ParentUuid  string           `gorm:"type:varchar(36);not null;default:'';comment:父级uuid;" json:"parent_uuid"`
		Parent      *RbacMenuModel   `gorm:"foreignKey:parent_uuid;references:uuid;comment:所属父级;" json:"parent"`
		Subs        []*RbacMenuModel `gorm:"foreignKey:parent_uuid;references:uuid;comment:相关子集;" json:"subs"`
	}

	PivotRbacRoleAccountModel struct {
		RbacRoleUuid string `gorm:"type:varchar(36);not null;default:'';comment:角色uuid" json:"rbac_role_uuid"`
		AccountUuid  string `gorm:"type:varchar(36);not null;default:'';comment:用户uuid" json:"account_uuid"`
	}

	PivotRbacRoleRbacPermissionModel struct {
		RbacRoleUuid       string `gorm:"type:varchar(36);not null;default:'';comment:角色uuid" json:"rbac_role_uuid"`
		RbacPermissionUuid string `gorm:"type:varchar(36);not null;default:'';comment:权限uuid" json:"rbac_permission_uuid"`
	}

	PivotRbacRoleRbacMenuModel struct {
		RbacRoleUuid string `gorm:"type:varchar(36);not null;default:'';comment:角色uuid" json:"rbac_role_uuid"`
		RbacMenuUuid string `gorm:"type:varchar(36);not null;default:'';comment:菜单uuid" json:"rbac_menu_uuid"`
	}
)

// NewRbacRoleModel 创建一个新的 RBAC 角色模型
func NewRbacRoleModel() *MySqlModel {
	return NewMySqlModel().SetModel(&RbacRoleModel{})
}

// TableName 角色表名称
func (RbacRoleModel) TableName() string {
	return "rbac_roles"
}

// GetListByQuery 根据Query获取角色列表
func (receiver RbacRoleModel) GetListByQuery(ctx *gin.Context) *gorm.DB {
	return NewRbacRoleModel().
		SetWheresEqual("be_enable").
		SetWheresDateBetween("created_at", "updated_at", "deleted_at").
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

// NewRbacPermissionModel 返回一个新的 RbacPermissionModel 模型实例化的指针
func NewRbacPermissionModel() *MySqlModel {
	return NewMySqlModel().SetModel(&RbacPermissionModel{})
}

// TableName 权限表名称
func (RbacPermissionModel) TableName() string {
	return "rbac_permissions"
}

// GetListByQuery 根据Query获取权限列表
func (receiver RbacPermissionModel) GetListByQuery(ctx *gin.Context) *gorm.DB {
	return NewRbacPermissionModel().
		SetWheresEqual("be_enable").
		SetWheresDateBetween("created_at", "updated_at", "deleted_at").
		SetWheresExtraHasValue(map[string]func(string, *gorm.DB) *gorm.DB{
			"name": func(value string, db *gorm.DB) *gorm.DB {
				return db.Where(fmt.Sprintf("rp.name like '%%%s%%'", value))
			},
			"rbac_role_uuid": func(value string, db *gorm.DB) *gorm.DB {
				return db.
					Joins("left join pivot_rbac_roles__rbac_permissions prrrp on rp.uuid = prrrp.rbac_permission_uuid").
					Joins("left join rbac_roles rr on prrrp.rbac_role_uuid = rr.uuid").
					Where("rr.uuid =?", value)
			},
		}).
		SetWheresExtraHasValues(map[string]func([]string, *gorm.DB) *gorm.DB{
			"names[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("rp.name in (?)", values)
			},
			"rbac_role_uuids[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.
					Joins("left join pivot_rbac_roles__rbac_permissions prrrp on rp.uuid = prrrp.rbac_permission_uuid").
					Joins("left join rbac_roles rr on prrrp.rbac_role_uuid = rr.uuid").
					Where("rr.uuid in (?)", values)
			},
		}).
		SetCtx(ctx).
		GetDbUseQuery("").
		Table("rbac_permissions as rp")
}

// NewRbacMenuModel 返回一个新的 RbacMenuModel 模型实例指针
func NewRbacMenuModel() *MySqlModel {
	return NewMySqlModel().SetModel(&RbacMenuModel{})
}

// TableName 菜单表名称
func (RbacMenuModel) TableName() string {
	return "rbac_menus"
}

// GetListByQuery 根据Query获取菜单列表
func (receiver RbacMenuModel) GetListByQuery(ctx *gin.Context) *gorm.DB {
	var (
		notHasSubs = ctx.Query("not_has_subs")
		subs       = make(map[string]map[string]string)
		subUuids   []string
	)
	subUuids = make([]string, 0)
	if notHasSubs != "" {
		subs = receiver.GetSubUuidsByParentUuid(notHasSubs)
		for _, sub := range subs {
			subUuids = append(subUuids, sub["uuid"])
		}
	}

	return NewRbacMenuModel().
		SetWheresEqual("be_enable").
		SetWheresDateBetween("created_at", "updated_at", "deleted_at").
		SetWheresExtraHasValue(map[string]func(string, *gorm.DB) *gorm.DB{
			"name": func(value string, db *gorm.DB) *gorm.DB {
				return db.Where(fmt.Sprintf("name like '%%%s%%'", value))
			},
			"uri": func(value string, db *gorm.DB) *gorm.DB {
				return db.Where("uri", value)
			},
			"not_has_subs": func(value string, db *gorm.DB) *gorm.DB {
				if len(subUuids) == 0 {
					return db
				}
				return db.Where("uuid not in ?", subUuids)
			},
		}).
		SetWheresExtraHasValues(map[string]func([]string, *gorm.DB) *gorm.DB{
			"names[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("name in (?)", values)
			},
			"uris[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("uri in (?)", values)
			},
		}).
		SetCtx(ctx).
		GetDbUseQuery("").
		Table("rbac_menus as rm")
}

// GetSubUuidsByParentUuid 根据父级uuid获取所有子集uuid
func (receiver RbacMenuModel) GetSubUuidsByParentUuid(parentUuid string) map[string]map[string]string {
	if rows := database.ExecSql(strings.Join([]string{
		"WITH RECURSIVE cte AS (",
		"SELECT uuid, name, NULL AS parent_uuid",
		"FROM rbac_menus",
		"WHERE parent_uuid = ?",
		"AND deleted_at IS NULL",
		"UNION ALL",
		"SELECT m.uuid, m.name, c.parent_uuid",
		"FROM rbac_menus m INNER JOIN cte c ON m.parent_uuid = c.uuid",
		"WHERE m.deleted_at IS NULL",
		")",
		"SELECT uuid, name FROM cte",
	}, "\r\n"), parentUuid); rows != nil {
		var (
			subs = make(map[string]map[string]string)
			err  error
		)
		for rows.Next() {
			var (
				uuid string
				name string
			)
			if err = rows.Scan(&uuid, &name); err != nil {
				wrongs.ThrowForbidden(err.Error())
			}
			subs[uuid] = map[string]string{
				"uuid": uuid,
				"name": name,
			}
		}
		return subs
	}
	return map[string]map[string]string{}
}

// NewPivotRbacRoleAccountModel 返回一个新的 PivotRbacRoleAccountModel 模型实例
func NewPivotRbacRoleAccountModel() *MySqlModel {
	return NewMySqlModel().SetModel(&PivotRbacRoleAccountModel{})
}

// TableName 角色与用户对应关系表名称
func (PivotRbacRoleAccountModel) TableName() string {
	return "pivot_rbac_roles__accounts"
}

// BindRbacRoles 绑定角色与用户
func (PivotRbacRoleAccountModel) BindRbacRoles(rbacRole *RbacRoleModel, accounts []*AccountModel) {
	database.NewGormLauncher().GetConn("").Table("pivot_rbac_roles__accounts").Where("rbac_role_uuid", rbacRole.Uuid).Delete(nil)
	if len(accounts) > 0 {
		for _, account := range accounts {
			NewPivotRbacRoleAccountModel().
				GetDb("").
				Create(&PivotRbacRoleAccountModel{
					RbacRoleUuid: rbacRole.Uuid,
					AccountUuid:  account.Uuid,
				})
		}
	}
}

// NewPivotRbacRoleRbacPermissionModel 返回一个新的 PivotRbacRoleRbacPermissionModel 模型的实例。
func NewPivotRbacRoleRbacPermissionModel() *MySqlModel {
	return NewMySqlModel().SetModel(&PivotRbacRoleRbacPermissionModel{})
}

// TableName 角色与权限对应关系表名称
func (PivotRbacRoleRbacPermissionModel) TableName() string {
	return "pivot_rbac_roles__rbac_permissions"
}

// BindRbacRoles 绑定角色与权限
func (PivotRbacRoleRbacPermissionModel) BindRbacRoles(rbacPermission *RbacPermissionModel, rbacRoles []*RbacRoleModel) {
	database.NewGormLauncher().GetConn("").Table("pivot_rbac_roles__rbac_permissions").Where("rbac_permission_uuid", rbacPermission.Uuid).Delete(nil)
	if len(rbacRoles) > 0 {
		for _, rbacRole := range rbacRoles {
			NewPivotRbacRoleRbacPermissionModel().
				GetDb("").
				Create(&PivotRbacRoleRbacPermissionModel{
					RbacRoleUuid:       rbacRole.Uuid,
					RbacPermissionUuid: rbacPermission.Uuid,
				})
		}
	}
}

// NewPivotRbacRoleRbacMenuModel 返回一个新的 PivotRbacRoleRbacMenuModel 模型的实例。
func NewPivotRbacRoleRbacMenuModel() *MySqlModel {
	return NewMySqlModel().SetModel(&PivotRbacRoleRbacMenuModel{})
}

// TableName 角色与菜单对应关系表名称
func (PivotRbacRoleRbacMenuModel) TableName() string {
	return "pivot_rbac_roles__rbac_menus"
}

// BindRbacRoles 绑定角色与菜单
func (PivotRbacRoleRbacMenuModel) BindRbacRoles(rbacMenu *RbacMenuModel, rbacRoles []*RbacRoleModel) {
	var pivotRbacRoleRbacMenus = make([]*PivotRbacRoleRbacMenuModel, 0)

	database.NewGormLauncher().GetConn("").Table("pivot_rbac_roles__rbac_menus").Where("rbac_menu_uuid = ?", rbacMenu.Uuid).Delete(nil)

	if len(rbacRoles) > 0 {
		for _, rbacRole := range rbacRoles {
			pivotRbacRoleRbacMenus = append(pivotRbacRoleRbacMenus, &PivotRbacRoleRbacMenuModel{
				RbacRoleUuid: rbacRole.Uuid,
				RbacMenuUuid: rbacMenu.Uuid,
			})
		}
	}

	NewPivotRbacRoleRbacMenuModel().GetDb("").Create(&pivotRbacRoleRbacMenus)
}
