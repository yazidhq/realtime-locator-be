package websocket

import "github.com/google/uuid"

type LocationMessage struct {
	Type      string    `json:"type"`
	UserID    uuid.UUID `json:"user_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

type ChatMessage struct {
	Type            string    `json:"type"`
	UserRecieverID  uuid.UUID `json:"user_receiver_id"`
	Message 		string 	  `json:"message"`
}

type UserStatusMessage struct {
	Type	string 	  `json:"type"`
    UserID	uuid.UUID `json:"user_id"`
	Online	bool 	  `json:"online"`
}

type OnlineUsersMessage struct {
	Type  string      `json:"type"`
	Users []uuid.UUID `json:"users"`
}