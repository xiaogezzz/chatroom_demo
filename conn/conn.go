package conn

import (
	"chatroom/libs"
	"chatroom/message"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connMap sync.Map

const JOIN = "上线啦"
const LEAVE = "下线啦"

const TEXT_TYPE = "text_type"
const STATUS_TYPE = "status_type"

type connInfo struct {
	Uid      string
	Gravatar string
	Conn     *websocket.Conn
	mutex    sync.Mutex
}

func Connection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	uid := r.FormValue("uid")
	gravatar := libs.UrlSize(uid, 32)

	msg := message.NewMessage(
		uid,
		uid+JOIN,
		gravatar,
		time.Now().Format("2006-01-02 15:04:05"),
		STATUS_TYPE)
	data, _ := msg.Encode()

	go Broadcast(data)

	connInfo := &connInfo{Uid: uid, Gravatar: gravatar, Conn: conn}
	connMap.Store(uid, connInfo)

	Receive(connInfo)
}

func Broadcast(data []byte) {
	connMap.Range(func(key, value interface{}) bool {
		tmpval, _ := value.(*connInfo)
		tmpval.mutex.Lock()
		defer tmpval.mutex.Unlock()

		if err := tmpval.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println(err)
		}
		return true
	})
}

func Receive(connInfo *connInfo) {
	for {
		var data []byte
		_, p, err := connInfo.Conn.ReadMessage()
		if err != nil {
			connMap.Delete(connInfo.Uid)
			msg := message.NewMessage(
				connInfo.Uid,
				connInfo.Uid+LEAVE,
				connInfo.Gravatar,
				time.Now().Format("2006-01-02 15:04:05"),
				STATUS_TYPE)
			data, _ = msg.Encode()
			Broadcast(data)
			return
		}
		msg := message.NewMessage(
			connInfo.Uid,
			string(p),
			connInfo.Gravatar,
			time.Now().Format("2006-01-02 15:04:05"),
			TEXT_TYPE)
		data, _ = msg.Encode()

		go Broadcast(data)
	}
}
