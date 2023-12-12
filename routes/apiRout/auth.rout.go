package apiRout

import (
	"jericho-gin/controllers"
	"jericho-gin/middlewares"

	"github.com/gin-gonic/gin"
)

// AuthRout 权鉴路由
type AuthRout struct{}

// NewAuthRout 构造函数
func NewAuthRout() AuthRout {
	return AuthRout{}
}

// Load 加载路由
func (AuthRout) Load(engine *gin.Engine) {
	// 无需登录
	noNeedLogin := engine.Group(
		"api/auth",
		// middlewares.CheckAuth(),
		// middlewares.CheckPermission(),
	)
	{
		// 登陆
		noNeedLogin.POST("login", controllers.NewAuthCtrl().PostLogin)
		// 注册
		noNeedLogin.POST("register", controllers.NewAuthCtrl().PostRegister)
	}

	// 需要登陆
	needLogin := engine.Group(
		"api/auth",
		middlewares.CheckAuth(),
	)
	{
		// 获取当前用户菜单
		needLogin.GET("menus", controllers.NewAuthCtrl().GetMenus)
		// 修改密码
		needLogin.PUT("updatePassword", controllers.NewAuthCtrl().PutUpdatePassword)
	}
}
