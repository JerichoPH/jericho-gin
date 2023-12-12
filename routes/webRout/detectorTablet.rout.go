package webRout

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DetectorTabletRout struct{}

func (DetectorTabletRout) Load(engine *gin.Engine) {
	engine.LoadHTMLGlob("templates/DetectorTablet/index.html")
	//engine.Static("/detectorTablet", "templates/DetectorTablet")
	engine.StaticFS("/detectorTablet", http.Dir("templates/DetectorTablet"))
}
