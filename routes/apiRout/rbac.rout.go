package apiRout

import (
	"jericho-gin/controllers"
	"jericho-gin/middlewares"

	"github.com/gin-gonic/gin"
)

// RbacRout 路由
type RbacRout struct{}

// NewRbacRout 构造函数
func NewRbacRout() RbacRout {
	return RbacRout{}
}

// Load 加载路由
func (RbacRout) Load(engine *gin.Engine) {
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
			rbacRoleRouter.POST("", controllers.NewRbacRoleCtrl().Store)
			// 删除
			rbacRoleRouter.DELETE("/:uuid", controllers.NewRbacRoleCtrl().Delete)
			// 编辑
			rbacRoleRouter.PUT("/:uuid", controllers.NewRbacRoleCtrl().Update)
			// 详情
			rbacRoleRouter.GET("/:uuid", controllers.NewRbacRoleCtrl().Detail)
			// 列表
			rbacRoleRouter.GET("", controllers.NewRbacRoleCtrl().List)
			// jquery-dataTable数据列表
			rbacRoleRouter.GET(".jdt", controllers.NewRbacRoleCtrl().ListJdt)
		}

		// 权限
		rbacPermissionRouter := r.Group(
			"permission",
			middlewares.CheckAuth(),
			middlewares.CheckPermission(),
		)
		{
			// 新建
			rbacPermissionRouter.POST("", controllers.NewRbacPermissionCtrl().Store)
			// 删除
			rbacPermissionRouter.DELETE("/:uuid", controllers.NewRbacPermissionCtrl().Delete)
			// 编辑
			rbacPermissionRouter.PUT("/:uuid", controllers.NewRbacPermissionCtrl().Update)
			// 详情
			rbacPermissionRouter.GET("/:uuid", controllers.NewRbacPermissionCtrl().Detail)
			// 列表
			rbacPermissionRouter.GET("", controllers.NewRbacPermissionCtrl().List)
			// jquery-dataTable数据列表
			rbacPermissionRouter.GET(".jdt", controllers.NewRbacPermissionCtrl().ListJdt)
		}

		rbacMenuRouter := r.Group(
			"menu",
			middlewares.CheckAuth(),
			middlewares.CheckPermission(),
		)
		{
			// 新建
			rbacMenuRouter.POST("", controllers.NewRbacMenuCtrl().Store)
			// 删除
			rbacMenuRouter.DELETE("/:uuid", controllers.NewRbacMenuCtrl().Delete)
			// 编辑
			rbacMenuRouter.PUT("/:uuid", controllers.NewRbacMenuCtrl().Update)
			// 详情
			rbacMenuRouter.GET("/:uuid", controllers.NewRbacMenuCtrl().Detail)
			// 列表
			rbacMenuRouter.GET("", controllers.NewRbacMenuCtrl().List)
			// jquery-dataTable数据列表
			rbacMenuRouter.GET(".jdt", controllers.NewRbacMenuCtrl().ListJdt)
		}
	}

}
