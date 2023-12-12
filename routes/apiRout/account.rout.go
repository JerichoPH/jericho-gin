package apiRout

import (
	"jericho-gin/controllers"
	"jericho-gin/middlewares"

	"github.com/gin-gonic/gin"
)

// AccountRout 用户路由
type AccountRout struct{}

// NewAccountRout 构造函数
func NewAccountRout() AccountRout {
	return AccountRout{}
}

// Load 加载路由
func (AccountRout) Load(engine *gin.Engine) {
	r := engine.Group(
		"api/account",
		middlewares.CheckAuth(),
		middlewares.CheckPermission(),
	)
	{
		// 新建
		r.POST("", controllers.NewAccountCtrl().Store)

		// 删除
		r.DELETE("/:uuid", controllers.NewAccountCtrl().Delete)

		// 编辑
		r.PUT("/:uuid", controllers.NewAccountCtrl().Update)

		// 详情
		r.GET("/:uuid", controllers.NewAccountCtrl().Detail)

		// 列表
		r.GET("", controllers.NewAccountCtrl().List)

		// jquery-dataTable分页列表
		r.GET(".jdt", controllers.NewAccountCtrl().ListJdt)

		// 修改密码
		r.PUT(":uuid/password", controllers.NewAccountCtrl().PutUpdatePassword)
	}
}
