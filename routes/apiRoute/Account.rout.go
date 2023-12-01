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
		r.POST("account", func(ctx *gin.Context) { controllers.NewAccountController().Store(ctx) })

		// 删除
		r.DELETE("account/:uuid", func(ctx *gin.Context) { controllers.NewAccountController().Delete(ctx) })

		// 编辑
		r.PUT("account.update/:uuid/update", func(ctx *gin.Context) { controllers.NewAccountController().Update(ctx) })

		// 详情
		r.GET("account/:uuid", func(ctx *gin.Context) { controllers.NewAccountController().Detail(ctx) })

		// 列表
		r.GET("account", func(ctx *gin.Context) { controllers.NewAccountController().List(ctx) })

		// jquery-dataTable分页列表
		r.GET("account.jdt", func(ctx *gin.Context) { controllers.NewAccountController().ListJdt(ctx) })
	}
}
