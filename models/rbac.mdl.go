package models

type (
	RbacRoleModel struct {
		GormModel
		Name            string                 `gorm:"unique;type:varchar(64);not null;comment:角色名称;"`
		Accounts        []*AccountModel        `gorm:"many2many:pivot_rbac_roles__accounts;foreignKey:uuid;joinForeignKey:rbac_rold_uuid;references:uuid;joinReferences:account_uuid;"`
		RbacPermissions []*RbacPermissionModel `gorm:"many2many:pivot_rbac_roles__rbac_permissions;foreignKey:uuid;joinForeignKey:rbac_rold_uuid;references:uuid;joinReferences:rbac_permission_uuid;"`
		RbacMenus       []*RbacMenuModel       `gorm:"many2many:pivot_rbac_roles__rbac_menus;foreignKey:uuid;joinForeignKey:rbac_rold_uuid;references:uuid;joinReferences:rbac_menu_uuid;"`
	}

	RbacPermissionModel struct {
		GormModel
		Name        string           `gorm:"unique;type:varchar(64);not null;comment:权限名称;"`
		Description string           `gorm:"type:text;comment:权限描述;"`
		Uri         string           `grom:"type:varchar(255);not null;default:'';comment:权限所属路由;"`
		RbacRoles   []*RbacRoleModel `gorm:"many2many:pivot_rbac_roles__rbac_permissions;foreignKey:uuid;joinForeignKey:rbac_permission_uuid;references:uuid;joinReferences:rbac_role_uuid;"`
	}

	RbacMenuModel struct {
		GormModel
		Name        string           `gorm:"unique;type:varchar(64);not null;comment:菜单名称"`
		SubTitle    string           `gorm:"type:varchar(255);not null;default:'';comment:菜单副标题"`
		Description string           `gorm:"type:text;comment:菜单描述"`
		Uri         string           `grom:"type:varchar(255);not null;default:'';comment:菜单所属路由"`
		RbacRoles   []*RbacRoleModel `gorm:"many2many:pivot_rbac_roles__rbac_menus;foreignKey:uuid;joinForeignKey:rbac_menu_uuid;references:uuid;joinReferences:rbac_role_uuid;"`
	}

	PivotRbacRoleAccountModel struct {
		// 角色 uuid，类型为 varchar(36)，不可为空，默认为空，注释为 "角色uuid"
		RbacRoleUuid string `grom:"type:varchar(36);not null;default:'';comment:角色uuid"`
		// 用户 uuid，类型为 varchar(36)，不可为空，默认为空，注释为 "用户uuid"
		AccountUuid string `grom:"type:varchar(36);not null;default:'';comment:用户uuid"`
	}

	PivotRbacRoleRbacPermissionModel struct {
		// 角色 uuid，类型为 varchar(36)，不可为空，默认为空，注释为 "角色uuid"
		RbacRoleUuid string `grom:"type:varchar(36);not null;default:'';comment:角色uuid"`
		// 权限 uuid，类型为 varchar(36)，不可为空，默认为空，注释为 "权限uuid"
		RbacPermissionUuid string `grom:"type:varchar(36);not null;default:'';comment:权限uuid"`
	}

	PivotRbacRoleRbacMenuModel struct {
		// 角色 uuid，类型为 varchar(36)，不可为空，默认为空，注释为 "角色uuid"
		RbacRoleUuid string `grom:"type:varchar(36);not null;default:'';comment:角色uuid"`
		// 菜单 uuid，类型为 varchar(36)，不可为空，默认为空，注释为 "菜单uuid"
		RbacMenuUuid string `grom:"type:varchar(36);not null;default:'';comment:菜单uuid"`
	}
)

// TableName 角色表名称
func (RbacRoleModel) TableName() string {
	return "rbac_roles"
}

// NewRbacRoleModel 创建一个新的 RBAC 角色模型
func NewRbacRoleModel() *GormModel {
	return NewGorm().SetModel(&RbacRoleModel{})
}

// TableName 权限表名称
func (RbacPermissionModel) TableName() string {
	return "rbac_permissions"
}

// NewRbacPermissionModel 返回一个新的 RbacPermissionModel 模型实例化的指针
func NewRbacPermissionModel() *GormModel {
	return NewGorm().SetModel(&RbacPermissionModel{})
}

// TableName 菜单表名称
func (RbacMenuModel) TableName() string {
	return "rbac_menus"
}

// NewRbacMenuModel 返回一个新的 RbacMenuModel 模型实例指针
func NewRbacMenuModel() *GormModel {
	return NewGorm().SetModel(&RbacMenuModel{})
}

// TableName 角色与用户对应关系表名称
func (PivotRbacRoleAccountModel) TableName() string {
	return "pivot_rbac_rolse__accounts"
}

// NewPivotRbacRoleAccountModel 返回一个新的 PivotRbacRoleAccountModel 模型实例
func NewPivotRbacRoleAccountModel() *GormModel {
	return NewGorm().SetModel(&PivotRbacRoleAccountModel{})
}

// TableName 角色与权限对应关系表名称
func (PivotRbacRoleRbacPermissionModel) TableName() string {
	return "pivot_rbac_roles__rbac_permissions"
}

// NewPivotRbacRoleRbacPermissionModel 返回一个新的 PivotRbacRoleRbacPermissionModel 模型的实例。
func NewPivotRbacRoleRbacPermissionModel() *GormModel {
	return NewGorm().SetModel(&PivotRbacRoleRbacPermissionModel{})
}

// TableName 角色与菜单对应关系表名称
func (PivotRbacRoleRbacMenuModel) TableName() string {
	return "pivot_rbac_role__rbac_menus"
}

// NewPivotRbacRoleRbacMenuModel 返回一个新的 PivotRbacRoleRbacMenuModel 模型的实例。
func NewPivotRbacRoleRbacMenuModel() *GormModel {
	return NewGorm().SetModel(&PivotRbacRoleRbacMenuModel{})
}
