package controllers

import (
	"fmt"
	"jericho-gin/models"
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
		Name string `json:"name"`
	}

	// RbacPermissionController 权限控制器
	RbacPermissionController struct{}

	// RbacPermissionStoreForm 权限表单
	RbacPermissionStoreForm struct {
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Uri           string   `json:"uri"`
		RbacRoleUuids []string `json:"rbac_role_uuids"`
		rbacRoles     []*models.RbacRoleModel
	}

	// RbacMenuController 菜单控制器
	RbacMenuController struct{}

	// RbacMenuStoreForm 菜单表单
	RbacMenuStoreForm struct {
		Name          string `json:"name"`
		SubTitle      string `json:"sub_title"`
		Description   string `json:"description"`
		Uri           string `json:"uri"`
		ParentUuid    string `json:"parent_uuid"`
		parentMenu    *models.RbacMenuModel
		RbacRoleUuids []string `json:"rbac_role_uuids"`
		rbacRoles     []*models.RbacRoleModel
	}
)

// ShouldBind 角色表单绑定
func (receiver RbacRoleStoreForm) ShouldBind(ctx *gin.Context) RbacRoleStoreForm {
	var err error

	if ctx.Request.Method == "PUT" && ctx.Param("uuid") == "" {
		wrongs.ThrowValidate("角色编号不能为空")
	}
	if err = ctx.ShouldBind(&receiver); err != nil {
		if len(receiver.Name) == 0 {
			wrongs.ThrowValidate("角色名称不能为空")
		}
	}

	return receiver
}

// ShouldBind 权限表单绑定
func (receiver RbacPermissionStoreForm) ShouldBind(ctx *gin.Context) RbacPermissionStoreForm {
	if ctx.Request.Method == "PUT" && ctx.Param("uuid") == "" {
		wrongs.ThrowValidate("权限编号不能为空")
	}
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}
	if receiver.Name == "" {
		wrongs.ThrowValidate("权限名称必填")
	}
	if receiver.Uri == "" {
		wrongs.ThrowValidate("权限路由必填")
	}
	if len(receiver.RbacRoleUuids) > 0 {
		models.NewRbacRoleModel().GetDb("").Where("uuid in ?", receiver.RbacRoleUuids).Find(&receiver.rbacRoles)
	}

	return receiver
}

// ShouldBind 菜单表单绑定
func (receiver RbacMenuStoreForm) ShouldBind(ctx *gin.Context) RbacMenuStoreForm {
	var ret *gorm.DB

	if ctx.Request.Method == "PUT" && ctx.Param("uuid") == "" {
		wrongs.ThrowValidate("菜单编号不能为空")
	}
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}
	if receiver.Name == "" {
		wrongs.ThrowValidate("菜单名称必填")
	}
	if receiver.ParentUuid != "" {
		ret = models.NewRbacMenuModel().GetDb("").Where("uuid = ?", receiver.ParentUuid).First(&receiver.parentMenu)
		wrongs.ThrowWhenEmpty(ret, fmt.Sprintf("父级菜单（%s）", receiver.ParentUuid))
	}
	if len(receiver.RbacRoleUuids) > 0 {
		models.NewRbacRoleModel().GetDb("").Where("uuid in ?", receiver.RbacRoleUuids).Find(&receiver.rbacRoles)
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
	wrongs.ThrowWhenRepeat(ret, "角色名称")

	// 新建
	rbacRole := &models.RbacRoleModel{
		MySqlModel: models.MySqlModel{Uuid: uuid.NewV4().String()},
		Name:       form.Name,
	}
	if ret = models.NewRbacRoleModel().
		GetDb("").
		Create(&rbacRole); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Created(map[string]any{"rbac_role": rbacRole}).ToGinResponse())
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
	wrongs.ThrowWhenEmpty(ret, "角色")

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
	wrongs.ThrowWhenRepeat(ret, "角色名称")

	// 查询
	ret = models.NewRbacRoleModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&rbacRole)
	wrongs.ThrowWhenEmpty(ret, "角色")

	// 编辑
	rbacRole.Name = form.Name
	if ret = models.NewRbacRoleModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		Save(&rbacRole); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Updated(map[string]any{"rbac_role": rbacRole}).ToGinResponse())
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
	wrongs.ThrowWhenEmpty(ret, "角色")

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Datum(map[string]any{"rbac_role": rbacRole}).ToGinResponse())
}

// List 列表
func (receiver RbacRoleController) List(ctx *gin.Context) {
	var rbacRoles []*models.RbacRoleModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForPager(
				models.RbacRoleModel{}.GetListByQuery(ctx),
				func(db *gorm.DB) map[string]interface{} {
					db.Find(&rbacRoles)
					return map[string]any{"rbac_roles": rbacRoles}
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
				models.RbacRoleModel{}.GetListByQuery(ctx),
				func(db *gorm.DB) map[string]interface{} {
					db.Find(&rbacRoles)
					return map[string]any{"rbac_roles": rbacRoles}
				},
			).
			ToGinResponse(),
	)
}

// NewRbacPermissionController 构造函数
func NewRbacPermissionController() *RbacPermissionController {
	return &RbacPermissionController{}
}

// Store 新建
func (RbacPermissionController) Store(ctx *gin.Context) {
	var (
		ret    *gorm.DB
		repeat models.RbacPermissionModel
	)

	// 表单
	form := RbacPermissionStoreForm{}.ShouldBind(ctx)

	// 查重
	ret = models.NewRbacPermissionModel().
		GetDb("").
		Where("name = ?", form.Name).
		First(&repeat)
	wrongs.ThrowWhenRepeat(ret, "权限名称")

	// 新建
	rbacPermission := &models.RbacPermissionModel{
		MySqlModel:  models.MySqlModel{Uuid: uuid.NewV4().String()},
		Name:        form.Name,
		Uri:         form.Uri,
		Description: &form.Description,
	}
	if ret = models.NewRbacPermissionModel().
		GetDb("").
		Create(&rbacPermission); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	// 绑定角色与权限
	models.PivotRbacRoleRbacPermissionModel{}.BindRbacRoles(rbacPermission, form.rbacRoles)

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Created(map[string]any{"rbac_permission": rbacPermission}).ToGinResponse())
}

// Delete 删除
func (RbacPermissionController) Delete(ctx *gin.Context) {
	var (
		ret            *gorm.DB
		rbacPermission models.RbacPermissionModel
	)

	// 查询
	ret = models.NewRbacPermissionModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&rbacPermission)
	wrongs.ThrowWhenEmpty(ret, "权限")

	// 删除
	if ret := models.NewRbacPermissionModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		Delete(&rbacPermission); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Deleted().ToGinResponse())
}

// Update 编辑
func (RbacPermissionController) Update(ctx *gin.Context) {
	var (
		ret                    *gorm.DB
		rbacPermission, repeat *models.RbacPermissionModel
	)

	// 表单
	form := RbacPermissionStoreForm{}.ShouldBind(ctx)

	// 查重
	ret = models.NewRbacPermissionModel().
		GetDb("").
		Where("name = ? and uuid <> ?", form.Name, ctx.Param("uuid")).
		First(&repeat)
	wrongs.ThrowWhenRepeat(ret, fmt.Sprintf("权限名称 %s %s", ctx.Param("uuid"), repeat.Uuid))

	// 查询
	ret = models.NewRbacPermissionModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&rbacPermission)
	wrongs.ThrowWhenEmpty(ret, "权限")

	// 编辑
	rbacPermission.Name = form.Name
	rbacPermission.Description = &form.Description
	rbacPermission.Uri = form.Uri
	if ret = models.NewRbacPermissionModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		Save(&rbacPermission); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	// 绑定角色与权限
	models.PivotRbacRoleRbacPermissionModel{}.BindRbacRoles(rbacPermission, form.rbacRoles)

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Updated(map[string]any{"rbac_permission": rbacPermission}).ToGinResponse())
}

// Detail 详情
func (RbacPermissionController) Detail(ctx *gin.Context) {
	var (
		ret            *gorm.DB
		rbacPermission models.RbacPermissionModel
	)
	ret = models.NewRbacPermissionModel().
		SetCtx(ctx).
		GetDbUseQuery("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&rbacPermission)
	wrongs.ThrowWhenEmpty(ret, "权限")

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Datum(map[string]any{"rbac_permission": rbacPermission}).ToGinResponse())
}

// List 列表
func (receiver RbacPermissionController) List(ctx *gin.Context) {
	var rbacPermissions []*models.RbacPermissionModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForPager(
				models.RbacPermissionModel{}.GetListByQuery(ctx),
				func(db *gorm.DB) map[string]interface{} {
					db.Find(&rbacPermissions)
					return map[string]any{"rbac_permissions": rbacPermissions}
				},
			).
			ToGinResponse(),
	)
}

// ListJdt jquery-dataTable后端分页数据
func (receiver RbacPermissionController) ListJdt(ctx *gin.Context) {
	var rbacPermissions []*models.RbacPermissionModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForJqueryDataTable(
				models.RbacPermissionModel{}.GetListByQuery(ctx),
				func(db *gorm.DB) map[string]interface{} {
					db.Find(&rbacPermissions)
					return map[string]any{"rbac_permissions": rbacPermissions}
				},
			).
			ToGinResponse(),
	)
}

// NewRbacMenuController 构造函数
func NewRbacMenuController() *RbacMenuController {
	return &RbacMenuController{}
}

// Store 新建
func (RbacMenuController) Store(ctx *gin.Context) {
	var (
		ret    *gorm.DB
		repeat *models.RbacMenuModel
	)

	// 表单
	form := RbacMenuStoreForm{}.ShouldBind(ctx)

	// 查重
	ret = models.NewRbacMenuModel().
		GetDb("").
		Where("name = ?", form.Name).
		Where("parent_uuid = ?", form.ParentUuid).
		First(&repeat)
	wrongs.ThrowWhenRepeat(ret, "菜单名称")

	// 新建
	rbacMenu := &models.RbacMenuModel{
		MySqlModel:  models.MySqlModel{Uuid: uuid.NewV4().String()},
		Name:        form.Name,
		SubTitle:    form.SubTitle,
		Description: &form.Description,
		Uri:         form.Uri,
		ParentUuid:  form.ParentUuid,
	}
	if ret = models.NewRbacMenuModel().
		GetDb("").
		Create(&rbacMenu); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	// 绑定角色与菜单
	models.PivotRbacRoleRbacMenuModel{}.BindRbacRoles(rbacMenu, form.rbacRoles)

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Created(map[string]any{"rbac_menu": rbacMenu}).ToGinResponse())
}

// Delete 删除
func (RbacMenuController) Delete(ctx *gin.Context) {
	var (
		ret      *gorm.DB
		rbacMenu *models.RbacMenuModel
		subs     = make(map[string]map[string]string)
	)

	// 查询
	ret = models.NewRbacMenuModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&rbacMenu)
	wrongs.ThrowWhenEmpty(ret, "菜单")
	// 查询该菜单下是否存在子集
	subs = rbacMenu.GetSubUuidsByParentUuid(rbacMenu.Uuid)
	if len(subs) > 0 {
		wrongs.ThrowForbidden("该菜单下存在子集，请先删除子集")
	}

	// 删除
	if ret := models.NewRbacMenuModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		Delete(&rbacMenu); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Deleted().ToGinResponse())
}

// Update 编辑
func (RbacMenuController) Update(ctx *gin.Context) {
	var (
		ret              *gorm.DB
		rbacMenu, repeat *models.RbacMenuModel
		subs             map[string]map[string]string
	)

	// 表单
	form := RbacMenuStoreForm{}.ShouldBind(ctx)

	// 查重
	ret = models.NewRbacMenuModel().
		GetDb("").
		Where("name = ? and uuid <> ?", form.Name, ctx.Param("uuid")).
		Where("uuid <> ?", ctx.Param("uuid")).
		Where("parent_uuid <> ?", form.ParentUuid).Debug().
		First(&repeat)
	wrongs.ThrowWhenRepeat(ret, "菜单名称")

	// 查询
	ret = models.NewRbacMenuModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&rbacMenu)
	wrongs.ThrowWhenEmpty(ret, "菜单")
	if form.ParentUuid != "" {
		if rbacMenu.ParentUuid == form.ParentUuid {
			wrongs.ThrowValidate("父级菜单不能是自己")
		}
	}
	// 查询所有子菜单
	subs = models.RbacMenuModel{}.GetSubUuidsByParentUuid(rbacMenu.Uuid)
	if len(subs) > 0 && form.ParentUuid != "" {
		if sub, exist := subs[form.ParentUuid]; exist {
			wrongs.ThrowValidate("「%s」是「%s」的父级，不能将子集绑定为子集的父级", form.Name, sub["name"])
		}
	}

	// 编辑
	rbacMenu.Name = form.Name
	rbacMenu.SubTitle = form.SubTitle
	rbacMenu.Description = &form.Description
	rbacMenu.Uri = form.Uri
	rbacMenu.ParentUuid = form.ParentUuid
	if ret = models.NewRbacMenuModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		Save(&rbacMenu); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	// 绑定角色与菜单
	models.PivotRbacRoleRbacMenuModel{}.BindRbacRoles(rbacMenu, form.rbacRoles)

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Updated(map[string]any{"rbac_menu": rbacMenu}).ToGinResponse())
}

// Detail 详情
func (RbacMenuController) Detail(ctx *gin.Context) {
	var (
		ret      *gorm.DB
		rbacMenu models.RbacMenuModel
	)
	ret = models.NewRbacMenuModel().
		SetCtx(ctx).
		GetDbUseQuery("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&rbacMenu)
	wrongs.ThrowWhenEmpty(ret, "菜单")

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Datum(map[string]any{"rbac_menu": rbacMenu}).ToGinResponse())
}

// List 列表
func (receiver RbacMenuController) List(ctx *gin.Context) {
	var rbacMenus []*models.RbacMenuModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForPager(
				models.RbacMenuModel{}.GetListByQuery(ctx),
				func(db *gorm.DB) map[string]interface{} {
					db.Find(&rbacMenus)
					return map[string]any{"rbac_menus": rbacMenus}
				},
			).
			ToGinResponse(),
	)
}

// ListJdt jquery-dataTable后端分页数据
func (receiver RbacMenuController) ListJdt(ctx *gin.Context) {
	var rbacMenus []*models.RbacMenuModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForJqueryDataTable(
				models.RbacMenuModel{}.GetListByQuery(ctx),
				func(db *gorm.DB) map[string]interface{} {
					db.Find(&rbacMenus)
					return map[string]any{"rbac_menus": rbacMenus}
				},
			).
			ToGinResponse(),
	)
}
