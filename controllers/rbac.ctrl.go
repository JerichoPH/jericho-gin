package controllers

import (
	"jericho-gin/models"
	"jericho-gin/services"
	"jericho-gin/tools"
	"jericho-gin/wrongs"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type (
	// RbacRoleController 角色控制器
	RbacRoleController struct{}

	// RbacRoleStoreForm 角色表单
	RbacRoleStoreForm struct {
		Name string `json:"name" binding:"required"`
	}
)

// 角色表单绑定
func (receiver RbacRoleStoreForm) ShouldBind(ctx *gin.Context) RbacRoleStoreForm {
	var err error

	if err = ctx.ShouldBind(&receiver); err != nil {
		if len(receiver.Name) == 0 {
			wrongs.ThrowValidate("角色名称不能为空")
		}
	}

	return receiver
}

// NewRbacRoleController 构造函数
func NewRbacRoleController() *RbacRoleController {
	return &RbacRoleController{}
}

// Store 新建
func (RbacRoleController) Store(ctx *gin.Context) {
	var (
		ret    *gorm.DB
		repeat models.RbacRoleModel
	)

	// 表单
	form := (&RbacRoleStoreForm{}).ShouldBind(ctx)

	// 查重
	ret = models.NewRbacRoleModel().
		GetDb("").
		Where("name = ?", form.Name).
		First(&repeat)
	wrongs.ThrowWhenIsRepeat(ret, "角色名称")

	// 新建
	rbacRole := &models.RbacRoleModel{
		MysqlModel: models.MysqlModel{Uuid: uuid.NewV4().String()},
	}
	if ret = models.NewRbacRoleModel().
		GetDb("").
		Create(&rbacRole); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Created(map[string]any{"rbacRole": rbacRole}).ToGinResponse())
}

// Delete 删除
func (RbacRoleController) Delete(ctx *gin.Context) {
	var (
		ret      *gorm.DB
		rbacRole models.RbacRoleModel
	)

	// 查询
	ret = models.NewRbacRoleModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&rbacRole)
	wrongs.ThrowWhenIsEmpty(ret, "角色")

	// 删除
	if ret := models.NewRbacRoleModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		Delete(&rbacRole); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Deleted().ToGinResponse())
}

// Update 编辑
func (RbacRoleController) Update(ctx *gin.Context) {
	var (
		ret              *gorm.DB
		rbacRole, repeat models.RbacRoleModel
	)

	// 表单
	form := (&RbacRoleStoreForm{}).ShouldBind(ctx)

	// 查重
	ret = models.NewRbacRoleModel().
		GetDb("").
		Where("name = ? and uuid <> ?", form.Name, ctx.Param("uuid")).
		First(&repeat)
	wrongs.ThrowWhenIsRepeat(ret, "角色名称")

	// 查询
	ret = models.NewRbacRoleModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&rbacRole)
	wrongs.ThrowWhenIsEmpty(ret, "角色")

	// 编辑
	rbacRole.Name = form.Name
	if ret = models.NewRbacRoleModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		Save(&rbacRole); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Updated(map[string]any{"rbacRole": rbacRole}).ToGinResponse())
}

// Detail 详情
func (RbacRoleController) Detail(ctx *gin.Context) {
	var (
		ret      *gorm.DB
		rbacRole models.RbacRoleModel
	)
	ret = models.NewRbacRoleModel().
		SetCtx(ctx).
		GetDbUseQuery("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&rbacRole)
	wrongs.ThrowWhenIsEmpty(ret, "角色")

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Datum(map[string]any{"rbacRole": rbacRole}).ToGinResponse())
}

func (RbacRoleController) listUseQuery(ctx *gin.Context) *gorm.DB {
	return services.NewRbacRoleService(services.BaseService{Model: models.NewMySqlModel().SetModel(models.RbacRoleModel{}), Ctx: ctx}).GetListByQuery()
}

// List 列表
func (receiver RbacRoleController) List(ctx *gin.Context) {
	var rbacRoles []*models.RbacRoleModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForPager(
				receiver.listUseQuery(ctx),
				func(db *gorm.DB) map[string]interface{} {
					db.Find(&rbacRoles)
					return map[string]any{"rbacRoles": rbacRoles}
				},
			).
			ToGinResponse(),
	)
}

// ListJdt jquery-dataTable后端分页数据
func (receiver RbacRoleController) ListJdt(ctx *gin.Context) {
	var rbacRoles []*models.RbacRoleModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForJqueryDataTable(
				receiver.listUseQuery(ctx),
				func(db *gorm.DB) map[string]interface{} {
					db.Find(&rbacRoles)
					return map[string]any{"rbacRoles": rbacRoles}
				},
			).
			ToGinResponse(),
	)
}
