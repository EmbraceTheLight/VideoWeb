package WebSocket

import (
	DAO "VideoWeb/DAO/EntitySets"
	"VideoWeb/define"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type ClientConnection struct {
	hub1       *ServerHub
	UserID     string
	Conn       *websocket.Conn
	Send       chan *DAO.Message //向客户端发送消息
	CreateTime uint64
}

func NewConnection(UserID string, ctx *gin.Context) (*ClientConnection, error) {
	up := define.Upgrader
	conn, err := up.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("[NewConnection] Error in NewConnection!")
		log.Println(err)
		return nil, err
	}
	var newConn *ClientConnection

	newConn = &ClientConnection{
		hub1:       hub,
		UserID:     UserID,
		Conn:       conn,
		Send:       make(chan *DAO.Message),
		CreateTime: uint64(time.Now().Unix()),
	}
	newConn.hub1.register <- newConn
	return newConn, nil
}

func (c *ClientConnection) ReadFromClient() {
	//defer func() {
	//	c.hub1.unregister <- c
	//}()
	//设置读取Pong消息的超时时间以及设置Pong消息的处理函数
	c.Conn.SetReadDeadline(time.Now().Add(define.PongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(define.PongWait)); return nil })
	for {
		var msg = new(DAO.Message)
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ReadFromClient] IsUnexpectedCloseError error: %v", err)
			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("[ReadFromClient] Client closed the connection.")
			} else {
				log.Printf("[ReadFromClient] error: %v", err)
			}
			break
		}
		fmt.Printf("[ReadFromClient] Get Message: Title: %s, Body: %s\n", msg.MessageContent.Title, msg.MessageContent.Body)

	}
}

func (c *ClientConnection) WriteToClient() {
	ticker := time.NewTicker(define.PingPeriod)
	defer func() {
		ticker.Stop()
		c.hub1.unregister <- c
	}()
	for {
		select {
		case <-ticker.C:
			//设置写操作的超时时间，若过期，则 c.Conn.WriteMessage 会返回一个错误
			c.Conn.SetWriteDeadline(time.Now().Add(define.WriteWait))
			err := c.Conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("[WriteToClient] error:", err)
				return
			}
			fmt.Println("[WriteToClient] Send PingMessage success!", time.Now().Format("2006-01-02T15:04:05"))
		case msg, ok := <-c.Send:
			if !ok {
				fmt.Println("[WriteToClient] Client is Closed.")
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(define.WriteWait))
			fmt.Printf("[WriteToClient] msg from %v:%v:%s\n", c.UserID, c.CreateTime, msg)
			err := c.Conn.WriteJSON(msg)
			if err != nil {
				log.Printf("[WriteToClient] error at WriteMessage from %v:%v:%s\n", c.UserID, c.CreateTime, err)
				return
			}
		}
	}
}
