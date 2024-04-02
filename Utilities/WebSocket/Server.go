package WebSocket

import (
	"fmt"
	"log"
	"sync"
)

type ServerHub struct {
	// UserConnections 维护一个用户ID-->webSocket连接的映射关系
	UserConnections map[string]map[uint64]*ClientConnection
	register        chan *ClientConnection //注册新的webSocket
	unregister      chan *ClientConnection //注销新的webSocket
	mu              sync.Mutex
}

var Hub *ServerHub

// NewServerHub 生成新的ServerHub来管理全部WebSocket连接ClientConnection
func NewServerHub() *ServerHub {
	return &ServerHub{
		UserConnections: make(map[string]map[uint64]*ClientConnection),
		register:        make(chan *ClientConnection),
		unregister:      make(chan *ClientConnection),
		mu:              sync.Mutex{},
	}
}

// RegisterConnections 向ServerHub中添加新的client
func (s *ServerHub) RegisterConnections(client *ClientConnection) {
	fmt.Println("[RegisterConnections] starting registering......")
	s.mu.Lock()
	defer s.mu.Unlock()

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

// UnRegisterConnections 向ServerHub减小对应连接的引用计数，若计数为0，删除连接
func (s *ServerHub) UnRegisterConnections(client *ClientConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	close(s.UserConnections[client.UserID][client.CreateTime].Recv)
	close(s.UserConnections[client.UserID][client.CreateTime].Send)
	err := s.UserConnections[client.UserID][client.CreateTime].Conn.Close()
	if err != nil {
		log.Println("[UnregisterConnections] err at closing Conn:", err)
	}
	delete(s.UserConnections[client.UserID], client.CreateTime)
	for k, v := range s.UserConnections[client.UserID] {
		fmt.Println(k, "--->", v.Conn.RemoteAddr().String())
	}
	if len(s.UserConnections[client.UserID]) == 0 {
		delete(s.UserConnections, client.UserID)
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
