package webRoute

import "github.com/gin-gonic/gin"

type RouterHandle struct{}

func (RouterHandle) Register(engine *gin.Engine) {
	HomeRouter{}.Load(engine)           // 欢迎页
	DetectorTabletRouter{}.Load(engine) // 检测台旁边的平板
	CommandRouter{}.Load(engine)        // Command控制台
	WsTestRouter{}.Load(engine)         // web-socket-test
}
