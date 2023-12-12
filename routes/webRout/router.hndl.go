package webRout

import "github.com/gin-gonic/gin"

type RoutHandle struct{}

func (RoutHandle) Register(engine *gin.Engine) {
	HomeRout{}.Load(engine)           // 欢迎页
	DetectorTabletRout{}.Load(engine) // 检测台旁边的平板
	CommandRout{}.Load(engine)        // Command控制台
	WsTestRout{}.Load(engine)         // web-socket-test
}
