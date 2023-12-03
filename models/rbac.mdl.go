package models

type (
	RbacRoleModel struct {
		MysqlModel
		Name            string                 `gorm:"unique;type:varchar(64);not null;comment:角色名称;"`
		Accounts        []*AccountModel        `gorm:"many2many:pivot_rbac_roles__accounts;foreignKey:uuid;joinForeignKey:rbac_rold_uuid;references:uuid;joinReferences:account_uuid;"`
		RbacPermissions []*RbacPermissionModel `gorm:"many2many:pivot_rbac_roles__rbac_permissions;foreignKey:uuid;joinForeignKey:rbac_rold_uuid;references:uuid;joinReferences:rbac_permission_uuid;"`
		RbacMenus       []*RbacMenuModel       `gorm:"many2many:pivot_rbac_roles__rbac_menus;foreignKey:uuid;joinForeignKey:rbac_rold_uuid;references:uuid;joinReferences:rbac_menu_uuid;"`
	}

	RbacPermissionModel struct {
		MysqlModel
		Name        string           `gorm:"unique;type:varchar(64);not null;comment:权限名称;"`
		Description string           `gorm:"type:text;comment:权限描述;"`
		Uri         string           `grom:"type:varchar(255);not null;default:'';comment:权限所属路由;"`
		RbacRoles   []*RbacRoleModel `gorm:"many2many:pivot_rbac_roles__rbac_permissions;foreignKey:uuid;joinForeignKey:rbac_permission_uuid;references:uuid;joinReferences:rbac_role_uuid;"`
	}

	RbacMenuModel struct {
		MysqlModel
		Name        string           `gorm:"unique;type:varchar(64);not null;comment:菜单名称"`
		SubTitle    string           `gorm:"type:varchar(255);not null;default:'';comment:菜单副标题"`
		Description string           `gorm:"type:text;comment:菜单描述"`
		Uri         string           `grom:"type:varchar(255);not null;default:'';comment:菜单所属路由"`
		RbacRoles   []*RbacRoleModel `gorm:"many2many:pivot_rbac_roles__rbac_menus;foreignKey:uuid;joinForeignKey:rbac_menu_uuid;references:uuid;joinReferences:rbac_role_uuid;"`
	}

	PivotRbacRoleAccountModel struct {
		RbacRoleUuid string `grom:"type:varchar(36);not null;default:'';comment:角色uuid"`
		AccountUuid  string `grom:"type:varchar(36);not null;default:'';comment:用户uuid"`
	}

	PivotRbacRoleRbacPermissionModel struct {
		RbacRoleUuid       string `grom:"type:varchar(36);not null;default:'';comment:角色uuid"`
		RbacPermissionUuid string `grom:"type:varchar(36);not null;default:'';comment:权限uuid"`
	}

	PivotRbacRoleRbacMenuModel struct {
		RbacRoleUuid string `grom:"type:varchar(36);not null;default:'';comment:角色uuid"`
		RbacMenuUuid string `grom:"type:varchar(36);not null;default:'';comment:菜单uuid"`
	}
)

// TableName 角色表名称
func (RbacRoleModel) TableName() string {
	return "rbac_roles"
}

// NewRbacRoleModel 创建一个新的 RBAC 角色模型
func NewRbacRoleModel() *MysqlModel {
	return NewMySqlModel().SetModel(&RbacRoleModel{})
}

// TableName 权限表名称
func (RbacPermissionModel) TableName() string {
	return "rbac_permissions"
}

// NewRbacPermissionModel 返回一个新的 RbacPermissionModel 模型实例化的指针
func NewRbacPermissionModel() *MysqlModel {
	return NewMySqlModel().SetModel(&RbacPermissionModel{})
}

// TableName 菜单表名称
func (RbacMenuModel) TableName() string {
	return "rbac_menus"
}

// NewRbacMenuModel 返回一个新的 RbacMenuModel 模型实例指针
func NewRbacMenuModel() *MysqlModel {
	return NewMySqlModel().SetModel(&RbacMenuModel{})
}

// TableName 角色与用户对应关系表名称
func (PivotRbacRoleAccountModel) TableName() string {
	return "pivot_rbac_rolse__accounts"
}

// NewPivotRbacRoleAccountModel 返回一个新的 PivotRbacRoleAccountModel 模型实例
func NewPivotRbacRoleAccountModel() *MysqlModel {
	return NewMySqlModel().SetModel(&PivotRbacRoleAccountModel{})
}

// TableName 角色与权限对应关系表名称
func (PivotRbacRoleRbacPermissionModel) TableName() string {
	return "pivot_rbac_roles__rbac_permissions"
}

// NewPivotRbacRoleRbacPermissionModel 返回一个新的 PivotRbacRoleRbacPermissionModel 模型的实例。
func NewPivotRbacRoleRbacPermissionModel() *MysqlModel {
	return NewMySqlModel().SetModel(&PivotRbacRoleRbacPermissionModel{})
}

// TableName 角色与菜单对应关系表名称
func (PivotRbacRoleRbacMenuModel) TableName() string {
	return "pivot_rbac_role__rbac_menus"
}

// NewPivotRbacRoleRbacMenuModel 返回一个新的 PivotRbacRoleRbacMenuModel 模型的实例。
func NewPivotRbacRoleRbacMenuModel() *MysqlModel {
	return NewMySqlModel().SetModel(&PivotRbacRoleRbacMenuModel{})
}
