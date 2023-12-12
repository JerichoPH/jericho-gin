package apiRout

import (
	"jericho-gin/controllers"

	"github.com/gin-gonic/gin"
)

// TestRout 路由
type TestRout struct{}

// NewTestRout 构造函数
func NewTestRout() TestRout {
	return TestRout{}
}

// Load 加载路由
func (TestRout) Load(engine *gin.Engine) {
	r := engine.Group(
		"api/test",
		// middlewares.CheckAuthorization(),
		// middlewares.CheckPermission(),
	)
	{
		r.Any("sendToWebsocket", controllers.NewTestController().AnySendToWebsocket)
		r.Any("sendToTcpServer", controllers.NewTestController().AnySendToTcpServer)
		r.Any("sendToTcpClient", controllers.NewTestController().AnySendToTcpClient)
		r.Any("sendToKafkaClient", controllers.NewTestController().AnySendToKafkaClient)
	}
}
