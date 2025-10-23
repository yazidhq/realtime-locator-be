package websocket

import (
	"TeamTrackerBE/internal/domain/model"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWs(w http.ResponseWriter, r *http.Request, user *model.User) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	
	hub := GetHub()
	client := &Client{
		ID: r.RemoteAddr,
		UserID: user.ID,
		Role: string(user.Role),
		Conn: conn,
		Send: make(chan []byte, 256),
		HubRef: hub,
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
