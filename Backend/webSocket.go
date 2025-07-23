package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	clients    [2]*websocket.Conn
	clientsMux sync.Mutex
}

func (h *Hub) addClient(conn *websocket.Conn) (int, error) {
	h.clientsMux.Lock()
	defer h.clientsMux.Unlock()
	for i := 0; i < 2; i++ {
		if h.clients[i] == nil {
			h.clients[i] = conn
			return i, nil
		}
	}
	return -1, http.ErrServerClosed
}

func (h *Hub) removeClient(index int) {
	h.clientsMux.Lock()
	defer h.clientsMux.Unlock()
	h.clients[index] = nil
}

func (h *Hub) relayMessages(clientIndex int) {
	conn := h.clients[clientIndex]
	otherIndex := 1 - clientIndex

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			h.removeClient(clientIndex)
			return
		}

		if messageType == websocket.TextMessage {
			h.clientsMux.Lock()
			otherConn := h.clients[otherIndex]
			h.clientsMux.Unlock()

			if otherConn != nil {
				err = otherConn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("write error:", err)
				}
			}
		} else {
			log.Println("unexpected message type, ignoring")
		}
	}
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	clientIndex, err := hub.addClient(conn)
	if err != nil {
		log.Println("max clients reached")
		conn.WriteMessage(websocket.TextMessage, []byte("Server full"))
		conn.Close()
		return
	}

	log.Printf("Client %d connected", clientIndex)
	defer conn.Close()

	hub.relayMessages(clientIndex)
	log.Printf("Client %d disconnected", clientIndex)
}
