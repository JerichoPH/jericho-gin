package services

import (
	"jericho-gin/models"

	"github.com/gin-gonic/gin"
)

type (
	BaseService struct {
		Model      *models.MySqlModel
		Ctx        *gin.Context
		DbConnName string
	}
)
