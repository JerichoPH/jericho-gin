package wrongs

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func RecoverHandler(c *gin.Context) {
	defer func() {
		if reco := recover(); reco != nil {
			// 打印错误堆栈信息
			log.Printf("panic: %v\n", reco)

			// 判断错误类型
			switch fmt.Sprintf("%T", reco) {
			case "*validator.Validationerrors":
				// 表单验证错误
				c.JSON(NewInCorrect().Validate("", errorToString(reco)).ToGinResponse())
			case "*wrongs.ValidateWrong":
				c.JSON(NewInCorrect().Validate(errorToString(reco), map[string]string{}).ToGinResponse())
			case "*wrongs.ForbiddenWrong":
				// 禁止操作
				c.JSON(NewInCorrect().Forbidden(errorToString(reco)).ToGinResponse())
			case "*wrongs.EmptyWrong":
				// 空数据
				c.JSON(NewInCorrect().Empty(errorToString(reco)).ToGinResponse())
			case "*wrongs.UnAuthWrong":
				// 未授权
				c.JSON(NewInCorrect().UnAuthorization(errorToString(reco)).ToGinResponse())
			case "*wrongs.UnLoginWrong":
				// 未登录
				c.JSON(NewInCorrect().ErrUnLogin().ToGinResponse())
			default:
				// 其他错误
				c.JSON(NewInCorrect().Accident(errorToString(reco), reco).ToGinResponse())
				debug.PrintStack() // 打印堆栈信息
			}

			c.Abort()
		}
	}()

	c.Next()
}

// recover错误，转string
func errorToString(recover interface{}) string {
	switch errorType := recover.(type) {
	case error:
		return errorType.Error()
	default:
		return recover.(string)
	}
}
