package services

import (
	"jericho-go/models"

	"github.com/gin-gonic/gin"
)

type (
	BaseService struct {
		Model      *models.GormModel
		Ctx        *gin.Context
		DbConnName string
	}
)
