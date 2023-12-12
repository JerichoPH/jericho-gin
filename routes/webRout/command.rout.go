package webRout

import (
	"jericho-gin/controllers"

	"github.com/gin-gonic/gin"
)

type CommandRout struct{}

func (CommandRout) Load(engine *gin.Engine) {
	r := engine.Group("command")
	{
		// ExcelHelper类演示
		r.GET("excelHelperDemo", controllers.NewCommandCtrl().ExcelHelperDemo)

		// 初始化数据库
		r.GET("initData", controllers.NewCommandCtrl().InitData)

	}
}
