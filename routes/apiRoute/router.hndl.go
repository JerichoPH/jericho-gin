package apiRoute

import (
	"github.com/gin-gonic/gin"
)

// RouterHandle 分组路由
type RouterHandle struct{}

// Register 组册路由
func (RouterHandle) Register(engine *gin.Engine) {
	NewTestRouter().Load(engine)    // 测试
	NewAuthRouter().Load(engine)    // 权鉴
	NewAccountRouter().Load(engine) // 用户
	NewRbacRouter().Load(engine)    // 权限管理
}