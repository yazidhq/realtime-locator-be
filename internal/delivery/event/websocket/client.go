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

// Heartbeat timings
const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func (c *Client) ReadPump() {
	defer func() {
		c.HubRef.Unregister <- c
		c.Conn.Close()
	}()

	// Configure read deadline and pong handler to detect dead clients
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			// Send ping periodically to keep the connection alive and detect dead peers
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
