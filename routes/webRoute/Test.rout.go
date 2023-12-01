package webRoute

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WsTestRouter struct{}

func (WsTestRouter) Load(engine *gin.Engine) {

	r := engine.Group("test")
	{
		r.GET("ws", func(ctx *gin.Context) {
			engine.LoadHTMLFiles("templates/Test/ws.html")
			ctx.HTML(http.StatusOK, "ws.html", map[string]any{})
		})

		r.GET("sse", func(ctx *gin.Context) {
			engine.LoadHTMLFiles("templates/Test/sse.html")
			ctx.HTML(http.StatusOK, "sse.html", map[string]any{})
		})
	}

}
