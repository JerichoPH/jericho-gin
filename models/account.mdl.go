package models

import (
	"jericho-gin/database"
	"strings"
)

type (
	// AccountModel 用户模型
	AccountModel struct {
		MySqlModel
		Username  string           `gorm:"unique;type:varchar(64);not null;comment:账号;" json:"username"`
		Password  string           `gorm:"type:varchar(255);not null;comment:密码;" json:"-"`
		Nickname  string           `gorm:"unique;type:varchar(64);not null;comment:昵称;" json:"nickname"`
		BeAdmin   bool             `gorm:"type:boolean;not null;default:0;comment:是否是管理员" json:"be_admin"`
		RbacRoles []*RbacRoleModel `gorm:"many2many:pivot_rbac_roles__accounts;foreignKey:uuid;joinForeignKey:account_uuid;references:uuid;joinReferences:rbac_role_uuid;" json:"rbac_roles"`
	}
)

func (AccountModel) TableName() string {
	return "accounts"
}

func NewAccountModel() *MySqlModel {
	return NewMySqlModel().SetModel(AccountModel{})
}

// GetPermissionUuids 获取当前用户所有权限
func (receiver AccountModel) GetPermissionUuids() (rbacPermissionUuids []string) {
	database.NewGormLauncher().
		GetConn("").
		Table("accounts as a").
		Select(strings.Join([]string{
			"DISTINCT rp.uuid",
		}, ",")).
		Joins("join pivot_rbac_roles__accounts prra on a.uuid = prra.account_uuid").
		Joins("join pivot_rbac_roles__rbac_permissions prrrp on prra.rbac_role_uuid = prrrp.rbac_role_uuid").
		Joins("join rbac_permissions rp on prrrp.rbac_permission_uuid = rp.uuid").
		Where("a.uuid = ?", receiver.Uuid).
		Pluck("uuid", &rbacPermissionUuids)

	return
}
