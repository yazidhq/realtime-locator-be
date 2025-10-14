package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWs(h *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &Client{
		ID:     r.RemoteAddr,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		HubRef: h,
	}

	h.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
