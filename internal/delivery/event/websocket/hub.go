package websocket

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type Hub struct {
	Clients    map[uuid.UUID]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

var hubInstance *Hub

func GetHub() *Hub {
	if hubInstance == nil {
		hubInstance = &Hub{
			Clients:    make(map[uuid.UUID]*Client),
			Broadcast:  make(chan []byte),
			Register:   make(chan *Client),
			Unregister: make(chan *Client),
		}

		go hubInstance.Run()
	}

	return hubInstance
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if client == nil || client.UserID == uuid.Nil {
				log.Println("Register: invalid client or empty UserID")
				continue
			}

			if existing, ok := h.Clients[client.UserID]; ok {
				log.Println("User", client.UserID, "reconnected")
				existing.Conn.Close()
				close(existing.Send)
				delete(h.Clients, client.UserID)
			}

			h.Clients[client.UserID] = client
			log.Println("Client connected:", client.UserID)

		case client := <-h.Unregister:
			if client == nil || client.UserID == uuid.Nil {
				continue
			}

			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
				log.Println("Client disconnected:", client.UserID)
			}

		case message := <-h.Broadcast:
			var locMsg LocationMessage
			if err := json.Unmarshal(message, &locMsg); err != nil {
				log.Println("Invalid message:", err)
				continue
			}

			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					if client.Conn != nil {
						client.Conn.Close()
					}
					close(client.Send)
					delete(h.Clients, client.UserID)
				}
			}
		}
	}
}
