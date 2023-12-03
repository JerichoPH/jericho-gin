package models

type (
	// AccountModel 用户模型
	AccountModel struct {
		GormModel
		Username  string           `gorm:"unique;type:varchar(64);not null;comment:账号;" json:"username"`
		Password  string           `gorm:"type:varchar(255);not null;comment:密码;" json:"password"`
		Nickname  string           `gorm:"unique;type:varchar(64);not null;comment:昵称;" json:"nickname"`
		RbacRoles []*RbacRoleModel `gorm:"many2many:pivot_rbac_roles__accounts;foreignKey:uuid;joinForeignKey:account_uuid;references:uuid;joinReferences:rbac_rold_uuid;" json:"rbac_roles"`
	}
)

func NewAccountModel() *GormModel {
	return NewGorm().SetModel(AccountModel{})
}

func (AccountModel) TableName() string {
	return "accounts"
}
