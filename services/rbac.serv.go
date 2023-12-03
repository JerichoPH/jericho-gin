package services

import "gorm.io/gorm"

type (
	// RbacRoleService 角色服务
	RbacRoleService struct{ BaseService }
)

// NewRbacRoleService 返回一个RbacRoleService的实例。
// 参数baseService是一个BaseService类型的参数，用于传递基本服务。
func NewRbacRoleService(baseService BaseService) *RbacRoleService {
	return &RbacRoleService{BaseService: baseService}
}

// GetListByQuery 根据查询条件获取列表
func (receiver RbacRoleService) GetListByQuery() *gorm.DB {
	return receiver.
		Model.
		SetWheresEqual("name").
		SetWheresExtraHasValues(map[string]func([]string, *gorm.DB) *gorm.DB{
			"name": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("name in (?)", values)
			},
		}).
		GetDb("").
		Table("rbac_roles as rr")
}
