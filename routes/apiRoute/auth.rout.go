package apiRoute

import (
	"jericho-gin/controllers"

	"github.com/gin-gonic/gin"
)

// AuthRouter 权鉴路由
type AuthRouter struct{}

// NewAuthRouter 构造函数
func NewAuthRouter() AuthRouter {
	return AuthRouter{}
}

// Load 加载路由
func (AuthRouter) Load(engine *gin.Engine) {
	r := engine.Group(
		"api/auth",
		// middlewares.CheckJwt(),
		// middlewares.CheckPermission(),
	)
	{
		// 登陆
		r.POST("login", controllers.NewAuthController().PostLogin)

		// 注册
		r.POST("register", controllers.NewAuthController().PostRegister)
	}
}
