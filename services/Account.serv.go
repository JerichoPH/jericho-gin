package services

import "gorm.io/gorm"

// AccountService 用户服务
type AccountService struct{ BaseService }

// NewAccountService 构造函数
func NewAccountService(baseService BaseService) *AccountService {
	return &AccountService{BaseService: baseService}
}

// GetListByQuery 通过Query获取列表
func (receiver *AccountService) GetListByQuery() *gorm.DB {
	return receiver.
		Model.
		SetWheresEqual("open_id", "work_area_unique_code", "rank").
		SetWheresFuzzy(map[string]string{ // 模糊查询字段
			"account":  "a.account like ?",
			"nickname": "a.nickname like ?",
		}).
		SetWheresDateBetween("created_at", "updated_at").
		SetWheresExtraHasValues(map[string]func([]string, *gorm.DB) *gorm.DB{
			"ranks[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("a.rank in ?", values)
			},
		}).
		SetCtx(receiver.Ctx).
		GetDbUseQuery("").
		Table("accounts as a")
}
