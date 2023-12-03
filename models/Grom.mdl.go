package models

import (
	"fmt"
	"jericho-gin/database"
	"jericho-gin/wrongs"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GormModel 基础模型
type GormModel struct {
	Id                       uint64         `gorm:"primaryKey;autoIncrement;type:bigint unsigned auto_increment;" json:"id"`
	CreatedAt                time.Time      `gorm:"<-:create;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;" json:"created_at,omitempty"`
	UpdatedAt                time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间;" json:"updated_at,omitempty"`
	DeletedAt                gorm.DeletedAt `gorm:"index;type:datetime" json:"deleted_at"`
	Uuid                     string         `gorm:"unique;type:varchar(36);not null;comment:uuid;" json:"uuid"`
	Sort                     int64          `gorm:"type:bigint;default:0;comment:排序;" json:"sort"`
	ctx                      *gin.Context
	preloads                 []string
	selects                  []string
	omits                    []string
	whereFields              []string
	notWhereFields           []string
	orWhereFields            []string
	ignoreFields             []string
	distinctFieldNames       []string
	wheres                   map[string]any
	notWheres                map[string]any
	orWheres                 map[string]any
	wheresExtra              map[string]func(string, *gorm.DB) *gorm.DB
	wheresExtraExist         map[string]func(string, *gorm.DB) *gorm.DB
	wheresExtraExists        map[string]func([]string, *gorm.DB) *gorm.DB
	scopes                   []func(*gorm.DB) *gorm.DB
	whereDateBetween         []string
	whereDatetimeBetween     []string
	whereBetween             []string
	whereIntBetween          []string
	whereFloatBetween        []string
	whereIn                  map[string]string
	whereFuzzyQueryCondition map[string]string
	model                    interface{}
}

// JoinString 生成join语句
func JoinString(slaveTableName, slaveShortTableName, slaveFieldName, masterFieldName string) string {
	return fmt.Sprintf("join %s %s on %s = %s", slaveTableName, slaveShortTableName, slaveFieldName, masterFieldName)
}

// LeftJoinString 生成left join语句
func LeftJoinString(slaveTableName, slaveShortTableName, slaveFieldName, masterFieldName string) string {
	return "left " + JoinString(slaveTableName, slaveShortTableName, slaveFieldName, masterFieldName)
}

// RightJoinString 生成right join语句
func RightJoinString(slaveTableName, slaveShortTableName, slaveFieldName, masterFieldName string) string {
	return "right " + JoinString(slaveTableName, slaveShortTableName, slaveFieldName, masterFieldName)
}

// NewGorm 构造函数
func NewGorm() *GormModel {
	return &GormModel{}
}

// demoFindOne 获取单条数据演示
func (receiver *GormModel) demoFindOne() {
	var b GormModel
	ret := receiver.
		SetModel(GormModel{}).
		SetWheres(map[string]any{}).
		SetNotWheres(map[string]any{}).
		GetDb("").
		First(b)
	wrongs.ThrowWhenIsEmpty(ret, "XX")
}

// demoFind 获取多条数据演示
func (receiver *GormModel) demoFind() {
	var b GormModel
	var ctx *gin.Context
	receiver.
		SetModel(GormModel{}).
		SetWheresEqual("uuid").
		SetWheresHasKey(map[string]func(string, *gorm.DB) *gorm.DB{
			"name": func(fieldName string, db *gorm.DB) *gorm.DB {
				if queryValue, exist := ctx.GetQuery(fieldName); exist {
					db = db.Where(fieldName+" like ?", fmt.Sprintf("%%%s%%", queryValue))
				}
				return db
			},
		}).
		SetWheresDateBetween("created_at", "updated_at").
		SetCtx(ctx).
		GetDbUseQuery("连接池名称：空字符串为默认").
		Find(&b)
}

// ScopeBeEnableTrue 启用（查询域）
func (*GormModel) ScopeBeEnableTrue(db *gorm.DB) *gorm.DB {
	return db.Where("be_enable is true")
}

// ScopeBeEnableFalse 不启用（查询域）
func (*GormModel) ScopeBeEnableFalse(db *gorm.DB) *gorm.DB {
	return db.Where("be_enable is false")
}

// SetCtx 设置Context
func (receiver *GormModel) SetCtx(ctx *gin.Context) *GormModel {
	receiver.ctx = ctx
	return receiver
}

// SetModel 设置使用的模型
func (receiver *GormModel) SetModel(model interface{}) *GormModel {
	receiver.model = model
	return receiver
}

// SetDistinct 设置不重复字段
func (receiver *GormModel) SetDistinct(distinctFieldNames ...string) *GormModel {
	receiver.distinctFieldNames = distinctFieldNames

	return receiver
}

// SetPreloads 设置Preloads
func (receiver *GormModel) SetPreloads(preloads ...string) *GormModel {
	receiver.preloads = preloads
	return receiver
}

// SetPreloadsByDefault 设置Preloads为默认
func (receiver *GormModel) SetPreloadsByDefault() *GormModel {
	receiver.preloads = []string{clause.Associations}
	return receiver
}

// SetSelects 设置Selects
func (receiver *GormModel) SetSelects(selects ...string) *GormModel {
	receiver.selects = selects
	return receiver
}

// SetOmits 设置Omits
func (receiver *GormModel) SetOmits(omits ...string) *GormModel {
	receiver.omits = omits
	return receiver
}

// SetWheresEqual 设置WhereFields
func (receiver *GormModel) SetWheresEqual(whereFields ...string) *GormModel {
	receiver.whereFields = whereFields
	return receiver
}

// SetNotWhereFields 设置NotWhereFields
func (receiver *GormModel) SetNotWhereFields(notWhereFields ...string) *GormModel {
	receiver.notWhereFields = notWhereFields
	return receiver
}

// SetOrWhereFields 设置OrWhere字段
func (receiver *GormModel) SetOrWhereFields(orWhereFields ...string) *GormModel {
	receiver.orWhereFields = orWhereFields
	return receiver
}

// SetIgnoreFields 设置IgnoreFields
func (receiver *GormModel) SetIgnoreFields(ignoreFields ...string) *GormModel {
	receiver.ignoreFields = ignoreFields
	return receiver
}

// SetWheres 通过Map设置Wheres
func (receiver *GormModel) SetWheres(wheres map[string]any) *GormModel {
	receiver.wheres = wheres
	return receiver
}

// SetNotWheres 设置NotWheres
func (receiver *GormModel) SetNotWheres(notWheres map[string]any) *GormModel {
	receiver.notWheres = notWheres
	return receiver
}

// SetOrWheres 设置Or条件
func (receiver *GormModel) SetOrWheres(orWheres map[string]any) *GormModel {
	receiver.orWheres = orWheres
	return receiver
}

// SetScopes 设置Scopes
func (receiver *GormModel) SetScopes(scopes ...func(*gorm.DB) *gorm.DB) *GormModel {
	receiver.scopes = scopes
	return receiver
}

// SetWheresHasKey 设置额外搜索条件字段
func (receiver *GormModel) SetWheresHasKey(wheresExtra map[string]func(string, *gorm.DB) *gorm.DB) *GormModel {
	receiver.wheresExtra = wheresExtra

	return receiver
}

// SetWheresExtraHasValue 当Query参数存在时设置额外搜索条件（单个条件）
func (receiver *GormModel) SetWheresExtraHasValue(wheresExtraExist map[string]func(string, *gorm.DB) *gorm.DB) *GormModel {
	receiver.wheresExtraExist = wheresExtraExist

	return receiver
}

// SetWheresExtraHasValues 当Query参数存在时设置额外搜索条件（多个条件）
func (receiver *GormModel) SetWheresExtraHasValues(wheresExtraExists map[string]func([]string, *gorm.DB) *gorm.DB) *GormModel {
	receiver.wheresExtraExists = wheresExtraExists

	return receiver
}

// SetWheresDateBetween 设置需要检查的日期范围字段
func (receiver *GormModel) SetWheresDateBetween(fieldNames ...string) *GormModel {
	receiver.whereDateBetween = fieldNames

	return receiver
}

// SetWheresBetween 设置需要范围查询的条件
func (receiver *GormModel) SetWheresBetween(fieldNames ...string) *GormModel {
	receiver.whereBetween = fieldNames

	return receiver
}

// SetWheresIntBetween 设置需要范围查询的条件（数字）
func (receiver *GormModel) SetWheresIntBetween(fieldNames ...string) *GormModel {
	receiver.whereIntBetween = fieldNames

	return receiver
}

// SetWheresFloatBetween 设置需要范围查询的条件（浮点）
func (receiver *GormModel) SetWheresFloatBetween(fieldNames ...string) *GormModel {
	receiver.whereFloatBetween = fieldNames

	return receiver
}

// SetWheresDatetimeBetween 设置需要检查的日期时间范围字段
func (receiver *GormModel) SetWheresDatetimeBetween(fieldNames ...string) *GormModel {
	receiver.whereDatetimeBetween = fieldNames

	return receiver
}

// SetWheresIn 设置自定义查询条件(in)
func (receiver *GormModel) SetWheresIn(condition map[string]string) *GormModel {
	receiver.whereIn = condition

	return receiver
}

// SetWheresFuzzy 设置模糊查询条件
func (receiver *GormModel) SetWheresFuzzy(condition map[string]string) *GormModel {
	receiver.whereFuzzyQueryCondition = condition

	return receiver
}

// BeforeCreate 插入数据前
func (receiver *GormModel) BeforeCreate(db *gorm.DB) (err error) {
	receiver.Uuid = uuid.NewV4().String()
	receiver.CreatedAt = time.Now()
	receiver.UpdatedAt = time.Now()
	return
}

// BeforeSave 修改数据前
func (receiver *GormModel) BeforeSave(db *gorm.DB) (err error) {
	receiver.UpdatedAt = time.Now()
	return
}

// GetDb 初始化
func (receiver *GormModel) GetDb(dbConnName string) (query *gorm.DB) {
	query = database.NewGormLauncher().GetConn(dbConnName)

	query = query.Where(receiver.wheres).Not(receiver.notWheres)

	if receiver.model != nil {
		query = query.Model(&receiver.model)
	}

	// 设置scopes
	if len(receiver.scopes) > 0 {
		query = query.Scopes(receiver.scopes...)
	}

	// 拼接preloads关系
	if len(receiver.preloads) > 0 {
		for _, v := range receiver.preloads {
			query = query.Preload(v)
		}
	}

	// 拼接distinct
	if len(receiver.distinctFieldNames) > 0 {
		query = query.Distinct(receiver.distinctFieldNames)
	}

	// 拼接selects字段
	if len(receiver.selects) > 0 {
		query = query.Select(receiver.selects)
	}

	// 拼接omits字段
	if len(receiver.omits) > 0 {
		query = query.Omit(receiver.omits...)
	}

	return query
}

// GetDbUseQuery 根据Query参数初始化
func (receiver *GormModel) GetDbUseQuery(dbConnName string) *gorm.DB {
	dbSession := receiver.GetDb(dbConnName)

	wheres := make(map[string]any)
	notWheres := make(map[string]any)
	orWheres := make(map[string]any)

	// 自动化处理⬇

	// 拼接需要跳过的字段
	ignoreFields := make(map[string]int32)
	if len(receiver.ignoreFields) > 0 {
		for _, v := range receiver.ignoreFields {
			ignoreFields[v] = 1
		}
	}

	// 拼接Where条件
	for _, whereField := range receiver.whereFields {
		if _, ok := ignoreFields[whereField]; !ok {
			if val, ok := receiver.ctx.GetQuery(whereField); ok {
				wheres[whereField] = val
			}
		}
	}

	// 拼接WhereMap条件
	if m, exist := receiver.ctx.GetQueryMap("__wheres__"); exist {
		for whereMapField, whereMapValue := range m {
			if _, ok := ignoreFields[whereMapField]; !ok {
				if whereMapValue != "" {
					wheres[whereMapField] = whereMapValue
				}
			}
		}
	}

	// 拼接NotWhere条件
	for _, notWhereField := range receiver.notWhereFields {
		if _, ok := ignoreFields[notWhereField]; !ok {
			if notWhereValue, ok := receiver.ctx.GetQuery(notWhereField); ok {
				notWheres[notWhereField] = notWhereValue
			}
		}
	}

	// 拼接NotWhereMap条件
	if m, exist := receiver.ctx.GetQueryMap("__not_wheres__"); exist {
		for notWhereMapField, notWhereMapValue := range m {
			if _, ok := ignoreFields[notWhereMapField]; !ok {
				if notWhereMapValue != "" {
					notWheres[notWhereMapField] = notWhereMapValue
				}
			}
		}
	}

	// 拼接OrWhere条件
	for _, orWhereField := range receiver.orWhereFields {
		if _, ok := ignoreFields[orWhereField]; !ok {
			if val, ok := receiver.ctx.GetQuery(orWhereField); ok {
				orWheres[orWhereField] = val
			}
		}
	}

	// 拼接OrWhereMap条件
	if m, exist := receiver.ctx.GetQueryMap("__or_wheres__"); exist {
		for orWhereMapField, orWhereMapValue := range m {
			if _, ok := ignoreFields[orWhereMapField]; !ok {
				if orWhereMapValue != "" {
					notWheres[orWhereMapField] = orWhereMapValue
				}
			}
		}
	}
	dbSession = dbSession.Where(wheres).Not(notWheres).Or(orWheres)

	// 拼接自主条件（like）
	if m, exist := receiver.ctx.GetQueryMap("__like__"); exist {
		for field, value := range m {
			if _, ok := ignoreFields[field]; !ok {
				if value != "" {
					dbSession = dbSession.Where(fmt.Sprintf("%s like ? %%%s%%", field, value))
				}
			}
		}
	}

	// 拼接自主条件（左like）
	if m, exist := receiver.ctx.GetQueryMap("__like_l__"); exist {
		for field, value := range m {
			if _, ok := ignoreFields[field]; !ok {
				if value != "" {
					dbSession = dbSession.Where(fmt.Sprintf("%s like ? %%%s", field, value))
				}
			}
		}
	}

	// 拼接自主条件（右like）
	if m, exist := receiver.ctx.GetQueryMap("__like_r__"); exist {
		for field, value := range m {
			if _, ok := ignoreFields[field]; !ok {
				if value != "" {
					dbSession = dbSession.Where(fmt.Sprintf("%s like ? %s%%", field, value))
				}
			}
		}
	}

	// 拼接自主条件（等于）
	if m, exist := receiver.ctx.GetQueryMap("__eq__"); exist {
		for field, value := range m {
			if _, ok := ignoreFields[field]; !ok {
				if value != "" {
					dbSession = dbSession.Where(fmt.Sprintf("%s = ?", field), value)
				}
			}
		}
	}

	// 拼接自主条件（等于空字符串）
	if fields, exist := receiver.ctx.GetQueryArray("__eq_empty__[]"); exist {
		for _, field := range fields {
			if _, ok := ignoreFields[field]; !ok {
				dbSession = dbSession.Where(fmt.Sprintf("%s = ''", field))
			}
		}
	}

	// 拼接自主条件（不等于空字符串）
	if fields, exist := receiver.ctx.GetQueryArray("__neq_empty__[]"); exist {
		for _, field := range fields {
			if _, ok := ignoreFields[field]; !ok {
				dbSession = dbSession.Where(fmt.Sprintf("%s <> ''", field))
			}
		}
	}

	// 拼接自主条件（等于空）
	if fields, exist := receiver.ctx.GetQueryArray("__eq_null__[]"); exist {
		for _, field := range fields {
			if _, ok := ignoreFields[field]; !ok {
				dbSession = dbSession.Where(fmt.Sprintf("%s is null", field))
			}
		}
	}

	// 拼接自主条件（非空）
	if fields, exist := receiver.ctx.GetQueryArray("__neq_null__[]"); exist {
		for _, field := range fields {
			if _, ok := ignoreFields[field]; !ok {
				dbSession = dbSession.Where(fmt.Sprintf("%s is not null", field))
			}
		}
	}

	// 拼接自主条件（小于）
	if m, exist := receiver.ctx.GetQueryMap("__lt__"); exist {
		for field, value := range m {
			if _, ok := ignoreFields[field]; !ok {
				if value != "" {
					dbSession = dbSession.Where(fmt.Sprintf("%s < ?", field), value)
				}
			}
		}
	}

	// 拼接自主条件（大于）
	if m, exist := receiver.ctx.GetQueryMap("__gt__"); exist {
		for field, value := range m {
			if _, ok := ignoreFields[field]; !ok {
				if value != "" {
					dbSession = dbSession.Where(fmt.Sprintf("%s > ?", field), value)
				}
			}
		}
	}

	// 拼接自主条件（不等于）
	if m, exist := receiver.ctx.GetQueryMap("__neq__"); exist {
		for field, value := range m {
			if _, ok := ignoreFields[field]; !ok {
				if value != "" {
					dbSession = dbSession.Where(fmt.Sprintf("%s != ?", field), value)
				}
			}
		}
	}

	// 拼接自主条件（小于等于）
	if m, exist := receiver.ctx.GetQueryMap("__elt__"); exist {
		for field, value := range m {
			if _, ok := ignoreFields[field]; !ok {
				if value != "" {
					dbSession = dbSession.Where(fmt.Sprintf("%s <= ?", field), value)
				}
			}
		}
	}

	// 拼接自主条件（大于等于）
	if m, exist := receiver.ctx.GetQueryMap("__egt__"); exist {
		for field, value := range m {
			if _, ok := ignoreFields[field]; !ok {
				if value != "" {
					dbSession = dbSession.Where(fmt.Sprintf("%s >= ?", field), value)
				}
			}
		}
	}

	// 拼接自主条件（大于等于）
	if m, exist := receiver.ctx.GetQueryMap("__egt__"); exist {
		for field, value := range m {
			if _, ok := ignoreFields[field]; !ok {
				if value != "" {
					dbSession = dbSession.Where(fmt.Sprintf("%s >= ?", field), value)
				}
			}
		}
	}

	// 拼接自主条件（in）
	if m, exist := receiver.ctx.GetQueryMap("__in__"); exist {
		for field, values := range m {
			if _, ok := ignoreFields[field]; !ok {
				if len(strings.Split(values, ",")) > 0 {
					dbSession = dbSession.Where(fmt.Sprintf("%s in ?", field), values)
				}
			}
		}
	}

	// 拼接自主条件（between）
	if m, exist := receiver.ctx.GetQueryMap("__between__"); exist {
		for field, values := range m {
			if _, ok := ignoreFields[field]; !ok {
				values := strings.Split(values, ",")
				if len(values) > 0 {
					dbSession = dbSession.Where(fmt.Sprintf("%s between ? and ?", field), values[0], values[1])
				}
			}
		}
	}

	// 自主拼接条件（between int）
	if m, exist := receiver.ctx.GetQueryMap("__between_in__"); exist {
		for field, values := range m {
			if _, ok := ignoreFields[field]; !ok {
				values := strings.Split(values, ",")
				int1, err := strconv.Atoi(values[0])
				if err != nil {
					wrongs.ThrowValidate("%s 必须是数字，（当前值：%v）", field, values[0])
				}
				int2, err := strconv.Atoi(values[1])
				if err != nil {
					wrongs.ThrowValidate("%s 必须是数字，（当前值：%v）", field, values[1])
				}
				if len(values) > 0 {
					dbSession = dbSession.Where(fmt.Sprintf("%s between ? and ?", field), int1, int2)
				}
			}
		}
	}

	// 自主拼接条件（between float）
	if m, exist := receiver.ctx.GetQueryMap("__between_float__"); exist {
		for field, values := range m {
			if _, ok := ignoreFields[field]; !ok {
				values := strings.Split(values, ",")
				float1, err := strconv.ParseFloat(values[0], 64)
				if err != nil {
					wrongs.ThrowValidate("%s 必须是小数，（当前值：%v）", field, values[0])
				}
				float2, err := strconv.ParseFloat(values[1], 64)
				if err != nil {
					wrongs.ThrowValidate("%s 必须是小数，（当前值：%v）", field, values[1])
				}
				if len(values) > 0 {
					dbSession = dbSession.Where(fmt.Sprintf("%s between ? and ?", field), float1, float2)
				}
			}
		}
	}

	// 自主拼接preload
	if a, exist := receiver.ctx.GetQueryArray("__preloads__[]"); exist {
		for _, d := range a {
			dbSession = dbSession.Preload(d)
		}
	}

	// 手动处理⬇

	// 拼接额外搜索条件
	for fieldName, v := range receiver.wheresExtra {
		if _, ok := ignoreFields[fieldName]; !ok {
			if _, ok := receiver.ctx.GetQuery(fieldName); ok {
				dbSession = v(fieldName, dbSession)
			}
		}
	}

	// 拼接额外搜索条件（判断值是否存在，单个条件）
	for fieldName, v := range receiver.wheresExtraExist {
		if _, ok := ignoreFields[fieldName]; !ok {
			if value := receiver.ctx.Query(fieldName); value != "" {
				dbSession = v(value, dbSession)
			}
		}
	}

	// 拼接额外搜索条件（判断值是否存在，多个条件）
	for fieldName, v := range receiver.wheresExtraExists {
		if _, ok := ignoreFields[fieldName]; !ok {
			if values := receiver.ctx.QueryArray(fieldName); len(values) > 0 {
				dbSession = v(values, dbSession)
			}
		}
	}

	// 拼接日期范围查询
	if len(receiver.whereDateBetween) > 0 {
		for _, fieldName := range receiver.whereDateBetween {
			if value, exist := receiver.ctx.GetQuery(fieldName); exist {
				times := strings.Split(value, "~")
				originalAt, finishedAt := times[0], times[1]
				dbSession = dbSession.Where(fieldName+" between ? and ?", originalAt+" 00:00:00", finishedAt+" 23:59:59")
			}
		}
	}

	// 拼接日期时间范围查询
	if len(receiver.whereDatetimeBetween) > 0 {
		for _, fieldName := range receiver.whereDateBetween {
			if value, exist := receiver.ctx.GetQuery(fieldName); exist {
				times := strings.Split(value, "~")
				originalAt, finishedAt := times[0], times[1]
				dbSession = dbSession.Where(fieldName+" between ? and ?", originalAt, finishedAt)
			}
		}
	}

	// 拼接between范围查询条件
	if len(receiver.whereBetween) > 0 {
		for _, fieldName := range receiver.whereBetween {
			if values, exist := receiver.ctx.GetQuery(fieldName); exist {
				values := strings.Split(values, ",")
				if len(values) == 2 {
					dbSession = dbSession.Where(fmt.Sprintf("%s between ? and ?", fieldName), values[0], values[1])
				}
			}
		}
	}

	// 拼接int between范围查询条件
	if len(receiver.whereIntBetween) > 0 {
		for _, fieldName := range receiver.whereIntBetween {
			if values, exist := receiver.ctx.GetQuery(fieldName); exist {
				values := strings.Split(values, ",")
				if len(values) == 2 {
					int1, err := strconv.Atoi(values[0])
					if err != nil {
						wrongs.ThrowValidate("%s 必须是数字（当前值：%v）", fieldName, values[0])
					}
					int2, err := strconv.Atoi(values[1])
					if err != nil {
						wrongs.ThrowValidate("%s 必须是数字（当前值：%v）", fieldName, values[1])
					}

					dbSession = dbSession.Where(fmt.Sprintf("%s between ? and ?", fieldName), int1, int2)
				}
			}
		}
	}

	// 拼接flat between范围查询条件
	if len(receiver.whereFloatBetween) > 0 {
		for _, fieldName := range receiver.whereFloatBetween {
			if values, exist := receiver.ctx.GetQuery(fieldName); exist {
				values := strings.Split(values, ",")
				if len(values) == 2 {
					float1, err := strconv.ParseFloat(values[0], 64)
					if err != nil {
						wrongs.ThrowValidate("%s 必须是小数（当前值：%v）", fieldName, values[0])
					}
					float2, err := strconv.ParseFloat(values[0], 64)
					if err != nil {
						wrongs.ThrowValidate("%s 必须是小数（当前值：%v）", fieldName, values[1])
					}

					dbSession = dbSession.Where(fmt.Sprintf("%s between ? and ?", fieldName), float1, float2)
				}
			}
		}
	}

	// 拼接in查询条件
	if len(receiver.whereIn) > 0 {
		for field, condition := range receiver.whereIn {
			if values, exist := receiver.ctx.GetQueryArray(field); exist {
				if len(values) > 0 {
					dbSession = dbSession.Where(condition, values)
				}
			}
		}
	}

	// 拼接模糊查询条件
	if len(receiver.whereFuzzyQueryCondition) > 0 {
		for field, condition := range receiver.whereFuzzyQueryCondition {
			if value, exist := receiver.ctx.GetQuery(field); exist {
				if value != "" {
					dbSession = dbSession.Where(condition, "%"+value+"%")
				}
			}
		}
	}

	// 排序
	if order, ok := receiver.ctx.GetQuery("__order__"); ok {
		dbSession.Order(order)
	}

	// 指定排序
	if orderField, ok := receiver.ctx.GetQueryArray("__order_field__"); ok {
		if len(orderField) > 2 {
			orderFieldV := "'" + strings.Join(orderField[1:], "','") + "'"
			dbSession.Order(fmt.Sprintf("field (%s,%s)", orderField[0], orderFieldV))
		}
	}

	return dbSession
}
