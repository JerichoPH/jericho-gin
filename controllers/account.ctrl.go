package controllers

import (
	"jericho-gin/models"
	"jericho-gin/services"
	"jericho-gin/tools"
	"jericho-gin/wrongs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	// AccountController 用户控制器
	AccountController struct{}
	// AccountStoreForm 用户表单
	AccountStoreForm struct{}
)

// NewAccountController 构造函数
func NewAccountController() *AccountController {
	return &AccountController{}
}

// ShouldBind 表单绑定
func (receiver AccountStoreForm) ShouldBind(ctx *gin.Context) AccountStoreForm {
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}

	return receiver
}

// Store 新建
func (AccountController) Store(ctx *gin.Context) {
	var (
		ret *gorm.DB
		// repeat models.AccountModel
	)

	// 新建
	account := &models.AccountModel{}
	if ret = models.NewMySqlModel().SetModel(models.AccountModel{}).
		GetDb("").
		Create(&account); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Created(map[string]any{"account": account}).ToGinResponse())
}

// Delete 删除
func (AccountController) Delete(ctx *gin.Context) {
	var (
		ret     *gorm.DB
		account models.AccountModel
	)

	// 查询
	ret = models.NewMySqlModel().SetModel(models.AccountModel{}).
		SetWheres(map[string]any{"uuid": ctx.Param("uuid")}).
		GetDb("").
		First(&account)
	wrongs.ThrowWhenIsEmpty(ret, "用户")

	// 删除
	if ret := models.NewMySqlModel().SetModel(models.AccountModel{}).GetDb("").Delete(&account); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Deleted().ToGinResponse())
}

// Update 编辑
func (AccountController) Update(ctx *gin.Context) {
	var (
		ret     *gorm.DB
		account models.AccountModel
		// repeat  models.AccountModel
	)

	// 表单
	// form := new(accountStoreForm).ShouldBind(ctx)

	// 查询
	ret = models.NewMySqlModel().SetModel(models.AccountModel{}).
		SetWheres(map[string]any{"uuid": ctx.Param("uuid")}).
		GetDb("").
		First(&account)
	wrongs.ThrowWhenIsEmpty(ret, "用户")

	// 编辑
	if ret = models.NewMySqlModel().SetModel(models.AccountModel{}).
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		Save(&account); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Updated(map[string]any{"account": account}).ToGinResponse())
}

// Detail 详情
func (AccountController) Detail(ctx *gin.Context) {
	var (
		ret     *gorm.DB
		account models.AccountModel
	)
	ret = models.NewMySqlModel().SetModel(models.AccountModel{}).
		SetWheres(map[string]any{"uuid": ctx.Param("uuid")}).
		GetDb("").
		First(&account)
	wrongs.ThrowWhenIsEmpty(ret, "用户")

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Datum(map[string]any{"account": account}).ToGinResponse())
}

func (AccountController) listByQuery(ctx *gin.Context) *gorm.DB {
	return services.NewAccountService(services.BaseService{Model: models.NewMySqlModel().SetModel(models.AccountModel{}), Ctx: ctx}).GetListByQuery()
}

// List 列表
func (receiver AccountController) List(ctx *gin.Context) {
	var accounts []models.AccountModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForPager(
				receiver.listByQuery(ctx),
				func(db *gorm.DB) map[string]any {
					db.Find(&accounts)
					return map[string]any{"accounts": accounts}
				},
			).ToGinResponse(),
	)
}

// ListJdt jquery-dataTable分页列表
func (receiver AccountController) ListJdt(ctx *gin.Context) {
	var accounts []models.AccountModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForJqueryDataTable(
				receiver.listByQuery(ctx),
				func(db *gorm.DB) map[string]any {
					db.Find(&accounts)
					return map[string]any{"accounts": accounts}
				},
			).ToGinResponse(),
	)
}
