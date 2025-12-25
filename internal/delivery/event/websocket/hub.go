package websocket

import (
	_repo "TeamTrackerBE/internal/domain/repository"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BroadcastMessage struct {
    SenderID uuid.UUID
    IsAdmin  bool
    Payload  []byte
}

type Hub struct {
	Clients      map[uuid.UUID]*Client
	Broadcast    chan BroadcastMessage
	Register     chan *Client
	Unregister   chan *Client
	flushTickers map[uuid.UUID]*time.Ticker
	locationRepo *_repo.LocationRepository
	Online 		 map[uuid.UUID]bool
	offlineGrace map[uuid.UUID]time.Time
	offlineTicker *time.Ticker
}

var hubInstance *Hub

func GetHub() *Hub {
	if hubInstance == nil {
		hubInstance = &Hub{
			Clients: make(map[uuid.UUID]*Client),
			Broadcast: make(chan BroadcastMessage),
			Register: make(chan *Client),
			Unregister: make(chan *Client),
			flushTickers: make(map[uuid.UUID]*time.Ticker),
			Online: make(map[uuid.UUID]bool),
			offlineGrace: make(map[uuid.UUID]time.Time),
			offlineTicker: time.NewTicker(1 * time.Second),
		}

		go hubInstance.Run()
	}

	return hubInstance
}

func LocationRepository(repo *_repo.LocationRepository) {
	h := GetHub()
	h.locationRepo = repo
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if client == nil || client.UserID == uuid.Nil {
				continue
			}

			if existing, ok := h.Clients[client.UserID]; ok {
				existing.Conn.Close()
				safeCloseSend(existing.Send)
				delete(h.Clients, client.UserID)
			}

			h.Clients[client.UserID] = client

			delete(h.offlineGrace, client.UserID)

			h.Online[client.UserID] = true
            status := UserStatusMessage{
                Type:   "user_status",
                UserID: client.UserID,
                Online: true,
            }

            if b, err := json.Marshal(status); err == nil {
                for _, c := range h.Clients {
                    select {
                    case c.Send <- b:
                    default:
						if c.Conn != nil {
							c.Conn.Close()
						}
						safeCloseSend(c.Send)
                        delete(h.Clients, c.UserID)
                        h.StopFlushLoop(c.UserID, true)
                    }
                }
            }

		case client := <-h.Unregister:
			if client == nil || client.UserID == uuid.Nil {
				continue
			}

			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				safeCloseSend(client.Send)
			}

			// Schedule offline after a grace period to survive quick refreshes
			uid := client.UserID
			h.offlineGrace[uid] = time.Now().Add(10 * time.Second)

		case message := <-h.Broadcast:
			if message.IsAdmin {
				var chatMsg ChatMessage
				if err := json.Unmarshal(message.Payload, &chatMsg); err != nil {
					continue
				}

				targetClient, ok := h.Clients[chatMsg.UserRecieverID]
				if !ok {
					continue
				}

				select {
				case targetClient.Send <- message.Payload:
				default:
					if targetClient.Conn != nil {
						targetClient.Conn.Close()
					}
					safeCloseSend(targetClient.Send)
					delete(h.Clients, targetClient.UserID)
					h.StopFlushLoop(targetClient.UserID, true)
				}

				continue
			}

			var locMsg LocationMessage
			if err := json.Unmarshal(message.Payload, &locMsg); err != nil {
				continue
			}

			locMsg.UserID = message.SenderID

			h.CacheLocation(locMsg)

			for _, client := range h.Clients {
				select {
				case client.Send <- message.Payload:
				default:
					if client.Conn != nil {
						client.Conn.Close()
					}

					safeCloseSend(client.Send)
					delete(h.Clients, client.UserID)
					h.StopFlushLoop(client.UserID, true)
				}
				// Periodically check for users whose grace period expired
				select {
				case <-h.offlineTicker.C:
					now := time.Now()
					for uid, until := range h.offlineGrace {
						if now.After(until) {
							// If user hasn't re-registered, mark offline and broadcast
							if _, still := h.Clients[uid]; !still {
								h.Online[uid] = false
								status := UserStatusMessage{
									Type:   "user_status",
									UserID: uid,
									Online: false,
								}
								if b, err := json.Marshal(status); err == nil {
									for _, c := range h.Clients {
										select {
										case c.Send <- b:
										default:
											if c.Conn != nil {
												c.Conn.Close()
											}
											safeCloseSend(c.Send)
											delete(h.Clients, c.UserID)
											h.StopFlushLoop(c.UserID, true)
										}
									}
								}
								h.StopFlushLoop(uid, true)
							}
							delete(h.offlineGrace, uid)
						}
					}
				default:
					// no-op
				}
			}
		}
	}
}

func safeCloseSend(ch chan []byte) {
	if ch == nil {
		return
	}
	defer func() {
		_ = recover()
	}()
	close(ch)
}

// ForceOffline immediately marks a user offline and broadcasts the status, closing any active connection and cancelling any grace period.
func (h *Hub) ForceOffline(userID uuid.UUID) {
	if userID == uuid.Nil {
		return
	}

	// Cancel any pending grace period
	delete(h.offlineGrace, userID)

	// Close and remove active client connection if exists
	if client, ok := h.Clients[userID]; ok {
		if client.Conn != nil {
			client.Conn.Close()
		}
		safeCloseSend(client.Send)
		delete(h.Clients, userID)
		h.StopFlushLoop(userID, true)
	}

	h.Online[userID] = false
	status := UserStatusMessage{
		Type:   "user_status",
		UserID: userID,
		Online: false,
	}

	if b, err := json.Marshal(status); err == nil {
		for _, c := range h.Clients {
			select {
			case c.Send <- b:
			default:
				if c.Conn != nil {
					c.Conn.Close()
				}
				safeCloseSend(c.Send)
				delete(h.Clients, c.UserID)
				h.StopFlushLoop(c.UserID, true)
			}
		}
	}
}