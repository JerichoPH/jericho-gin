package tools

import (
	"strings"
)

// Operation 计算工具
type Operation struct {
}

// NewOperation 初始化计算工具
func NewOperation() Operation {
	return Operation{}
}

// Ternary 三元表达式
func (receiver Operation) Ternary(condition bool, trueValue interface{}, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

// StringIsTrue 判断字符串内的内容是否代表真(如果严格代表真,则返回 true)
func (receiver Operation) StringIsTrue(value string) bool {
	var isTrue bool
	isTrue = false

	switch strings.ToLower(value) {
	case "on":
		isTrue = true
	case "yes":
		isTrue = true
	case "1":
		isTrue = true
	case "true":
		isTrue = true
	}

	return isTrue
}

// StringIsFalse 判断字符串内的内容是否代表假(如果严格代表假,则返回true)
func (receiver Operation) StringIsFalse(value string) bool {
	var isFalse bool
	isFalse = false

	switch strings.ToLower(value) {
	case "off":
		isFalse = true
	case "no":
		isFalse = true
	case "0":
		isFalse = true
	case "false":
		isFalse = true
	}

	return isFalse
}
