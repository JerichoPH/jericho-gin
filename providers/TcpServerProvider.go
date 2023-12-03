package providers

import (
	"encoding/json"
	"io"
	"jericho-gin/tools"
	"jericho-gin/types"
	"jericho-gin/wrongs"
	"log"
	"net"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	tcpClients        = make(map[net.Conn]bool)
	tcpClientMutex    = &sync.Mutex{}
	tcpUuidToAddrDict = make(map[string]string)
	tcpAddrToUuidDict = make(map[string]string)
)

func TcpAddClient(conn net.Conn) {
	tcpClientMutex.Lock()
	defer tcpClientMutex.Unlock()
	tcpClients[conn] = true
}

func TcpRemoveClient(conn net.Conn) {
	tcpClientMutex.Lock()
	defer tcpClientMutex.Unlock()

	log.Printf("[tcp-server-debug] [关闭连接] %s %s\n", tcpAddrToUuidDict[conn.RemoteAddr().String()], conn.RemoteAddr().String())

	delete(tcpClients, conn)
	delete(tcpUuidToAddrDict, tcpAddrToUuidDict[conn.RemoteAddr().String()])
	delete(tcpAddrToUuidDict, conn.RemoteAddr().String())

	log.Printf("[tcp-server-debug] [剩余连接] %v\n", tcpClients)
}

// TcpServerSendMessageByAddr 通过用户地址发送消息
func TcpServerSendMessageByAddr(message, addr string) {
	tcpClientMutex.Lock()
	defer tcpClientMutex.Unlock()
	for conn := range tcpClients {
		if conn.RemoteAddr().String() == addr {
			_, err := conn.Write([]byte(message))
			if err != nil {
				log.Printf("[tcp-server-error] 发送消息失败 %s %s\n", addr, message)
			}
			return
		}
	}
}

// TcpServerSendMessageByUuid 通过用户编号发送消息
func TcpServerSendMessageByUuid(message, uuid string) {
	TcpServerSendMessageByAddr(message, tcpUuidToAddrDict[uuid])
}

// TcpServerSendMessageToAllClient 群发消息
func TcpServerSendMessageToAllClient(message string) {
	tcpClientMutex.Lock()
	defer tcpClientMutex.Unlock()
	for conn := range tcpClients {
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Printf("[tcp-server-error] 发送消息失败 %s %s\n", conn.RemoteAddr().String(), message)
		}
		return
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		TcpRemoveClient(conn)
		err := conn.Close()
		if err != nil {
			log.Printf("[tcp-server-error] [关闭链接失败] %s\n", err.Error())
		}
	}()

	newUuid := uuid.NewV4().String()
	_, err := io.WriteString(conn, tools.NewCorrectWithBusiness("链接成功", "connection-success", "").Datum(map[string]any{"uuid": newUuid}).ToJsonStr())
	tcpUuidToAddrDict[newUuid] = conn.RemoteAddr().String()
	tcpAddrToUuidDict[conn.RemoteAddr().String()] = newUuid

	if err != nil {
		log.Printf("[tcp-server-error] [发送消息失败] %s %s\n", conn.RemoteAddr().String(), "链接成功")
	}
	for {
		buf := make([]byte, 256)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		message := string(buf[:n])

		business := &types.StdBusiness{}
		err = json.Unmarshal(buf[:n], business)
		if err != nil {
			log.Printf("[tcp-server-error] [解析业务失败] %s\n", message)

			TcpServerSendMessageByAddr(wrongs.NewInCorrectWithBusniess("error").Error("业务解析失败", map[string]any{"request_content": message}).ToJsonStr(), conn.RemoteAddr().String())
		}

		log.Printf("[tcp-server-debug] [接收客户端消息] %s %v\n", conn.RemoteAddr().String(), business)

		switch business.BusinessType {
		case "ping":
			log.Printf("[tcp-server-debug] [%s] %s\n", business.BusinessType, message)

			TcpServerSendMessageByAddr(tools.NewCorrectWithBusiness("pong", "pong", "").Datum(map[string]any{"time": time.Now().Unix()}).ToJsonStr(), conn.RemoteAddr().String())
		case "authorization/bindUserUuid":
			log.Printf("[tcp-server-debug] [%s] 绑定用户uuid\n", business.BusinessType)

			tcpAddrToUuidDict[conn.RemoteAddr().String()] = business.Content["uuid"].(string)
			tcpUuidToAddrDict[business.Content["uuid"].(string)] = conn.RemoteAddr().String()

			TcpServerSendMessageByAddr(tools.NewCorrectWithBusiness("绑定成功", business.BusinessType, "").Datum(map[string]any{}).ToJsonStr(), conn.RemoteAddr().String())
		}
	}
}

func TcpServerHandler(tcpServerAddr string) {
	log.Printf("[tcp-server-debug] [启动TCP服务] %s\n", tcpServerAddr)
	listener, err := net.Listen("tcp", tcpServerAddr)
	if err != nil {
		log.Fatalf("[tcp-server-error] [启动TCP服务失败] %v\n", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("[tcp-server-error] [建立连接TCP端口失败] %v", err)
		} else {
			log.Printf("[tcp-server-debug] [建立连接TCP端口成功] 监听：%s\n", conn.RemoteAddr().String())
			TcpAddClient(conn)
			go handleConnection(conn)
		}
	}
}
