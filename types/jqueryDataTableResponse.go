package types

import "encoding/json"

// JqueryDataTableResponse 标准响应结构
type JqueryDataTableResponse struct {
	Draw            string
	Start           int
	Length          int
	Content         MapStringToAny
	Max             *int64
	RecordsTotal    *int64
	RecordsFiltered *int64
}

// ToMap 转map[string]interface
func (receiver JqueryDataTableResponse) ToMap() MapStringToAny {
	return map[string]interface{}{
		"content":         receiver.Content,
		"draw":            receiver.Draw,
		"start":           receiver.Start,
		"length":          receiver.Length,
		"max":             receiver.Max,
		"recordsTotal":    receiver.RecordsTotal,
		"recordsFiltered": receiver.RecordsFiltered,
	}
}

// ToJson 转json
func (receiver JqueryDataTableResponse) ToJson() []byte {
	data := receiver.ToMap()
	if marshal, err := json.Marshal(data); err != nil {
		return []byte{}
	} else {
		return marshal
	}
}

// ToJsonStr 转json字符串
func (receiver JqueryDataTableResponse) ToJsonStr() string {
	return string(receiver.ToJson())
}

// ToGinResponse 转gin响应格式 int+map[string]interface{}
func (receiver JqueryDataTableResponse) ToGinResponse() (int, MapStringToAny) {
	return 200, receiver.ToMap()
}
