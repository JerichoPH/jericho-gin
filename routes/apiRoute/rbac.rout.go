package apiRoute

import (
	"jericho-gin/controllers"
	"jericho-gin/middlewares"

	"github.com/gin-gonic/gin"
)

// RbacRouter 路由
type RbacRouter struct{}

// NewRbacRouter 构造函数
func NewRbacRouter() RbacRouter {
	return RbacRouter{}
}

// Load 加载路由
func (RbacRouter) Load(engine *gin.Engine) {
	r := engine.Group("/api/rbac")
	{
		// 角色
		rbacRoleRouter := r.Group(
			"role",
			middlewares.CheckAuth(),
			middlewares.CheckPermission(),
		)
		{
			// 新建
			rbacRoleRouter.POST("", controllers.NewRbacRoleController().Store)
			// 删除
			rbacRoleRouter.DELETE("/:uuid", controllers.NewRbacRoleController().Delete)
			// 编辑
			rbacRoleRouter.PUT("/:uuid", controllers.NewRbacRoleController().Update)
			// 详情
			rbacRoleRouter.GET("/:uuid", controllers.NewRbacRoleController().Detail)
			// 列表
			rbacRoleRouter.GET("", controllers.NewRbacRoleController().List)
			// jquery-dataTable数据列表
			rbacRoleRouter.GET(".jdt", controllers.NewRbacRoleController().ListJdt)
		}

		// 权限
		rbacPermissionRouter := r.Group(
			"permission",
			middlewares.CheckAuth(),
			middlewares.CheckPermission(),
		)
		{
			// 新建
			rbacPermissionRouter.POST("", controllers.NewRbacPermissionController().Store)
			// 删除
			rbacPermissionRouter.DELETE("/:uuid", controllers.NewRbacPermissionController().Delete)
			// 编辑
			rbacPermissionRouter.PUT("/:uuid", controllers.NewRbacPermissionController().Update)
			// 详情
			rbacPermissionRouter.GET("/:uuid", controllers.NewRbacPermissionController().Detail)
			// 列表
			rbacPermissionRouter.GET("", controllers.NewRbacPermissionController().List)
			// jquery-dataTable数据列表
			rbacPermissionRouter.GET(".jdt", controllers.NewRbacPermissionController().ListJdt)
		}

		rbacMenuRouter := r.Group(
			"menu",
			middlewares.CheckAuth(),
			middlewares.CheckPermission(),
		)
		{
			// 新建
			rbacMenuRouter.POST("", controllers.NewRbacMenuController().Store)
			// 删除
			rbacMenuRouter.DELETE("/:uuid", controllers.NewRbacMenuController().Delete)
			// 编辑
			rbacMenuRouter.PUT("/:uuid", controllers.NewRbacMenuController().Update)
			// 详情
			rbacMenuRouter.GET("/:uuid", controllers.NewRbacMenuController().Detail)
			// 列表
			rbacMenuRouter.GET("", controllers.NewRbacMenuController().List)
			// jquery-dataTable数据列表
			rbacMenuRouter.GET(".jdt", controllers.NewRbacMenuController().ListJdt)
		}
	}

}
