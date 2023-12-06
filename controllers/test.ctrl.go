package controllers

import (
	"fmt"
	"jericho-gin/providers"
	"jericho-gin/settings"
	"jericho-gin/tools"
	"jericho-gin/wrongs"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

type (
	// TestController 测试控制器
	TestController struct{}

	// 发送消息表单
	sendMessageForm struct {
		ReceiverUuid   string `json:"receiver_uuid"`
		MessageTitle   string `json:"message_title"`
		MessageContent any    `json:"message_content"`
	}
)

// NewTestController 构造函数
func NewTestController() *TestController {
	return &TestController{}
}

// ShouldBind 表单绑定
func (receiver sendMessageForm) ShouldBind(ctx *gin.Context) sendMessageForm {
	err := ctx.ShouldBind(&receiver)
	if err != nil {
		wrongs.ThrowValidate("表单绑定失败：%s", err.Error())
		return sendMessageForm{}
	}

	return receiver
}

// AnySendToWebsocket 接收并转发消息(websocket)
func (receiver TestController) AnySendToWebsocket(ctx *gin.Context) {
	form := sendMessageForm{}.ShouldBind(ctx)

	providers.WebsocketSendMessageByUuid(tools.NewCorrectWithBusiness(form.MessageTitle, "message", "").Datum(form.MessageContent).ToJsonStr(), form.ReceiverUuid)

	ctx.JSON(tools.NewCorrectWithGinContext("OK", ctx).Datum(map[string]any{"content": fmt.Sprintf("%s ==> %s %v", form.ReceiverUuid, form.MessageTitle, form.MessageContent)}).ToGinResponse())
}

// AnySendToTcpServer 接收并转发消息(tcp server)
func (receiver TestController) AnySendToTcpServer(ctx *gin.Context) {
	form := sendMessageForm{}.ShouldBind(ctx)

	providers.TcpServerSendMessageByUuid(tools.NewCorrectWithBusiness(form.MessageTitle, "message", "").Datum(form.MessageContent).ToJsonStr(), form.ReceiverUuid)

	ctx.JSON(tools.NewCorrectWithGinContext("OK", ctx).Datum(map[string]any{"content": fmt.Sprintf("%s ==> %s %v", form.ReceiverUuid, form.MessageTitle, form.MessageContent)}).ToGinResponse())
}

// AnySendToTcpClient 接收病转发消息(tcp client)
func (receiver TestController) AnySendToTcpClient(ctx *gin.Context) {
	form := sendMessageForm{}.ShouldBind(ctx)

	providers.TcpClientSendMessage(tools.NewCorrectWithBusiness(form.MessageTitle, "message", "").Datum(map[string]any{"content": fmt.Sprintf("%s %v", form.MessageTitle, form.MessageContent)}).ToJsonStr())

	ctx.JSON(tools.NewCorrectWithGinContext("OK", ctx).Datum(map[string]any{"content": fmt.Sprintf("%s %v", form.MessageTitle, form.MessageContent)}).ToGinResponse())
}

// AnySendToKafkaClient 发送消息到kafka
func (receiver TestController) AnySendToKafkaClient(ctx *gin.Context) {
	var (
		appSetting *ini.File
		topic      string
		err        error
	)

	appSetting = settings.NewSetting().App
	topic = appSetting.Section("kafka-server").Key("topic").MustString("")

	form := sendMessageForm{}.ShouldBind(ctx)

	if err = providers.KafkaSendMessage(topic, "message", tools.NewCorrectWithBusiness(form.MessageTitle, "message", "").Datum(map[string]any{"content": fmt.Sprintf("%s %v", form.MessageTitle, form.MessageContent)}).ToJsonStr()); err != nil {
		wrongs.ThrowForbidden("发送失败：%s", err.Error())
	}
}
