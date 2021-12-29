package conn

import (
	"chatroom/libs"
	"chatroom/message"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const JOIN = "上线啦"
const LEAVE = "下线啦"

const TEXT_TYPE = "text_type"
const STATUS_TYPE = "status_type"

type room struct {
	register   chan *client
	unregister chan *client
	clients    map[*client]bool
	broadcast  chan *message.Message
}

type client struct {
	uid      string
	gravatar string
	conn     *websocket.Conn
	send     chan *message.Message
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	chatroom = &room{
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
		broadcast:  make(chan *message.Message),
	}
)

func Connection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	uid := r.FormValue("uid")
	gravatar := libs.UrlSize(uid, 32)

	// 添加用户
	c := &client{
		uid:      uid,
		gravatar: gravatar,
		conn:     conn,
		send:     make(chan *message.Message, 256),
	}
	chatroom.register <- c

	go c.ReadMessage()
	go c.SendMessage()
}

// 读取消息
func (c *client) ReadMessage() {
	preMessageTime := int64(0)
	for {
		// 接收消息
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			c.conn.Close()
			chatroom.unregister <- c
			return
		}

		// 限制用户发送消息频率，每1秒只能发送一条消息
		curMessageTime := time.Now().Unix()
		if curMessageTime-preMessageTime < 0 {
			log.Println("消息发送过于频繁，请稍后再试")
			continue
		}
		preMessageTime = curMessageTime

		// 解析消息
		msg := message.NewMessage(
			c.uid,
			string(p),
			c.gravatar,
			time.Now().Format("2006-01-02 15:04:05"),
			TEXT_TYPE,
		)

		chatroom.broadcast <- msg
	}
}

// 发送消息
func (c *client) SendMessage() {
	for {
		m, _ := (<-c.send).Encode()
		if err := c.conn.WriteMessage(websocket.TextMessage, m); err != nil {
			c.conn.Close()
			chatroom.unregister <- c
			return
		}
	}
}

func init() {
	log.Println("初始化聊天室")
	go func() {
		for {
			select {
			case c := <-chatroom.register:
				chatroom.clients[c] = true
				log.Println(c.uid + " join")

				go func() {
					// 构造并发送上线消息
					msg := message.NewMessage(
						c.uid,
						c.uid+JOIN,
						c.gravatar,
						time.Now().Format("2006-01-02 15:04:05"),
						STATUS_TYPE)
					chatroom.broadcast <- msg
				}()

			case c := <-chatroom.unregister:
				if _, ok := chatroom.clients[c]; ok {
					delete(chatroom.clients, c)
					close(c.send)
					log.Println(c.uid + " leave")
				}
			case m := <-chatroom.broadcast:
				for c := range chatroom.clients {
					select {
					case c.send <- m:
					default:
						close(c.send)
						delete(chatroom.clients, c)
					}
				}
			}
		}
	}()
}
