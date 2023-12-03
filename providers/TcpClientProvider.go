package providers

import (
	"encoding/json"
	"jericho-gin/tools"
	"jericho-gin/types"
	"log"
	"net"
	"time"
)

var TcpClient net.Conn

func TcpClientHandler(tcpClientAddr string) {
	var err error
	TcpClient, err = net.DialTimeout("tcp", tcpClientAddr, 5*time.Second)
	if err != nil {
		log.Fatalf("[tcp-client-error] [启动TCP客户端失败] %v\n", err)
	}
	defer func(TcpClient net.Conn) {
		if err = TcpClient.Close(); err != nil {
			log.Printf("[tcp-client-error] [关闭TCP客户端失败] %v\n", err)
		}
	}(TcpClient)

	// 从TCP客户端接收响应数据并返回给HTTP客户端
	buf := make([]byte, 1024)
	business := &types.StdBusiness{}
	n, err := TcpClient.Read(buf)
	message := string(buf[:n])
	if err != nil {
		log.Printf("[tcp-client-error] [接收消息失败] %v\n", err)
		return
	}

	if err = json.Unmarshal(buf[:n], business); err != nil {
		log.Printf("[tcp-client-error] [解析消息失败] %v\n", err)
		return
	}
	switch business.BusinessType {
	case "ping":
		log.Printf("[tcp-client-debug] [%s] %s\n", business.BusinessType, message)
	}

	log.Printf("[tcp-client-debug] [接收消息] %s", message)

	for {
		TcpClientSendMessage(tools.NewCorrectWithBusiness("ping", "ping", "").Datum(map[string]any{"time": time.Now().Unix()}).ToJsonStr())
		time.Sleep(5 * time.Second)
	}
}

// TcpClientSendMessage 发送消息到TCP服务器端
func TcpClientSendMessage(message string) {
	_, err := TcpClient.Write([]byte(message))
	if err != nil {
		log.Fatalf("[tcp-client-service] [发送消息失败] %v\n", err)
		return
	}
}
