package service

import (
	"VideoWeb/Utilities/WebSocket"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

// CreateWebSocket 建立websocket连接
func CreateWebSocket(c *gin.Context) {
	UserID := c.Param("UserID")
	conn, err := WebSocket.NewConnection(UserID, c)
	if err != nil {
		log.Println("[CreateWebSocket] Error creating WebSocket connection: ", err)
		return
	}
	go conn.WriteToClient()
	go conn.ReadFromClient()
	fmt.Println("[CreateWebSocket] Create WebSocket Successfully!")
}
