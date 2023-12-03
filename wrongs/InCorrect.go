package wrongs

import (
	"jericho-gin/types"
	"sync"
)

type inCorrect struct {
	businessType string
	stdResponse  types.StdResponse
}

var ins *inCorrect
var once sync.Once

func NewInCorrect() *inCorrect {
	once.Do(func() { ins = &inCorrect{} })
	return ins
}

func NewInCorrectWithBusniess(businessType string) *inCorrect {
	once.Do(func() { ins = &inCorrect{} })
	ins.businessType = businessType
	return ins
}

func (receiver inCorrect) UnAuthorization(msg string) types.StdResponse {
	if msg == "" {
		msg = "未授权"
	}

	return types.StdResponse{
		Msg:          msg,
		Content:      nil,
		Status:       406,
		ErrorCode:    1,
		BusinessType: receiver.businessType,
	}
}

func (receiver inCorrect) ErrUnLogin() types.StdResponse {
	return types.StdResponse{
		Msg:          "未登录",
		Content:      nil,
		Status:       401,
		ErrorCode:    2,
		BusinessType: receiver.businessType,
	}
}

func (receiver inCorrect) Forbidden(msg string) types.StdResponse {
	if msg == "" {
		msg = "禁止操作"
	}

	return types.StdResponse{
		Msg:          msg,
		Content:      nil,
		Status:       403,
		ErrorCode:    3,
		BusinessType: receiver.businessType,
	}
}

func (receiver inCorrect) Empty(msg string) types.StdResponse {
	if msg == "" {
		msg = "数不存在"
	}

	return types.StdResponse{
		Msg:          msg,
		Content:      nil,
		Status:       404,
		ErrorCode:    4,
		BusinessType: receiver.businessType,
	}
}

func (receiver inCorrect) Validate(msg string, content any) types.StdResponse {
	if msg == "" {
		msg = "表单验证错误"
	}

	return types.StdResponse{
		Msg:          msg,
		Content:      content,
		Status:       421,
		ErrorCode:    5,
		BusinessType: receiver.businessType,
	}
}

func (receiver inCorrect) Error(msg string, content types.MapStringToAny) types.StdResponse {
	if msg == "" {
		msg = "错误"
	}

	return types.StdResponse{
		Msg:          msg,
		Content:      content,
		Status:       400,
		ErrorCode:    6,
		BusinessType: receiver.businessType,
	}
}

func (receiver inCorrect) Accident(msg string, err any) types.StdResponse {
	if msg == "" {
		msg = "意外错误"
	}

	return types.StdResponse{
		Msg:          msg,
		Content:      nil,
		Status:       500,
		ErrorCode:    -1,
		BusinessType: receiver.businessType,
	}
}
