package apiRoute

import (
	"github.com/gin-gonic/gin"
	"jericho-go/controllers"
)

// AuthorizationRouter 权鉴路由
type AuthorizationRouter struct{}

// NewAuthorizationRouter 构造函数
func NewAuthorizationRouter() AuthorizationRouter {
	return AuthorizationRouter{}
}

// Load 加载路由
func (AuthorizationRouter) Load(engine *gin.Engine) {
	r := engine.Group(
		"api/authorization",
		// middlewares.CheckJwt(),
		// middlewares.CheckPermission(),
	)
	{
		// 登陆
		r.POST("login", controllers.NewAuthorizationController().Login)

		// 注册
		r.POST("register", controllers.NewAuthorizationController().Register)
	}
}
