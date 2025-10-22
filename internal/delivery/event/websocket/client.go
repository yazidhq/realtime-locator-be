package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	UserID uuid.UUID
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
			log.Println("read error: ", err)
			break
		}

		var loc LocationMessage
		if err := json.Unmarshal(msg, &loc); err != nil {
			log.Println("invalid message: ", err)
			continue
		}

		c.HubRef.Broadcast <- msg
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("write error: ", err)
			return
		}
	}
}
