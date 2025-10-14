package websocket

import (
	"log"
)

// Hub mewakili satu grup
type Hub struct {
	GroupID    string
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// Semua grup disimpan di map global
var hubs = make(map[string]*Hub)

// GetHub mengambil atau membuat hub baru untuk group tertentu
func GetHub(groupID string) *Hub {
	if h, ok := hubs[groupID]; ok {
		return h
	}

	h := &Hub{
		GroupID:    groupID,
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}

	hubs[groupID] = h
	go h.Run()
	return h
}

// Jalankan loop untuk handle event register, unregister, broadcast
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if existing, ok := h.Clients[client.UserID]; ok {
				log.Println("User", client.UserID, "reconnected to group", h.GroupID)
				existing.Conn.Close()
			}
			h.Clients[client.UserID] = client
			log.Println("Client connected:", client.UserID, "to group", h.GroupID)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
				log.Println("Client disconnected:", client.UserID, "from group", h.GroupID)
			}

		case message := <-h.Broadcast:
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.UserID)
				}
			}
		}
	}
}
