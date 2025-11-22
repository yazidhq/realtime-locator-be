package websocket

import (
	_repo "TeamTrackerBE/internal/domain/repository"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BroadcastMessage struct {
    SenderID uuid.UUID
    IsAdmin  bool
    Payload  []byte
}

type Hub struct {
	Clients      map[uuid.UUID]*Client
	Broadcast    chan BroadcastMessage
	Register     chan *Client
	Unregister   chan *Client
	flushTickers map[uuid.UUID]*time.Ticker
	locationRepo *_repo.LocationRepository
}

var hubInstance *Hub

func GetHub() *Hub {
	if hubInstance == nil {
		hubInstance = &Hub{
			Clients:    make(map[uuid.UUID]*Client),
			Broadcast:  make(chan BroadcastMessage),
			Register:   make(chan *Client),
			Unregister: make(chan *Client),
			flushTickers: make(map[uuid.UUID]*time.Ticker),
		}

		go hubInstance.Run()
	}

	return hubInstance
}

func LocationRepository(repo *_repo.LocationRepository) {
	h := GetHub()
	h.locationRepo = repo
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if client == nil || client.UserID == uuid.Nil {
				continue
			}

			if existing, ok := h.Clients[client.UserID]; ok {
				existing.Conn.Close()
				close(existing.Send)
				delete(h.Clients, client.UserID)
			}

			h.Clients[client.UserID] = client

		case client := <-h.Unregister:
			if client == nil || client.UserID == uuid.Nil {
				continue
			}

			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}

			h.StopFlushLoop(client.UserID, true)

		case message := <-h.Broadcast:
			if message.IsAdmin {
				var chatMsg ChatMessage
				if err := json.Unmarshal(message.Payload, &chatMsg); err != nil {
					continue
				}

				targetClient, ok := h.Clients[chatMsg.UserRecieverID]
				if !ok {
					continue
				}

				select {
				case targetClient.Send <- message.Payload:
				default:
					if targetClient.Conn != nil {
						targetClient.Conn.Close()
					}
					close(targetClient.Send)
					delete(h.Clients, targetClient.UserID)
					h.StopFlushLoop(targetClient.UserID, true)
				}

				continue
			}

			var locMsg LocationMessage
			if err := json.Unmarshal(message.Payload, &locMsg); err != nil {
				continue
			}

			locMsg.UserID = message.SenderID

			h.CacheLocation(locMsg)

			for _, client := range h.Clients {
				select {
				case client.Send <- message.Payload:
				default:
					if client.Conn != nil {
						client.Conn.Close()
					}

					close(client.Send)
					delete(h.Clients, client.UserID)
					h.StopFlushLoop(client.UserID, true)
				}
			}
		}
	}
}