package commands

import (
	"fmt"
	"jericho-gin/database"
	"jericho-gin/models"
	"jericho-gin/tools"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// UpgradeCommand 程序升级
type UpgradeCommand struct{}

// NewUpgradeCommand 构造函数
func NewUpgradeCommand() UpgradeCommand {
	return UpgradeCommand{}
}

func (UpgradeCommand) init() []string {
	var (
		rbacPermissions = make([]*models.RbacPermissionModel, 0)
		rbacMenus       = make([]*models.RbacMenuModel, 0)
		ret             *gorm.DB
	)

	std := tools.NewStdoutHelper("初始化项目")

	std.EchoLineDebug("初始化权限")
	database.ExecSql("truncate table rbac_permissions")
	std.EchoLineSuccess("截断权限表成功")
	// 创建权限
	for _, rbacPermissionDatum := range []map[string]string{
		{"name": "用户列表", "method": "get", "uri": "/api/account"},
		{"name": "用户详情", "method": "get", "uri": "/api/account/:uuid"},
		{"name": "新建用户", "method": "post", "uri": "/api/account"},
		{"name": "编辑用户", "method": "put", "uri": "/api/account/:uuid"},
		{"name": "删除用户", "method": "delete", "uri": "/api/account/:uuid"},
		{"name": "角色列表", "method": "get", "uri": "/api/rbac/role"},
		{"name": "角色详情", "method": "get", "uri": "/api/rbac/role/:uuid"},
		{"name": "新建角色", "method": "post", "uri": "/api/rbac/role"},
		{"name": "编辑角色", "method": "put", "uri": "/api/rbac/role/:uuid"},
		{"name": "删除角色", "method": "delete", "uri": "/api/rbac/role/:uuid"},
		{"name": "权限列表", "method": "get", "uri": "/api/rbac/permission"},
		{"name": "权限详情", "method": "get", "uri": "/api/rbac/permission/:uuid"},
		{"name": "新建权限", "method": "post", "uri": "/api/rbac/permission"},
		{"name": "编辑权限", "method": "put", "uri": "/api/rbac/permission/:uuid"},
		{"name": "删除权限", "method": "delete", "uri": "/api/rbac/permission/:uuid"},
		{"name": "菜单列表", "method": "get", "uri": "/api/rbac/menu"},
		{"name": "菜单详情", "method": "get", "uri": "/api/rbac/menu/:uuid"},
		{"name": "新建菜单", "method": "post", "uri": "/api/rbac/menu"},
		{"name": "编辑菜单", "method": "put", "uri": "/api/rbac/menu/:uuid"},
		{"name": "删除菜单", "method": "delete", "uri": "/api/rbac/menu/:uuid"},
	} {
		rbacPermissions = append(rbacPermissions, &models.RbacPermissionModel{
			MySqlModel:  models.MySqlModel{Uuid: uuid.NewV4().String(), BeEnable: true},
			Name:        rbacPermissionDatum["name"],
			Uri:         rbacPermissionDatum["uri"],
			Method:      rbacPermissionDatum["method"],
			Description: nil,
		})
	}
	if ret = models.NewRbacPermissionModel().GetDb("").Create(&rbacPermissions); ret.Error != nil {
		std.EchoLineWrong(fmt.Sprintf("错误：%v", ret.Error.Error()))
	}
	std.EchoLineSuccess("成功")

	std.EchoLineDebug("初始化菜单")
	database.ExecSql("truncate table rbac_menus")
	std.EchoLineSuccess("截断菜单表成功")
	for _, rbacMenuDatum := range []map[string]string{
		{"name": "用户列表", "uri": "/account"},
		{"name": "角色列表", "uri": "/rbac/role"},
		{"name": "权限列表", "uri": "/rbac/permission"},
		{"name": "权限列表", "uri": "/rbac/permission"},
	} {
		rbacMenus = append(rbacMenus, &models.RbacMenuModel{
			Name: rbacMenuDatum["name"],
			Uri:  rbacMenuDatum["uri"],
		})
	}
	if ret = models.NewRbacMenuModel().GetDb("").Create(&rbacMenus); ret.Error != nil {
		std.EchoLineWrong(fmt.Sprintf("错误：%v", ret.Error.Error()))
	}
	std.EchoLineSuccess("成功")

	return []string{}
}

// Do 执行命令
func (receiver UpgradeCommand) Do(params []string) []string {
	switch params[0] {
	case "init":
		return receiver.init()
	}

	return []string{"执行完成"}
}
