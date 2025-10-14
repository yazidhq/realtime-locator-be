package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client mewakili satu koneksi websocket
type Client struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
	HubRef *Hub
}

// ReadPump menerima pesan dari user
func (c *Client) ReadPump() {
	defer func() {
		c.HubRef.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		// Validasi JSON (optional)
		var loc LocationMessage
		if err := json.Unmarshal(msg, &loc); err != nil {
			log.Println("invalid message:", err)
			continue
		}

		// Broadcast ke semua user di grup yang sama
		c.HubRef.Broadcast <- msg
	}
}

// WritePump mengirim pesan ke user
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("write error:", err)
			return
		}
	}
}
