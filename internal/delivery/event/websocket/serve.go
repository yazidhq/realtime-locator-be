package websocket

import (
	"TeamTrackerBE/internal/domain/model"
	"encoding/json"
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

	// Send list of all online users to the new client
	onlineUsers := hub.GetAllOnlineUsers()
	onlineMsg := OnlineUsersMessage{
		Type: "online_users_list",
		Users: onlineUsers,
	}

	if b, err := json.Marshal(onlineMsg); err == nil {
		select {
		case client.Send <- b:
		default:
			if client.Conn != nil {
				client.Conn.Close()
			}
			safeCloseSend(client.Send)
		}
	}

	go client.WritePump()
	go client.ReadPump()
}
