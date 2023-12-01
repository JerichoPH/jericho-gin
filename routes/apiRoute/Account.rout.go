package apiRoute

import (
	"jericho-go/controllers"
	"jericho-go/middlewares"

	"github.com/gin-gonic/gin"
)

// AccountRouter 用户路由
type AccountRouter struct{}

// NewAccountRouter 构造函数
func NewAccountRouter() AccountRouter {
	return AccountRouter{}
}

// Load 加载路由
func (AccountRouter) Load(engine *gin.Engine) {
	r := engine.Group(
		"api",
		middlewares.CheckAuthorization(),
		// middlewares.CheckPermission(),
	)
	{
		// 新建
		r.POST("account", controllers.NewAccountController().Store)

		// 删除
		r.DELETE("account/:uuid", controllers.NewAccountController().Delete)

		// 编辑
		r.PUT("account.update/:uuid/update", controllers.NewAccountController().Update)

		// 详情
		r.GET("account/:uuid", controllers.NewAccountController().Detail)

		// 列表
		r.GET("account", controllers.NewAccountController().List)

		// jquery-dataTable分页列表
		r.GET("account.jdt", controllers.NewAccountController().ListJdt)
	}
}
