package websocket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWs(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader websocket error: ", err)
		return
	}

	hub := GetHub()
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
