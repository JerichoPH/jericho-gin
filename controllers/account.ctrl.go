package controllers

import (
	"jericho-gin/models"
	"jericho-gin/tools"
	"jericho-gin/wrongs"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type (
	// AccountCtrl 用户控制器
	AccountCtrl struct{}
	// AccountStoreForm 新建用户表单
	AccountStoreForm struct {
		Username             string   `json:"username"`
		Nickname             string   `json:"nickname"`
		Password             string   `json:"password"`
		PasswordConfirmation string   `json:"password_confirmation"`
		RbacRoleUuids        []string `json:"rbac_role_uuids"`
		rbacRoles            []*models.RbacRoleModel
	}
	// AccountUpdateForm 编辑用户表单
	AccountUpdateForm struct {
		Username      string   `json:"username"`
		Nickname      string   `json:"nickname"`
		RbacRoleUuids []string `json:"rbac_role_uuids"`
		rbacRoles     []*models.RbacRoleModel
	}

	// AccountUpdatePasswordForm 修改密码表单
	AccountUpdatePasswordForm struct {
		OldPassword          string `json:"old_password"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
	}
)

// NewAccountCtrl 构造函数
func NewAccountCtrl() *AccountCtrl {
	return &AccountCtrl{}
}

// ShouldBind 新建用户表单绑定
func (receiver AccountStoreForm) ShouldBind(ctx *gin.Context) AccountStoreForm {
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}
	if len(receiver.Username) < 2 {
		wrongs.ThrowValidate("账号不能小于2位")
	}
	if len(receiver.Nickname) < 2 {
		wrongs.ThrowValidate("昵称不能小于2位")
	}
	if len(receiver.Password) < 6 {
		wrongs.ThrowValidate("密码不能小于6位")
	}
	if receiver.Password != receiver.PasswordConfirmation {
		wrongs.ThrowValidate("两次密码不一致")
	}
	if len(receiver.RbacRoleUuids) > 0 {
		models.NewRbacRoleModel().GetDb("").Where("uuid in ?", receiver.RbacRoleUuids).Find(&receiver.rbacRoles)
	}

	return receiver
}

// ShouldBind 编辑用户表单绑定
func (receiver AccountUpdateForm) ShouldBind(ctx *gin.Context) AccountUpdateForm {
	if ctx.Param("uuid") == "" {
		wrongs.ThrowValidate("用户编号不能为空")
	}
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}
	if len(receiver.Username) < 2 {
		wrongs.ThrowValidate("账号不能小于2位")
	}
	if len(receiver.Nickname) < 2 {
		wrongs.ThrowValidate("昵称不能小于2位")
	}
	if len(receiver.RbacRoleUuids) > 0 {
		models.NewRbacRoleModel().GetDb("").Where("uuid in ?", receiver.RbacRoleUuids).Find(&receiver.rbacRoles)
	}

	return receiver
}

// ShouldBind 修改用户密码表单绑定
func (receiver AccountUpdatePasswordForm) ShouldBind(ctx *gin.Context) AccountUpdatePasswordForm {
	if ctx.Param("uuid") == "" {
		wrongs.ThrowValidate("用户编号不能为空")
	}

	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}
	if len(receiver.OldPassword) < 6 {
		wrongs.ThrowValidate("原密码不能小于6位")
	}
	if len(receiver.Password) < 6 {
		wrongs.ThrowValidate("新密码不能小于6位")
	}
	if receiver.Password != receiver.PasswordConfirmation {
		wrongs.ThrowValidate("两次密码不一致")
	}

	return receiver
}

// Store 新建
func (AccountCtrl) Store(ctx *gin.Context) {
	var (
		ret    *gorm.DB
		repeat *models.AccountModel
	)

	// 表单绑定
	form := AccountStoreForm{}.ShouldBind(ctx)

	// 查重
	ret = models.NewAccountModel().GetDb("").Where("username = ?", form.Username).First(&repeat)
	wrongs.ThrowWhenRepeat(ret, "账号被占用")
	ret = models.NewAccountModel().GetDb("").Where("nickname = ?", form.Nickname).First(&repeat)
	wrongs.ThrowWhenRepeat(ret, "昵称被占用")

	// 新建
	account := &models.AccountModel{
		MySqlModel: models.MySqlModel{Uuid: uuid.NewV4().String()},
		Username:   form.Username,
		Nickname:   form.Nickname,
		Password:   tools.GeneratePassword(form.Password),
		BeAdmin:    false,
	}
	if ret = models.
		NewAccountModel().
		GetDb("").
		Create(&account); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	// 用户绑定角色
	models.AccountModel{}.BindRbacRoles(account, form.rbacRoles)

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Created(map[string]any{"account": account}).ToGinResponse())
}

// Delete 删除
func (AccountCtrl) Delete(ctx *gin.Context) {
	var (
		ret     *gorm.DB
		account models.AccountModel
	)

	// 查询
	ret = models.
		NewAccountModel().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&account)
	wrongs.ThrowWhenEmpty(ret, "用户")

	// 删除
	if ret := models.NewAccountModel().GetDb("").Where("uuid = ?", account.Uuid).Delete(&account); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Deleted().ToGinResponse())
}

// Update 编辑
func (AccountCtrl) Update(ctx *gin.Context) {
	var (
		ret     *gorm.DB
		account *models.AccountModel
		// repeat  *models.AccountModel
	)

	// 表单
	form := new(AccountUpdateForm).ShouldBind(ctx)

	// 查重
	ret = models.NewAccountModel().GetDb("").Where("username = ?", form.Username).Where("uuid <> ?", ctx.Param("uuid")).First(nil)
	wrongs.ThrowWhenRepeat(ret, "账号被占用")
	ret = models.NewAccountModel().GetDb("").Where("nickname = ?", form.Nickname).Where("uuid <> ?", ctx.Param("uuid")).First(nil)
	wrongs.ThrowWhenRepeat(ret, "昵称被占用")

	// 查询
	ret = models.NewAccountModel().GetDb("").Where("uuid", ctx.Param("uuid")).First(&account)
	wrongs.ThrowWhenEmpty(ret, "用户")

	// 编辑
	account.Username = form.Username
	account.Nickname = form.Nickname
	if ret = models.NewAccountModel().GetDb("").Where("uuid = ?", ctx.Param("uuid")).Save(&account); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	// 用户绑定角色
	models.AccountModel{}.BindRbacRoles(account, form.rbacRoles)

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Updated(map[string]any{"account": account}).ToGinResponse())
}

// Detail 详情
func (AccountCtrl) Detail(ctx *gin.Context) {
	var (
		ret     *gorm.DB
		account models.AccountModel
	)
	ret = models.NewAccountModel().SetCtx(ctx).GetDbUseQuery("").Where("uuid", ctx.Param("uuid")).First(&account)
	wrongs.ThrowWhenEmpty(ret, "用户")

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Datum(map[string]any{"account": account}).ToGinResponse())
}

// List 列表
func (receiver AccountCtrl) List(ctx *gin.Context) {
	var accounts []models.AccountModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForPager(
				models.AccountModel{}.GetListByQuery(ctx),
				func(db *gorm.DB) map[string]any {
					db.Find(&accounts)
					return map[string]any{"accounts": accounts}
				},
			).ToGinResponse(),
	)
}

// ListJdt jquery-dataTable分页列表
func (receiver AccountCtrl) ListJdt(ctx *gin.Context) {
	var accounts []models.AccountModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForJqueryDataTable(
				models.AccountModel{}.GetListByQuery(ctx),
				func(db *gorm.DB) map[string]any {
					db.Find(&accounts)
					return map[string]any{"accounts": accounts}
				},
			).ToGinResponse(),
	)
}

// PutUpdatePassword 修改密码
func (recevier AccountCtrl) PutUpdatePassword(ctx *gin.Context) {
	var (
		ret     *gorm.DB
		account *models.AccountModel
	)

	form := AccountUpdatePasswordForm{}.ShouldBind(ctx)

	ret = models.NewAccountModel().GetDb("").Where("uuid = ?", ctx.Param("uuid")).First(&account)
	wrongs.ThrowWhenEmpty(ret, "用户")

	// 验证密码
	tools.CheckPassword(form.OldPassword, account.Password)

	account.Password = tools.GeneratePassword(form.Password)
	models.NewAccountModel().GetDb("").Where("uuid = ?", ctx.Param("uuid")).Save(&account)

	ctx.JSON(tools.NewCorrectWithGinContext("修改成功", ctx).Blank().ToGinResponse())
}
