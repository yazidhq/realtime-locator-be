package websocket

import "github.com/google/uuid"

type LocationMessage struct {
	UserID    uuid.UUID `json:"user_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}