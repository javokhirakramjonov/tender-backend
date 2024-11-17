package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"tender-backend/db"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust this for security
	},
}

type Client struct {
	UserID int64
	Conn   *websocket.Conn
	Send   chan []byte
}

type Server struct {
	Clients   map[*Client]bool
	Register  chan *Client
	Broadcast chan []byte
	mu        sync.Mutex
	ns        *NotificationService
}

func NewNotificationServer() *Server {
	return &Server{
		Clients:  make(map[*Client]bool),
		Register: make(chan *Client),
		ns:       NewNotificationService(db.DB),
	}
}

func (s *Server) Run() {
	go func() {
		s.ns.ConsumeNotifications()
	}()

	for client := range s.Register {
		s.mu.Lock()
		s.Clients[client] = true
		s.mu.Unlock()
		log.Println("Client registered with user_id: ", client.UserID)
		go func() {
			err := s.ns.PublishNotDeliveredNotificationsForUser(client.UserID)

			if err != nil {
				log.Printf("Failed to publish notifications for user: %v", err)
			}
		}()
	}
}

func (s *Server) HandleConnection(c *gin.Context) {
	// Upgrade the HTTP connection to a WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	userId := c.GetInt64("user_id")

	client := &Client{
		Conn:   conn,
		Send:   make(chan []byte),
		UserID: userId,
	}
	s.Register <- client
}
