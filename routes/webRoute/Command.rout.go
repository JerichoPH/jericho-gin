package webRoute

import (
	"jericho-gin/controllers"

	"github.com/gin-gonic/gin"
)

type CommandRouter struct{}

func (CommandRouter) Load(engine *gin.Engine) {
	r := engine.Group("command")
	{
		// ExcelHelper类演示
		r.GET("excelHelperDemo", controllers.NewCommandController().ExcelHelperDemo)

		// 初始化数据库
		r.GET("initData", controllers.NewCommandController().InitData)

	}
}
