package apiRout

import (
	"github.com/gin-gonic/gin"
)

// RoutHandle 分组路由
type RoutHandle struct{}

// Register 组册路由
func (RoutHandle) Register(engine *gin.Engine) {
	NewTestRout().Load(engine)    // 测试
	NewAuthRout().Load(engine)    // 权鉴
	NewAccountRout().Load(engine) // 用户
	NewRbacRout().Load(engine)    // 权限管理
}
