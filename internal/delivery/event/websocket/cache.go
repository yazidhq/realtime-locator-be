package websocket

import (
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/utils"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

func (h *Hub) CacheLocation(loc LocationMessage) {
	if !utils.IsCacheAvailable() {
		return
	}

	if loc.UserID == uuid.Nil {
		return
	}

	key := h.RedisKey(loc.UserID)
	payload, err := json.Marshal(loc)
	if err != nil {
		return
	}

	if err := utils.ListRPush(key, string(payload)); err != nil {
		log.Printf("redis rpush failed for %s: %v", key, err)
		return
	}

	if ttl, err := utils.TTL(key); err == nil && (ttl < 0) {
		_ = utils.Expire(key, 60*time.Second)
	}

	if _, ok := h.flushTickers[loc.UserID]; !ok {
		h.StartFlushLoop(loc.UserID)
	}
}

func (h *Hub) StartFlushLoop(userID uuid.UUID) {
	tk := time.NewTicker(30 * time.Second)
	h.flushTickers[userID] = tk

	go func(uid uuid.UUID, t *time.Ticker) {
		for range t.C {
			h.FlushUserLocations(uid)
		}
	}(userID, tk)
}

func (h *Hub) StopFlushLoop(userID uuid.UUID, flushNow bool) {
	if tk, ok := h.flushTickers[userID]; ok {
		tk.Stop()
		delete(h.flushTickers, userID)
		if flushNow {
			h.FlushUserLocations(userID)
		}
	}
}

func (h *Hub) FlushUserLocations(userID uuid.UUID) {
	if !utils.IsCacheAvailable() || h.locationRepo == nil || userID == uuid.Nil {
		return
	}
	key := h.RedisKey(userID)

	values, err := utils.ListLRange(key, 0, -1)
	if err != nil {
		log.Printf("redis lrange failed for %s: %v", key, err)
		return
	}
	if len(values) == 0 {
		return
	}

	records := make([]model.Location, 0, len(values))
	for _, raw := range values {
		var msg LocationMessage
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			continue
		}

		msg.UserID = userID
		records = append(records, model.Location{
			UserID:    msg.UserID,
			Latitude:  msg.Latitude,
			Longitude: msg.Longitude,
		})
	}
	if len(records) == 0 {
		_ = utils.DeleteCache(key)
		return
	}

	if err := h.locationRepo.BulkCreate(records); err != nil {
		log.Printf("bulk insert failed for user %s: %v", userID.String(), err)
		_ = utils.Expire(key, 60*time.Second)
		return
	}

	if err := utils.DeleteCache(key); err != nil {
		log.Printf("redis del failed for %s: %v", key, err)
	}
}

func (h *Hub) RedisKey(userID uuid.UUID) string {
	return "ws:locations:" + userID.String()
}
