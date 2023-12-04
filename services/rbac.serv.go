package services

import (
	"fmt"
	"gorm.io/gorm"
)

type (
	// RbacRoleService 角色服务
	RbacRoleService struct{ BaseService }

	// RbacPermissionService 权限服务
	RbacPermissionService struct{ BaseService }
)

func NewRbacRoleService(baseService BaseService) *RbacRoleService {
	return &RbacRoleService{BaseService: baseService}
}

func (receiver RbacRoleService) GetListByQuery() *gorm.DB {
	return receiver.
		Model.
		SetWheresExtraHasValue(map[string]func(string, *gorm.DB) *gorm.DB{
			"name": func(value string, db *gorm.DB) *gorm.DB {
				return db.Where(fmt.Sprintf("name like %%%s%%", value))
			},
		}).
		SetWheresExtraHasValues(map[string]func([]string, *gorm.DB) *gorm.DB{
			"names[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("name in (?)", values)
			},
		}).
		GetDb("").
		Table("rbac_roles as rr")
}

func NewRbacPermissionService(baseService BaseService) *RbacPermissionService {
	return &RbacPermissionService{BaseService: baseService}
}

func (receiver RbacPermissionService) GetListByQuery() *gorm.DB {
	return receiver.
		Model.
		SetWheresExtraHasValue(map[string]func(string, *gorm.DB) *gorm.DB{
			"name": func(value string, db *gorm.DB) *gorm.DB {
				return db.Where(fmt.Sprintf("name like %%%s%%", value))
			},
			"rbac_role_uuid": func(value string, db *gorm.DB) *gorm.DB {
				return db.Where("rr.uuid =?", value)
			},
		}).
		SetWheresExtraHasValues(map[string]func([]string, *gorm.DB) *gorm.DB{
			"names[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("name in (?)", values)
			},
			"rbac_role_uuids[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("rr.uuid in (?)", values)
			},
		}).
		GetDb("").
		Table("rbac_permissions as rp").
		Joins("join pivot_rbac_roles__rbac_permissions prrrp on rp.uuid = prrrp.rbac_permission_uuid").
		Joins("join rbac_roles rr on prrrp.rbac_role_uuid = rr.uuid")
}
