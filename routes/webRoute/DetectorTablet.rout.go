package webRoute

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DetectorTabletRouter struct{}

func (DetectorTabletRouter) Load(engine *gin.Engine) {
	engine.LoadHTMLGlob("templates/DetectorTablet/index.html")
	//engine.Static("/detectorTablet", "templates/DetectorTablet")
	engine.StaticFS("/detectorTablet", http.Dir("templates/DetectorTablet"))
}
