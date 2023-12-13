package controllers

import (
	"jericho-gin/database"
	"jericho-gin/models"
	"jericho-gin/tools"
	"jericho-gin/wrongs"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type (
	// AuthCtrl 权鉴控制器
	AuthCtrl struct{}
	// AuthRegisterForm 注册表单
	AuthRegisterForm struct {
		Username             string `json:"username" binding:"required"`
		Password             string `json:"password" binding:"required"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required"`
		Nickname             string `json:"nickname" binding:"required"`
	}
	// AuthLoginForm 登录表单
	AuthLoginForm struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)

// NewAuthCtrl 构造函数
func NewAuthCtrl() *AuthCtrl {
	return &AuthCtrl{}
}

// ShouldBind 绑定表单（注册）
func (receiver AuthRegisterForm) ShouldBind(ctx *gin.Context) AuthRegisterForm {
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}
	if receiver.Username == "" {
		wrongs.ThrowValidate("账号必填")
	}
	if receiver.Password == "" {
		wrongs.ThrowValidate("密码必填")
	}
	if len(receiver.Password) < 6 || len(receiver.Password) > 18 {
		wrongs.ThrowValidate("密码不可小于6位或大于18位")
	}
	if receiver.Password != receiver.PasswordConfirmation {
		wrongs.ThrowValidate("两次密码输入不一致")
	}

	return receiver
}

// ShouldBind 绑定表单（登陆）
func (receiver AuthLoginForm) ShouldBind(ctx *gin.Context) AuthLoginForm {
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}
	if receiver.Username == "" {
		wrongs.ThrowValidate("账号必填")
	}
	if receiver.Password == "" {
		wrongs.ThrowValidate("密码必填")
	}
	if len(receiver.Password) < 6 || len(receiver.Password) > 18 {
		wrongs.ThrowValidate("密码不可小于6位或大于18位")
	}

	return receiver
}

// PostRegister 注册
func (AuthCtrl) PostRegister(ctx *gin.Context) {
	// 表单验证
	form := AuthRegisterForm{}.ShouldBind(ctx)

	// 检查重复项（用户名）
	var repeat models.AccountModel
	var ret *gorm.DB
	ret = models.NewAccountModel().GetDb("").Where("username", form.Username).First(&repeat)
	wrongs.ThrowWhenRepeat(ret, "用户名")
	ret = models.NewAccountModel().GetDb("").Where("nickname", form.Nickname).First(&repeat)
	wrongs.ThrowWhenRepeat(ret, "昵称")

	// 保存新用户
	account := &models.AccountModel{
		MySqlModel: models.MySqlModel{Uuid: uuid.NewV4().String()},
		Username:   form.Username,
		Password:   tools.GeneratePassword(form.Password),
		Nickname:   form.Nickname,
	}
	if ret = models.NewAccountModel().GetDb("").Create(&account); ret.Error != nil {
		wrongs.ThrowForbidden("创建失败：" + ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("注册成功", ctx).Created(map[string]any{"account": account}).ToGinResponse())
}

// PostLogin 登录
func (AuthCtrl) PostLogin(ctx *gin.Context) {
	// 表单验证
	form := AuthLoginForm{}.ShouldBind(ctx)

	var (
		account models.AccountModel
		ret     *gorm.DB
	)

	// 获取用户
	ret = models.NewAccountModel().GetDb("").Where("username", form.Username).First(&account)
	wrongs.ThrowWhenEmpty(ret, "用户")

	// 验证密码
	if !tools.CheckPassword(form.Password, account.Password) {
		wrongs.ThrowUnAuth("账号或密码错误")
	}

	// 生成Jwt
	if token, err := tools.GenerateJwt(
		account.Id,
		account.Username,
		account.Nickname,
		account.Uuid,
	); err != nil {
		// 生成jwt错误
		wrongs.ThrowForbidden(err.Error())
	} else {
		ctx.JSON(tools.NewCorrectWithGinContext("登陆成功", ctx).Datum(map[string]any{
			"token": token,
			"account": map[string]any{
				"id":       account.Id,
				"username": account.Username,
				"nickname": account.Nickname,
				"uuid":     account.Uuid,
			},
		}).ToGinResponse())
	}
}

// GetMenus 获取当前账号菜单
func (AuthCtrl) GetMenus(ctx *gin.Context) {
	var (
		account       models.AccountModel
		rbacMenus     = make([]*models.RbacMenuModel, 0)
		rbacMenuUuids = make([]string, 0)
	)

	account = tools.GetAuth(ctx).(models.AccountModel)

	if account.BeAdmin {
		models.NewRbacMenuModel().GetDb("").Find(&rbacMenus)
	} else {
		database.
			NewGormLauncher().
			GetConn("").
			Table("rbac_menus m").
			Select("distinct row (m.uuid)").
			Joins("join pivot_rbac_roles__rbac_menus rm on m.uuid = rm.rbac_menu_uuid").
			Joins("join rbac_roles r on rm.rbac_role_uuid = r.uuid").
			Joins("join pivot_rbac_roles__accounts ra on r.uuid = ra.rbac_role_uuid").
			Joins("join accounts a on ra.account_uuid = a.uuid").
			Where("a.account_uuid =?", account.Uuid).
			Find(&rbacMenuUuids)

		models.NewRbacMenuModel().GetDb("").Where("uuid in ?", rbacMenuUuids).Find(&rbacMenus)
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Datum(map[string]any{"rbac_menus": rbacMenus}).ToGinResponse())
}

// PutUpdatePassword 修改密码
func (AuthCtrl) PutUpdatePassword(ctx *gin.Context) {
	var (
		ret            *gorm.DB
		account        *models.AccountModel
		currentAccount models.AccountModel
	)

	currentAccount = tools.GetAuth(ctx).(models.AccountModel)

	form := AccountUpdatePasswordForm{}.ShouldBind(ctx)

	ret = models.NewAccountModel().GetDb("").Where("uuid = ?", currentAccount.Uuid).First(&account)
	wrongs.ThrowWhenEmpty(ret, "用户")

	// 验证密码
	tools.CheckPassword(form.OldPassword, account.Password)

	account.Password = tools.GeneratePassword(form.Password)
	models.NewAccountModel().GetDb("").Where("uuid = ?", ctx.Param("uuid")).Save(&account)

	ctx.JSON(tools.NewCorrectWithGinContext("修改成功", ctx).Blank().ToGinResponse())
}
