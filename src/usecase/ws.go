package usecase

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WsUsecase interface {
	TextMessage() error
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

func TextMessage(w http.ResponseWriter, r *http.Request) error {

	broadcast := func(messageType int, message []byte) {
		clientsMu.Lock()
		defer clientsMu.Unlock()
		for conn := range clients {
			if err := conn.WriteMessage(messageType, message); err != nil {
				fmt.Println("Broadcast error:", err)
				delete(clients, conn)
				conn.Close()
			}
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 新しいクライアントをマップに追加
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, conn) // クライアントが切断された場合はマップから削除
		clientsMu.Unlock()
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		fmt.Println("Received message:", string(message))
		broadcast(messageType, message)
	}
}
