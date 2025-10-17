package websocket

import (
	"TeamTrackerBE/internal/domain/model"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWs meng-handle koneksi baru ke WebSocket
func ServeWs(w http.ResponseWriter, r *http.Request, groupID, userID string, group *model.Group) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader websocket error: ", err)
		return
	}

	hub := GetHub(groupID, group)

	client := &Client{
		ID:     r.RemoteAddr,
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		HubRef: hub,
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
