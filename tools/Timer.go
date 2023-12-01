package tools

import (
	"fmt"
	"jericho-go/types"
	"os"
	"time"
)

// Timer 计时器
// 使用方法说明：
// c := make(chan uint64)
// go sum(c)
// go tools.NewTimer(5).Ticker()
// runtime.Gosched()
// println(<-c)
type Timer struct {
	originalTime time.Time
	maxTime      float64
	needPrint    bool
}

// NewTimer 初始化计时器
func NewTimer(maxTime float64) *Timer {
	return &Timer{originalTime: time.Now(), maxTime: maxTime}
}

// Ticker 计时器，超出时间自动停止
func (receiver *Timer) Ticker() {
	for {
		if time.Now().Sub(receiver.originalTime).Seconds() > receiver.maxTime {

			println(fmt.Sprintf("程序超时，放弃执行『结束进程』。开始时间：『%s』。结束时间：『%s』超时时间：『%f』秒", receiver.originalTime.Format(string(types.TIME_FORMAT_DATETIME)), NewTime().SetTimeNowAdd8Hour().ToDateTimeString(), receiver.maxTime))
			os.Exit(0)
		} else {
			if receiver.needPrint {
				println(time.Now().Sub(receiver.originalTime).Seconds())
			}
		}
	}
}

// NeedPrint 设置是否需要打印
func (receiver *Timer) NeedPrint() *Timer {
	receiver.needPrint = true

	return receiver
}
