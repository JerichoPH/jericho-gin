package webRout

import (
	"fmt"
	"jericho-gin/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeRout struct{}

func (HomeRout) Load(engine *gin.Engine) {
	r := engine.Group("")
	{
		r.GET("", func(ctx *gin.Context) {
			engine.LoadHTMLGlob("templates/Home/*")
			version := settings.NewSetting().App.Section("app").Key("version").MustString("")
			subVersion := settings.NewSetting().App.Section("app").Key("sub-version").MustString("")
			if subVersion != "" {
				subVersion = fmt.Sprintf("(%s)", subVersion)
			}
			ctx.HTML(http.StatusOK, "index.html", map[string]string{
				"version":    version,
				"subVersion": subVersion,
			})
		})
	}
}
