package providers

import (
	"encoding/json"
	"jericho-go/tools"
	"jericho-go/types"
	"jericho-go/wrongs"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	websocketClients        = make(map[*websocket.Conn]bool)
	websocketClientMutex    = &sync.Mutex{}
	websocketUuidToAddrDict = make(map[string]string, 0)
	websocketAddrToUuidDict = make(map[string]string, 0)
)

func WebsocketRemoveClient(conn *websocket.Conn) {
	websocketClientMutex.Lock()
	defer websocketClientMutex.Unlock()

	log.Printf("[websocket-debug] [关闭链接] %s %s\n", websocketAddrToUuidDict[conn.RemoteAddr().String()], conn.RemoteAddr().String())

	delete(websocketClients, conn)
	delete(websocketUuidToAddrDict, websocketAddrToUuidDict[conn.RemoteAddr().String()])
	delete(websocketAddrToUuidDict, conn.RemoteAddr().String())

	log.Printf("[websocket-debug] [剩余链接] %v\n", websocketClients)
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebsocketSendMessageByAddr 通过用户地址发送消息
func WebsocketSendMessageByAddr(message, addr string) {
	websocketClientMutex.Lock()
	defer websocketClientMutex.Unlock()
	for ws := range websocketClients {
		if ws.RemoteAddr().String() == addr {
			err := ws.WriteMessage(1, []byte(message))
			if err != nil {
				log.Printf("[websocket-error] 发送消息失败 %s %s\n", addr, message)
				break
			}
			return
		}
	}
}

// WebsocketSendMessageByUuid 通过用户编号发送消息
func WebsocketSendMessageByUuid(message, uuid string) {
	log.Printf("[websocket-debug] [通过 uuid 发送消息] %s %s %s\n", uuid, websocketUuidToAddrDict[uuid], message)
	WebsocketSendMessageByAddr(message, websocketUuidToAddrDict[uuid])
}

func WebsocketAddClient(conn *websocket.Conn) {
	websocketClientMutex.Lock()
	defer websocketClientMutex.Unlock()
	websocketClients[conn] = true

	log.Printf("[websocket-debug] [链接成功，当前连接池] %v\n", websocketClients)
}

func WebsocketHandler(ctx *gin.Context) {
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("[websocket-error] %v\n", err)
		return
	}

	defer func() {
		WebsocketRemoveClient(ws)
		if err = ws.Close(); err != nil {
			log.Printf("[websocket-error] [关闭链接失败] %s\n", err.Error())
		}
	}()

	// newUuid := uuid.NewV4().String()
	if err = ws.WriteMessage(1, []byte(tools.NewCorrectWithBusiness("连接成功", "connection-success", "").Datum(nil).ToJsonStr())); err != nil {
		log.Printf("[websocket-error] [发送消息失败] %v\n", err)
	}
	// websocketUuidToAddrDict[newUuid] = ws.RemoteAddr().String()
	// websocketAddrToUuidDict[ws.RemoteAddr().String()] = newUuid
	WebsocketAddClient(ws)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("[websocket-error] [读取消息失败] %s\n", ws.RemoteAddr().String())
			break
		}

		business := &types.StdBusiness{}
		err = json.Unmarshal(message, business)
		if err != nil {
			log.Printf("[websocket-error] [解析业务失败] %s\n", message)

			WebsocketSendMessageByAddr(wrongs.NewInCorrectWithBusniess("error").Error("业务解析失败", map[string]any{"request_content": message}).ToJsonStr(), ws.RemoteAddr().String())
		}

		switch business.BusinessType {
		case "echo":
			log.Printf("[websocket-debug] [%s] %s\n", business.BusinessType, message)
			WebsocketSendMessageByAddr(tools.NewCorrectWithBusiness("echo", "echo", "").Datum(business.Content).ToJsonStr(), ws.RemoteAddr().String())
		case "ping":
			// log.Printf("[websocket-debug] [%s] %s\n", business.BusinessType, message)
			WebsocketSendMessageByAddr(tools.NewCorrectWithBusiness("pong", "pong", "").Datum(map[string]any{"time": time.Now().Unix()}).ToJsonStr(), ws.RemoteAddr().String())
		case "authorization/bindUserUuid":
			log.Printf("[websocket-debug] [%s] 绑定用户uuid %s %s %v\n", business.BusinessType, business.Content["uuid"], ws.RemoteAddr().String(), websocketClients)

			websocketAddrToUuidDict[ws.RemoteAddr().String()] = business.Content["uuid"].(string)
			websocketUuidToAddrDict[business.Content["uuid"].(string)] = ws.RemoteAddr().String()

			WebsocketSendMessageByAddr(tools.NewCorrectWithBusiness("绑定成功", business.BusinessType, "").Datum(map[string]any{}).ToJsonStr(), ws.RemoteAddr().String())
		}
	}
}
