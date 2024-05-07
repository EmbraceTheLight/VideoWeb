package WebSocket

import (
	"fmt"
	"log"
)

type ServerHub struct {
	// UserConnections 维护一个用户ID-->webSocket连接的映射关系
	UserConnections map[string]map[uint64]*ClientConnection
	register        chan *ClientConnection //注册新的webSocket
	unregister      chan *ClientConnection //注销新的webSocket
}

var hub *ServerHub

// NewServerHub 生成新的ServerHub来管理全部WebSocket连接ClientConnection
func NewServerHub() *ServerHub {
	hub = &ServerHub{
		UserConnections: make(map[string]map[uint64]*ClientConnection),
		register:        make(chan *ClientConnection),
		unregister:      make(chan *ClientConnection),
	}
	return hub
}

// RegisterConnections 向ServerHub中添加新的client
func (s *ServerHub) RegisterConnections(client *ClientConnection) {
	fmt.Println("[RegisterConnections] starting registering......")
	//map没有初始化，则初始化一个空的map
	if _, ok := s.UserConnections[client.UserID]; !ok {
		s.UserConnections[client.UserID] = make(map[uint64]*ClientConnection)
	}
	s.UserConnections[client.UserID][client.CreateTime] = client
	fmt.Println("[RegisterConnections] User ", client.UserID, "has cntConns below: ")
	for k, v := range s.UserConnections[client.UserID] {
		fmt.Println(k, ":", v.Conn.RemoteAddr().String())
	}
	fmt.Println("[RegisterConnections] Create connection Successfully!")
}

func cleanConn(client *ClientConnection) error {
	if _, ok := hub.UserConnections[client.UserID][client.CreateTime]; !ok {
		fmt.Println("[CleanConn] User ", client.UserID, " has already been cleaned.")
		return nil
	}
	close(hub.UserConnections[client.UserID][client.CreateTime].Send)
	err := hub.UserConnections[client.UserID][client.CreateTime].Conn.Close()
	if err != nil {
		return err
	}
	//err = DAO.DeleteUserRoomDocumentByUserID(client.UserID)
	//if err != nil {
	//	return err
	//}
	delete(hub.UserConnections[client.UserID], client.CreateTime)
	if len(hub.UserConnections[client.UserID]) == 0 {
		delete(hub.UserConnections, client.UserID)
	}
	return nil
}

// UnRegisterConnections 删除指定的client
func (s *ServerHub) UnRegisterConnections(client *ClientConnection) {
	//清除该链接的一切有关记录
	err := cleanConn(client)
	if err != nil {
		log.Println("[UnregisterConnections] err at closing Conn:", err)
	}
	fmt.Println("[UnregisterConnections] delete successfully.")
}

func (s *ServerHub) Run() {
	for {
		select {
		case conn := <-s.register:
			s.RegisterConnections(conn)
		case conn := <-s.unregister:
			s.UnRegisterConnections(conn)
		}
	}
}
