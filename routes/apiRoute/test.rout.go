package apiRoute

import (
	"jericho-gin/controllers"

	"github.com/gin-gonic/gin"
)

// TestRouter 路由
type TestRouter struct{}

// NewTestRouter 构造函数
func NewTestRouter() TestRouter {
	return TestRouter{}
}

// Load 加载路由
func (TestRouter) Load(engine *gin.Engine) {
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
