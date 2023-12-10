package apiRoute

import (
	"jericho-gin/controllers"
	"jericho-gin/middlewares"

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
		"api/account",
		middlewares.CheckAuth(),
		middlewares.CheckPermission(),
	)
	{
		// 新建
		r.POST("", controllers.NewAccountController().Store)

		// 删除
		r.DELETE("/:uuid", controllers.NewAccountController().Delete)

		// 编辑
		r.PUT("/:uuid", controllers.NewAccountController().Update)

		// 详情
		r.GET("/:uuid", controllers.NewAccountController().Detail)

		// 列表
		r.GET("", controllers.NewAccountController().List)

		// jquery-dataTable分页列表
		r.GET(".jdt", controllers.NewAccountController().ListJdt)

		// 修改密码
		r.PUT(":uuid/password", controllers.NewAccountController().PutUpdatePassword)
	}
}
