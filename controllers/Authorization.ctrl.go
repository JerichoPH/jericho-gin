package controllers

import (
	"jericho-go/models"
	"jericho-go/tools"
	"jericho-go/wrongs"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	// AuthorizationController 权鉴控制器
	AuthorizationController struct{}
	// authorizationRegisterForm 注册表单
	authorizationRegisterForm struct {
		Username             string `json:"username" binding:"required"`
		Password             string `json:"password" binding:"required"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required"`
		Nickname             string `json:"nickname" binding:"required"`
	}
	// authorizationLoginForm 登录表单
	authorizationLoginForm struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)

// NewAuthorizationController 构造函数
func NewAuthorizationController() *AuthorizationController {
	return &AuthorizationController{}
}

// ShouldBind 绑定表单（注册）
func (receiver authorizationRegisterForm) ShouldBind(ctx *gin.Context) authorizationRegisterForm {
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
func (receiver authorizationLoginForm) ShouldBind(ctx *gin.Context) authorizationLoginForm {
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

// Register 注册
func (AuthorizationController) Register(ctx *gin.Context) {
	// 表单验证
	form := authorizationRegisterForm{}.ShouldBind(ctx)

	// 检查重复项（用户名）
	var repeat models.AccountModel
	var ret *gorm.DB
	ret = models.NewAccountModel().GetDb("").Where("username", form.Username).First(&repeat)
	wrongs.ThrowWhenIsRepeat(ret, "用户名")
	ret = models.NewAccountModel().GetDb("").Where("nickname", form.Nickname).First(&repeat)
	wrongs.ThrowWhenIsRepeat(ret, "昵称")

	// 密码加密
	bytes, _ := bcrypt.GenerateFromPassword([]byte(form.Password), 14)

	// 保存新用户
	account := &models.AccountModel{
		GormModel: models.GormModel{Uuid: uuid.NewV4().String()},
		Username:  form.Username,
		Password:  string(bytes),
		Nickname:  form.Nickname,
	}
	if ret = models.NewAccountModel().GetDb("").Create(&account); ret.Error != nil {
		wrongs.ThrowForbidden("创建失败：" + ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("注册成功", ctx).Created(map[string]any{"account": account}).ToGinResponse())
}

// Login 登录
func (AuthorizationController) Login(ctx *gin.Context) {
	// 表单验证
	form := authorizationLoginForm{}.ShouldBind(ctx)

	var (
		account models.AccountModel
		ret     *gorm.DB
	)

	// 获取用户
	ret = models.NewAccountModel().GetDb("").Where("username", form.Username).First(&account)
	wrongs.ThrowWhenIsEmpty(ret, "用户")

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(form.Password)); err != nil {
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

// GetMenus 获取当前用户菜单
// func (AuthorizationController) GetMenus(ctx *gin.Context) {
//	var ret *gorm.GetDb
//	if accountUuid, exists := ctx.Get(tools.ACCOUNT_OPEN_ID); !exists {
//		wrongs.ThrowUnLogin("用户未登录")
//	} else {
//		// 获取当前用户信息
//		var account models.AccountModel
//		ret = models.NewGorm().SetModel(models.AccountModel{}).
//			SetWheres(map[string]any{"uuid": accountUuid}).
//			SetPreloads("RbacRoles", "RbacRoles.Menus").
//			GetDb("",nil).
//			FindOneUseQuery(&account)
//		if !wrongs.ThrowWhenIsEmpty(ret, "") {
//			wrongs.ThrowUnLogin("当前令牌指向用户不存在")
//		}
//
//		var menus []models.MenuModel
//		models.NewGorm().SetModel(models.MenuModel{}).
//			GetDb("",nil).
//			Joins("join pivot_rbac_role_and_menus prram on menus.uuid = prram.menu_uuid").
//			Joins("join rbac_roles r on prram.rbac_role_uuid = r.uuid").
//			Joins("join pivot_rbac_role_and_accounts prraa on r.uuid = prraa.rbac_role_uuid").
//			Joins("join accounts a on prraa.account_uuid = a.uuid").
//			Where("a.uuid = ?", account.GormModel.Uuid).
//			Where("menus.deleted_at is null").
//			Where("menus.parent_uuid = ''").
//			Order("menus.sort asc").
//			Order("menus.id asc").
//			Preload("Subs").
//			Find(&menus)
//
//		ctx.JSON(tools.CorrectInit("", ctx).Datum(map[string]any{"menus": menus}))
//	}
// }
