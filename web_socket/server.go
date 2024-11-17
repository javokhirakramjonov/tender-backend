package web_socket

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

var clients = make(map[int64]*websocket.Conn) // map[user_id]connection
var lock sync.Mutex

func SendNotification(userID int64, message []byte) error {
	lock.Lock()
	conn, exists := clients[userID]
	lock.Unlock()

	if exists {
		return conn.WriteMessage(websocket.TextMessage, message)
	}

	return errors.New("client is not online")
}
