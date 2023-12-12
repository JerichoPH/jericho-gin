package controllers

import (
	"jericho-gin/providers"
	"jericho-gin/tools"
	"jericho-gin/wrongs"

	"github.com/gin-gonic/gin"
)

type (
	// WebsocketCtrl websocket控制器
	WebsocketCtrl      struct{}
	webSocketStoreForm struct {
		ReceiverUuid string `json:"receiver_uuid"`
		Message      any    `json:"message"`
	}
)

// ShouldBind 表单绑定
func (receiver webSocketStoreForm) ShouldBind(ctx *gin.Context) webSocketStoreForm {
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}

	if receiver.ReceiverUuid == "" {
		wrongs.ThrowValidate("接收人uuid不能为空")
	}
	if receiver.Message == nil {
		wrongs.ThrowValidate("消息不能为空")
	}

	return receiver
}

// NewWebsocketCtrl 初始化websocket控制
func NewWebsocketCtrl() *WebsocketCtrl {
	return &WebsocketCtrl{}
}

// SendTo 发送消息
func (receiver WebsocketCtrl) SendTo(ctx *gin.Context) {
	form := (&webSocketStoreForm{}).ShouldBind(ctx)

	// 发送到websocket客户端
	providers.WebsocketSendMessageByUuid(tools.NewCorrectWithBusiness("", "message", "").Datum(form.Message).ToJsonStr(), form.ReceiverUuid)
}
