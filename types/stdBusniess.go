package types

import (
	"encoding/json"
)

type (
	// StdBusiness tcp服务器端、tcp客户端、websocket业务消息格式
	StdBusiness struct {
		MessageId    string                 `json:"message_id"`
		TraceId      string                 `json:"trace_id"`
		BusinessType string                 `json:"business_type"`
		Content      map[string]interface{} `json:"content"`
	}
)

// ToMap 转map
func (receiver StdBusiness) ToMap() MapStringToAny {
	return MapStringToAny{
		"message_id":    receiver.MessageId,
		"trace_id":      receiver.TraceId,
		"content":       receiver.Content,
		"business_type": receiver.BusinessType,
	}
}

// ToJson 转json
func (receiver StdBusiness) ToJson() []byte {
	data := receiver.ToMap()
	if marshal, err := json.Marshal(data); err != nil {
		return []byte{}
	} else {
		return marshal
	}
}

// ToJsonStr 转json字符串
func (receiver StdBusiness) ToJsonStr() string {
	return string(receiver.ToJson())
}
