package WebSocket

import (
	"VideoWeb/define"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type ClientConnection struct {
	UserID     string
	Conn       *websocket.Conn
	Recv       chan *define.Message //暂存消息
	Send       chan *define.Message //发送消息
	CreateTime uint64
	//CntConn int
}

func NewConnection(UserID string, ctx *gin.Context) (*ClientConnection, error) {
	conn, err := define.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("[NewConnection] Error in NewConnection!")
		log.Println(err)
		return nil, err
	}
	var newConn *ClientConnection

	newConn = &ClientConnection{
		UserID:     UserID,
		Conn:       conn,
		Recv:       make(chan *define.Message, 99),
		Send:       make(chan *define.Message, 10),
		CreateTime: uint64(time.Now().Unix()),
	}
	Hub.register <- newConn
	return newConn, nil
}

func (c *ClientConnection) ReadFromConn() {
	defer func() {
		Hub.unregister <- c
	}()
	fmt.Println("[ReadFromConn] Starting read......")
	c.Conn.SetReadDeadline(time.Now().Add(define.PongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(define.PongWait)); return nil })
	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ReadFromConn] IsUnexpectedCloseError error: %v", err)
			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("[ReadFromConn] Client closed the connection.")
			} else {
				log.Printf("[ReadFromConn] error: %v", err)
			}
			break
		}
		fmt.Println("[ReadFromConn] Get Message: ", string(message), "messageType: ", messageType)
		c.Conn.WriteMessage(1, []byte("111111"))
	}
}

func (c *ClientConnection) WriteToConn() {
	ticker := time.NewTicker(define.PingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			//设置写操作的超时时间，若过期，则 c.Conn.WriteMessage 会返回一个错误
			c.Conn.SetWriteDeadline(time.Now().Add(define.WriteWait))
			err := c.Conn.WriteMessage(websocket.PingMessage, nil)
			fmt.Printf("[WriteToConn] Tick-tick:")
			if err != nil {
				log.Println("[WriteToConn] error:", err)
				return
			}
			fmt.Println("[WriteToConn] Send PingMessage success!", time.Now().Format("2006-01-02T15:04:05"))
		case msg, ok := <-c.Send:
			if !ok {
				fmt.Println("[WriteToConn] Client is Closed.")
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(define.WriteWait))
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			fmt.Println("[WriteToConn] msg:", msg)
			if msg != nil {
				w.Write([]byte(msg.Title + " " + msg.Body))
			}
			if err := w.Close(); err != nil {
				return
			}
			//c.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Title+" "+msg.Body))
		}
	}
}
