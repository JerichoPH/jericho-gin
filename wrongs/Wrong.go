package wrongs

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type Wrong struct{ ErrorMessage string }

// EmptyWrong 空数据异常
type EmptyWrong struct{ Wrong }

// Error 获取异常信息
//
//	@receiver receiver
//	@return string
func (receiver *Wrong) Error() string {
	return receiver.ErrorMessage
}

// ValidateWrong 表单验证错误
type ValidateWrong struct{ Wrong }

// ThrowValidate 421错误
//
//	@param text
func ThrowValidate(text string, more ...interface{}) {
	panic(&ValidateWrong{Wrong{ErrorMessage: fmt.Sprintf(text, more...)}})
}

// ThrowEmpty 404错误
//
//	@param text
//	@return error
func ThrowEmpty(text string, more ...interface{}) {
	panic(&EmptyWrong{Wrong{ErrorMessage: fmt.Sprintf(text, more...)}})
}

// ForbiddenWrong 操作错误
type ForbiddenWrong struct{ Wrong }

// ThrowForbidden 403错误
//
//	@param text
//	@return error
func ThrowForbidden(text string, more ...interface{}) {
	panic(&ForbiddenWrong{Wrong{ErrorMessage: fmt.Sprintf(text, more...)}})
}

// UnAuthWrong 未授权异常
type UnAuthWrong struct{ Wrong }

// ThrowUnAuth 未授权错误
//
//	@param text
//	@return error
func ThrowUnAuth(text string, more ...interface{}) {
	panic(&UnAuthWrong{Wrong{ErrorMessage: fmt.Sprintf(text, more...)}})
}

// UnLoginWrong 未登录异常
type UnLoginWrong struct{ Wrong }

// ThrowUnLogin 未登录错误
//
//	@param text
func ThrowUnLogin(text string, more ...interface{}) {
	panic(&UnLoginWrong{Wrong{ErrorMessage: fmt.Sprintf(text, more...)}})
}

// ThrowWhenIsNotInt 文字转整型
//
//	@param v
//	@param errMsg
//	@return intValue
func ThrowWhenIsNotInt(strValue string, errorMessage string) (intValue int) {
	intValue, err := strconv.Atoi(strValue)
	if err != nil && errorMessage != "" {
		ThrowForbidden(errorMessage)
	}
	return
}

// ThrowWhenIsNotUint 文字转无符号整型
//
//	@param v
//	@param errMsg
//	@return uintValue
func ThrowWhenIsNotUint(strValue string, errorMessage string) (uintValue uint) {
	intValue := ThrowWhenIsNotInt(strValue, errorMessage)
	uintValue = uint(intValue)
	return
}

// ThrowWhenIsEmpty 当数据库返回空则报错
//
//	@param db
//	@param name
//	@return bool
func ThrowWhenIsEmpty(db *gorm.DB, errorField string) bool {
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			if errorField != "" {
				ThrowEmpty(errorField + "不存在")
				return true
			} else {
				return true
			}
		}
	}
	return false
}

// ThrowWhenIsRepeat 当数据库返回不空则报错
//
//	@param db
//	@param name
//	@return bool
func ThrowWhenIsRepeat(db *gorm.DB, errorField string) bool {
	if db.Error == nil {
		if errorField != "" {
			ThrowForbidden(errorField + "重复")
			return false
		} else {
			return false
		}
	}
	return true
}
