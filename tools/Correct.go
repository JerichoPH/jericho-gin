package tools

import (
	"jericho-gin/wrongs"
	"math"
	"strconv"
	"sync"

	uuid "github.com/satori/go.uuid"

	"jericho-gin/types"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type (
	Correct struct {
		msg          string
		messageId    string
		traceId      string
		context      *gin.Context
		businessType string
	}

	Response struct {
		DB        *gorm.DB
		PagerFunc func(db *gorm.DB) map[string]any
		Msg       string
		Ctx       *gin.Context
		Context   map[string]any
		Method    types.ResponseMethod
	}
)

var ins *Correct
var once sync.Once

// AsGin 作为Gin格式返回
func (receiver Response) AsGin() (int, map[string]any) {
	switch receiver.Method {
	case types.RESPONSE_METHOD_OK:
		return NewCorrectWithGinContext(receiver.Msg, receiver.Ctx).Ok().ToGinResponse()
	case types.RESPONSE_METHOD_DATUM:
		return NewCorrectWithGinContext(receiver.Msg, receiver.Ctx).Datum(receiver.Context).ToGinResponse()
	case types.RESPONSE_METHOD_CREATED:
		return NewCorrectWithGinContext(receiver.Msg, receiver.Ctx).Created(receiver.Context).ToGinResponse()
	case types.RESPONSE_METHOD_UPDATED:
		return NewCorrectWithGinContext(receiver.Msg, receiver.Ctx).Updated(receiver.Context).ToGinResponse()
	case types.RESPONSE_METHOD_DELETED:
		return NewCorrectWithGinContext(receiver.Msg, receiver.Ctx).Deleted().ToGinResponse()
	case types.RESPONSE_METHOD_PAGER:
		return NewCorrectWithGinContext(receiver.Msg, receiver.Ctx).DataForPager(receiver.DB, receiver.PagerFunc).ToGinResponse()
	case types.RESPONSE_METHOD_JQUERY_DATA_PAGER:
		return NewCorrectWithGinContext(receiver.Msg, receiver.Ctx).DataForJqueryDataTable(receiver.DB, receiver.PagerFunc).ToGinResponse()
	default:
		wrongs.ThrowForbidden("响应序列化器错误")
		return 0, nil
	}
}

// AsBusiness 作为Business格式返回
func (receiver Response) AsBusiness(businessType, traceId string) map[string]any {
	switch receiver.Method {
	case types.RESPONSE_METHOD_OK:
		return NewCorrectWithBusiness(receiver.Msg, businessType, traceId).Ok().ToMap()
	case types.RESPONSE_METHOD_DATUM:
		return NewCorrectWithBusiness(receiver.Msg, businessType, traceId).Datum(receiver.Context).ToMap()
	case types.RESPONSE_METHOD_CREATED:
		return NewCorrectWithBusiness(receiver.Msg, businessType, traceId).Created(receiver.Context).ToMap()
	case types.RESPONSE_METHOD_UPDATED:
		return NewCorrectWithBusiness(receiver.Msg, businessType, traceId).Updated(receiver.Context).ToMap()
	case types.RESPONSE_METHOD_DELETED:
		return NewCorrectWithBusiness(receiver.Msg, businessType, traceId).Deleted().ToMap()
	case types.RESPONSE_METHOD_PAGER:
		return NewCorrectWithBusiness(receiver.Msg, businessType, traceId).DataForPager(receiver.DB, receiver.PagerFunc).ToMap()
	default:
		wrongs.ThrowForbidden("响应序列化器错误")
		return nil
	}
}

// NewCorrectWithGinContext 正确返回值
func NewCorrectWithGinContext(msg string, ctx *gin.Context) *Correct {
	once.Do(func() { ins = &Correct{msg: ""} })
	ins.msg = msg
	ins.context = ctx
	ins.messageId = uuid.NewV4().String()
	return ins
}

// NewCorrectWithBusiness 返回正确值（不使用gin context）
func NewCorrectWithBusiness(msg, businessType, traceId string) *Correct {
	once.Do(func() { ins = &Correct{msg: ""} })
	ins.msg = msg
	ins.businessType = businessType
	ins.traceId = traceId
	ins.messageId = uuid.NewV4().String()
	return ins
}

// Ok 操作成功
func (receiver *Correct) Ok() types.StdResponse {
	if receiver.msg == "" {
		receiver.msg = "OK"
	}

	return types.StdResponse{
		MessageId:    receiver.messageId,
		TraceId:      receiver.traceId,
		Msg:          receiver.msg,
		Content:      0,
		Status:       200,
		ErrorCode:    0,
		Pagination:   nil,
		BusinessType: receiver.businessType,
	}
}

// Datum 读取成功
func (receiver *Correct) Datum(content any) types.StdResponse {
	if receiver.msg == "" {
		receiver.msg = "OK"
	}

	return types.StdResponse{
		MessageId:    receiver.messageId,
		TraceId:      receiver.traceId,
		Msg:          receiver.msg,
		Content:      content,
		Status:       200,
		ErrorCode:    0,
		Pagination:   nil,
		BusinessType: receiver.businessType,
	}
}

// Blank 返回空响应
func (receiver *Correct) Blank() types.StdResponse {
	return receiver.Datum(map[string]any{})
}

// DataForPager 返回分页数据
func (receiver *Correct) DataForPager(db *gorm.DB, read func(*gorm.DB) map[string]any) types.StdResponse {
	var count int64

	if receiver.msg == "" {
		receiver.msg = "OK"
	}

	pageStr := receiver.context.Query("__page__")
	limitStr := receiver.context.DefaultQuery("__limit__", "1000")
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)

	if pageStr != "" {
		db.Count(&count)

		return types.StdResponse{
			MessageId: receiver.messageId,
			TraceId:   receiver.traceId,
			Msg:       receiver.msg,
			Content:   read(db.Offset((page - 1) * limit).Limit(limit)),
			Status:    200,
			ErrorCode: 0,
			Pagination: map[string]interface{}{
				"page_curr": page,
				"page_last": math.Ceil(float64(count/int64(limit))) + 1,
				"page_prev": page - 1,
				"page_next": page + 1,
				"page_size": limit,
				"count":     count,
			},
			BusinessType: receiver.businessType,
		}
	} else {
		if limit <= 0 {
			return types.StdResponse{
				MessageId:    receiver.messageId,
				Msg:          receiver.msg,
				Content:      read(db),
				Status:       200,
				ErrorCode:    0,
				Pagination:   nil,
				BusinessType: receiver.businessType,
			}
		} else {
			return types.StdResponse{
				MessageId:    receiver.messageId,
				Msg:          receiver.msg,
				Content:      read(db.Limit(limit)),
				Status:       200,
				ErrorCode:    0,
				Pagination:   nil,
				BusinessType: receiver.businessType,
			}
		}
	}
}

// DataForJqueryDataTable 返回jquery-dataTable格式分页数据
func (receiver *Correct) DataForJqueryDataTable(db *gorm.DB, read func(*gorm.DB) map[string]any) types.JqueryDataTableResponse {
	var count int64

	if receiver.msg == "" {
		receiver.msg = "OK"
	}

	startStr := receiver.context.DefaultQuery("start", "1")
	limitStr := receiver.context.DefaultQuery("length", "1000")

	limit, _ := strconv.Atoi(limitStr)
	start, _ := strconv.Atoi(startStr)
	db.Count(&count)

	return types.JqueryDataTableResponse{
		Content:         read(db.Offset(start).Limit(limit)),
		Draw:            receiver.context.DefaultQuery("draw", "1"),
		Start:           start,
		Length:          limit,
		Max:             &count,
		RecordsTotal:    &count,
		RecordsFiltered: &count,
	}
}

// Created 新建成功
func (receiver *Correct) Created(content any) types.StdResponse {
	if receiver.msg == "" {
		receiver.msg = "新建成功"
	}

	return types.StdResponse{
		MessageId:    receiver.messageId,
		TraceId:      receiver.traceId,
		Msg:          receiver.msg,
		Content:      content,
		Status:       201,
		ErrorCode:    0,
		Pagination:   nil,
		BusinessType: receiver.businessType,
	}
}

// Updated 更新成功
func (receiver *Correct) Updated(content any) types.StdResponse {
	if receiver.msg == "" {
		receiver.msg = "编辑成功"
	}

	return types.StdResponse{
		MessageId:    receiver.messageId,
		TraceId:      receiver.traceId,
		Msg:          receiver.msg,
		Content:      content,
		Status:       202,
		ErrorCode:    0,
		Pagination:   nil,
		BusinessType: receiver.businessType,
	}
}

// Deleted 删除成功
func (receiver *Correct) Deleted() types.StdResponse {
	if receiver.msg == "" {
		receiver.msg = "删除成功"
	}
	return types.StdResponse{
		MessageId:    receiver.messageId,
		TraceId:      receiver.traceId,
		Msg:          receiver.msg,
		Content:      nil,
		Status:       204,
		ErrorCode:    0,
		Pagination:   nil,
		BusinessType: receiver.businessType,
	}
}
