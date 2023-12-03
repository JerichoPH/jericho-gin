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
	roleRouter := engine.Group(
		"api/rbac",
		middlewares.CheckAuth(),
	// middlewares.CheckPermission(),
	)
	{
		// 新建
		roleRouter.POST("role", controllers.NewRbacRoleController().Store)
		// 删除
		roleRouter.DELETE("role/:uuid", controllers.NewRbacRoleController().Delete)
		// 编辑
		roleRouter.PUT("role/:uuid", controllers.NewRbacRoleController().Update)
		// 详情
		roleRouter.GET("role/:uuid", controllers.NewRbacRoleController().Detail)
		// 列表
		roleRouter.GET("role", controllers.NewRbacRoleController().List)
		// jquery-dataTable数据列表
		roleRouter.GET("role.jdt", controllers.NewRbacRoleController().ListJdt)
	}
}
