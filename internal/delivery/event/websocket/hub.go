package websocket

import (
	"TeamTrackerBE/internal/domain/model"
	"encoding/json"
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
func GetHub(groupID string, group *model.Group) *Hub {
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
	go h.Run(group)
	// fmt.Println("Message", h.Broadcast)
	return h
}

// Jalankan loop untuk handle event register, unregister, broadcast
func (h *Hub) Run(group *model.Group) {
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
			var locMsg LocationMessage
			if err := json.Unmarshal(message, &locMsg); err != nil {
				log.Println("Invalid message:", err)
				continue
			}

			if len(group.RadiusArea) == 0 {
				log.Printf("group have not radius yet: %s", group.ID)
				continue
			}

			var area RadiusArea
			if err := json.Unmarshal(group.RadiusArea, &area); err != nil {
				log.Printf("failed to parse radius area: %v", err)
				continue
			}

			distance := Distance(locMsg.Latitude, locMsg.Longitude, area.CenterLat, area.CenterLon)

			warn := CrossedRadiusMessage{
				Type: "warning",
				UserID: locMsg.UserID,
				Message: "the user has crossed the radius boundary",
			}

			warnMsg, _ := json.Marshal(warn)

			for _, client := range h.Clients {
				if distance > area.Radius {
					select {
					case client.Send <- warnMsg:
					default:
						close(client.Send)
						delete(h.Clients, client.UserID)
					}
				} else {
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
}
