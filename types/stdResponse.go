package types

import (
	"encoding/json"
)

// StdResponse 标准响应结构
type StdResponse struct {
	MessageId    string
	TraceId      string
	Msg          string
	Content      any
	Status       int
	ErrorCode    int
	Pagination   any
	BusinessType string
}

// ToMap 转map
func (receiver StdResponse) ToMap() MapStringToAny {
	return MapStringToAny{
		"message_id":    receiver.MessageId,
		"trace_id":      receiver.TraceId,
		"msg":           receiver.Msg,
		"content":       receiver.Content,
		"status":        receiver.Status,
		"error_code":    receiver.ErrorCode,
		"pagination":    receiver.Pagination,
		"business_type": receiver.BusinessType,
	}
}

// ToJson 转json
func (receiver StdResponse) ToJson() []byte {
	data := receiver.ToMap()
	if marshal, err := json.Marshal(data); err != nil {
		return []byte{}
	} else {
		return marshal
	}
}

// ToJsonStr 转json字符串
func (receiver StdResponse) ToJsonStr() string {
	return string(receiver.ToJson())
}

// ToGinResponse 转gin响应格式 int+map[string]interface{}
func (receiver StdResponse) ToGinResponse() (int, MapStringToAny) {
	return receiver.Status, receiver.ToMap()
}
