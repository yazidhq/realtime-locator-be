package websocket

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	UserID uuid.UUID
	Role   string
	Conn   *websocket.Conn
	Send   chan []byte
	HubRef *Hub
}

func (c *Client) ReadPump() {
	defer func() {
		c.HubRef.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var meta struct{
			Type string `json:"type"`
		}
		
		if err := json.Unmarshal(msg, &meta); err != nil {
			continue
		}

		isAdmin := strings.EqualFold(c.Role, "admin")

		switch meta.Type {
		case "user_location":
			if isAdmin {
				continue
			}

			c.HubRef.Broadcast <- BroadcastMessage{
				SenderID: c.UserID,
				IsAdmin:  false,
				Payload:  msg,
			}

		case "user_message":
			if !isAdmin {
				continue
			}

			c.HubRef.Broadcast <- BroadcastMessage{
				SenderID: c.UserID,
				IsAdmin:  true,
				Payload:  msg,
			}

		default:
			continue
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
