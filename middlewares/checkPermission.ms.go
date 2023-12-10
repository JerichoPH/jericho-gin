package middlewares

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jericho-gin/models"
	"jericho-gin/tools"
	"jericho-gin/wrongs"
)

func CheckPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			account               models.AccountModel
			rbacPermissionUuids   []string
			currentUri            string
			currentRbacPermission *models.RbacPermissionModel
			ret                   *gorm.DB
			yes                   = false
		)
		// 获取当前用户
		account = tools.GetAuth(ctx).(models.AccountModel)
		if !account.BeAdmin {
			rbacPermissionUuids = account.GetPermissionUuids()

			// 获取当前路由
			currentUri = ctx.Request.URL.Path

			// 查询当前路由是否存在权限
			ret = models.NewRbacPermissionModel().GetDb("").Where("uri", currentUri).First(&currentRbacPermission)
			wrongs.ThrowWhenIsEmpty(ret, "当前路由对应权限")

			// 检查当前路由是否合法
			for _, uuid := range rbacPermissionUuids {
				if uuid == currentRbacPermission.Uuid {
					yes = true
					break
				}
			}

			if !yes {
				wrongs.ThrowUnAuth("权限不足")
			}
		}
		
		ctx.Next()
	}
}
