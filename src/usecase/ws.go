package usecase

import (
	"encoding/json"
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

type Coordinate struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func Location(w http.ResponseWriter, r *http.Request) error {

	broadcast := func(message []byte) {
		clientsMu.Lock()
		defer clientsMu.Unlock()
		for conn := range clients {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
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
		_, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		var loc Coordinate
		if err := json.Unmarshal(message, &loc); err != nil {
			fmt.Println("Error decoding message:", err)
			continue
		}

		fmt.Printf("Received location: %+v\n", loc)

		encodedMessage, err := json.Marshal(loc)
		if err != nil {
			fmt.Println("Error encoding message:", err)
			continue
		}

		broadcast(encodedMessage)
	}
}
